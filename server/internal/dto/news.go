package dto

import "time"

// NewsResponse represents a news entry in API responses.
type NewsResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Content     string    `json:"content,omitempty"`
	ImagePath   string    `json:"image_path"`
	VideoPath   string    `json:"video_path,omitempty"`
	VideoPaths  []string  `json:"video_paths,omitempty"`
	ImagePaths  []string  `json:"image_paths,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateNewsRequest represents a create news request.
type CreateNewsRequest struct {
	Title       string `json:"title" form:"title" binding:"required"`
	Description string `json:"description" form:"description"`
	Content     string `json:"content" form:"content"`
	Status      string `json:"status" form:"status"`
}

// UpdateNewsRequest represents an update news request.
type UpdateNewsRequest struct {
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
	Content     string `json:"content" form:"content"`
	Status      string `json:"status" form:"status"`
}
