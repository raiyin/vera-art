package models

import (
	_ "github.com/mattn/go-sqlite3"
)

type Illustration struct {
	Id       int    `json:"id"`
	StrId    string `json:"str_id"`
	Dir      string `json:"dir"`
	NameRu   string `json:"name_ru"`
	NameEn   string `json:"name_en"`
	BaseId   string `json:"base_id"`
	Year     int    `json:"year"`
	ImgCount int    `json:"img_count"`
	Desc     string `json:"desc"`
	BaseRu   string `json:"base_ru"`
	BaseEn   string `json:"base_en"`
}
