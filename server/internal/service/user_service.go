package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
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
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		slog.Error("UserService.GetProfile: failed to get user",
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}
	slog.Debug("UserService.GetProfile: profile retrieved",
		"user_id", userID,
	)
	return user, nil
}

// UpdateProfile updates a user's profile.
func (s *UserService) UpdateProfile(ctx context.Context, userID int64, name string) error {
	if err := validator.ValidateName(name); err != nil {
		slog.Warn("UserService.UpdateProfile: invalid name",
			"user_id", userID,
			"name", name,
			"error", err,
		)
		return err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		slog.Error("UserService.UpdateProfile: failed to get user",
			"user_id", userID,
			"error", err,
		)
		return err
	}

	user.Name = name
	if err := s.userRepo.Update(ctx, user); err != nil {
		slog.Error("UserService.UpdateProfile: failed to update user",
			"user_id", userID,
			"name", name,
			"error", err,
		)
		return err
	}

	slog.Info("UserService.UpdateProfile: profile updated",
		"user_id", userID,
		"name", name,
	)
	return nil
}

// UploadAvatar uploads a user's avatar.
func (s *UserService) UploadAvatar(ctx context.Context, userID int64, filename string, reader io.Reader) (string, error) {
	if err := validator.ValidateFileExtension(filename, validator.ImageExtensions); err != nil {
		slog.Warn("UserService.UploadAvatar: invalid file extension",
			"user_id", userID,
			"filename", filename,
			"error", err,
		)
		return "", err
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		slog.Error("UserService.UploadAvatar: failed to get user",
			"user_id", userID,
			"error", err,
		)
		return "", err
	}

	// Delete old avatar if exists
	if user.AvatarPath != "" {
		if err := s.fileRepo.Delete(ctx, user.AvatarPath); err != nil {
			slog.Warn("UserService.UploadAvatar: failed to delete old avatar",
				"user_id", userID,
				"avatar_path", user.AvatarPath,
				"error", err,
			)
		}
	}

	ext := filepath.Ext(filename)
	avatarFilename := fmt.Sprintf("avatar_%d%s", userID, ext)
	avatarPath := filepath.Join(s.avatarDir, avatarFilename)

	if err := s.fileRepo.Save(ctx, avatarPath, reader); err != nil {
		slog.Error("UserService.UploadAvatar: failed to save avatar file",
			"user_id", userID,
			"avatar_path", avatarPath,
			"error", err,
		)
		return "", err
	}

	user.AvatarPath = avatarPath
	if err := s.userRepo.Update(ctx, user); err != nil {
		slog.Error("UserService.UploadAvatar: failed to update user with avatar path",
			"user_id", userID,
			"avatar_path", avatarPath,
			"error", err,
		)
		return "", err
	}

	slog.Info("UserService.UploadAvatar: avatar uploaded",
		"user_id", userID,
		"avatar_path", avatarPath,
	)
	return avatarPath, nil
}

// DeleteAvatar deletes a user's avatar.
func (s *UserService) DeleteAvatar(ctx context.Context, userID int64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		slog.Error("UserService.DeleteAvatar: failed to get user",
			"user_id", userID,
			"error", err,
		)
		return err
	}

	if user.AvatarPath != "" {
		if err := s.fileRepo.Delete(ctx, user.AvatarPath); err != nil {
			slog.Warn("UserService.DeleteAvatar: failed to delete avatar file",
				"user_id", userID,
				"avatar_path", user.AvatarPath,
				"error", err,
			)
		}
		user.AvatarPath = ""
		if err := s.userRepo.Update(ctx, user); err != nil {
			slog.Error("UserService.DeleteAvatar: failed to update user after avatar deletion",
				"user_id", userID,
				"error", err,
			)
			return err
		}
		slog.Info("UserService.DeleteAvatar: avatar deleted",
			"user_id", userID,
		)
		return nil
	}

	slog.Debug("UserService.DeleteAvatar: no avatar to delete",
		"user_id", userID,
	)
	return nil
}

// ServeAvatar returns the avatar path for a user.
func (s *UserService) ServeAvatar(ctx context.Context, userID int64) (string, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		slog.Error("UserService.ServeAvatar: failed to get user",
			"user_id", userID,
			"error", err,
		)
		return "", err
	}
	if user.AvatarPath == "" {
		slog.Warn("UserService.ServeAvatar: no avatar found",
			"user_id", userID,
		)
		return "", domain.ErrNotFound
	}
	slog.Debug("UserService.ServeAvatar: avatar path retrieved",
		"user_id", userID,
		"avatar_path", user.AvatarPath,
	)
	return user.AvatarPath, nil
}

// ListUsers lists users with optional filtering.
func (s *UserService) ListUsers(ctx context.Context, filter domain.UserFilter) ([]domain.User, int, error) {
	users, total, err := s.userRepo.List(ctx, filter)
	if err != nil {
		slog.Error("UserService.ListUsers: failed to list users",
			"filter", filter,
			"error", err,
		)
		return nil, 0, err
	}
	slog.Debug("UserService.ListUsers: users listed",
		"count", len(users),
		"total", total,
	)
	return users, total, nil
}

// DeleteUser deletes a user by ID.
func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		slog.Error("UserService.DeleteUser: failed to delete user",
			"user_id", id,
			"error", err,
		)
		return err
	}
	slog.Info("UserService.DeleteUser: user deleted",
		"user_id", id,
	)
	return nil
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
	works, total, err := s.workRepo.List(ctx, filter)
	if err != nil {
		slog.Error("GalleryService.GetWorks: failed to list works",
			"filter", filter,
			"error", err,
		)
		return nil, 0, err
	}
	slog.Debug("GalleryService.GetWorks: works listed",
		"count", len(works),
		"total", total,
	)
	return works, total, nil
}

// GetWorkByID retrieves a work by ID.
func (s *GalleryService) GetWorkByID(ctx context.Context, id int64) (*domain.Work, error) {
	work, err := s.workRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("GalleryService.GetWorkByID: failed to get work",
			"work_id", id,
			"error", err,
		)
		return nil, err
	}
	slog.Debug("GalleryService.GetWorkByID: work retrieved",
		"work_id", id,
		"title", work.Title,
	)
	return work, nil
}

// CreateWork creates a new work.
func (s *GalleryService) CreateWork(ctx context.Context, work *domain.Work, filename string, reader io.Reader) error {
	if work.Title == "" {
		slog.Warn("GalleryService.CreateWork: empty title")
		return domain.ErrInvalidInput
	}

	// Save image
	imagePath := filepath.Join(s.imagesDir, fmt.Sprintf("work_%d%s", time.Now().UnixNano(), filepath.Ext(filename)))
	if err := s.fileRepo.Save(ctx, imagePath, reader); err != nil {
		slog.Error("GalleryService.CreateWork: failed to save image",
			"title", work.Title,
			"image_path", imagePath,
			"error", err,
		)
		return err
	}
	work.ImagePath = imagePath

	if err := s.workRepo.Create(ctx, work); err != nil {
		slog.Error("GalleryService.CreateWork: failed to create work",
			"title", work.Title,
			"image_path", imagePath,
			"error", err,
		)
		return err
	}

	if len(work.MaterialIDs) > 0 {
		if err := s.workRepo.SetMaterials(ctx, work.ID, work.MaterialIDs); err != nil {
			slog.Error("GalleryService.CreateWork: failed to set materials",
				"work_id", work.ID,
				"error", err,
			)
			return err
		}
	}

	if len(work.BaseIDs) > 0 {
		if err := s.workRepo.SetBases(ctx, work.ID, work.BaseIDs); err != nil {
			slog.Error("GalleryService.CreateWork: failed to set bases",
				"work_id", work.ID,
				"error", err,
			)
			return err
		}
	}

	slog.Info("GalleryService.CreateWork: work created",
		"work_id", work.ID,
		"title", work.Title,
	)
	return nil
}

// UpdateWork updates a work.
func (s *GalleryService) UpdateWork(ctx context.Context, work *domain.Work, filename string, reader io.Reader) error {
	existing, err := s.workRepo.GetByID(ctx, work.ID)
	if err != nil {
		slog.Error("GalleryService.UpdateWork: failed to get existing work",
			"work_id", work.ID,
			"error", err,
		)
		return err
	}

	if reader != nil {
		// Delete old image
		if existing.ImagePath != "" {
			if err := s.fileRepo.Delete(ctx, existing.ImagePath); err != nil {
				slog.Warn("GalleryService.UpdateWork: failed to delete old image",
					"work_id", work.ID,
					"image_path", existing.ImagePath,
					"error", err,
				)
			}
		}
		imagePath := filepath.Join(s.imagesDir, fmt.Sprintf("work_%d%s", time.Now().UnixNano(), filepath.Ext(filename)))
		if err := s.fileRepo.Save(ctx, imagePath, reader); err != nil {
			slog.Error("GalleryService.UpdateWork: failed to save new image",
				"work_id", work.ID,
				"image_path", imagePath,
				"error", err,
			)
			return err
		}
		work.ImagePath = imagePath
	} else {
		work.ImagePath = existing.ImagePath
	}

	if err := s.workRepo.Update(ctx, work); err != nil {
		slog.Error("GalleryService.UpdateWork: failed to update work",
			"work_id", work.ID,
			"title", work.Title,
			"error", err,
		)
		return err
	}

	if err := s.workRepo.SetMaterials(ctx, work.ID, work.MaterialIDs); err != nil {
		slog.Error("GalleryService.UpdateWork: failed to set materials",
			"work_id", work.ID,
			"error", err,
		)
		return err
	}

	if err := s.workRepo.SetBases(ctx, work.ID, work.BaseIDs); err != nil {
		slog.Error("GalleryService.UpdateWork: failed to set bases",
			"work_id", work.ID,
			"error", err,
		)
		return err
	}

	slog.Info("GalleryService.UpdateWork: work updated",
		"work_id", work.ID,
		"title", work.Title,
	)
	return nil
}

// DeleteWork deletes a work.
func (s *GalleryService) DeleteWork(ctx context.Context, id int64) error {
	work, err := s.workRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("GalleryService.DeleteWork: failed to get work",
			"work_id", id,
			"error", err,
		)
		return err
	}

	if work.ImagePath != "" {
		if err := s.fileRepo.Delete(ctx, work.ImagePath); err != nil {
			slog.Warn("GalleryService.DeleteWork: failed to delete work image",
				"work_id", id,
				"image_path", work.ImagePath,
				"error", err,
			)
		}
	}

	if err := s.workRepo.Delete(ctx, id); err != nil {
		slog.Error("GalleryService.DeleteWork: failed to delete work",
			"work_id", id,
			"title", work.Title,
			"error", err,
		)
		return err
	}

	slog.Info("GalleryService.DeleteWork: work deleted",
		"work_id", id,
		"title", work.Title,
	)
	return nil
}

// GetSales retrieves sales with filtering.
func (s *GalleryService) GetSales(ctx context.Context, filter domain.SaleFilter) ([]domain.Sale, int, error) {
	sales, total, err := s.saleRepo.List(ctx, filter)
	if err != nil {
		slog.Error("GalleryService.GetSales: failed to list sales",
			"filter", filter,
			"error", err,
		)
		return nil, 0, err
	}
	slog.Debug("GalleryService.GetSales: sales listed",
		"count", len(sales),
		"total", total,
	)
	return sales, total, nil
}

// GetSaleByID retrieves a sale by ID.
func (s *GalleryService) GetSaleByID(ctx context.Context, id int64) (*domain.Sale, error) {
	sale, err := s.saleRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("GalleryService.GetSaleByID: failed to get sale",
			"sale_id", id,
			"error", err,
		)
		return nil, err
	}
	slog.Debug("GalleryService.GetSaleByID: sale retrieved",
		"sale_id", id,
		"title", sale.Title,
	)
	return sale, nil
}

// CreateSale creates a new sale.
func (s *GalleryService) CreateSale(ctx context.Context, sale *domain.Sale, filename string, reader io.Reader) error {
	if sale.Title == "" {
		slog.Warn("GalleryService.CreateSale: empty title")
		return domain.ErrInvalidInput
	}

	imagePath := filepath.Join(s.imagesDir, fmt.Sprintf("sale_%d%s", time.Now().UnixNano(), filepath.Ext(filename)))
	if err := s.fileRepo.Save(ctx, imagePath, reader); err != nil {
		slog.Error("GalleryService.CreateSale: failed to save image",
			"title", sale.Title,
			"image_path", imagePath,
			"error", err,
		)
		return err
	}
	sale.ImagePath = imagePath

	if err := s.saleRepo.Create(ctx, sale); err != nil {
		slog.Error("GalleryService.CreateSale: failed to create sale",
			"title", sale.Title,
			"image_path", imagePath,
			"error", err,
		)
		return err
	}

	if len(sale.MaterialIDs) > 0 {
		if err := s.saleRepo.SetMaterials(ctx, sale.ID, sale.MaterialIDs); err != nil {
			slog.Error("GalleryService.CreateSale: failed to set materials",
				"sale_id", sale.ID,
				"error", err,
			)
			return err
		}
	}

	if len(sale.BaseIDs) > 0 {
		if err := s.saleRepo.SetBases(ctx, sale.ID, sale.BaseIDs); err != nil {
			slog.Error("GalleryService.CreateSale: failed to set bases",
				"sale_id", sale.ID,
				"error", err,
			)
			return err
		}
	}

	slog.Info("GalleryService.CreateSale: sale created",
		"sale_id", sale.ID,
		"title", sale.Title,
	)
	return nil
}

// UpdateSale updates a sale.
func (s *GalleryService) UpdateSale(ctx context.Context, sale *domain.Sale, filename string, reader io.Reader) error {
	existing, err := s.saleRepo.GetByID(ctx, sale.ID)
	if err != nil {
		slog.Error("GalleryService.UpdateSale: failed to get existing sale",
			"sale_id", sale.ID,
			"error", err,
		)
		return err
	}

	if reader != nil {
		if existing.ImagePath != "" {
			if err := s.fileRepo.Delete(ctx, existing.ImagePath); err != nil {
				slog.Warn("GalleryService.UpdateSale: failed to delete old sale image",
					"sale_id", sale.ID,
					"image_path", existing.ImagePath,
					"error", err,
				)
			}
		}
		imagePath := filepath.Join(s.imagesDir, fmt.Sprintf("sale_%d%s", time.Now().UnixNano(), filepath.Ext(filename)))
		if err := s.fileRepo.Save(ctx, imagePath, reader); err != nil {
			slog.Error("GalleryService.UpdateSale: failed to save new image",
				"sale_id", sale.ID,
				"image_path", imagePath,
				"error", err,
			)
			return err
		}
		sale.ImagePath = imagePath
	} else {
		sale.ImagePath = existing.ImagePath
	}

	if err := s.saleRepo.Update(ctx, sale); err != nil {
		slog.Error("GalleryService.UpdateSale: failed to update sale",
			"sale_id", sale.ID,
			"title", sale.Title,
			"error", err,
		)
		return err
	}

	if err := s.saleRepo.SetMaterials(ctx, sale.ID, sale.MaterialIDs); err != nil {
		slog.Error("GalleryService.UpdateSale: failed to set materials",
			"sale_id", sale.ID,
			"error", err,
		)
		return err
	}

	if err := s.saleRepo.SetBases(ctx, sale.ID, sale.BaseIDs); err != nil {
		slog.Error("GalleryService.UpdateSale: failed to set bases",
			"sale_id", sale.ID,
			"error", err,
		)
		return err
	}

	slog.Info("GalleryService.UpdateSale: sale updated",
		"sale_id", sale.ID,
		"title", sale.Title,
	)
	return nil
}

// DeleteSale deletes a sale.
func (s *GalleryService) DeleteSale(ctx context.Context, id int64) error {
	sale, err := s.saleRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("GalleryService.DeleteSale: failed to get sale",
			"sale_id", id,
			"error", err,
		)
		return err
	}

	if sale.ImagePath != "" {
		if err := s.fileRepo.Delete(ctx, sale.ImagePath); err != nil {
			slog.Warn("GalleryService.DeleteSale: failed to delete sale image",
				"sale_id", id,
				"image_path", sale.ImagePath,
				"error", err,
			)
		}
	}

	if err := s.saleRepo.Delete(ctx, id); err != nil {
		slog.Error("GalleryService.DeleteSale: failed to delete sale",
			"sale_id", id,
			"title", sale.Title,
			"error", err,
		)
		return err
	}

	slog.Info("GalleryService.DeleteSale: sale deleted",
		"sale_id", id,
		"title", sale.Title,
	)
	return nil
}
