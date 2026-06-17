package domain

import "time"

type Sale struct {
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
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	MaterialIDs []int64   `json:"material_ids,omitempty"`
	BaseIDs     []int64   `json:"base_ids,omitempty"`
}

type SaleFilter struct {
	Status     string
	MaterialID int64
	BaseID     int64
	Query      string
	SortBy     string
	SortOrder  string
	Page       int
	Limit      int
}
