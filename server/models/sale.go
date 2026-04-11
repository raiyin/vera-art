package models

import (
	_ "github.com/mattn/go-sqlite3"
)

type Sale struct {
	Id           int      `json:"id"`
	Dir          string   `json:"dir"`
	Width        int      `json:"width"`
	Height       int      `json:"height"`
	Year         int      `json:"year"`
	Price        int      `json:"price"`
	NameRu       string   `json:"name_ru"`
	NameEn       string   `json:"name_en"`
	BaseId       int      `json:"base_id"`
	StrId        string   `json:"str_id"`
	Descr        string   `json:"descr"`
	Images       []string `json:"images"`
	MaterialsIds []int    `json:"materials_ids"`
}
