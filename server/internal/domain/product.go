package domain

import "time"

type Product struct {
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

type ProductFilter struct {
	CategoryID    int64
	Status        string
	IsDigital     *bool
	IsMasterClass *bool
	Query         string
	SortBy        string
	SortOrder     string
	Page          int
	Limit         int
}

type ProductCategory struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}
