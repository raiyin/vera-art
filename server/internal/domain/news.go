package domain

import "time"

type News struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Content     string    `json:"content,omitempty"`
	ImagePath   string    `json:"image_path"`
	VideoPath   string    `json:"video_path,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type NewsFilter struct {
	Status string
	Query  string
	Page   int
	Limit  int
}
