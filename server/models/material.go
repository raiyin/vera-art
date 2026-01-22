package models

import (
	_ "github.com/mattn/go-sqlite3"
)

type Material struct {
	Id         int    `json:"id"`
	MaterialRu string `json:"material_ru"`
	MaterialEn string `json:"material_en"`
}
