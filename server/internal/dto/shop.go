package dto

import "time"

// ProductResponse represents a product in API responses.
type ProductResponse struct {
	ID                   int64     `json:"id"`
	Type                 string    `json:"type"`
	TitleRu              string    `json:"title_ru"`
	TitleEn              string    `json:"title_en"`
	DescriptionRu        string    `json:"description_ru,omitempty"`
	DescriptionEn        string    `json:"description_en,omitempty"`
	ShortDescriptionRu   string    `json:"short_description_ru,omitempty"`
	ShortDescriptionEn   string    `json:"short_description_en,omitempty"`
	Price                int64     `json:"price"`
	DurationDays         *int      `json:"duration_days,omitempty"`
	ThumbnailURL         string    `json:"thumbnail_url,omitempty"`
	VideoURL             string    `json:"video_url,omitempty"`
	Status               string    `json:"status"`
	Difficulty           string    `json:"difficulty,omitempty"`
	TotalLessons         int       `json:"total_lessons"`
	TotalDurationMinutes int       `json:"total_duration_minutes"`
	CategoryID           *int64    `json:"category_id,omitempty"`
	InstructorID         *int64    `json:"instructor_id,omitempty"`
	Tags                 []string  `json:"tags,omitempty"`
	PrerequisitesRu      string    `json:"prerequisites_ru,omitempty"`
	PrerequisitesEn      string    `json:"prerequisites_en,omitempty"`
	LearningOutcomesRu   string    `json:"learning_outcomes_ru,omitempty"`
	LearningOutcomesEn   string    `json:"learning_outcomes_en,omitempty"`
	CertificateAvailable bool      `json:"certificate_available"`
	MaxStudents          *int      `json:"max_students,omitempty"`
	StartDate            *string   `json:"start_date,omitempty"`
	Language             string    `json:"language"`
	IsFeatured           bool      `json:"is_featured"`
	ViewCount            int       `json:"view_count"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// CreateProductRequest represents a create product request.
type CreateProductRequest struct {
	Type                 string   `json:"type" binding:"required"`
	TitleRu              string   `json:"title_ru" binding:"required"`
	TitleEn              string   `json:"title_en" binding:"required"`
	DescriptionRu        string   `json:"description_ru,omitempty"`
	DescriptionEn        string   `json:"description_en,omitempty"`
	ShortDescriptionRu   string   `json:"short_description_ru,omitempty"`
	ShortDescriptionEn   string   `json:"short_description_en,omitempty"`
	Price                int64    `json:"price" binding:"required"`
	DurationDays         *int     `json:"duration_days,omitempty"`
	ThumbnailURL         string   `json:"thumbnail_url,omitempty"`
	VideoURL             string   `json:"video_url,omitempty"`
	Status               string   `json:"status"`
	Difficulty           string   `json:"difficulty,omitempty"`
	TotalLessons         int      `json:"total_lessons"`
	TotalDurationMinutes int      `json:"total_duration_minutes"`
	CategoryID           *int64   `json:"category_id,omitempty"`
	InstructorID         *int64   `json:"instructor_id,omitempty"`
	Tags                 []string `json:"tags,omitempty"`
	PrerequisitesRu      string   `json:"prerequisites_ru,omitempty"`
	PrerequisitesEn      string   `json:"prerequisites_en,omitempty"`
	LearningOutcomesRu   string   `json:"learning_outcomes_ru,omitempty"`
	LearningOutcomesEn   string   `json:"learning_outcomes_en,omitempty"`
	CertificateAvailable bool     `json:"certificate_available"`
	MaxStudents          *int     `json:"max_students,omitempty"`
	StartDate            *string  `json:"start_date,omitempty"`
	Language             string   `json:"language"`
	IsFeatured           bool     `json:"is_featured"`
}

// UpdateProductRequest represents an update product request.
type UpdateProductRequest struct {
	Type                 *string  `json:"type,omitempty"`
	TitleRu              *string  `json:"title_ru,omitempty"`
	TitleEn              *string  `json:"title_en,omitempty"`
	DescriptionRu        *string  `json:"description_ru,omitempty"`
	DescriptionEn        *string  `json:"description_en,omitempty"`
	ShortDescriptionRu   *string  `json:"short_description_ru,omitempty"`
	ShortDescriptionEn   *string  `json:"short_description_en,omitempty"`
	Price                *int64   `json:"price,omitempty"`
	DurationDays         *int     `json:"duration_days,omitempty"`
	ThumbnailURL         *string  `json:"thumbnail_url,omitempty"`
	VideoURL             *string  `json:"video_url,omitempty"`
	Status               *string  `json:"status,omitempty"`
	Difficulty           *string  `json:"difficulty,omitempty"`
	TotalLessons         *int     `json:"total_lessons,omitempty"`
	TotalDurationMinutes *int     `json:"total_duration_minutes,omitempty"`
	CategoryID           *int64   `json:"category_id,omitempty"`
	InstructorID         *int64   `json:"instructor_id,omitempty"`
	Tags                 []string `json:"tags,omitempty"`
	PrerequisitesRu      *string  `json:"prerequisites_ru,omitempty"`
	PrerequisitesEn      *string  `json:"prerequisites_en,omitempty"`
	LearningOutcomesRu   *string  `json:"learning_outcomes_ru,omitempty"`
	LearningOutcomesEn   *string  `json:"learning_outcomes_en,omitempty"`
	CertificateAvailable *bool    `json:"certificate_available,omitempty"`
	MaxStudents          *int     `json:"max_students,omitempty"`
	StartDate            *string  `json:"start_date,omitempty"`
	Language             *string  `json:"language,omitempty"`
	IsFeatured           *bool    `json:"is_featured,omitempty"`
}

// CategoryResponse represents a product category in API responses.
type CategoryResponse struct {
	ID            int64     `json:"id"`
	NameRu        string    `json:"name_ru"`
	NameEn        string    `json:"name_en"`
	Slug          string    `json:"slug"`
	DescriptionRu string    `json:"description_ru,omitempty"`
	DescriptionEn string    `json:"description_en,omitempty"`
	SortOrder     int       `json:"sort_order"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateCategoryRequest represents a create category request.
type CreateCategoryRequest struct {
	NameRu        string `json:"name_ru" binding:"required"`
	NameEn        string `json:"name_en" binding:"required"`
	Slug          string `json:"slug" binding:"required"`
	DescriptionRu string `json:"description_ru,omitempty"`
	DescriptionEn string `json:"description_en,omitempty"`
	SortOrder     int    `json:"sort_order"`
	IsActive      bool   `json:"is_active"`
}

// UpdateCategoryRequest represents an update category request.
type UpdateCategoryRequest struct {
	NameRu        *string `json:"name_ru,omitempty"`
	NameEn        *string `json:"name_en,omitempty"`
	Slug          *string `json:"slug,omitempty"`
	DescriptionRu *string `json:"description_ru,omitempty"`
	DescriptionEn *string `json:"description_en,omitempty"`
	SortOrder     *int    `json:"sort_order,omitempty"`
	IsActive      *bool   `json:"is_active,omitempty"`
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
