package dto

import (
	"strings"
	"time"
)

// WorkResponse represents a work in API responses.
type WorkResponse struct {
	ID          int64    `json:"id"`
	StrID       string   `json:"str_id"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Year        int      `json:"year"`
	NameRu      string   `json:"name_ru"`
	NameEn      string   `json:"name_en"`
	BaseID      int64    `json:"base_id"`
	DescrRu     string   `json:"descr_ru,omitempty"`
	DescrEn     string   `json:"descr_en,omitempty"`
	WorkPath    string   `json:"work_path"`
	Dir         string   `json:"dir"`
	Images      []string `json:"images,omitempty"`
	MaterialIDs []int64  `json:"material_ids,omitempty"`
}

// UpdateWorkResponse represents a work in update API responses.
type UpdateWorkResponse struct {
	ID          int64    `json:"id"`
	StrID       string   `json:"str_id"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Year        int      `json:"year"`
	NameRu      string   `json:"name_ru"`
	NameEn      string   `json:"name_en"`
	BaseID      int64    `json:"base_id"`
	DescrRu     string   `json:"descr_ru,omitempty"`
	DescrEn     string   `json:"descr_en,omitempty"`
	WorkPath    string   `json:"work_path"`
	Dir         string   `json:"dir"`
	Images      []string `json:"images,omitempty"`
	MaterialIDs []int64  `json:"material_ids,omitempty"`
}

// CreateWorkRequest represents a create work request.
type CreateWorkRequest struct {
	StrID       string  `json:"str_id" form:"str_id"`
	Width       int     `json:"width" form:"width"`
	Height      int     `json:"height" form:"height"`
	Year        int     `json:"year" form:"year"`
	NameRu      string  `json:"name_ru" form:"name_ru" binding:"required"`
	NameEn      string  `json:"name_en" form:"name_en" binding:"required"`
	BaseID      int64   `json:"base_id" form:"base_id"`
	DescrRu     string  `json:"descr_ru" form:"descr_ru"`
	DescrEn     string  `json:"descr_en" form:"descr_en"`
	MaterialIDs []int64 `json:"material_ids" form:"material_ids"`
}

// UpdateWorkRequest represents an update work request.
type UpdateWorkRequest struct {
	StrID       string   `json:"str_id" form:"str_id"`
	Width       int      `json:"width" form:"width"`
	Height      int      `json:"height" form:"height"`
	Year        int      `json:"year" form:"year"`
	NameRu      string   `json:"name_ru" form:"name_ru"`
	NameEn      string   `json:"name_en" form:"name_en"`
	BaseID      int64    `json:"base_id" form:"base_id"`
	DescrRu     string   `json:"descr_ru" form:"descr_ru"`
	DescrEn     string   `json:"descr_en" form:"descr_en"`
	MaterialIDs []int64  `json:"material_ids" form:"material_ids"`
	Images      []string `json:"images" form:"images"`
}

// SplitImages splits a semicolon-separated images string into a slice.
func SplitImages(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ";")
}

// JoinImages joins a slice of image filenames into a semicolon-separated string.
func JoinImages(imgs []string) string {
	return strings.Join(imgs, ";")
}

// SaleResponse represents a sale in API responses.
type SaleResponse struct {
	ID          int64     `json:"id"`
	NameRu      string    `json:"name_ru"`
	NameEn      string    `json:"name_en"`
	DescrRu     string    `json:"descr_ru,omitempty"`
	DescrEn     string    `json:"descr_en,omitempty"`
	ImagePath   string    `json:"image_path"`
	SalePath    string    `json:"sale_path"`
	Dir         string    `json:"dir"`
	Images      []string  `json:"images"`
	Price       float64   `json:"price"`
	Year        int       `json:"year,omitempty"`
	Technique   string    `json:"technique,omitempty"`
	Width       int       `json:"width,omitempty"`
	Height      int       `json:"height,omitempty"`
	Status      string    `json:"status"`
	SortOrder   int       `json:"sort_order"`
	Sold        bool      `json:"sold"`
	MaterialIDs []int64   `json:"material_ids,omitempty"`
	BaseIDs     []int64   `json:"base_ids,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UpdateSaleResponse represents a sale in update API responses.
type UpdateSaleResponse struct {
	ID          int64     `json:"id"`
	NameRu      string    `json:"name_ru"`
	NameEn      string    `json:"name_en"`
	DescrRu     string    `json:"descr_ru,omitempty"`
	DescrEn     string    `json:"descr_en,omitempty"`
	ImagePath   string    `json:"image_path"`
	SalePath    string    `json:"sale_path"`
	Dir         string    `json:"dir"`
	Images      []string  `json:"images"`
	Price       float64   `json:"price"`
	Year        int       `json:"year,omitempty"`
	Technique   string    `json:"technique,omitempty"`
	Width       int       `json:"width,omitempty"`
	Height      int       `json:"height,omitempty"`
	Status      string    `json:"status"`
	SortOrder   int       `json:"sort_order"`
	Sold        bool      `json:"sold"`
	MaterialIDs []int64   `json:"material_ids,omitempty"`
	BaseIDs     []int64   `json:"base_ids,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateSaleRequest represents a create sale request.
type CreateSaleRequest struct {
	NameRu      string  `json:"name_ru" form:"name_ru" binding:"required"`
	NameEn      string  `json:"name_en" form:"name_en" binding:"required"`
	DescrRu     string  `json:"descr_ru" form:"descr_ru"`
	DescrEn     string  `json:"descr_en" form:"descr_en"`
	Price       float64 `json:"price" form:"price" binding:"required"`
	Year        int     `json:"year" form:"year"`
	Technique   string  `json:"technique" form:"technique"`
	Width       int     `json:"width" form:"width"`
	Height      int     `json:"height" form:"height"`
	Status      string  `json:"status" form:"status"`
	SortOrder   int     `json:"sort_order" form:"sort_order"`
	Sold        bool    `json:"sold" form:"sold"`
	MaterialIDs []int64 `json:"material_ids" form:"material_ids"`
	BaseIDs     []int64 `json:"base_ids" form:"base_ids"`
}

// UpdateSaleRequest represents an update sale request.
type UpdateSaleRequest struct {
	NameRu      string  `json:"name_ru" form:"name_ru"`
	NameEn      string  `json:"name_en" form:"name_en"`
	DescrRu     string  `json:"descr_ru" form:"descr_ru"`
	DescrEn     string  `json:"descr_en" form:"descr_en"`
	Price       float64 `json:"price" form:"price"`
	Year        int     `json:"year" form:"year"`
	Technique   string  `json:"technique" form:"technique"`
	Width       int     `json:"width" form:"width"`
	Height      int     `json:"height" form:"height"`
	Status      string  `json:"status" form:"status"`
	SortOrder   int     `json:"sort_order" form:"sort_order"`
	Sold        bool    `json:"sold" form:"sold"`
	MaterialIDs []int64 `json:"material_ids" form:"material_ids"`
	BaseIDs     []int64 `json:"base_ids" form:"base_ids"`
}
