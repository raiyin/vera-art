package dto

import "time"

// LessonResponse represents a lesson in API responses.
type LessonResponse struct {
	ID              int64     `json:"id"`
	ProductID       int64     `json:"product_id"`
	Title           string    `json:"title"`
	Description     string    `json:"description,omitempty"`
	Content         string    `json:"content,omitempty"`
	VideoURL        string    `json:"video_url,omitempty"`
	Resources       []string  `json:"resources,omitempty"`
	DurationMinutes int       `json:"duration_minutes,omitempty"`
	SortOrder       int       `json:"sort_order"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CreateLessonRequest represents a create lesson request.
type CreateLessonRequest struct {
	ProductID       int64    `json:"product_id" binding:"required"`
	Title           string   `json:"title" binding:"required"`
	Description     string   `json:"description"`
	Content         string   `json:"content"`
	VideoURL        string   `json:"video_url"`
	Resources       []string `json:"resources"`
	DurationMinutes int      `json:"duration_minutes"`
	SortOrder       int      `json:"sort_order"`
	Status          string   `json:"status"`
}

// UpdateLessonRequest represents an update lesson request.
type UpdateLessonRequest struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Content         string   `json:"content"`
	VideoURL        string   `json:"video_url"`
	Resources       []string `json:"resources"`
	DurationMinutes int      `json:"duration_minutes"`
	SortOrder       int      `json:"sort_order"`
	Status          string   `json:"status"`
}

// LearningProgressResponse represents learning progress in API responses.
type LearningProgressResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	LessonID  int64     `json:"lesson_id"`
	ProductID int64     `json:"product_id"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateLessonProgressRequest represents an update lesson progress request.
type UpdateLessonProgressRequest struct {
	Completed bool `json:"completed" binding:"required"`
}

// MasterClassResponse represents a master class in API responses.
type MasterClassResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price"`
	ImagePath   string    `json:"image_path"`
	VideoURL    string    `json:"video_url,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateMasterClassRequest represents a create master class request.
type CreateMasterClassRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImagePath   string  `json:"image_path"`
	VideoURL    string  `json:"video_url"`
	Status      string  `json:"status"`
}

// UpdateMasterClassRequest represents an update master class request.
type UpdateMasterClassRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImagePath   string  `json:"image_path"`
	VideoURL    string  `json:"video_url"`
	Status      string  `json:"status"`
}
