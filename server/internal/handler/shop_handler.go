package handler

import (
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
		Query:      c.Query("q"),
		SortBy:     c.Query("sort_by"),
		SortOrder:  c.Query("sort_order"),
		Page:       ParseIntQuery(c, "page", 1),
		Limit:      ParseIntQuery(c, "limit", 20),
	}

	products, total, err := h.shopService.GetProducts(c.Request.Context(), filter)
	if err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.ProductResponse, len(products))
	for i, p := range products {
		responses[i] = dto.ProductResponse{
			ID:              p.ID,
			Title:           p.Title,
			Slug:            p.Slug,
			Description:     p.Description,
			FullDescription: p.FullDescription,
			Price:           p.Price,
			OldPrice:        p.OldPrice,
			ImagePath:       p.ImagePath,
			CategoryID:      p.CategoryID,
			CategoryName:    p.CategoryName,
			Status:          p.Status,
			IsDigital:       p.IsDigital,
			IsMasterClass:   p.IsMasterClass,
			SortOrder:       p.SortOrder,
			Tags:            p.Tags,
			CreatedAt:       p.CreatedAt,
			UpdatedAt:       p.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"products": responses,
		"total":    total,
	})
}

// GetProductBySlug returns a product by slug.
func (h *ShopHandler) GetProductBySlug(c *gin.Context) {
	slug := c.Param("slug")

	product, err := h.shopService.GetProductBySlug(c.Request.Context(), slug)
	if err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.ProductResponse{
		ID:              product.ID,
		Title:           product.Title,
		Slug:            product.Slug,
		Description:     product.Description,
		FullDescription: product.FullDescription,
		Price:           product.Price,
		OldPrice:        product.OldPrice,
		ImagePath:       product.ImagePath,
		CategoryID:      product.CategoryID,
		CategoryName:    product.CategoryName,
		Status:          product.Status,
		IsDigital:       product.IsDigital,
		IsMasterClass:   product.IsMasterClass,
		SortOrder:       product.SortOrder,
		Tags:            product.Tags,
		CreatedAt:       product.CreatedAt,
		UpdatedAt:       product.UpdatedAt,
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
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.ProductResponse{
		ID:              product.ID,
		Title:           product.Title,
		Slug:            product.Slug,
		Description:     product.Description,
		FullDescription: product.FullDescription,
		Price:           product.Price,
		OldPrice:        product.OldPrice,
		ImagePath:       product.ImagePath,
		CategoryID:      product.CategoryID,
		CategoryName:    product.CategoryName,
		Status:          product.Status,
		IsDigital:       product.IsDigital,
		IsMasterClass:   product.IsMasterClass,
		SortOrder:       product.SortOrder,
		Tags:            product.Tags,
		CreatedAt:       product.CreatedAt,
		UpdatedAt:       product.UpdatedAt,
	})
}

// CreateProduct creates a new product.
func (h *ShopHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if !BindJSON(c, &req) {
		return
	}

	product := &domain.Product{
		Title:           req.Title,
		Slug:            req.Slug,
		Description:     req.Description,
		FullDescription: req.FullDescription,
		Price:           req.Price,
		OldPrice:        req.OldPrice,
		ImagePath:       req.ImagePath,
		CategoryID:      req.CategoryID,
		Status:          req.Status,
		IsDigital:       req.IsDigital,
		IsMasterClass:   req.IsMasterClass,
		SortOrder:       req.SortOrder,
		Tags:            req.Tags,
	}

	if err := h.shopService.CreateProduct(c.Request.Context(), product); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusCreated, dto.ProductResponse{
		ID:              product.ID,
		Title:           product.Title,
		Slug:            product.Slug,
		Description:     product.Description,
		FullDescription: product.FullDescription,
		Price:           product.Price,
		OldPrice:        product.OldPrice,
		ImagePath:       product.ImagePath,
		CategoryID:      product.CategoryID,
		CategoryName:    product.CategoryName,
		Status:          product.Status,
		IsDigital:       product.IsDigital,
		IsMasterClass:   product.IsMasterClass,
		SortOrder:       product.SortOrder,
		Tags:            product.Tags,
		CreatedAt:       product.CreatedAt,
		UpdatedAt:       product.UpdatedAt,
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
		ID:              id,
		Title:           req.Title,
		Slug:            req.Slug,
		Description:     req.Description,
		FullDescription: req.FullDescription,
		Price:           req.Price,
		OldPrice:        req.OldPrice,
		ImagePath:       req.ImagePath,
		CategoryID:      req.CategoryID,
		Status:          req.Status,
		IsDigital:       req.IsDigital,
		IsMasterClass:   req.IsMasterClass,
		SortOrder:       req.SortOrder,
		Tags:            req.Tags,
	}

	if err := h.shopService.UpdateProduct(c.Request.Context(), product); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.ProductResponse{
		ID:              product.ID,
		Title:           product.Title,
		Slug:            product.Slug,
		Description:     product.Description,
		FullDescription: product.FullDescription,
		Price:           product.Price,
		OldPrice:        product.OldPrice,
		ImagePath:       product.ImagePath,
		CategoryID:      product.CategoryID,
		CategoryName:    product.CategoryName,
		Status:          product.Status,
		IsDigital:       product.IsDigital,
		IsMasterClass:   product.IsMasterClass,
		SortOrder:       product.SortOrder,
		Tags:            product.Tags,
		CreatedAt:       product.CreatedAt,
		UpdatedAt:       product.UpdatedAt,
	})
}

// DeleteProduct deletes a product.
func (h *ShopHandler) DeleteProduct(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	if err := h.shopService.DeleteProduct(c.Request.Context(), id); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// GetCategories returns a list of categories.
func (h *ShopHandler) GetCategories(c *gin.Context) {
	categories, err := h.shopService.GetCategories(c.Request.Context())
	if err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.CategoryResponse, len(categories))
	for i, cat := range categories {
		responses[i] = dto.CategoryResponse{
			ID:        cat.ID,
			Name:      cat.Name,
			Slug:      cat.Slug,
			SortOrder: cat.SortOrder,
			CreatedAt: cat.CreatedAt,
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
		Name:      req.Name,
		Slug:      req.Slug,
		SortOrder: req.SortOrder,
	}

	if err := h.shopService.CreateCategory(c.Request.Context(), category); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusCreated, dto.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		SortOrder: category.SortOrder,
		CreatedAt: category.CreatedAt,
	})
}

// UpdateCategory updates a category.
func (h *ShopHandler) UpdateCategory(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	var req struct {
		Name      string `json:"name"`
		Slug      string `json:"slug"`
		SortOrder int    `json:"sort_order"`
	}
	if !BindJSON(c, &req) {
		return
	}

	category := &domain.ProductCategory{
		ID:        id,
		Name:      req.Name,
		Slug:      req.Slug,
		SortOrder: req.SortOrder,
	}

	if err := h.shopService.UpdateCategory(c.Request.Context(), category); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		SortOrder: category.SortOrder,
		CreatedAt: category.CreatedAt,
	})
}

// DeleteCategory deletes a category.
func (h *ShopHandler) DeleteCategory(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	if err := h.shopService.DeleteCategory(c.Request.Context(), id); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}

// GetPromoCodes returns a list of promo codes.
func (h *ShopHandler) GetPromoCodes(c *gin.Context) {
	codes, total, err := h.shopService.GetPromoCodes(c.Request.Context())
	if err != nil {
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
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

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
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promo code updated successfully"})
}

// DeletePromoCode deletes a promo code.
func (h *ShopHandler) DeletePromoCode(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	if err := h.shopService.DeletePromoCode(c.Request.Context(), id); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Promo code deleted successfully"})
}

// ValidatePromoCode validates a promo code.
func (h *ShopHandler) ValidatePromoCode(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Code is required",
		})
		return
	}

	promo, err := h.shopService.ValidatePromoCode(c.Request.Context(), code)
	if err != nil {
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
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

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
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review updated successfully"})
}

// DeleteReview deletes a review.
func (h *ShopHandler) DeleteReview(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	if err := h.shopService.DeleteReview(c.Request.Context(), id); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

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
