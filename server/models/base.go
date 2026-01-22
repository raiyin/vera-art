package models

import (
	_ "github.com/mattn/go-sqlite3"
)

type Base struct {
	Id     int    `json:"id"`
	BaseRu string `json:"base_ru"`
	BaseEn string `json:"base_en"`
}
