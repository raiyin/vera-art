package models

import (
	"encoding/json"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Lesson struct {
	Id                 int        `json:"id"`
	ProductId          int        `json:"product_id"`
	TitleRu            string     `json:"title_ru"`
	TitleEn            string     `json:"title_en"`
	DescriptionRu      string     `json:"description_ru,omitempty"`
	DescriptionEn      string     `json:"description_en,omitempty"`
	ContentType        string     `json:"content_type"` // "video", "text", "pdf", "quiz", "assignment"
	ContentUrl         string     `json:"content_url,omitempty"`
	DurationMinutes    int        `json:"duration_minutes"`
	SortOrder          int        `json:"sort_order"`
	IsPreview          bool       `json:"is_preview"`
	Resources          []Resource `json:"resources,omitempty"`
	HomeworkRu         string     `json:"homework_ru,omitempty"`
	HomeworkEn         string     `json:"homework_en,omitempty"`
	EstimatedStudyTime *int       `json:"estimated_study_time,omitempty"`
	IsRequired         bool       `json:"is_required"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type Resource struct {
	Type  string `json:"type"` // "pdf", "link", "image", "audio"
	Url   string `json:"url"`
	Title string `json:"title"`
}

// ResourcesJSON implements sql.Scanner and driver.Valuer for JSON array
type ResourcesJSON []Resource

func (r *ResourcesJSON) Scan(value interface{}) error {
	if value == nil {
		*r = []Resource{}
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(b, r)
}

func (r ResourcesJSON) Value() (interface{}, error) {
	if r == nil {
		return "[]", nil
	}
	return json.Marshal(r)
}

// LessonCreateRequest - запрос на создание урока
type LessonCreateRequest struct {
	ProductId          int        `json:"product_id" binding:"required"`
	TitleRu            string     `json:"title_ru" binding:"required"`
	TitleEn            string     `json:"title_en" binding:"required"`
	DescriptionRu      string     `json:"description_ru,omitempty"`
	DescriptionEn      string     `json:"description_en,omitempty"`
	ContentType        string     `json:"content_type" binding:"required,oneof=video text pdf quiz assignment"`
	ContentUrl         string     `json:"content_url,omitempty"`
	DurationMinutes    int        `json:"duration_minutes,omitempty"`
	SortOrder          int        `json:"sort_order,omitempty"`
	IsPreview          bool       `json:"is_preview,omitempty"`
	Resources          []Resource `json:"resources,omitempty"`
	HomeworkRu         string     `json:"homework_ru,omitempty"`
	HomeworkEn         string     `json:"homework_en,omitempty"`
	EstimatedStudyTime *int       `json:"estimated_study_time,omitempty"`
	IsRequired         bool       `json:"is_required,omitempty"`
}

// LessonUpdateRequest - запрос на обновление урока
type LessonUpdateRequest struct {
	TitleRu            *string    `json:"title_ru,omitempty"`
	TitleEn            *string    `json:"title_en,omitempty"`
	DescriptionRu      *string    `json:"description_ru,omitempty"`
	DescriptionEn      *string    `json:"description_en,omitempty"`
	ContentType        *string    `json:"content_type,omitempty" binding:"omitempty,oneof=video text pdf quiz assignment"`
	ContentUrl         *string    `json:"content_url,omitempty"`
	DurationMinutes    *int       `json:"duration_minutes,omitempty"`
	SortOrder          *int       `json:"sort_order,omitempty"`
	IsPreview          *bool      `json:"is_preview,omitempty"`
	Resources          []Resource `json:"resources,omitempty"`
	HomeworkRu         *string    `json:"homework_ru,omitempty"`
	HomeworkEn         *string    `json:"homework_en,omitempty"`
	EstimatedStudyTime *int       `json:"estimated_study_time,omitempty"`
	IsRequired         *bool      `json:"is_required,omitempty"`
}
