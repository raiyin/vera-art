package models

import (
	_ "github.com/mattn/go-sqlite3"
)

type Work struct {
	Id     int    `json:"id"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Year   int    `json:"year"`
	NameRu string `json:"name_ru"`
	NameEn string `json:"name_en"`
	BaseId int    `json:"base_id"`
	StrId  string `json:"str_id"`
	Descr  string `json:"descr"`
	Type   int    `json:"type"`
	Images string `json:"images"`
}
