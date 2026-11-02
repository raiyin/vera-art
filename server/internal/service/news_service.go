package service

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/port"
)

type NewsService struct {
	newsRepo port.NewsRepository
	fileRepo port.FileRepository
	newsDir  string
}

func NewNewsService(newsRepo port.NewsRepository, fileRepo port.FileRepository, newsDir string) *NewsService {
	return &NewsService{
		newsRepo: newsRepo,
		fileRepo: fileRepo,
		newsDir:  newsDir,
	}
}

func (s *NewsService) GetNews(ctx context.Context, filter domain.NewsFilter) ([]domain.News, int, error) {
	news, total, err := s.newsRepo.List(ctx, filter)
	if err != nil {
		slog.Error("NewsService.GetNews: failed", "filter", filter, "error", err)
		return nil, 0, err
	}
	return news, total, nil
}

func (s *NewsService) GetNewsByID(ctx context.Context, id string) (*domain.News, error) {
	news, err := s.newsRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("NewsService.GetNewsByID: failed", "news_id", id, "error", err)
		return nil, err
	}
	return news, nil
}

func (s *NewsService) CreateNews(ctx context.Context, news *domain.News, mainImageFile *domain.UploadedFile, imageFiles, videoFiles []domain.UploadedFile) error {
	if news.TitleRu == "" || news.TitleEn == "" {
		return domain.ErrInvalidInput
	}

	if news.DateTime == "" {
		news.DateTime = time.Now().Format("2006-01-02")
	}

	news.ID = s.generateID(news.DateTime)
	news.Dir = s.generateDir(news.DateTime)

	if mainImageFile != nil {
		if err := s.saveFile(ctx, mainImageFile, news.Dir); err != nil {
			return err
		}
		news.MainImage = mainImageFile.Filename
	}

	var savedImageNames []string
	for _, f := range imageFiles {
		if err := s.saveFile(ctx, &f, news.Dir); err != nil {
			continue
		}
		savedImageNames = append(savedImageNames, f.Filename)
	}
	if len(savedImageNames) > 0 {
		news.Images = savedImageNames
	}

	var savedVideoNames []string
	for _, f := range videoFiles {
		dirname := strings.TrimSuffix(f.Filename, filepath.Ext(f.Filename))
		videoSubDir := filepath.Join(news.Dir, "videos", dirname)
		if err := s.saveFile(ctx, &f, videoSubDir); err != nil {
			continue
		}
		savedVideoNames = append(savedVideoNames, dirname)
	}
	if len(savedVideoNames) > 0 {
		news.Videos = savedVideoNames
	}

	if err := s.newsRepo.Create(ctx, news); err != nil {
		slog.Error("NewsService.CreateNews: failed to create", "error", err)
		return err
	}

	slog.Info("NewsService.CreateNews: created", "news_id", news.ID, "title_ru", news.TitleRu)
	return nil
}

func (s *NewsService) UpdateNews(ctx context.Context, news *domain.News, mainImageFile *domain.UploadedFile, imageFiles, videoFiles []domain.UploadedFile) error {
	existing, err := s.newsRepo.GetByID(ctx, news.ID)
	if err != nil {
		slog.Error("NewsService.UpdateNews: not found", "news_id", news.ID, "error", err)
		return err
	}

	news.Dir = existing.Dir

	if mainImageFile != nil {
		if existing.MainImage != "" {
			_ = s.fileRepo.Delete(ctx, s.fullPath(existing.Dir, existing.MainImage))
		}
		if err := s.saveFile(ctx, mainImageFile, existing.Dir); err != nil {
			return err
		}
		news.MainImage = mainImageFile.Filename
	} else if news.MainImage == "" {
		// The user removed the image, delete the stored file as well.
		if existing.MainImage != "" {
			_ = s.fileRepo.Delete(ctx, s.fullPath(existing.Dir, existing.MainImage))
		}
		news.MainImage = ""
	} else {
		news.MainImage = existing.MainImage
	}

	// Only the files the user left on the preview must stay on the server.
	if news.Images == nil {
		news.Images = existing.Images
	} else {
		news.Images = s.reconcileFiles(ctx, existing.Dir, existing.Images, news.Images, imageFiles, false)
	}

	if news.Videos == nil {
		news.Videos = existing.Videos
	} else {
		news.Videos = s.reconcileFiles(ctx, existing.Dir, existing.Videos, news.Videos, videoFiles, true)
	}

	if err := s.newsRepo.Update(ctx, news); err != nil {
		slog.Error("NewsService.UpdateNews: failed to update", "news_id", news.ID, "error", err)
		return err
	}

	slog.Info("NewsService.UpdateNews: updated", "news_id", news.ID)
	return nil
}

// reconcileFiles keeps only the files the user left on the preview: it deletes
// previously stored files that are no longer desired and saves newly uploaded ones.
func (s *NewsService) reconcileFiles(ctx context.Context, dir string, existingNames, desiredNames []string, newFiles []domain.UploadedFile, isVideo bool) []string {
	keep := make(map[string]bool)
	result := make([]string, 0, len(desiredNames)+len(newFiles))
	for _, name := range desiredNames {
		name = strings.TrimSpace(name)
		if name == "" || keep[name] {
			continue
		}
		keep[name] = true
		result = append(result, name)
	}

	for _, oldName := range existingNames {
		if keep[oldName] {
			continue
		}
		if isVideo {
			_ = s.fileRepo.RemoveDir(ctx, s.fullPath(dir, "videos", oldName))
		} else {
			_ = s.fileRepo.Delete(ctx, s.fullPath(dir, oldName))
		}
	}

	for _, f := range newFiles {
		var storedName string
		if isVideo {
			storedName = strings.TrimSuffix(f.Filename, filepath.Ext(f.Filename))
			videoSubDir := filepath.Join(dir, "videos", storedName)
			if err := s.saveFile(ctx, &f, videoSubDir); err != nil {
				continue
			}
		} else {
			if err := s.saveFile(ctx, &f, dir); err != nil {
				continue
			}
			storedName = f.Filename
		}
		if storedName == "" || keep[storedName] {
			continue
		}
		keep[storedName] = true
		result = append(result, storedName)
	}

	return result
}

func (s *NewsService) DeleteNews(ctx context.Context, id string) error {
	news, err := s.newsRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("NewsService.DeleteNews: not found", "news_id", id, "error", err)
		return err
	}

	if news.Dir != "" {
		_ = s.fileRepo.RemoveDir(ctx, s.fullPath(news.Dir))
	}

	if err := s.newsRepo.Delete(ctx, id); err != nil {
		slog.Error("NewsService.DeleteNews: failed to delete", "news_id", id, "error", err)
		return err
	}

	slog.Info("NewsService.DeleteNews: deleted", "news_id", id)
	return nil
}

func (s *NewsService) BulkDeleteNews(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := s.DeleteNews(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (s *NewsService) saveFile(ctx context.Context, file *domain.UploadedFile, subDir string) error {
	return s.fileRepo.Save(ctx, s.fullPath(subDir, file.Filename), file.Reader)
}

// fullPath resolves URL-style content paths (e.g. "/content/news/2025/01/01/")
// against the absolute news directory.
func (s *NewsService) fullPath(elem ...string) string {
	rel := strings.TrimLeft(filepath.Join(elem...), "/")
	return filepath.Join(s.newsDir, rel)
}

func (s *NewsService) generateID(datetime string) string {
	parts := strings.Split(datetime, "T")
	datePart := parts[0]
	id := strings.ReplaceAll(datePart, "-", "")
	if id == "" {
		id = time.Now().Format("20060102")
	}
	return id
}

func (s *NewsService) generateDir(datetime string) string {
	parts := strings.Split(datetime, "T")
	datePart := parts[0]
	date, err := time.Parse("2006-01-02", datePart)
	if err != nil {
		date = time.Now()
	}
	return fmt.Sprintf("/content/news/%s/%s/%s/",
		date.Format("2006"),
		date.Format("01"),
		date.Format("02"),
	)
}
