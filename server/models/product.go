package models

import (
	"encoding/json"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Product struct {
	Id                   int        `json:"id"`
	Type                 string     `json:"type"` // "course" or "masterclass"
	TitleRu              string     `json:"title_ru"`
	TitleEn              string     `json:"title_en"`
	DescriptionRu        string     `json:"description_ru,omitempty"`
	DescriptionEn        string     `json:"description_en,omitempty"`
	ShortDescriptionRu   string     `json:"short_description_ru,omitempty"`
	ShortDescriptionEn   string     `json:"short_description_en,omitempty"`
	Price                int        `json:"price"` // in kopecks
	DurationDays         *int       `json:"duration_days,omitempty"`
	ThumbnailUrl         string     `json:"thumbnail_url,omitempty"`
	VideoUrl             string     `json:"video_url,omitempty"`
	Status               string     `json:"status"`               // "draft", "published", "archived"
	Difficulty           string     `json:"difficulty,omitempty"` // "beginner", "intermediate", "advanced"
	TotalLessons         int        `json:"total_lessons"`
	TotalDurationMinutes int        `json:"total_duration_minutes"`
	CategoryId           *int       `json:"category_id,omitempty"`
	InstructorId         *int       `json:"instructor_id,omitempty"`
	Tags                 []string   `json:"tags,omitempty"`
	PrerequisitesRu      string     `json:"prerequisites_ru,omitempty"`
	PrerequisitesEn      string     `json:"prerequisites_en,omitempty"`
	LearningOutcomesRu   string     `json:"learning_outcomes_ru,omitempty"`
	LearningOutcomesEn   string     `json:"learning_outcomes_en,omitempty"`
	CertificateAvailable bool       `json:"certificate_available"`
	MaxStudents          *int       `json:"max_students,omitempty"`
	StartDate            *time.Time `json:"start_date,omitempty"`
	Language             string     `json:"language"` // "ru", "en", "both"
	IsFeatured           bool       `json:"is_featured"`
	ViewCount            int        `json:"view_count"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`

	// Joined fields (optional)
	Category   *ProductCategory `json:"category,omitempty"`
	Instructor *User            `json:"instructor,omitempty"`
}

// TagsJSON implements sql.Scanner and driver.Valuer for JSON array
type TagsJSON []string

func (t *TagsJSON) Scan(value interface{}) error {
	if value == nil {
		*t = []string{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, t)
}

func (t TagsJSON) Value() (interface{}, error) {
	if t == nil {
		return "[]", nil
	}
	return json.Marshal(t)
}

// ProductCreateRequest - запрос на создание продукта
type ProductCreateRequest struct {
	Type                 string     `json:"type" binding:"required,oneof=course masterclass"`
	TitleRu              string     `json:"title_ru" binding:"required"`
	TitleEn              string     `json:"title_en" binding:"required"`
	DescriptionRu        string     `json:"description_ru,omitempty"`
	DescriptionEn        string     `json:"description_en,omitempty"`
	ShortDescriptionRu   string     `json:"short_description_ru,omitempty"`
	ShortDescriptionEn   string     `json:"short_description_en,omitempty"`
	Price                int        `json:"price" binding:"required,min=0"`
	DurationDays         *int       `json:"duration_days,omitempty"`
	ThumbnailUrl         string     `json:"thumbnail_url,omitempty"`
	VideoUrl             string     `json:"video_url,omitempty"`
	Status               string     `json:"status,omitempty" binding:"omitempty,oneof=draft published archived"`
	Difficulty           string     `json:"difficulty,omitempty" binding:"omitempty,oneof=beginner intermediate advanced"`
	CategoryId           *int       `json:"category_id,omitempty"`
	InstructorId         *int       `json:"instructor_id,omitempty"`
	Tags                 []string   `json:"tags,omitempty"`
	PrerequisitesRu      string     `json:"prerequisites_ru,omitempty"`
	PrerequisitesEn      string     `json:"prerequisites_en,omitempty"`
	LearningOutcomesRu   string     `json:"learning_outcomes_ru,omitempty"`
	LearningOutcomesEn   string     `json:"learning_outcomes_en,omitempty"`
	CertificateAvailable bool       `json:"certificate_available,omitempty"`
	MaxStudents          *int       `json:"max_students,omitempty"`
	StartDate            *time.Time `json:"start_date,omitempty"`
	Language             string     `json:"language,omitempty" binding:"omitempty,oneof=ru en both"`
	IsFeatured           bool       `json:"is_featured,omitempty"`
}

// ProductUpdateRequest - запрос на обновление продукта
type ProductUpdateRequest struct {
	TitleRu              *string    `json:"title_ru,omitempty"`
	TitleEn              *string    `json:"title_en,omitempty"`
	DescriptionRu        *string    `json:"description_ru,omitempty"`
	DescriptionEn        *string    `json:"description_en,omitempty"`
	ShortDescriptionRu   *string    `json:"short_description_ru,omitempty"`
	ShortDescriptionEn   *string    `json:"short_description_en,omitempty"`
	Price                *int       `json:"price,omitempty" binding:"omitempty,min=0"`
	DurationDays         *int       `json:"duration_days,omitempty"`
	ThumbnailUrl         *string    `json:"thumbnail_url,omitempty"`
	VideoUrl             *string    `json:"video_url,omitempty"`
	Status               *string    `json:"status,omitempty" binding:"omitempty,oneof=draft published archived"`
	Difficulty           *string    `json:"difficulty,omitempty" binding:"omitempty,oneof=beginner intermediate advanced"`
	CategoryId           *int       `json:"category_id,omitempty"`
	InstructorId         *int       `json:"instructor_id,omitempty"`
	Tags                 []string   `json:"tags,omitempty"`
	PrerequisitesRu      *string    `json:"prerequisites_ru,omitempty"`
	PrerequisitesEn      *string    `json:"prerequisites_en,omitempty"`
	LearningOutcomesRu   *string    `json:"learning_outcomes_ru,omitempty"`
	LearningOutcomesEn   *string    `json:"learning_outcomes_en,omitempty"`
	CertificateAvailable *bool      `json:"certificate_available,omitempty"`
	MaxStudents          *int       `json:"max_students,omitempty"`
	StartDate            *time.Time `json:"start_date,omitempty"`
	Language             *string    `json:"language,omitempty" binding:"omitempty,oneof=ru en both"`
	IsFeatured           *bool      `json:"is_featured,omitempty"`
}

// ProductFilter - фильтры для поиска продуктов
type ProductFilter struct {
	Type       string `form:"type,omitempty"`
	CategoryId *int   `form:"category_id,omitempty"`
	Status     string `form:"status,omitempty"`
	Difficulty string `form:"difficulty,omitempty"`
	Language   string `form:"language,omitempty"`
	IsFeatured *bool  `form:"is_featured,omitempty"`
	Search     string `form:"search,omitempty"`
	Limit      int    `form:"limit,omitempty"`
	Offset     int    `form:"offset,omitempty"`
}
