package service

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/port"
	"github.com/raiyin/artserver/pkg/validator"
)

// NewsService implements port.NewsService.
type NewsService struct {
	newsRepo port.NewsRepository
	fileRepo port.FileRepository
	newsDir  string
}

// NewNewsService creates a new NewsService.
func NewNewsService(newsRepo port.NewsRepository, fileRepo port.FileRepository, newsDir string) *NewsService {
	return &NewsService{
		newsRepo: newsRepo,
		fileRepo: fileRepo,
		newsDir:  newsDir,
	}
}

// GetNews retrieves news entries with filtering.
func (s *NewsService) GetNews(ctx context.Context, filter domain.NewsFilter) ([]domain.News, int, error) {
	return s.newsRepo.List(ctx, filter)
}

// GetNewsByID retrieves a news entry by ID.
func (s *NewsService) GetNewsByID(ctx context.Context, id int64) (*domain.News, error) {
	return s.newsRepo.GetByID(ctx, id)
}

// CreateNews creates a new news entry.
func (s *NewsService) CreateNews(ctx context.Context, news *domain.News, imageFile *domain.UploadedFile, videoFile *domain.UploadedFile) error {
	if news.Title == "" {
		return domain.ErrInvalidInput
	}

	newsDir := filepath.Join(s.newsDir, fmt.Sprintf("news_%d", time.Now().UnixNano()))

	if imageFile != nil {
		if err := validator.ValidateFileExtension(imageFile.Filename, validator.ImageExtensions); err != nil {
			return err
		}
		imagePath := filepath.Join(newsDir, imageFile.Filename)
		if err := s.fileRepo.Save(ctx, imagePath, imageFile.Reader); err != nil {
			return err
		}
		news.ImagePath = imagePath
	}

	if videoFile != nil {
		if err := validator.ValidateFileExtension(videoFile.Filename, validator.VideoExtensions); err != nil {
			return err
		}
		videoPath := filepath.Join(newsDir, videoFile.Filename)
		if err := s.fileRepo.Save(ctx, videoPath, videoFile.Reader); err != nil {
			return err
		}
		news.VideoPath = videoPath
	}

	return s.newsRepo.Create(ctx, news)
}

// UpdateNews updates a news entry.
func (s *NewsService) UpdateNews(ctx context.Context, news *domain.News, imageFile *domain.UploadedFile, videoFile *domain.UploadedFile) error {
	existing, err := s.newsRepo.GetByID(ctx, news.ID)
	if err != nil {
		return err
	}

	newsDir := filepath.Dir(existing.ImagePath)
	if newsDir == "." {
		newsDir = filepath.Join(s.newsDir, fmt.Sprintf("news_%d", news.ID))
	}

	if imageFile != nil {
		if existing.ImagePath != "" {
			if err := s.fileRepo.Delete(ctx, existing.ImagePath); err != nil {
				slog.Warn("NewsService.UpdateNews: failed to delete old image",
					"news_id", news.ID,
					"image_path", existing.ImagePath,
					"error", err,
				)
			}
		}
		imagePath := filepath.Join(newsDir, imageFile.Filename)
		if err := s.fileRepo.Save(ctx, imagePath, imageFile.Reader); err != nil {
			return err
		}
		news.ImagePath = imagePath
	} else {
		news.ImagePath = existing.ImagePath
	}

	if videoFile != nil {
		if existing.VideoPath != "" {
			if err := s.fileRepo.Delete(ctx, existing.VideoPath); err != nil {
				slog.Warn("NewsService.UpdateNews: failed to delete old video",
					"news_id", news.ID,
					"video_path", existing.VideoPath,
					"error", err,
				)
			}
		}
		videoPath := filepath.Join(newsDir, videoFile.Filename)
		if err := s.fileRepo.Save(ctx, videoPath, videoFile.Reader); err != nil {
			return err
		}
		news.VideoPath = videoPath
	} else {
		news.VideoPath = existing.VideoPath
	}

	return s.newsRepo.Update(ctx, news)
}

// DeleteNews deletes a news entry.
func (s *NewsService) DeleteNews(ctx context.Context, id int64) error {
	news, err := s.newsRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if news.ImagePath != "" {
		if err := s.fileRepo.Delete(ctx, news.ImagePath); err != nil {
			slog.Warn("NewsService.DeleteNews: failed to delete news image",
				"news_id", id,
				"image_path", news.ImagePath,
				"error", err,
			)
		}
	}
	if news.VideoPath != "" {
		if err := s.fileRepo.Delete(ctx, news.VideoPath); err != nil {
			slog.Warn("NewsService.DeleteNews: failed to delete news video",
				"news_id", id,
				"video_path", news.VideoPath,
				"error", err,
			)
		}
	}

	return s.newsRepo.Delete(ctx, id)
}

// Ensure interface compliance.
var _ port.NewsService = (*NewsService)(nil)
