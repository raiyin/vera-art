package service

import (
	"context"
	"log/slog"
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
	products, total, err := s.productRepo.List(ctx, filter)
	if err != nil {
		slog.Error("ShopService.GetProducts: failed to list products",
			"filter", filter,
			"error", err,
		)
		return nil, 0, err
	}
	slog.Debug("ShopService.GetProducts: products listed",
		"count", len(products),
		"total", total,
	)
	return products, total, nil
}

func (s *ShopService) GetProductBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	product, err := s.productRepo.GetBySlug(ctx, slug)
	if err != nil {
		slog.Error("ShopService.GetProductBySlug: failed to get product by slug",
			"slug", slug,
			"error", err,
		)
		return nil, err
	}
	slog.Debug("ShopService.GetProductBySlug: product retrieved",
		"slug", slug,
		"product_id", product.ID,
	)
	return product, nil
}

func (s *ShopService) GetProductByID(ctx context.Context, id int64) (*domain.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("ShopService.GetProductByID: failed to get product by ID",
			"product_id", id,
			"error", err,
		)
		return nil, err
	}
	slog.Debug("ShopService.GetProductByID: product retrieved",
		"product_id", id,
		"title_ru", product.TitleRu,
	)
	return product, nil
}

func (s *ShopService) CreateProduct(ctx context.Context, product *domain.Product) error {
	if product.TitleRu == "" && product.TitleEn == "" {
		slog.Warn("ShopService.CreateProduct: missing title",
			"title_ru", product.TitleRu,
			"title_en", product.TitleEn,
		)
		return domain.ErrInvalidInput
	}
	if err := s.productRepo.Create(ctx, product); err != nil {
		slog.Error("ShopService.CreateProduct: failed to create product",
			"title_ru", product.TitleRu,
			"title_en", product.TitleEn,
			"error", err,
		)
		return err
	}
	slog.Info("ShopService.CreateProduct: product created",
		"product_id", product.ID,
		"title_ru", product.TitleRu,
		"title_en", product.TitleEn,
	)
	return nil
}

func (s *ShopService) UpdateProduct(ctx context.Context, product *domain.Product) error {
	if err := s.productRepo.Update(ctx, product); err != nil {
		slog.Error("ShopService.UpdateProduct: failed to update product",
			"product_id", product.ID,
			"title_ru", product.TitleRu,
			"error", err,
		)
		return err
	}
	slog.Info("ShopService.UpdateProduct: product updated",
		"product_id", product.ID,
		"title_ru", product.TitleRu,
	)
	return nil
}

func (s *ShopService) DeleteProduct(ctx context.Context, id int64) error {
	if err := s.productRepo.Delete(ctx, id); err != nil {
		slog.Error("ShopService.DeleteProduct: failed to delete product",
			"product_id", id,
			"error", err,
		)
		return err
	}
	slog.Info("ShopService.DeleteProduct: product deleted",
		"product_id", id,
	)
	return nil
}

// Categories

func (s *ShopService) GetCategories(ctx context.Context) ([]domain.ProductCategory, error) {
	categories, err := s.categoryRepo.List(ctx)
	if err != nil {
		slog.Error("ShopService.GetCategories: failed to list categories",
			"error", err,
		)
		return nil, err
	}
	slog.Debug("ShopService.GetCategories: categories listed",
		"count", len(categories),
	)
	return categories, nil
}

func (s *ShopService) CreateCategory(ctx context.Context, category *domain.ProductCategory) error {
	if category.NameRu == "" && category.NameEn == "" {
		slog.Warn("ShopService.CreateCategory: empty name")
		return domain.ErrInvalidInput
	}
	if err := s.categoryRepo.Create(ctx, category); err != nil {
		slog.Error("ShopService.CreateCategory: failed to create category",
			"name_ru", category.NameRu,
			"error", err,
		)
		return err
	}
	slog.Info("ShopService.CreateCategory: category created",
		"category_id", category.ID,
		"name_ru", category.NameRu,
	)
	return nil
}

func (s *ShopService) UpdateCategory(ctx context.Context, category *domain.ProductCategory) error {
	if err := s.categoryRepo.Update(ctx, category); err != nil {
		slog.Error("ShopService.UpdateCategory: failed to update category",
			"category_id", category.ID,
			"name_ru", category.NameRu,
			"error", err,
		)
		return err
	}
	slog.Info("ShopService.UpdateCategory: category updated",
		"category_id", category.ID,
		"name_ru", category.NameRu,
	)
	return nil
}

func (s *ShopService) DeleteCategory(ctx context.Context, id int64) error {
	if err := s.categoryRepo.Delete(ctx, id); err != nil {
		slog.Error("ShopService.DeleteCategory: failed to delete category",
			"category_id", id,
			"error", err,
		)
		return err
	}
	slog.Info("ShopService.DeleteCategory: category deleted",
		"category_id", id,
	)
	return nil
}

// Promo Codes

func (s *ShopService) GetPromoCodes(ctx context.Context) ([]domain.PromoCode, int, error) {
	codes, total, err := s.promoRepo.List(ctx)
	if err != nil {
		slog.Error("ShopService.GetPromoCodes: failed to list promo codes",
			"error", err,
		)
		return nil, 0, err
	}
	slog.Debug("ShopService.GetPromoCodes: promo codes listed",
		"count", len(codes),
		"total", total,
	)
	return codes, total, nil
}

func (s *ShopService) CreatePromoCode(ctx context.Context, code *domain.PromoCode) error {
	if code.Code == "" || code.DiscountPercent <= 0 {
		slog.Warn("ShopService.CreatePromoCode: invalid promo code data",
			"code", code.Code,
			"discount_percent", code.DiscountPercent,
		)
		return domain.ErrInvalidInput
	}
	if err := s.promoRepo.Create(ctx, code); err != nil {
		slog.Error("ShopService.CreatePromoCode: failed to create promo code",
			"code", code.Code,
			"discount_percent", code.DiscountPercent,
			"error", err,
		)
		return err
	}
	slog.Info("ShopService.CreatePromoCode: promo code created",
		"promo_id", code.ID,
		"code", code.Code,
		"discount_percent", code.DiscountPercent,
	)
	return nil
}

func (s *ShopService) UpdatePromoCode(ctx context.Context, code *domain.PromoCode) error {
	if err := s.promoRepo.Update(ctx, code); err != nil {
		slog.Error("ShopService.UpdatePromoCode: failed to update promo code",
			"promo_id", code.ID,
			"code", code.Code,
			"error", err,
		)
		return err
	}
	slog.Info("ShopService.UpdatePromoCode: promo code updated",
		"promo_id", code.ID,
		"code", code.Code,
	)
	return nil
}

func (s *ShopService) DeletePromoCode(ctx context.Context, id int64) error {
	if err := s.promoRepo.Delete(ctx, id); err != nil {
		slog.Error("ShopService.DeletePromoCode: failed to delete promo code",
			"promo_id", id,
			"error", err,
		)
		return err
	}
	slog.Info("ShopService.DeletePromoCode: promo code deleted",
		"promo_id", id,
	)
	return nil
}

func (s *ShopService) ValidatePromoCode(ctx context.Context, code string) (*domain.PromoCode, error) {
	promo, err := s.promoRepo.GetByCode(ctx, code)
	if err != nil {
		slog.Warn("ShopService.ValidatePromoCode: promo code not found",
			"code", code,
			"error", err,
		)
		return nil, domain.ErrPromoCodeInvalid
	}

	if !promo.IsActive {
		slog.Warn("ShopService.ValidatePromoCode: promo code is inactive",
			"code", code,
			"promo_id", promo.ID,
		)
		return nil, domain.ErrPromoCodeInvalid
	}

	if time.Now().After(promo.ExpiresAt) {
		slog.Warn("ShopService.ValidatePromoCode: promo code expired",
			"code", code,
			"promo_id", promo.ID,
			"expires_at", promo.ExpiresAt,
		)
		return nil, domain.ErrPromoCodeInvalid
	}

	if promo.MaxUses > 0 && promo.CurrentUses >= promo.MaxUses {
		slog.Warn("ShopService.ValidatePromoCode: promo code usage limit reached",
			"code", code,
			"promo_id", promo.ID,
			"current_uses", promo.CurrentUses,
			"max_uses", promo.MaxUses,
		)
		return nil, domain.ErrPromoCodeUsed
	}

	slog.Info("ShopService.ValidatePromoCode: promo code validated",
		"code", code,
		"promo_id", promo.ID,
		"discount_percent", promo.DiscountPercent,
	)
	return promo, nil
}

// Reviews

func (s *ShopService) GetReviews(ctx context.Context) ([]domain.Review, int, error) {
	reviews, total, err := s.reviewRepo.List(ctx)
	if err != nil {
		slog.Error("ShopService.GetReviews: failed to list reviews",
			"error", err,
		)
		return nil, 0, err
	}
	slog.Debug("ShopService.GetReviews: reviews listed",
		"count", len(reviews),
		"total", total,
	)
	return reviews, total, nil
}

func (s *ShopService) GetReviewsByProduct(ctx context.Context, productID int64) ([]domain.Review, error) {
	reviews, err := s.reviewRepo.ListByProduct(ctx, productID)
	if err != nil {
		slog.Error("ShopService.GetReviewsByProduct: failed to list reviews by product",
			"product_id", productID,
			"error", err,
		)
		return nil, err
	}
	slog.Debug("ShopService.GetReviewsByProduct: reviews listed by product",
		"product_id", productID,
		"count", len(reviews),
	)
	return reviews, nil
}

func (s *ShopService) CreateReview(ctx context.Context, review *domain.Review) error {
	if review.Rating < 1 || review.Rating > 5 {
		slog.Warn("ShopService.CreateReview: invalid rating",
			"rating", review.Rating,
			"product_id", review.ProductID,
			"user_id", review.UserID,
		)
		return domain.ErrInvalidInput
	}
	if err := s.reviewRepo.Create(ctx, review); err != nil {
		slog.Error("ShopService.CreateReview: failed to create review",
			"product_id", review.ProductID,
			"user_id", review.UserID,
			"rating", review.Rating,
			"error", err,
		)
		return err
	}
	slog.Info("ShopService.CreateReview: review created",
		"review_id", review.ID,
		"product_id", review.ProductID,
		"user_id", review.UserID,
		"rating", review.Rating,
	)
	return nil
}

func (s *ShopService) UpdateReview(ctx context.Context, review *domain.Review) error {
	if err := s.reviewRepo.Update(ctx, review); err != nil {
		slog.Error("ShopService.UpdateReview: failed to update review",
			"review_id", review.ID,
			"product_id", review.ProductID,
			"user_id", review.UserID,
			"error", err,
		)
		return err
	}
	slog.Info("ShopService.UpdateReview: review updated",
		"review_id", review.ID,
		"product_id", review.ProductID,
		"user_id", review.UserID,
	)
	return nil
}

func (s *ShopService) DeleteReview(ctx context.Context, id int64) error {
	if err := s.reviewRepo.Delete(ctx, id); err != nil {
		slog.Error("ShopService.DeleteReview: failed to delete review",
			"review_id", id,
			"error", err,
		)
		return err
	}
	slog.Info("ShopService.DeleteReview: review deleted",
		"review_id", id,
	)
	return nil
}

// Ensure interface compliance.
var _ port.ShopService = (*ShopService)(nil)
