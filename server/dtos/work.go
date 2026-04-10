package dtos

import (
	_ "github.com/mattn/go-sqlite3"
)

type GetWorkDto struct {
	Id          int      `json:"id"`
	StrId       string   `json:"str_id"`
	Dir         string   `json:"dir"`
	NameRu      string   `json:"name_ru"`
	NameEn      string   `json:"name_en"`
	Year        int      `json:"year"`
	Descr       string   `json:"descr"`
	BaseRu      string   `json:"base_ru"`
	BaseEn      string   `json:"base_en"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Type        int      `json:"type"`
	Images      []string `json:"images"`
	MaterialsEn []string `json:"materials_en"`
	MaterialsRu []string `json:"materials_ru"`
}
