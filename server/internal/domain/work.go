package domain

import "time"

type Work struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	ImagePath   string    `json:"image_path"`
	Year        int       `json:"year,omitempty"`
	Technique   string    `json:"technique,omitempty"`
	Width       int       `json:"width,omitempty"`
	Height      int       `json:"height,omitempty"`
	Status      string    `json:"status"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	MaterialIDs []int64   `json:"material_ids,omitempty"`
	BaseIDs     []int64   `json:"base_ids,omitempty"`
}

type WorkFilter struct {
	Status     string
	MaterialID int64
	BaseID     int64
	Query      string
	SortBy     string
	SortOrder  string
	Page       int
	Limit      int
}
