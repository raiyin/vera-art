package dto

import "time"

// WorkResponse represents a work in API responses.
type WorkResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	ImagePath   string    `json:"image_path"`
	Year        int       `json:"year,omitempty"`
	Technique   string    `json:"technique,omitempty"`
	Size        string    `json:"size,omitempty"`
	Status      string    `json:"status"`
	SortOrder   int       `json:"sort_order"`
	MaterialIDs []int64   `json:"material_ids,omitempty"`
	BaseIDs     []int64   `json:"base_ids,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UpdateWorkResponse represents a work in update API responses.
type UpdateWorkResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	ImagePath   string    `json:"image_path"`
	Year        int       `json:"year,omitempty"`
	Technique   string    `json:"technique,omitempty"`
	Size        string    `json:"size,omitempty"`
	Status      string    `json:"status"`
	SortOrder   int       `json:"sort_order"`
	MaterialIDs []int64   `json:"material_ids,omitempty"`
	BaseIDs     []int64   `json:"base_ids,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateWorkRequest represents a create work request.
type CreateWorkRequest struct {
	Title       string  `json:"title" form:"title" binding:"required"`
	Description string  `json:"description" form:"description"`
	Year        int     `json:"year" form:"year"`
	Technique   string  `json:"technique" form:"technique"`
	Size        string  `json:"size" form:"size"`
	Status      string  `json:"status" form:"status"`
	SortOrder   int     `json:"sort_order" form:"sort_order"`
	MaterialIDs []int64 `json:"material_ids" form:"material_ids"`
	BaseIDs     []int64 `json:"base_ids" form:"base_ids"`
}

// UpdateWorkRequest represents an update work request.
type UpdateWorkRequest struct {
	Title       string  `json:"title" form:"title"`
	Description string  `json:"description" form:"description"`
	Year        int     `json:"year" form:"year"`
	Technique   string  `json:"technique" form:"technique"`
	Size        string  `json:"size" form:"size"`
	Status      string  `json:"status" form:"status"`
	SortOrder   int     `json:"sort_order" form:"sort_order"`
	MaterialIDs []int64 `json:"material_ids" form:"material_ids"`
	BaseIDs     []int64 `json:"base_ids" form:"base_ids"`
}

// SaleResponse represents a sale in API responses.
type SaleResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	ImagePath   string    `json:"image_path"`
	Price       float64   `json:"price"`
	OldPrice    float64   `json:"old_price,omitempty"`
	Year        int       `json:"year,omitempty"`
	Technique   string    `json:"technique,omitempty"`
	Size        string    `json:"size,omitempty"`
	Status      string    `json:"status"`
	SortOrder   int       `json:"sort_order"`
	Sold        bool      `json:"sold"`
	MaterialIDs []int64   `json:"material_ids,omitempty"`
	BaseIDs     []int64   `json:"base_ids,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UpdateSaleResponse represents a sale in update API responses.
type UpdateSaleResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	ImagePath   string    `json:"image_path"`
	Price       float64   `json:"price"`
	OldPrice    float64   `json:"old_price,omitempty"`
	Year        int       `json:"year,omitempty"`
	Technique   string    `json:"technique,omitempty"`
	Size        string    `json:"size,omitempty"`
	Status      string    `json:"status"`
	SortOrder   int       `json:"sort_order"`
	Sold        bool      `json:"sold"`
	MaterialIDs []int64   `json:"material_ids,omitempty"`
	BaseIDs     []int64   `json:"base_ids,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateSaleRequest represents a create sale request.
type CreateSaleRequest struct {
	Title       string  `json:"title" form:"title" binding:"required"`
	Description string  `json:"description" form:"description"`
	Price       float64 `json:"price" form:"price" binding:"required"`
	OldPrice    float64 `json:"old_price" form:"old_price"`
	Year        int     `json:"year" form:"year"`
	Technique   string  `json:"technique" form:"technique"`
	Size        string  `json:"size" form:"size"`
	Status      string  `json:"status" form:"status"`
	SortOrder   int     `json:"sort_order" form:"sort_order"`
	Sold        bool    `json:"sold" form:"sold"`
	MaterialIDs []int64 `json:"material_ids" form:"material_ids"`
	BaseIDs     []int64 `json:"base_ids" form:"base_ids"`
}

// UpdateSaleRequest represents an update sale request.
type UpdateSaleRequest struct {
	Title       string  `json:"title" form:"title"`
	Description string  `json:"description" form:"description"`
	Price       float64 `json:"price" form:"price"`
	OldPrice    float64 `json:"old_price" form:"old_price"`
	Year        int     `json:"year" form:"year"`
	Technique   string  `json:"technique" form:"technique"`
	Size        string  `json:"size" form:"size"`
	Status      string  `json:"status" form:"status"`
	SortOrder   int     `json:"sort_order" form:"sort_order"`
	Sold        bool    `json:"sold" form:"sold"`
	MaterialIDs []int64 `json:"material_ids" form:"material_ids"`
	BaseIDs     []int64 `json:"base_ids" form:"base_ids"`
}
