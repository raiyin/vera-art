package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/dto"
	"github.com/raiyin/artserver/internal/port"
	"github.com/raiyin/artserver/pkg/apperror"
)

// ShopHandler handles shop (products, categories, promo codes, reviews) HTTP requests.
type ShopHandler struct {
	shopService port.ShopService
}

// NewShopHandler creates a new ShopHandler.
func NewShopHandler(shopService port.ShopService) *ShopHandler {
	return &ShopHandler{shopService: shopService}
}

// GetProducts returns a list of products.
func (h *ShopHandler) GetProducts(c *gin.Context) {
	filter := domain.ProductFilter{
		CategoryID: ParseIntQueryAsInt64(c, "category_id"),
		Status:     c.Query("status"),
		Type:       c.Query("type"),
		Difficulty: c.Query("difficulty"),
		Query:      c.Query("q"),
		SortBy:     c.Query("sort_by"),
		SortOrder:  c.Query("sort_order"),
		Page:       ParseIntQuery(c, "page", 1),
		Limit:      ParseIntQuery(c, "limit", 20),
	}

	products, total, err := h.shopService.GetProducts(c.Request.Context(), filter)
	if err != nil {
		slog.Error("GetProducts: failed to list products",
			"error", err,
			"filter", filter,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = dto.ProductResponse{
			ID:                   p.ID,
			Type:                 p.Type,
			TitleRu:              p.TitleRu,
			TitleEn:              p.TitleEn,
			DescriptionRu:        p.DescriptionRu,
			DescriptionEn:        p.DescriptionEn,
			ShortDescriptionRu:   p.ShortDescriptionRu,
			ShortDescriptionEn:   p.ShortDescriptionEn,
			Price:                p.Price,
			DurationDays:         p.DurationDays,
			ThumbnailURL:         p.ThumbnailURL,
			VideoURL:             p.VideoURL,
			Status:               p.Status,
			Difficulty:           p.Difficulty,
			TotalLessons:         p.TotalLessons,
			TotalDurationMinutes: p.TotalDurationMinutes,
			CategoryID:           p.CategoryID,
			InstructorID:         p.InstructorID,
			Tags:                 p.Tags,
			PrerequisitesRu:      p.PrerequisitesRu,
			PrerequisitesEn:      p.PrerequisitesEn,
			LearningOutcomesRu:   p.LearningOutcomesRu,
			LearningOutcomesEn:   p.LearningOutcomesEn,
			CertificateAvailable: p.CertificateAvailable,
			MaxStudents:          p.MaxStudents,
			StartDate:            p.StartDate,
			Language:             p.Language,
			IsFeatured:           p.IsFeatured,
			ViewCount:            p.ViewCount,
			CreatedAt:            p.CreatedAt,
			UpdatedAt:            p.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"products": responses,
		"total":    total,
	})
}

// GetProductBySlug returns a product by slug.
func (h *ShopHandler) GetProductBySlug(c *gin.Context) {
	slug := c.Param("product_slug")

	product, err := h.shopService.GetProductBySlug(c.Request.Context(), slug)
	if err != nil {
		slog.Error("GetProductBySlug: failed to get product",
			"slug", slug,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.ProductResponse{
		ID:                   product.ID,
		Type:                 product.Type,
		TitleRu:              product.TitleRu,
		TitleEn:              product.TitleEn,
		DescriptionRu:        product.DescriptionRu,
		DescriptionEn:        product.DescriptionEn,
		ShortDescriptionRu:   product.ShortDescriptionRu,
		ShortDescriptionEn:   product.ShortDescriptionEn,
		Price:                product.Price,
		DurationDays:         product.DurationDays,
		ThumbnailURL:         product.ThumbnailURL,
		VideoURL:             product.VideoURL,
		Status:               product.Status,
		Difficulty:           product.Difficulty,
		TotalLessons:         product.TotalLessons,
		TotalDurationMinutes: product.TotalDurationMinutes,
		CategoryID:           product.CategoryID,
		InstructorID:         product.InstructorID,
		Tags:                 product.Tags,
		PrerequisitesRu:      product.PrerequisitesRu,
		PrerequisitesEn:      product.PrerequisitesEn,
		LearningOutcomesRu:   product.LearningOutcomesRu,
		LearningOutcomesEn:   product.LearningOutcomesEn,
		CertificateAvailable: product.CertificateAvailable,
		MaxStudents:          product.MaxStudents,
		StartDate:            product.StartDate,
		Language:             product.Language,
		IsFeatured:           product.IsFeatured,
		ViewCount:            product.ViewCount,
		CreatedAt:            product.CreatedAt,
		UpdatedAt:            product.UpdatedAt,
	})
}

// GetProductByID returns a product by ID.
func (h *ShopHandler) GetProductByID(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	product, err := h.shopService.GetProductByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("GetProductByID: failed to get product",
			"product_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.ProductResponse{
		ID:                   product.ID,
		Type:                 product.Type,
		TitleRu:              product.TitleRu,
		TitleEn:              product.TitleEn,
		DescriptionRu:        product.DescriptionRu,
		DescriptionEn:        product.DescriptionEn,
		ShortDescriptionRu:   product.ShortDescriptionRu,
		ShortDescriptionEn:   product.ShortDescriptionEn,
		Price:                product.Price,
		DurationDays:         product.DurationDays,
		ThumbnailURL:         product.ThumbnailURL,
		VideoURL:             product.VideoURL,
		Status:               product.Status,
		Difficulty:           product.Difficulty,
		TotalLessons:         product.TotalLessons,
		TotalDurationMinutes: product.TotalDurationMinutes,
		CategoryID:           product.CategoryID,
		InstructorID:         product.InstructorID,
		Tags:                 product.Tags,
		PrerequisitesRu:      product.PrerequisitesRu,
		PrerequisitesEn:      product.PrerequisitesEn,
		LearningOutcomesRu:   product.LearningOutcomesRu,
		LearningOutcomesEn:   product.LearningOutcomesEn,
		CertificateAvailable: product.CertificateAvailable,
		MaxStudents:          product.MaxStudents,
		StartDate:            product.StartDate,
		Language:             product.Language,
		IsFeatured:           product.IsFeatured,
		ViewCount:            product.ViewCount,
		CreatedAt:            product.CreatedAt,
		UpdatedAt:            product.UpdatedAt,
	})
}

// CreateProduct creates a new product.
func (h *ShopHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if !BindJSON(c, &req) {
		return
	}

	product := &domain.Product{
		Type:                 req.Type,
		TitleRu:              req.TitleRu,
		TitleEn:              req.TitleEn,
		DescriptionRu:        req.DescriptionRu,
		DescriptionEn:        req.DescriptionEn,
		ShortDescriptionRu:   req.ShortDescriptionRu,
		ShortDescriptionEn:   req.ShortDescriptionEn,
		Price:                req.Price,
		DurationDays:         req.DurationDays,
		ThumbnailURL:         req.ThumbnailURL,
		VideoURL:             req.VideoURL,
		Status:               req.Status,
		Difficulty:           req.Difficulty,
		TotalLessons:         req.TotalLessons,
		TotalDurationMinutes: req.TotalDurationMinutes,
		CategoryID:           req.CategoryID,
		InstructorID:         req.InstructorID,
		Tags:                 req.Tags,
		PrerequisitesRu:      req.PrerequisitesRu,
		PrerequisitesEn:      req.PrerequisitesEn,
		LearningOutcomesRu:   req.LearningOutcomesRu,
		LearningOutcomesEn:   req.LearningOutcomesEn,
		CertificateAvailable: req.CertificateAvailable,
		MaxStudents:          req.MaxStudents,
		StartDate:            req.StartDate,
		Language:             req.Language,
		IsFeatured:           req.IsFeatured,
	}

	if err := h.shopService.CreateProduct(c.Request.Context(), product); err != nil {
		slog.Error("CreateProduct: failed to create product",
			"title_ru", req.TitleRu,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Product created successfully",
		"product_id", product.ID,
		"title_ru", product.TitleRu,
	)
	c.JSON(http.StatusCreated, dto.ProductResponse{
		ID:                   product.ID,
		Type:                 product.Type,
		TitleRu:              product.TitleRu,
		TitleEn:              product.TitleEn,
		DescriptionRu:        product.DescriptionRu,
		DescriptionEn:        product.DescriptionEn,
		ShortDescriptionRu:   product.ShortDescriptionRu,
		ShortDescriptionEn:   product.ShortDescriptionEn,
		Price:                product.Price,
		DurationDays:         product.DurationDays,
		ThumbnailURL:         product.ThumbnailURL,
		VideoURL:             product.VideoURL,
		Status:               product.Status,
		Difficulty:           product.Difficulty,
		TotalLessons:         product.TotalLessons,
		TotalDurationMinutes: product.TotalDurationMinutes,
		CategoryID:           product.CategoryID,
		InstructorID:         product.InstructorID,
		Tags:                 product.Tags,
		PrerequisitesRu:      product.PrerequisitesRu,
		PrerequisitesEn:      product.PrerequisitesEn,
		LearningOutcomesRu:   product.LearningOutcomesRu,
		LearningOutcomesEn:   product.LearningOutcomesEn,
		CertificateAvailable: product.CertificateAvailable,
		MaxStudents:          product.MaxStudents,
		StartDate:            product.StartDate,
		Language:             product.Language,
		IsFeatured:           product.IsFeatured,
		ViewCount:            product.ViewCount,
		CreatedAt:            product.CreatedAt,
		UpdatedAt:            product.UpdatedAt,
	})
}

// UpdateProduct updates a product.
func (h *ShopHandler) UpdateProduct(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateProductRequest
	if !BindJSON(c, &req) {
		return
	}

	product := &domain.Product{
		ID: id,
	}

	if req.Type != nil {
		product.Type = *req.Type
	}
	if req.TitleRu != nil {
		product.TitleRu = *req.TitleRu
	}
	if req.TitleEn != nil {
		product.TitleEn = *req.TitleEn
	}
	if req.DescriptionRu != nil {
		product.DescriptionRu = *req.DescriptionRu
	}
	if req.DescriptionEn != nil {
		product.DescriptionEn = *req.DescriptionEn
	}
	if req.ShortDescriptionRu != nil {
		product.ShortDescriptionRu = *req.ShortDescriptionRu
	}
	if req.ShortDescriptionEn != nil {
		product.ShortDescriptionEn = *req.ShortDescriptionEn
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.DurationDays != nil {
		product.DurationDays = req.DurationDays
	}
	if req.ThumbnailURL != nil {
		product.ThumbnailURL = *req.ThumbnailURL
	}
	if req.VideoURL != nil {
		product.VideoURL = *req.VideoURL
	}
	if req.Status != nil {
		product.Status = *req.Status
	}
	if req.Difficulty != nil {
		product.Difficulty = *req.Difficulty
	}
	if req.TotalLessons != nil {
		product.TotalLessons = *req.TotalLessons
	}
	if req.TotalDurationMinutes != nil {
		product.TotalDurationMinutes = *req.TotalDurationMinutes
	}
	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	}
	if req.InstructorID != nil {
		product.InstructorID = req.InstructorID
	}
	if req.PrerequisitesRu != nil {
		product.PrerequisitesRu = *req.PrerequisitesRu
	}
	if req.PrerequisitesEn != nil {
		product.PrerequisitesEn = *req.PrerequisitesEn
	}
	if req.LearningOutcomesRu != nil {
		product.LearningOutcomesRu = *req.LearningOutcomesRu
	}
	if req.LearningOutcomesEn != nil {
		product.LearningOutcomesEn = *req.LearningOutcomesEn
	}
	if req.CertificateAvailable != nil {
		product.CertificateAvailable = *req.CertificateAvailable
	}
	if req.MaxStudents != nil {
		product.MaxStudents = req.MaxStudents
	}
	if req.StartDate != nil {
		product.StartDate = req.StartDate
	}
	if req.Language != nil {
		product.Language = *req.Language
	}
	if req.IsFeatured != nil {
		product.IsFeatured = *req.IsFeatured
	}
	if req.Tags != nil {
		product.Tags = req.Tags
	}

	if err := h.shopService.UpdateProduct(c.Request.Context(), product); err != nil {
		slog.Error("UpdateProduct: failed to update product",
			"product_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Product updated successfully",
		"product_id", id,
	)
	c.JSON(http.StatusOK, dto.ProductResponse{
		ID:                   product.ID,
		Type:                 product.Type,
		TitleRu:              product.TitleRu,
		TitleEn:              product.TitleEn,
		DescriptionRu:        product.DescriptionRu,
		DescriptionEn:        product.DescriptionEn,
		ShortDescriptionRu:   product.ShortDescriptionRu,
		ShortDescriptionEn:   product.ShortDescriptionEn,
		Price:                product.Price,
		DurationDays:         product.DurationDays,
		ThumbnailURL:         product.ThumbnailURL,
		VideoURL:             product.VideoURL,
		Status:               product.Status,
		Difficulty:           product.Difficulty,
		TotalLessons:         product.TotalLessons,
		TotalDurationMinutes: product.TotalDurationMinutes,
		CategoryID:           product.CategoryID,
		InstructorID:         product.InstructorID,
		Tags:                 product.Tags,
		PrerequisitesRu:      product.PrerequisitesRu,
		PrerequisitesEn:      product.PrerequisitesEn,
		LearningOutcomesRu:   product.LearningOutcomesRu,
		LearningOutcomesEn:   product.LearningOutcomesEn,
		CertificateAvailable: product.CertificateAvailable,
		MaxStudents:          product.MaxStudents,
		StartDate:            product.StartDate,
		Language:             product.Language,
		IsFeatured:           product.IsFeatured,
		ViewCount:            product.ViewCount,
		CreatedAt:            product.CreatedAt,
		UpdatedAt:            product.UpdatedAt,
	})
}

// DeleteProduct deletes a product.
func (h *ShopHandler) DeleteProduct(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	if err := h.shopService.DeleteProduct(c.Request.Context(), id); err != nil {
		slog.Error("DeleteProduct: failed to delete product",
			"product_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Product deleted successfully",
		"product_id", id,
	)
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// GetCategories returns a list of categories.
func (h *ShopHandler) GetCategories(c *gin.Context) {
	categories, err := h.shopService.GetCategories(c.Request.Context())
	if err != nil {
		slog.Error("GetCategories: failed to list categories",
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = dto.CategoryResponse{
			ID:            cat.ID,
			NameRu:        cat.NameRu,
			NameEn:        cat.NameEn,
			Slug:          cat.Slug,
			DescriptionRu: cat.DescriptionRu,
			DescriptionEn: cat.DescriptionEn,
			SortOrder:     cat.SortOrder,
			IsActive:      cat.IsActive,
			CreatedAt:     cat.CreatedAt,
			UpdatedAt:     cat.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"categories": responses})
}

// CreateCategory creates a new category.
func (h *ShopHandler) CreateCategory(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if !BindJSON(c, &req) {
		return
	}

	category := &domain.ProductCategory{
		NameRu:        req.NameRu,
		NameEn:        req.NameEn,
		Slug:          req.Slug,
		DescriptionRu: req.DescriptionRu,
		DescriptionEn: req.DescriptionEn,
		SortOrder:     req.SortOrder,
		IsActive:      req.IsActive,
	}

	if err := h.shopService.CreateCategory(c.Request.Context(), category); err != nil {
		slog.Error("CreateCategory: failed to create category",
			"name_ru", req.NameRu,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Category created successfully",
		"category_id", category.ID,
		"name_ru", category.NameRu,
	)
	c.JSON(http.StatusCreated, dto.CategoryResponse{
		ID:            category.ID,
		NameRu:        category.NameRu,
		NameEn:        category.NameEn,
		Slug:          category.Slug,
		DescriptionRu: category.DescriptionRu,
		DescriptionEn: category.DescriptionEn,
		SortOrder:     category.SortOrder,
		IsActive:      category.IsActive,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
	})
}

// UpdateCategory updates a category.
func (h *ShopHandler) UpdateCategory(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateCategoryRequest
	if !BindJSON(c, &req) {
		return
	}

	category := &domain.ProductCategory{
		ID: id,
	}
	if req.NameRu != nil {
		category.NameRu = *req.NameRu
	}
	if req.NameEn != nil {
		category.NameEn = *req.NameEn
	}
	if req.Slug != nil {
		category.Slug = *req.Slug
	}
	if req.DescriptionRu != nil {
		category.DescriptionRu = *req.DescriptionRu
	}
	if req.DescriptionEn != nil {
		category.DescriptionEn = *req.DescriptionEn
	}
	if req.SortOrder != nil {
		category.SortOrder = *req.SortOrder
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := h.shopService.UpdateCategory(c.Request.Context(), category); err != nil {
		slog.Error("UpdateCategory: failed to update category",
			"category_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Category updated successfully",
		"category_id", id,
	)
	c.JSON(http.StatusOK, dto.CategoryResponse{
		ID:            category.ID,
		NameRu:        category.NameRu,
		NameEn:        category.NameEn,
		Slug:          category.Slug,
		DescriptionRu: category.DescriptionRu,
		DescriptionEn: category.DescriptionEn,
		SortOrder:     category.SortOrder,
		IsActive:      category.IsActive,
		CreatedAt:     category.CreatedAt,
		UpdatedAt:     category.UpdatedAt,
	})
}

// DeleteCategory deletes a category.
func (h *ShopHandler) DeleteCategory(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	if err := h.shopService.DeleteCategory(c.Request.Context(), id); err != nil {
		slog.Error("DeleteCategory: failed to delete category",
			"category_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Category deleted successfully",
		"category_id", id,
	)
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// GetPromoCodes returns a list of promo codes.
func (h *ShopHandler) GetPromoCodes(c *gin.Context) {
	codes, total, err := h.shopService.GetPromoCodes(c.Request.Context())
	if err != nil {
		slog.Error("GetPromoCodes: failed to list promo codes",
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.PromoCodeResponse, len(codes))
	for i, code := range codes {
		responses[i] = dto.PromoCodeResponse{
			ID:              code.ID,
			Code:            code.Code,
			DiscountPercent: code.DiscountPercent,
			MaxUses:         code.MaxUses,
			CurrentUses:     code.CurrentUses,
			ExpiresAt:       code.ExpiresAt,
			IsActive:        code.IsActive,
			CreatedAt:       code.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"promo_codes": responses,
		"total":       total,
	})
}

// CreatePromoCode creates a new promo code.
func (h *ShopHandler) CreatePromoCode(c *gin.Context) {
	var req dto.CreatePromoCodeRequest
	if !BindJSON(c, &req) {
		return
	}

	code := &domain.PromoCode{
		Code:            req.Code,
		DiscountPercent: req.DiscountPercent,
		MaxUses:         req.MaxUses,
		ExpiresAt:       req.ExpiresAt,
		IsActive:        req.IsActive,
	}

	if err := h.shopService.CreatePromoCode(c.Request.Context(), code); err != nil {
		slog.Error("CreatePromoCode: failed to create promo code",
			"code", req.Code,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Promo code created successfully",
		"promo_code_id", code.ID,
		"code", code.Code,
	)
	c.JSON(http.StatusCreated, dto.PromoCodeResponse{
		ID:              code.ID,
		Code:            code.Code,
		DiscountPercent: code.DiscountPercent,
		MaxUses:         code.MaxUses,
		CurrentUses:     code.CurrentUses,
		ExpiresAt:       code.ExpiresAt,
		IsActive:        code.IsActive,
		CreatedAt:       code.CreatedAt,
	})
}

// UpdatePromoCode updates a promo code.
func (h *ShopHandler) UpdatePromoCode(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	var req struct {
		Code            string  `json:"code"`
		DiscountPercent float64 `json:"discount_percent"`
		MaxUses         int     `json:"max_uses"`
		ExpiresAt       string  `json:"expires_at"`
		IsActive        *bool   `json:"is_active"`
	}
	if !BindJSON(c, &req) {
		return
	}

	code := &domain.PromoCode{
		ID: id,
	}

	if req.Code != "" {
		code.Code = req.Code
	}
	if req.DiscountPercent > 0 {
		code.DiscountPercent = req.DiscountPercent
	}
	if req.MaxUses > 0 {
		code.MaxUses = req.MaxUses
	}
	if req.IsActive != nil {
		code.IsActive = *req.IsActive
	}

	if err := h.shopService.UpdatePromoCode(c.Request.Context(), code); err != nil {
		slog.Error("UpdatePromoCode: failed to update promo code",
			"promo_code_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Promo code updated successfully",
		"promo_code_id", id,
	)
	c.JSON(http.StatusOK, gin.H{"message": "Promo code updated successfully"})
}

// DeletePromoCode deletes a promo code.
func (h *ShopHandler) DeletePromoCode(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	if err := h.shopService.DeletePromoCode(c.Request.Context(), id); err != nil {
		slog.Error("DeletePromoCode: failed to delete promo code",
			"promo_code_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Promo code deleted successfully",
		"promo_code_id", id,
	)
	c.JSON(http.StatusOK, gin.H{"message": "Promo code deleted successfully"})
}

// ValidatePromoCode validates a promo code.
func (h *ShopHandler) ValidatePromoCode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		slog.Warn("ValidatePromoCode: code is required")
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Code is required",
		})
		return
	}

	promo, err := h.shopService.ValidatePromoCode(c.Request.Context(), code)
	if err != nil {
		slog.Warn("ValidatePromoCode: validation failed",
			"code", code,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.PromoCodeResponse{
		ID:              promo.ID,
		Code:            promo.Code,
		DiscountPercent: promo.DiscountPercent,
		MaxUses:         promo.MaxUses,
		CurrentUses:     promo.CurrentUses,
		ExpiresAt:       promo.ExpiresAt,
		IsActive:        promo.IsActive,
		CreatedAt:       promo.CreatedAt,
	})
}

// GetReviews returns a list of reviews.
func (h *ShopHandler) GetReviews(c *gin.Context) {
	reviews, total, err := h.shopService.GetReviews(c.Request.Context())
	if err != nil {
		slog.Error("GetReviews: failed to list reviews",
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.ReviewResponse, len(reviews))
	for i, r := range reviews {
		responses[i] = dto.ReviewResponse{
			ID:        r.ID,
			UserID:    r.UserID,
			ProductID: r.ProductID,
			Rating:    r.Rating,
			Text:      r.Text,
			Status:    r.Status,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"reviews": responses,
		"total":   total,
	})
}

// GetReviewsByProduct returns reviews for a product.
func (h *ShopHandler) GetReviewsByProduct(c *gin.Context) {
	productID, ok := ParseInt64Param(c, "product_id")
	if !ok {
		return
	}

	reviews, err := h.shopService.GetReviewsByProduct(c.Request.Context(), productID)
	if err != nil {
		slog.Error("GetReviewsByProduct: failed to get reviews",
			"product_id", productID,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.ReviewResponse, len(reviews))
	for i, r := range reviews {
		responses[i] = dto.ReviewResponse{
			ID:        r.ID,
			UserID:    r.UserID,
			ProductID: r.ProductID,
			Rating:    r.Rating,
			Text:      r.Text,
			Status:    r.Status,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"reviews": responses})
}

// CreateReview creates a new review.
func (h *ShopHandler) CreateReview(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req dto.CreateReviewRequest
	if !BindJSON(c, &req) {
		return
	}

	review := &domain.Review{
		UserID:    userID,
		ProductID: req.ProductID,
		Rating:    req.Rating,
		Text:      req.Text,
		Status:    "pending",
	}

	if err := h.shopService.CreateReview(c.Request.Context(), review); err != nil {
		slog.Error("CreateReview: failed to create review",
			"user_id", userID,
			"product_id", req.ProductID,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Review created successfully",
		"review_id", review.ID,
		"user_id", userID,
	)
	c.JSON(http.StatusCreated, dto.ReviewResponse{
		ID:        review.ID,
		UserID:    review.UserID,
		ProductID: review.ProductID,
		Rating:    review.Rating,
		Text:      review.Text,
		Status:    review.Status,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	})
}

// UpdateReview updates a review.
func (h *ShopHandler) UpdateReview(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	var req struct {
		Rating int    `json:"rating"`
		Text   string `json:"text"`
		Status string `json:"status"`
	}
	if !BindJSON(c, &req) {
		return
	}

	review := &domain.Review{
		ID:     id,
		Rating: req.Rating,
		Text:   req.Text,
		Status: req.Status,
	}

	if err := h.shopService.UpdateReview(c.Request.Context(), review); err != nil {
		slog.Error("UpdateReview: failed to update review",
			"review_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Review updated successfully",
		"review_id", id,
	)
	c.JSON(http.StatusOK, gin.H{"message": "Review updated successfully"})
}

// DeleteReview deletes a review.
func (h *ShopHandler) DeleteReview(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	if err := h.shopService.DeleteReview(c.Request.Context(), id); err != nil {
		slog.Error("DeleteReview: failed to delete review",
			"review_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Review deleted successfully",
		"review_id", id,
	)
	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}

// ParseIntQueryAsInt64 parses an int64 query parameter.
func ParseIntQueryAsInt64(c *gin.Context, name string) int64 {
	valStr := c.Query(name)
	if valStr == "" {
		return 0
	}
	val, err := parseInt64(valStr)
	if err != nil {
		return 0
	}
	return val
}

func parseInt64(s string) (int64, error) {
	var result int64
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, nil
		}
		result = result*10 + int64(c-'0')
	}
	return result, nil
}
