package domain

import "time"

type Product struct {
	ID                   int64     `json:"id"`
	Type                 string    `json:"type"` // 'course' or 'masterclass'
	TitleRu              string    `json:"title_ru"`
	TitleEn              string    `json:"title_en"`
	DescriptionRu        string    `json:"description_ru,omitempty"`
	DescriptionEn        string    `json:"description_en,omitempty"`
	ShortDescriptionRu   string    `json:"short_description_ru,omitempty"`
	ShortDescriptionEn   string    `json:"short_description_en,omitempty"`
	Price                int64     `json:"price"` // в копейках
	DurationDays         *int      `json:"duration_days,omitempty"`
	ThumbnailURL         string    `json:"thumbnail_url,omitempty"`
	VideoURL             string    `json:"video_url,omitempty"`
	Status               string    `json:"status"` // 'draft', 'published', 'archived'
	Difficulty           string    `json:"difficulty,omitempty"`
	TotalLessons         int       `json:"total_lessons"`
	TotalDurationMinutes int       `json:"total_duration_minutes"`
	CategoryID           *int64    `json:"category_id,omitempty"`
	InstructorID         *int64    `json:"instructor_id,omitempty"`
	Tags                 []string  `json:"tags,omitempty"`
	PrerequisitesRu      string    `json:"prerequisites_ru,omitempty"`
	PrerequisitesEn      string    `json:"prerequisites_en,omitempty"`
	LearningOutcomesRu   string    `json:"learning_outcomes_ru,omitempty"`
	LearningOutcomesEn   string    `json:"learning_outcomes_en,omitempty"`
	CertificateAvailable bool      `json:"certificate_available"`
	MaxStudents          *int      `json:"max_students,omitempty"`
	StartDate            *string   `json:"start_date,omitempty"`
	Language             string    `json:"language"`
	IsFeatured           bool      `json:"is_featured"`
	ViewCount            int       `json:"view_count"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type ProductFilter struct {
	CategoryID int64
	Status     string
	Type       string // 'course' or 'masterclass'
	Difficulty string
	IsFeatured *bool
	Query      string
	SortBy     string
	SortOrder  string
	Page       int
	Limit      int
}

type ProductCategory struct {
	ID            int64     `json:"id"`
	NameRu        string    `json:"name_ru"`
	NameEn        string    `json:"name_en"`
	Slug          string    `json:"slug"`
	DescriptionRu string    `json:"description_ru,omitempty"`
	DescriptionEn string    `json:"description_en,omitempty"`
	SortOrder     int       `json:"sort_order"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
