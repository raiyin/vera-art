package dto

import "time"

// ProductResponse represents a product in API responses.
type ProductResponse struct {
	ID              int64     `json:"id"`
	Title           string    `json:"title"`
	Slug            string    `json:"slug"`
	Description     string    `json:"description,omitempty"`
	FullDescription string    `json:"full_description,omitempty"`
	Price           float64   `json:"price"`
	OldPrice        float64   `json:"old_price,omitempty"`
	ImagePath       string    `json:"image_path"`
	CategoryID      int64     `json:"category_id"`
	CategoryName    string    `json:"category_name,omitempty"`
	Status          string    `json:"status"`
	IsDigital       bool      `json:"is_digital"`
	IsMasterClass   bool      `json:"is_master_class"`
	SortOrder       int       `json:"sort_order"`
	Tags            []string  `json:"tags,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CreateProductRequest represents a create product request.
type CreateProductRequest struct {
	Title           string   `json:"title" binding:"required"`
	Slug            string   `json:"slug" binding:"required"`
	Description     string   `json:"description"`
	FullDescription string   `json:"full_description"`
	Price           float64  `json:"price" binding:"required"`
	OldPrice        float64  `json:"old_price"`
	ImagePath       string   `json:"image_path"`
	CategoryID      int64    `json:"category_id"`
	Status          string   `json:"status"`
	IsDigital       bool     `json:"is_digital"`
	IsMasterClass   bool     `json:"is_master_class"`
	SortOrder       int      `json:"sort_order"`
	Tags            []string `json:"tags"`
}

// UpdateProductRequest represents an update product request.
type UpdateProductRequest struct {
	Title           string   `json:"title"`
	Slug            string   `json:"slug"`
	Description     string   `json:"description"`
	FullDescription string   `json:"full_description"`
	Price           float64  `json:"price"`
	OldPrice        float64  `json:"old_price"`
	ImagePath       string   `json:"image_path"`
	CategoryID      int64    `json:"category_id"`
	Status          string   `json:"status"`
	IsDigital       bool     `json:"is_digital"`
	IsMasterClass   bool     `json:"is_master_class"`
	SortOrder       int      `json:"sort_order"`
	Tags            []string `json:"tags"`
}

// CategoryResponse represents a product category in API responses.
type CategoryResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateCategoryRequest represents a create category request.
type CreateCategoryRequest struct {
	Name      string `json:"name" binding:"required"`
	Slug      string `json:"slug" binding:"required"`
	SortOrder int    `json:"sort_order"`
}

// PromoCodeResponse represents a promo code in API responses.
type PromoCodeResponse struct {
	ID              int64     `json:"id"`
	Code            string    `json:"code"`
	DiscountPercent float64   `json:"discount_percent"`
	MaxUses         int       `json:"max_uses"`
	CurrentUses     int       `json:"current_uses"`
	ExpiresAt       time.Time `json:"expires_at"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
}

// CreatePromoCodeRequest represents a create promo code request.
type CreatePromoCodeRequest struct {
	Code            string    `json:"code" binding:"required"`
	DiscountPercent float64   `json:"discount_percent" binding:"required"`
	MaxUses         int       `json:"max_uses"`
	ExpiresAt       time.Time `json:"expires_at" binding:"required"`
	IsActive        bool      `json:"is_active"`
}

// ReviewResponse represents a review in API responses.
type ReviewResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	Rating    int       `json:"rating"`
	Text      string    `json:"text,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateReviewRequest represents a create review request.
type CreateReviewRequest struct {
	ProductID int64  `json:"product_id" binding:"required"`
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Text      string `json:"text"`
}
