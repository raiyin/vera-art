package dtos

import (
	_ "github.com/mattn/go-sqlite3"
)

type SaleResponse struct {
	Id          int      `json:"id"`
	StrId       string   `json:"str_id"`
	Dir         string   `json:"dir"`
	NameRu      string   `json:"name_ru"`
	NameEn      string   `json:"name_en"`
	Year        int      `json:"year"`
	DescrRu     string   `json:"descr_ru"`
	DescrEn     string   `json:"descr_en"`
	BaseRu      string   `json:"base_ru"`
	BaseEn      string   `json:"base_en"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Price       int      `json:"price"`
	Images      []string `json:"images"`
	MaterialsEn []string `json:"materials_en"`
	MaterialsRu []string `json:"materials_ru"`
}

// Служит для отображения информации о картине на продажу на странице редактирования.
type UpdateSaleResponse struct {
	Id           int      `json:"id"`
	StrId        string   `json:"str_id"`
	Dir          string   `json:"dir"`
	NameRu       string   `json:"name_ru"`
	NameEn       string   `json:"name_en"`
	BaseId       int      `json:"base_id"`
	Year         int      `json:"year"`
	DescrRu      string   `json:"descr_ru"`
	DescrEn      string   `json:"descr_en"`
	Width        int      `json:"width"`
	Height       int      `json:"height"`
	Price        int      `json:"price"`
	Images       []string `json:"images"`
	MaterialsIds []int    `json:"materials_ids"`
}

// UpdateSaleRequest is the DTO for updating a sale from frontend.
type UpdateSaleRequest struct {
	Id           int      `json:"id"`
	NameRu       string   `json:"name_ru"`
	NameEn       string   `json:"name_en"`
	BaseId       int      `json:"base_id"`
	Year         int      `json:"year"`
	DescrRu      string   `json:"descr_ru"`
	DescrEn      string   `json:"descr_en"`
	Width        int      `json:"width"`
	Height       int      `json:"height"`
	Price        int      `json:"price"`
	Images       []string `json:"images"`
	MaterialsIds []int    `json:"materials_ids"`
}

// CreateSaleRequest is the DTO for creating a sale from frontend.
type CreateSaleRequest struct {
	NameRu       string   `json:"name_ru"`
	NameEn       string   `json:"name_en"`
	BaseId       int      `json:"base_id"`
	Year         int      `json:"year"`
	DescrRu      string   `json:"descr_ru"`
	DescrEn      string   `json:"descr_en"`
	Width        int      `json:"width"`
	Height       int      `json:"height"`
	Price        int      `json:"price"`
	Images       []string `json:"images"`
	MaterialsIds []int    `json:"materials_ids"`
}
