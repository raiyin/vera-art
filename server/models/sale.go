package models

import (
	_ "github.com/mattn/go-sqlite3"
)

type Sale struct {
	Id      int    `json:"id"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	Year    int    `json:"year"`
	Price   int    `json:"price"`
	NameRu  string `json:"name_ru"`
	NameEn  string `json:"name_en"`
	BaseId  int    `json:"base_id"`
	StrId   string `json:"str_id"`
	DescrRu string `json:"descr_ru"`
	DescrEn string `json:"descr_en"`
	Images  string `json:"images"`
}
