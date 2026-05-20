package models

import (
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ProductCategory struct {
	Id            int       `json:"id"`
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

// ProductCategoryCreateRequest - запрос на создание категории
type ProductCategoryCreateRequest struct {
	NameRu        string `json:"name_ru" binding:"required"`
	NameEn        string `json:"name_en" binding:"required"`
	Slug          string `json:"slug" binding:"required"`
	DescriptionRu string `json:"description_ru,omitempty"`
	DescriptionEn string `json:"description_en,omitempty"`
	SortOrder     int    `json:"sort_order,omitempty"`
	IsActive      bool   `json:"is_active,omitempty"`
}

// ProductCategoryUpdateRequest - запрос на обновление категории
type ProductCategoryUpdateRequest struct {
	NameRu        string `json:"name_ru,omitempty"`
	NameEn        string `json:"name_en,omitempty"`
	Slug          string `json:"slug,omitempty"`
	DescriptionRu string `json:"description_ru,omitempty"`
	DescriptionEn string `json:"description_en,omitempty"`
	SortOrder     *int   `json:"sort_order,omitempty"`
	IsActive      *bool  `json:"is_active,omitempty"`
}
