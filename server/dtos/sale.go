package dtos

import (
	_ "github.com/mattn/go-sqlite3"
)

type GetSaleDto struct {
	Id          int      `json:"id"`
	Dir         string   `json:"dir"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Year        int      `json:"year"`
	Price       int      `json:"price"`
	NameRu      string   `json:"name_ru"`
	NameEn      string   `json:"name_en"`
	StrId       string   `json:"str_id"`
	Descr       string   `json:"descr"`
	Images      []string `json:"images"`
	BaseRu      string   `json:"base_ru"`
	BaseEn      string   `json:"base_en"`
	MaterialsEn []string `json:"materials_en"`
	MaterialsRu []string `json:"materials_ru"`
}

// EditSaleDto is the DTO for editing a sale from backend to frontend.
type EditSaleDto struct {
	Id           int      `json:"id"`
	Dir          string   `json:"dir"`
	Width        int      `json:"width"`
	Height       int      `json:"height"`
	Year         int      `json:"year"`
	Price        int      `json:"price"`
	NameRu       string   `json:"name_ru"`
	NameEn       string   `json:"name_en"`
	StrId        string   `json:"str_id"`
	BaseId       int      `json:"base_id"`
	Descr        string   `json:"descr"`
	Images       []string `json:"images"`
	MaterialsIds []int    `json:"materials_ids"`
}

// UpdateEditSaleDto is the DTO for updating a sale from frontend.
type UpdateSaleDto struct {
	Id           int      `json:"id"`
	Dir          string   `json:"dir"`
	Width        int      `json:"width"`
	Height       int      `json:"height"`
	Year         int      `json:"year"`
	Price        int      `json:"price"`
	NameRu       string   `json:"name_ru"`
	NameEn       string   `json:"name_en"`
	StrId        string   `json:"str_id"`
	BaseId       int      `json:"base_id"`
	Descr        string   `json:"descr"`
	Images       []string `json:"images"`
	MaterialsIds []int    `json:"materials_ids"`
}

// CreateSaleDto is the DTO for creating a sale from frontend.
type CreateSaleDto struct {
	Dir          string   `json:"dir"`
	Width        int      `json:"width"`
	Height       int      `json:"height"`
	Year         int      `json:"year"`
	Price        int      `json:"price"`
	NameRu       string   `json:"name_ru"`
	NameEn       string   `json:"name_en"`
	StrId        string   `json:"str_id"`
	BaseId       int      `json:"base_id"`
	Descr        string   `json:"descr"`
	Images       []string `json:"images"`
	MaterialsIds []int    `json:"materials_ids"`
}
