package service

import (
	"context"
	"log/slog"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/port"
)

// TagService implements port.TagService.
type TagService struct {
	tagRepo port.TagRepository
}

// NewTagService creates a new TagService.
func NewTagService(tagRepo port.TagRepository) *TagService {
	return &TagService{tagRepo: tagRepo}
}

func (s *TagService) GetTags(ctx context.Context) ([]domain.Tag, error) {
	tags, err := s.tagRepo.List(ctx)
	if err != nil {
		slog.Error("TagService.GetTags: failed to list tags",
			"error", err,
		)
		return nil, err
	}
	slog.Debug("TagService.GetTags: tags listed",
		"count", len(tags),
	)
	return tags, nil
}

func (s *TagService) CreateTag(ctx context.Context, tag *domain.Tag) error {
	if tag.NameRu == "" && tag.NameEn == "" {
		slog.Warn("TagService.CreateTag: empty name",
			"tag", tag,
		)
		return domain.ErrInvalidInput
	}
	if err := s.tagRepo.Create(ctx, tag); err != nil {
		slog.Error("TagService.CreateTag: failed to create tag",
			"tag_name_ru", tag.NameRu,
			"tag_name_en", tag.NameEn,
			"error", err,
		)
		return err
	}
	slog.Info("TagService.CreateTag: tag created",
		"tag_id", tag.ID,
		"tag_name_ru", tag.NameRu,
		"tag_name_en", tag.NameEn,
	)
	return nil
}

func (s *TagService) DeleteTag(ctx context.Context, id int64) error {
	if err := s.tagRepo.Delete(ctx, id); err != nil {
		slog.Error("TagService.DeleteTag: failed to delete tag",
			"tag_id", id,
			"error", err,
		)
		return err
	}
	slog.Info("TagService.DeleteTag: tag deleted",
		"tag_id", id,
	)
	return nil
}

var _ port.TagService = (*TagService)(nil)

// MaterialService implements port.MaterialService.
type MaterialService struct {
	materialRepo port.MaterialRepository
}

// NewMaterialService creates a new MaterialService.
func NewMaterialService(materialRepo port.MaterialRepository) *MaterialService {
	return &MaterialService{materialRepo: materialRepo}
}

func (s *MaterialService) GetMaterials(ctx context.Context) ([]domain.Material, error) {
	materials, err := s.materialRepo.List(ctx)
	if err != nil {
		slog.Error("MaterialService.GetMaterials: failed to list materials",
			"error", err,
		)
		return nil, err
	}
	slog.Debug("MaterialService.GetMaterials: materials listed",
		"count", len(materials),
	)
	return materials, nil
}

func (s *MaterialService) CreateMaterial(ctx context.Context, material *domain.Material) error {
	if material.NameRu == "" && material.NameEn == "" {
		slog.Warn("MaterialService.CreateMaterial: empty name",
			"material", material,
		)
		return domain.ErrInvalidInput
	}
	if err := s.materialRepo.Create(ctx, material); err != nil {
		slog.Error("MaterialService.CreateMaterial: failed to create material",
			"material_name_ru", material.NameRu,
			"material_name_en", material.NameEn,
			"error", err,
		)
		return err
	}
	slog.Info("MaterialService.CreateMaterial: material created",
		"material_id", material.ID,
		"material_name_ru", material.NameRu,
		"material_name_en", material.NameEn,
	)
	return nil
}

func (s *MaterialService) DeleteMaterial(ctx context.Context, id int64) error {
	if err := s.materialRepo.Delete(ctx, id); err != nil {
		slog.Error("MaterialService.DeleteMaterial: failed to delete material",
			"material_id", id,
			"error", err,
		)
		return err
	}
	slog.Info("MaterialService.DeleteMaterial: material deleted",
		"material_id", id,
	)
	return nil
}

var _ port.MaterialService = (*MaterialService)(nil)

// BaseService implements port.BaseService.
type BaseService struct {
	baseRepo port.BaseRepository
}

// NewBaseService creates a new BaseService.
func NewBaseService(baseRepo port.BaseRepository) *BaseService {
	return &BaseService{baseRepo: baseRepo}
}

func (s *BaseService) GetBases(ctx context.Context) ([]domain.Base, error) {
	bases, err := s.baseRepo.List(ctx)
	if err != nil {
		slog.Error("BaseService.GetBases: failed to list bases",
			"error", err,
		)
		return nil, err
	}
	slog.Debug("BaseService.GetBases: bases listed",
		"count", len(bases),
	)
	return bases, nil
}

func (s *BaseService) CreateBase(ctx context.Context, base *domain.Base) error {
	if base.NameRu == "" && base.NameEn == "" {
		slog.Warn("BaseService.CreateBase: empty name",
			"base", base,
		)
		return domain.ErrInvalidInput
	}
	if err := s.baseRepo.Create(ctx, base); err != nil {
		slog.Error("BaseService.CreateBase: failed to create base",
			"base_name_ru", base.NameRu,
			"base_name_en", base.NameEn,
			"error", err,
		)
		return err
	}
	slog.Info("BaseService.CreateBase: base created",
		"base_id", base.ID,
		"base_name_ru", base.NameRu,
		"base_name_en", base.NameEn,
	)
	return nil
}

func (s *BaseService) DeleteBase(ctx context.Context, id int64) error {
	if err := s.baseRepo.Delete(ctx, id); err != nil {
		slog.Error("BaseService.DeleteBase: failed to delete base",
			"base_id", id,
			"error", err,
		)
		return err
	}
	slog.Info("BaseService.DeleteBase: base deleted",
		"base_id", id,
	)
	return nil
}

var _ port.BaseService = (*BaseService)(nil)

// ConsentService implements port.ConsentService.
type ConsentService struct {
	consentRepo port.ConsentRepository
}

// NewConsentService creates a new ConsentService.
func NewConsentService(consentRepo port.ConsentRepository) *ConsentService {
	return &ConsentService{consentRepo: consentRepo}
}

func (s *ConsentService) RecordConsent(ctx context.Context, userID int64, consentType string, granted bool, ipAddress string) error {
	consent := &domain.UserConsent{
		UserID:      userID,
		ConsentType: consentType,
		Granted:     granted,
		IPAddress:   ipAddress,
	}
	if err := s.consentRepo.Create(ctx, consent); err != nil {
		slog.Error("ConsentService.RecordConsent: failed to record consent",
			"user_id", userID,
			"consent_type", consentType,
			"granted", granted,
			"error", err,
		)
		return err
	}
	slog.Info("ConsentService.RecordConsent: consent recorded",
		"user_id", userID,
		"consent_type", consentType,
		"granted", granted,
	)
	return nil
}

func (s *ConsentService) GetUserConsents(ctx context.Context, userID int64) ([]domain.UserConsent, error) {
	consents, err := s.consentRepo.ListByUser(ctx, userID)
	if err != nil {
		slog.Error("ConsentService.GetUserConsents: failed to list consents",
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}
	slog.Debug("ConsentService.GetUserConsents: consents listed",
		"user_id", userID,
		"count", len(consents),
	)
	return consents, nil
}

var _ port.ConsentService = (*ConsentService)(nil)

// MasterClassService implements port.MasterClassService.
type MasterClassService struct {
	mcRepo port.MasterClassRepository
}

// NewMasterClassService creates a new MasterClassService.
func NewMasterClassService(mcRepo port.MasterClassRepository) *MasterClassService {
	return &MasterClassService{mcRepo: mcRepo}
}

func (s *MasterClassService) GetMasterClasses(ctx context.Context) ([]domain.MasterClass, int, error) {
	classes, total, err := s.mcRepo.List(ctx)
	if err != nil {
		slog.Error("MasterClassService.GetMasterClasses: failed to list master classes",
			"error", err,
		)
		return nil, 0, err
	}
	slog.Debug("MasterClassService.GetMasterClasses: master classes listed",
		"count", len(classes),
		"total", total,
	)
	return classes, total, nil
}

func (s *MasterClassService) GetMasterClassByID(ctx context.Context, id int64) (*domain.MasterClass, error) {
	mc, err := s.mcRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("MasterClassService.GetMasterClassByID: failed to get master class",
			"master_class_id", id,
			"error", err,
		)
		return nil, err
	}
	slog.Debug("MasterClassService.GetMasterClassByID: master class retrieved",
		"master_class_id", id,
		"title", mc.Title,
	)
	return mc, nil
}

func (s *MasterClassService) CreateMasterClass(ctx context.Context, mc *domain.MasterClass) error {
	if mc.Title == "" {
		slog.Warn("MasterClassService.CreateMasterClass: empty title")
		return domain.ErrInvalidInput
	}
	if err := s.mcRepo.Create(ctx, mc); err != nil {
		slog.Error("MasterClassService.CreateMasterClass: failed to create master class",
			"title", mc.Title,
			"error", err,
		)
		return err
	}
	slog.Info("MasterClassService.CreateMasterClass: master class created",
		"master_class_id", mc.ID,
		"title", mc.Title,
	)
	return nil
}

func (s *MasterClassService) UpdateMasterClass(ctx context.Context, mc *domain.MasterClass) error {
	if err := s.mcRepo.Update(ctx, mc); err != nil {
		slog.Error("MasterClassService.UpdateMasterClass: failed to update master class",
			"master_class_id", mc.ID,
			"title", mc.Title,
			"error", err,
		)
		return err
	}
	slog.Info("MasterClassService.UpdateMasterClass: master class updated",
		"master_class_id", mc.ID,
		"title", mc.Title,
	)
	return nil
}

func (s *MasterClassService) DeleteMasterClass(ctx context.Context, id int64) error {
	if err := s.mcRepo.Delete(ctx, id); err != nil {
		slog.Error("MasterClassService.DeleteMasterClass: failed to delete master class",
			"master_class_id", id,
			"error", err,
		)
		return err
	}
	slog.Info("MasterClassService.DeleteMasterClass: master class deleted",
		"master_class_id", id,
	)
	return nil
}

var _ port.MasterClassService = (*MasterClassService)(nil)

// AdminService implements port.AdminService.
type AdminService struct {
	userRepo     port.UserRepository
	productRepo  port.ProductRepository
	paymentRepo  port.PaymentRepository
	purchaseRepo port.PurchaseRepository
}

// NewAdminService creates a new AdminService.
func NewAdminService(
	userRepo port.UserRepository,
	productRepo port.ProductRepository,
	paymentRepo port.PaymentRepository,
	purchaseRepo port.PurchaseRepository,
) *AdminService {
	return &AdminService{
		userRepo:     userRepo,
		productRepo:  productRepo,
		paymentRepo:  paymentRepo,
		purchaseRepo: purchaseRepo,
	}
}

func (s *AdminService) GetDashboardStats(ctx context.Context) (map[string]interface{}, error) {
	users, _, err := s.userRepo.List(ctx, domain.UserFilter{})
	if err != nil {
		slog.Error("AdminService.GetDashboardStats: failed to list users",
			"error", err,
		)
		return nil, err
	}

	products, _, err := s.productRepo.List(ctx, domain.ProductFilter{})
	if err != nil {
		slog.Error("AdminService.GetDashboardStats: failed to list products",
			"error", err,
		)
		return nil, err
	}

	payments, _, err := s.paymentRepo.List(ctx)
	if err != nil {
		slog.Error("AdminService.GetDashboardStats: failed to list payments",
			"error", err,
		)
		return nil, err
	}

	purchases, _, err := s.purchaseRepo.List(ctx)
	if err != nil {
		slog.Error("AdminService.GetDashboardStats: failed to list purchases",
			"error", err,
		)
		return nil, err
	}

	stats := map[string]interface{}{
		"total_users":     len(users),
		"total_products":  len(products),
		"total_payments":  len(payments),
		"total_purchases": len(purchases),
	}

	slog.Info("AdminService.GetDashboardStats: stats retrieved",
		"total_users", len(users),
		"total_products", len(products),
		"total_payments", len(payments),
		"total_purchases", len(purchases),
	)

	return stats, nil
}

var _ port.AdminService = (*AdminService)(nil)
