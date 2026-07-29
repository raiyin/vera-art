package domain

import "time"

type Sale struct {
	ID          int64     `json:"id"`
	NameRu      string    `json:"name_ru"`
	NameEn      string    `json:"name_en"`
	Description string    `json:"description,omitempty"`
	ImagePath   string    `json:"image_path"`
	SalePath    string    `json:"sale_path"`
	Price       float64   `json:"price"`
	Year        int       `json:"year,omitempty"`
	Technique   string    `json:"technique,omitempty"`
	Width       int       `json:"width,omitempty"`
	Height      int       `json:"height,omitempty"`
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
