package domain

import "time"

type Lesson struct {
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

type LessonFilter struct {
	ProductID int64
	Status    string
	Query     string
}

type LearningProgress struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	LessonID  int64     `json:"lesson_id"`
	ProductID int64     `json:"product_id"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MasterClass struct {
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
