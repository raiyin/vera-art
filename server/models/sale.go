package models

import (
	_ "github.com/mattn/go-sqlite3"
)

type Sale struct {
	Id       int    `json:"id"`
	Dir      string `json:"dir"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	Year     int    `json:"year"`
	Price    int    `json:"price"`
	NameRu   string `json:"name_ru"`
	NameEn   string `json:"name_en"`
	BaseId   string `json:"base_id"`
	StrId    string `json:"str_id"`
	ImgCount int    `json:"img_count"`
	BaseRu   string `json:"base_ru"`
	BaseEn   string `json:"base_en"`
}
