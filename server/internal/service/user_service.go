package service

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/port"
	"github.com/raiyin/artserver/pkg/validator"
)

// UserService implements port.UserService.
type UserService struct {
	userRepo  port.UserRepository
	fileRepo  port.FileRepository
	avatarDir string
}

// NewUserService creates a new UserService.
func NewUserService(userRepo port.UserRepository, fileRepo port.FileRepository, avatarDir string) *UserService {
	return &UserService{
		userRepo:  userRepo,
		fileRepo:  fileRepo,
		avatarDir: avatarDir,
	}
}

// GetProfile retrieves a user's profile.
func (s *UserService) GetProfile(ctx context.Context, userID int64) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

// UpdateProfile updates a user's profile.
func (s *UserService) UpdateProfile(ctx context.Context, userID int64, name string) error {
	if err := validator.ValidateName(name); err != nil {
		return err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Name = name
	return s.userRepo.Update(ctx, user)
}

// UploadAvatar uploads a user's avatar.
func (s *UserService) UploadAvatar(ctx context.Context, userID int64, filename string, reader io.Reader) (string, error) {
	if err := validator.ValidateFileExtension(filename, validator.ImageExtensions); err != nil {
		return "", err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	// Delete old avatar if exists
	if user.AvatarPath != "" {
		_ = s.fileRepo.Delete(ctx, user.AvatarPath)
	}

	ext := filepath.Ext(filename)
	avatarFilename := fmt.Sprintf("avatar_%d%s", userID, ext)
	avatarPath := filepath.Join(s.avatarDir, avatarFilename)

	if err := s.fileRepo.Save(ctx, avatarPath, reader); err != nil {
		return "", err
	}

	user.AvatarPath = avatarPath
	if err := s.userRepo.Update(ctx, user); err != nil {
		return "", err
	}

	return avatarPath, nil
}

// DeleteAvatar deletes a user's avatar.
func (s *UserService) DeleteAvatar(ctx context.Context, userID int64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.AvatarPath != "" {
		_ = s.fileRepo.Delete(ctx, user.AvatarPath)
		user.AvatarPath = ""
		return s.userRepo.Update(ctx, user)
	}

	return nil
}

// ServeAvatar returns the avatar path for a user.
func (s *UserService) ServeAvatar(ctx context.Context, userID int64) (string, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}
	if user.AvatarPath == "" {
		return "", domain.ErrNotFound
	}
	return user.AvatarPath, nil
}

// ListUsers lists users with optional filtering.
func (s *UserService) ListUsers(ctx context.Context, filter domain.UserFilter) ([]domain.User, int, error) {
	return s.userRepo.List(ctx, filter)
}

// DeleteUser deletes a user by ID.
func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	return s.userRepo.Delete(ctx, id)
}

// GalleryService implements port.GalleryService.
type GalleryService struct {
	workRepo  port.WorkRepository
	saleRepo  port.SaleRepository
	fileRepo  port.FileRepository
	imagesDir string
}

// NewGalleryService creates a new GalleryService.
func NewGalleryService(workRepo port.WorkRepository, saleRepo port.SaleRepository, fileRepo port.FileRepository, imagesDir string) *GalleryService {
	return &GalleryService{
		workRepo:  workRepo,
		saleRepo:  saleRepo,
		fileRepo:  fileRepo,
		imagesDir: imagesDir,
	}
}

// GetWorks retrieves works with filtering.
func (s *GalleryService) GetWorks(ctx context.Context, filter domain.WorkFilter) ([]domain.Work, int, error) {
	return s.workRepo.List(ctx, filter)
}

// GetWorkByID retrieves a work by ID.
func (s *GalleryService) GetWorkByID(ctx context.Context, id int64) (*domain.Work, error) {
	return s.workRepo.GetByID(ctx, id)
}

// CreateWork creates a new work.
func (s *GalleryService) CreateWork(ctx context.Context, work *domain.Work, filename string, reader io.Reader) error {
	if work.Title == "" {
		return domain.ErrInvalidInput
	}

	// Save image
	imagePath := filepath.Join(s.imagesDir, fmt.Sprintf("work_%d%s", time.Now().UnixNano(), filepath.Ext(filename)))
	if err := s.fileRepo.Save(ctx, imagePath, reader); err != nil {
		return err
	}
	work.ImagePath = imagePath

	return s.workRepo.Create(ctx, work)
}

// UpdateWork updates a work.
func (s *GalleryService) UpdateWork(ctx context.Context, work *domain.Work, filename string, reader io.Reader) error {
	existing, err := s.workRepo.GetByID(ctx, work.ID)
	if err != nil {
		return err
	}

	if reader != nil {
		// Delete old image
		if existing.ImagePath != "" {
			_ = s.fileRepo.Delete(ctx, existing.ImagePath)
		}
		imagePath := filepath.Join(s.imagesDir, fmt.Sprintf("work_%d%s", time.Now().UnixNano(), filepath.Ext(filename)))
		if err := s.fileRepo.Save(ctx, imagePath, reader); err != nil {
			return err
		}
		work.ImagePath = imagePath
	} else {
		work.ImagePath = existing.ImagePath
	}

	return s.workRepo.Update(ctx, work)
}

// DeleteWork deletes a work.
func (s *GalleryService) DeleteWork(ctx context.Context, id int64) error {
	work, err := s.workRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if work.ImagePath != "" {
		_ = s.fileRepo.Delete(ctx, work.ImagePath)
	}

	return s.workRepo.Delete(ctx, id)
}

// GetSales retrieves sales with filtering.
func (s *GalleryService) GetSales(ctx context.Context, filter domain.SaleFilter) ([]domain.Sale, int, error) {
	return s.saleRepo.List(ctx, filter)
}

// GetSaleByID retrieves a sale by ID.
func (s *GalleryService) GetSaleByID(ctx context.Context, id int64) (*domain.Sale, error) {
	return s.saleRepo.GetByID(ctx, id)
}

// CreateSale creates a new sale.
func (s *GalleryService) CreateSale(ctx context.Context, sale *domain.Sale, filename string, reader io.Reader) error {
	if sale.Title == "" {
		return domain.ErrInvalidInput
	}

	imagePath := filepath.Join(s.imagesDir, fmt.Sprintf("sale_%d%s", time.Now().UnixNano(), filepath.Ext(filename)))
	if err := s.fileRepo.Save(ctx, imagePath, reader); err != nil {
		return err
	}
	sale.ImagePath = imagePath

	return s.saleRepo.Create(ctx, sale)
}

// UpdateSale updates a sale.
func (s *GalleryService) UpdateSale(ctx context.Context, sale *domain.Sale, filename string, reader io.Reader) error {
	existing, err := s.saleRepo.GetByID(ctx, sale.ID)
	if err != nil {
		return err
	}

	if reader != nil {
		if existing.ImagePath != "" {
			_ = s.fileRepo.Delete(ctx, existing.ImagePath)
		}
		imagePath := filepath.Join(s.imagesDir, fmt.Sprintf("sale_%d%s", time.Now().UnixNano(), filepath.Ext(filename)))
		if err := s.fileRepo.Save(ctx, imagePath, reader); err != nil {
			return err
		}
		sale.ImagePath = imagePath
	} else {
		sale.ImagePath = existing.ImagePath
	}

	return s.saleRepo.Update(ctx, sale)
}

// DeleteSale deletes a sale.
func (s *GalleryService) DeleteSale(ctx context.Context, id int64) error {
	sale, err := s.saleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if sale.ImagePath != "" {
		_ = s.fileRepo.Delete(ctx, sale.ImagePath)
	}

	return s.saleRepo.Delete(ctx, id)
}
