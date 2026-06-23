package domain

import "time"

type Lesson struct {
	ID                 int64     `json:"id"`
	ProductID          int64     `json:"product_id"`
	TitleRu            string    `json:"title_ru"`
	TitleEn            string    `json:"title_en"`
	DescriptionRu      string    `json:"description_ru,omitempty"`
	DescriptionEn      string    `json:"description_en,omitempty"`
	ContentType        string    `json:"content_type"` // 'video', 'text', 'pdf', 'quiz', 'assignment'
	ContentURL         string    `json:"content_url,omitempty"`
	DurationMinutes    int       `json:"duration_minutes"`
	SortOrder          int       `json:"sort_order"`
	IsPreview          bool      `json:"is_preview"`
	Resources          []string  `json:"resources,omitempty"`
	HomeworkRu         string    `json:"homework_ru,omitempty"`
	HomeworkEn         string    `json:"homework_en,omitempty"`
	EstimatedStudyTime *int      `json:"estimated_study_time,omitempty"`
	IsRequired         bool      `json:"is_required"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type LessonFilter struct {
	ProductID int64
	Query     string
}

type LearningProgress struct {
	ID                   int64     `json:"id"`
	UserID               int64     `json:"user_id"`
	PurchaseID           int64     `json:"purchase_id"`
	LessonID             int64     `json:"lesson_id"`
	Completed            bool      `json:"completed"`
	CompletedAt          *time.Time `json:"completed_at,omitempty"`
	WatchDurationSeconds int       `json:"watch_duration_seconds"`
	LastPositionSeconds  int       `json:"last_position_seconds"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
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
