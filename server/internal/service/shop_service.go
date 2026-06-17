package service

import (
	"context"
	"time"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/port"
)

// ShopService implements port.ShopService.
type ShopService struct {
	productRepo  port.ProductRepository
	categoryRepo port.ProductCategoryRepository
	promoRepo    port.PromoCodeRepository
	reviewRepo   port.ReviewRepository
}

// NewShopService creates a new ShopService.
func NewShopService(
	productRepo port.ProductRepository,
	categoryRepo port.ProductCategoryRepository,
	promoRepo port.PromoCodeRepository,
	reviewRepo port.ReviewRepository,
) *ShopService {
	return &ShopService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		promoRepo:    promoRepo,
		reviewRepo:   reviewRepo,
	}
}

// Products

func (s *ShopService) GetProducts(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int, error) {
	return s.productRepo.List(ctx, filter)
}

func (s *ShopService) GetProductBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	return s.productRepo.GetBySlug(ctx, slug)
}

func (s *ShopService) GetProductByID(ctx context.Context, id int64) (*domain.Product, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *ShopService) CreateProduct(ctx context.Context, product *domain.Product) error {
	if product.Title == "" || product.Slug == "" {
		return domain.ErrInvalidInput
	}
	return s.productRepo.Create(ctx, product)
}

func (s *ShopService) UpdateProduct(ctx context.Context, product *domain.Product) error {
	return s.productRepo.Update(ctx, product)
}

func (s *ShopService) DeleteProduct(ctx context.Context, id int64) error {
	return s.productRepo.Delete(ctx, id)
}

// Categories

func (s *ShopService) GetCategories(ctx context.Context) ([]domain.ProductCategory, error) {
	return s.categoryRepo.List(ctx)
}

func (s *ShopService) CreateCategory(ctx context.Context, category *domain.ProductCategory) error {
	if category.Name == "" {
		return domain.ErrInvalidInput
	}
	return s.categoryRepo.Create(ctx, category)
}

func (s *ShopService) UpdateCategory(ctx context.Context, category *domain.ProductCategory) error {
	return s.categoryRepo.Update(ctx, category)
}

func (s *ShopService) DeleteCategory(ctx context.Context, id int64) error {
	return s.categoryRepo.Delete(ctx, id)
}

// Promo Codes

func (s *ShopService) GetPromoCodes(ctx context.Context) ([]domain.PromoCode, int, error) {
	return s.promoRepo.List(ctx)
}

func (s *ShopService) CreatePromoCode(ctx context.Context, code *domain.PromoCode) error {
	if code.Code == "" || code.DiscountPercent <= 0 {
		return domain.ErrInvalidInput
	}
	return s.promoRepo.Create(ctx, code)
}

func (s *ShopService) UpdatePromoCode(ctx context.Context, code *domain.PromoCode) error {
	return s.promoRepo.Update(ctx, code)
}

func (s *ShopService) DeletePromoCode(ctx context.Context, id int64) error {
	return s.promoRepo.Delete(ctx, id)
}

func (s *ShopService) ValidatePromoCode(ctx context.Context, code string) (*domain.PromoCode, error) {
	promo, err := s.promoRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, domain.ErrPromoCodeInvalid
	}

	if !promo.IsActive {
		return nil, domain.ErrPromoCodeInvalid
	}

	if time.Now().After(promo.ExpiresAt) {
		return nil, domain.ErrPromoCodeInvalid
	}

	if promo.MaxUses > 0 && promo.CurrentUses >= promo.MaxUses {
		return nil, domain.ErrPromoCodeUsed
	}

	return promo, nil
}

// Reviews

func (s *ShopService) GetReviews(ctx context.Context) ([]domain.Review, int, error) {
	return s.reviewRepo.List(ctx)
}

func (s *ShopService) GetReviewsByProduct(ctx context.Context, productID int64) ([]domain.Review, error) {
	return s.reviewRepo.ListByProduct(ctx, productID)
}

func (s *ShopService) CreateReview(ctx context.Context, review *domain.Review) error {
	if review.Rating < 1 || review.Rating > 5 {
		return domain.ErrInvalidInput
	}
	return s.reviewRepo.Create(ctx, review)
}

func (s *ShopService) UpdateReview(ctx context.Context, review *domain.Review) error {
	return s.reviewRepo.Update(ctx, review)
}

func (s *ShopService) DeleteReview(ctx context.Context, id int64) error {
	return s.reviewRepo.Delete(ctx, id)
}

// Ensure interface compliance.
var _ port.ShopService = (*ShopService)(nil)
