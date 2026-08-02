package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/dto"
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

func (s *UserService) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("UserService.GetUserByID: failed to get user",
			"user_id", id,
			"error", err,
		)
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateUserRole(ctx context.Context, id int64, role string) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("UserService.UpdateUserRole: failed to get user",
			"user_id", id,
			"error", err,
		)
		return err
	}
	user.Role = role
	if err := s.userRepo.Update(ctx, user); err != nil {
		slog.Error("UserService.UpdateUserRole: failed to update role",
			"user_id", id,
			"error", err,
		)
		return err
	}
	slog.Info("UserService.UpdateUserRole: role updated",
		"user_id", id,
		"role", role,
	)
	return nil
}

func (s *UserService) ToggleUserBlock(ctx context.Context, id int64) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("UserService.ToggleUserBlock: failed to get user",
			"user_id", id,
			"error", err,
		)
		return err
	}
	user.Blocked = !user.Blocked
	if err := s.userRepo.Update(ctx, user); err != nil {
		slog.Error("UserService.ToggleUserBlock: failed to toggle block",
			"user_id", id,
			"error", err,
		)
		return err
	}
	slog.Info("UserService.ToggleUserBlock: block toggled",
		"user_id", id,
		"blocked", user.Blocked,
	)
	return nil
}

// GalleryService implements port.GalleryService.
type GalleryService struct {
	workRepo    port.WorkRepository
	saleRepo    port.SaleRepository
	fileRepo    port.FileRepository
	imagesDir   string
	relWorksDir string
}

// NewGalleryService creates a new GalleryService.
func NewGalleryService(workRepo port.WorkRepository, saleRepo port.SaleRepository, fileRepo port.FileRepository, imagesDir, relWorksDir string) *GalleryService {
	return &GalleryService{
		workRepo:    workRepo,
		saleRepo:    saleRepo,
		fileRepo:    fileRepo,
		imagesDir:   imagesDir,
		relWorksDir: relWorksDir,
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
		"name_ru", work.NameRu,
	)
	return work, nil
}

// CreateWork creates a new work.
func (s *GalleryService) CreateWork(ctx context.Context, work *domain.Work, files []domain.UploadedFile) error {
	if work.NameRu == "" {
		slog.Warn("GalleryService.CreateWork: empty name_ru")
		return domain.ErrInvalidInput
	}

	if work.WorkPath == "" {
		work.WorkPath = s.buildWorkPath(work.StrID)
	}

	if len(files) > 0 {
		imageNames, err := s.saveWorkImages(ctx, work.WorkPath, "", files)
		if err != nil {
			slog.Error("GalleryService.CreateWork: failed to save image",
				"name_ru", work.NameRu,
				"error", err,
			)
			return err
		}
		work.Images = imageNames
	}

	if err := s.workRepo.Create(ctx, work); err != nil {
		slog.Error("GalleryService.CreateWork: failed to create work",
			"name_ru", work.NameRu,
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

	slog.Info("GalleryService.CreateWork: work created",
		"work_id", work.ID,
		"name_ru", work.NameRu,
	)
	return nil
}

// UpdateWork updates a work.
func (s *GalleryService) UpdateWork(ctx context.Context, work *domain.Work, files []domain.UploadedFile) error {
	existing, err := s.workRepo.GetByID(ctx, work.ID)
	if err != nil {
		slog.Error("GalleryService.UpdateWork: failed to get existing work",
			"work_id", work.ID,
			"error", err,
		)
		return err
	}

	workPath := existing.WorkPath
	if workPath == "" {
		workPath = s.buildWorkPath(work.StrID)
	}
	work.WorkPath = workPath

	// keptImages is the list of image filenames the user left in the preview.
	// When it is not provided, fall back to the existing list for backward compatibility.
	keptImages := work.Images
	if keptImages == "" {
		keptImages = existing.Images
	}

	// Delete files that were removed from the preview.
	s.removeDeletedImages(ctx, workPath, existing.Images, keptImages)

	// Save newly uploaded images and append them to the kept list.
	if len(files) > 0 {
		imageNames, err := s.saveWorkImages(ctx, workPath, keptImages, files)
		if err != nil {
			slog.Error("GalleryService.UpdateWork: failed to save image",
				"work_id", work.ID,
				"error", err,
			)
			return err
		}
		if keptImages != "" {
			work.Images = keptImages + ";" + imageNames
		} else {
			work.Images = imageNames
		}
	} else {
		work.Images = keptImages
	}

	if err := s.workRepo.Update(ctx, work); err != nil {
		slog.Error("GalleryService.UpdateWork: failed to update work",
			"work_id", work.ID,
			"name_ru", work.NameRu,
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

	slog.Info("GalleryService.UpdateWork: work updated",
		"work_id", work.ID,
		"name_ru", work.NameRu,
	)
	return nil
}

// removeDeletedImages deletes image files that exist on disk but are no longer present in the kept list.
func (s *GalleryService) removeDeletedImages(ctx context.Context, workPath, existingImages, keptImages string) {
	existing := dto.SplitImages(existingImages)
	kept := dto.SplitImages(keptImages)

	keptSet := make(map[string]struct{}, len(kept))
	for _, img := range kept {
		keptSet[img] = struct{}{}
	}

	for _, img := range existing {
		if img == "" {
			continue
		}
		if _, ok := keptSet[img]; ok {
			continue
		}
		relPath := strings.TrimLeft(filepath.Join(s.relWorksDir, workPath, img), "/")
		fullPath := filepath.Join(s.imagesDir, relPath)
		if err := s.fileRepo.Delete(ctx, fullPath); err != nil {
			slog.Warn("GalleryService.UpdateWork: failed to delete removed image",
				"image", img,
				"error", err,
			)
			continue
		}
		slog.Debug("GalleryService.UpdateWork: removed image",
			"image", img,
		)
	}
}

// buildWorkPath returns a directory path for a work based on its str_id.
func (s *GalleryService) buildWorkPath(strID string) string {
	dir := strings.TrimSpace(strID)
	if dir == "" {
		dir = fmt.Sprintf("work_%d", time.Now().UnixNano())
	}
	return strings.TrimSuffix(dir, "/") + "/"
}

// saveWorkImages saves uploaded images into a work's directory and returns the semicolon-separated image names.
func (s *GalleryService) saveWorkImages(ctx context.Context, workPath, currentImages string, files []domain.UploadedFile) (string, error) {
	current := currentImages
	var saved []string
	for _, f := range files {
		imageName, err := s.saveWorkImage(ctx, workPath, current, f.Filename, f.Reader)
		if err != nil {
			return "", err
		}
		saved = append(saved, imageName)
		if current == "" {
			current = imageName
		} else {
			current = current + ";" + imageName
		}
	}
	return strings.Join(saved, ";"), nil
}

// saveWorkImage saves an uploaded image into a work's directory and returns the image filename.
func (s *GalleryService) saveWorkImage(ctx context.Context, workPath, currentImages, filename string, reader io.Reader) (string, error) {
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".jpg"
	}
	imageName := s.nextWorkImageName(currentImages, ext)

	relPath := strings.TrimLeft(filepath.Join(s.relWorksDir, workPath, imageName), "/")
	fullPath := filepath.Join(s.imagesDir, relPath)
	if err := s.fileRepo.Save(ctx, fullPath, reader); err != nil {
		return "", err
	}
	return imageName, nil
}

// nextWorkImageName returns the next sequential image filename for a work directory.
func (s *GalleryService) nextWorkImageName(currentImages, ext string) string {
	count := 0
	maxNum := 0
	if currentImages != "" {
		imgs := strings.Split(currentImages, ";")
		count = len(imgs)
		for _, img := range imgs {
			n, err := strconv.Atoi(strings.TrimSuffix(img, filepath.Ext(img)))
			if err == nil && n > maxNum {
				maxNum = n
			}
		}
	}
	if maxNum < count {
		maxNum = count
	}
	return fmt.Sprintf("%d%s", maxNum+1, ext)
}

// DeleteWork deletes a work.
func (s *GalleryService) DeleteWork(ctx context.Context, id int64) error {
	if err := s.workRepo.Delete(ctx, id); err != nil {
		slog.Error("GalleryService.DeleteWork: failed to delete work",
			"work_id", id,
			"error", err,
		)
		return err
	}

	slog.Info("GalleryService.DeleteWork: work deleted",
		"work_id", id,
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
		"name_ru", sale.NameRu,
	)
	return sale, nil
}

// CreateSale creates a new sale.
func (s *GalleryService) CreateSale(ctx context.Context, sale *domain.Sale, filename string, reader io.Reader) error {
	if sale.NameRu == "" {
		slog.Warn("GalleryService.CreateSale: empty name_ru")
		return domain.ErrInvalidInput
	}

	imagePath := filepath.Join(s.imagesDir, fmt.Sprintf("sale_%d%s", time.Now().UnixNano(), filepath.Ext(filename)))
	if err := s.fileRepo.Save(ctx, imagePath, reader); err != nil {
		slog.Error("GalleryService.CreateSale: failed to save image",
			"name_ru", sale.NameRu,
			"image_path", imagePath,
			"error", err,
		)
		return err
	}
	sale.ImagePath = imagePath

	if err := s.saleRepo.Create(ctx, sale); err != nil {
		slog.Error("GalleryService.CreateSale: failed to create sale",
			"name_ru", sale.NameRu,
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
		"name_ru", sale.NameRu,
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
			"name_ru", sale.NameRu,
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
		"name_ru", sale.NameRu,
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
			"name_ru", sale.NameRu,
			"error", err,
		)
		return err
	}

	slog.Info("GalleryService.DeleteSale: sale deleted",
		"sale_id", id,
		"name_ru", sale.NameRu,
	)
	return nil
}

func (s *GalleryService) BulkDeleteWorks(ctx context.Context, ids []int64) error {
	for _, id := range ids {
		if err := s.DeleteWork(ctx, id); err != nil {
			slog.Error("GalleryService.BulkDeleteWorks: failed to delete work",
				"work_id", id,
				"error", err,
			)
			return err
		}
	}
	slog.Info("GalleryService.BulkDeleteWorks: works deleted",
		"count", len(ids),
	)
	return nil
}

func (s *GalleryService) BulkDeleteSales(ctx context.Context, ids []int64) error {
	for _, id := range ids {
		if err := s.DeleteSale(ctx, id); err != nil {
			slog.Error("GalleryService.BulkDeleteSales: failed to delete sale",
				"sale_id", id,
				"error", err,
			)
			return err
		}
	}
	slog.Info("GalleryService.BulkDeleteSales: sales deleted",
		"count", len(ids),
	)
	return nil
}
