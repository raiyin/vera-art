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
