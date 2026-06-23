package dto

import "time"

// LessonResponse represents a lesson in API responses.
type LessonResponse struct {
	ID                 int64     `json:"id"`
	ProductID          int64     `json:"product_id"`
	TitleRu            string    `json:"title_ru"`
	TitleEn            string    `json:"title_en"`
	DescriptionRu      string    `json:"description_ru,omitempty"`
	DescriptionEn      string    `json:"description_en,omitempty"`
	ContentType        string    `json:"content_type"`
	ContentURL         string    `json:"content_url,omitempty"`
	Resources          []string  `json:"resources,omitempty"`
	DurationMinutes    int       `json:"duration_minutes"`
	SortOrder          int       `json:"sort_order"`
	IsPreview          bool      `json:"is_preview"`
	HomeworkRu         string    `json:"homework_ru,omitempty"`
	HomeworkEn         string    `json:"homework_en,omitempty"`
	EstimatedStudyTime *int      `json:"estimated_study_time,omitempty"`
	IsRequired         bool      `json:"is_required"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// CreateLessonRequest represents a create lesson request.
type CreateLessonRequest struct {
	ProductID          int64    `json:"product_id" binding:"required"`
	TitleRu            string   `json:"title_ru" binding:"required"`
	TitleEn            string   `json:"title_en" binding:"required"`
	DescriptionRu      string   `json:"description_ru,omitempty"`
	DescriptionEn      string   `json:"description_en,omitempty"`
	ContentType        string   `json:"content_type"`
	ContentURL         string   `json:"content_url,omitempty"`
	Resources          []string `json:"resources,omitempty"`
	DurationMinutes    int      `json:"duration_minutes"`
	SortOrder          int      `json:"sort_order"`
	IsPreview          bool     `json:"is_preview"`
	HomeworkRu         string   `json:"homework_ru,omitempty"`
	HomeworkEn         string   `json:"homework_en,omitempty"`
	EstimatedStudyTime *int     `json:"estimated_study_time,omitempty"`
	IsRequired         bool     `json:"is_required"`
}

// UpdateLessonRequest represents an update lesson request.
type UpdateLessonRequest struct {
	TitleRu            *string  `json:"title_ru,omitempty"`
	TitleEn            *string  `json:"title_en,omitempty"`
	DescriptionRu      *string  `json:"description_ru,omitempty"`
	DescriptionEn      *string  `json:"description_en,omitempty"`
	ContentType        *string  `json:"content_type,omitempty"`
	ContentURL         *string  `json:"content_url,omitempty"`
	Resources          []string `json:"resources,omitempty"`
	DurationMinutes    *int     `json:"duration_minutes,omitempty"`
	SortOrder          *int     `json:"sort_order,omitempty"`
	IsPreview          *bool    `json:"is_preview,omitempty"`
	HomeworkRu         *string  `json:"homework_ru,omitempty"`
	HomeworkEn         *string  `json:"homework_en,omitempty"`
	EstimatedStudyTime *int     `json:"estimated_study_time,omitempty"`
	IsRequired         *bool    `json:"is_required,omitempty"`
}

// LearningProgressResponse represents learning progress in API responses.
type LearningProgressResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	LessonID  int64     `json:"lesson_id"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateProgressRequest represents a progress update request.
type UpdateProgressRequest struct {
	Completed bool `json:"completed"`
}
