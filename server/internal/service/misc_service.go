package service

import (
	"context"

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
	return s.tagRepo.List(ctx)
}

func (s *TagService) CreateTag(ctx context.Context, tag *domain.Tag) error {
	if tag.NameRu == "" && tag.NameEn == "" {
		return domain.ErrInvalidInput
	}
	return s.tagRepo.Create(ctx, tag)
}

func (s *TagService) DeleteTag(ctx context.Context, id int64) error {
	return s.tagRepo.Delete(ctx, id)
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
	return s.materialRepo.List(ctx)
}

func (s *MaterialService) CreateMaterial(ctx context.Context, material *domain.Material) error {
	if material.NameRu == "" && material.NameEn == "" {
		return domain.ErrInvalidInput
	}
	return s.materialRepo.Create(ctx, material)
}

func (s *MaterialService) DeleteMaterial(ctx context.Context, id int64) error {
	return s.materialRepo.Delete(ctx, id)
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
	return s.baseRepo.List(ctx)
}

func (s *BaseService) CreateBase(ctx context.Context, base *domain.Base) error {
	if base.NameRu == "" && base.NameEn == "" {
		return domain.ErrInvalidInput
	}
	return s.baseRepo.Create(ctx, base)
}

func (s *BaseService) DeleteBase(ctx context.Context, id int64) error {
	return s.baseRepo.Delete(ctx, id)
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
	return s.consentRepo.Create(ctx, consent)
}

func (s *ConsentService) GetUserConsents(ctx context.Context, userID int64) ([]domain.UserConsent, error) {
	return s.consentRepo.ListByUser(ctx, userID)
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
	return s.mcRepo.List(ctx)
}

func (s *MasterClassService) GetMasterClassByID(ctx context.Context, id int64) (*domain.MasterClass, error) {
	return s.mcRepo.GetByID(ctx, id)
}

func (s *MasterClassService) CreateMasterClass(ctx context.Context, mc *domain.MasterClass) error {
	if mc.Title == "" {
		return domain.ErrInvalidInput
	}
	return s.mcRepo.Create(ctx, mc)
}

func (s *MasterClassService) UpdateMasterClass(ctx context.Context, mc *domain.MasterClass) error {
	return s.mcRepo.Update(ctx, mc)
}

func (s *MasterClassService) DeleteMasterClass(ctx context.Context, id int64) error {
	return s.mcRepo.Delete(ctx, id)
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
		return nil, err
	}

	products, _, err := s.productRepo.List(ctx, domain.ProductFilter{})
	if err != nil {
		return nil, err
	}

	payments, _, err := s.paymentRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	purchases, _, err := s.purchaseRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_users":     len(users),
		"total_products":  len(products),
		"total_payments":  len(payments),
		"total_purchases": len(purchases),
	}

	return stats, nil
}

var _ port.AdminService = (*AdminService)(nil)
