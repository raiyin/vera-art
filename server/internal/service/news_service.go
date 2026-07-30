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

func (s *NewsService) CreateNews(ctx context.Context, news *domain.News, imgBackFile, imgBackfullFile *domain.UploadedFile, imageFiles, videoFiles []domain.UploadedFile) error {
	if news.TitleRu == "" || news.TitleEn == "" {
		return domain.ErrInvalidInput
	}

	if news.DateTime == "" {
		news.DateTime = time.Now().Format("2006-01-02")
	}

	news.ID = s.generateID(news.DateTime)
	news.Dir = s.generateDir(news.DateTime)

	if imgBackFile != nil {
		if err := s.saveFile(ctx, imgBackFile, news.Dir); err != nil {
			return err
		}
		news.ImgBack = imgBackFile.Filename
	}

	if imgBackfullFile != nil {
		if err := s.saveFile(ctx, imgBackfullFile, news.Dir); err != nil {
			return err
		}
		news.ImgBackfull = imgBackfullFile.Filename
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

func (s *NewsService) UpdateNews(ctx context.Context, news *domain.News, imgBackFile, imgBackfullFile *domain.UploadedFile, imageFiles, videoFiles []domain.UploadedFile) error {
	existing, err := s.newsRepo.GetByID(ctx, news.ID)
	if err != nil {
		slog.Error("NewsService.UpdateNews: not found", "news_id", news.ID, "error", err)
		return err
	}

	if imgBackFile != nil {
		if existing.ImgBack != "" {
			_ = s.fileRepo.Delete(ctx, filepath.Join(existing.Dir, existing.ImgBack))
		}
		if err := s.saveFile(ctx, imgBackFile, existing.Dir); err != nil {
			return err
		}
		news.ImgBack = imgBackFile.Filename
	} else {
		news.ImgBack = existing.ImgBack
	}

	if imgBackfullFile != nil {
		if existing.ImgBackfull != "" {
			_ = s.fileRepo.Delete(ctx, filepath.Join(existing.Dir, existing.ImgBackfull))
		}
		if err := s.saveFile(ctx, imgBackfullFile, existing.Dir); err != nil {
			return err
		}
		news.ImgBackfull = imgBackfullFile.Filename
	} else {
		news.ImgBackfull = existing.ImgBackfull
	}

	if len(imageFiles) > 0 {
		for _, oldImg := range existing.Images {
			_ = s.fileRepo.Delete(ctx, filepath.Join(existing.Dir, oldImg))
		}
		var names []string
		for _, f := range imageFiles {
			if err := s.saveFile(ctx, &f, existing.Dir); err != nil {
				continue
			}
			names = append(names, f.Filename)
		}
		news.Images = names
	} else {
		news.Images = existing.Images
	}

	if len(videoFiles) > 0 {
		for _, oldVideo := range existing.Videos {
			_ = s.fileRepo.RemoveDir(ctx, filepath.Join(existing.Dir, "videos", oldVideo))
		}
		var names []string
		for _, f := range videoFiles {
			dirname := strings.TrimSuffix(f.Filename, filepath.Ext(f.Filename))
			videoSubDir := filepath.Join(existing.Dir, "videos", dirname)
			if err := s.saveFile(ctx, &f, videoSubDir); err != nil {
				continue
			}
			names = append(names, dirname)
		}
		news.Videos = names
	} else {
		news.Videos = existing.Videos
	}

	news.Dir = existing.Dir

	if err := s.newsRepo.Update(ctx, news); err != nil {
		slog.Error("NewsService.UpdateNews: failed to update", "news_id", news.ID, "error", err)
		return err
	}

	slog.Info("NewsService.UpdateNews: updated", "news_id", news.ID)
	return nil
}

func (s *NewsService) DeleteNews(ctx context.Context, id string) error {
	news, err := s.newsRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("NewsService.DeleteNews: not found", "news_id", id, "error", err)
		return err
	}

	if news.Dir != "" {
		_ = s.fileRepo.RemoveDir(ctx, news.Dir)
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
	return s.fileRepo.Save(ctx, filepath.Join(subDir, file.Filename), file.Reader)
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
