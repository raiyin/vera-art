package models

import (
	_ "github.com/mattn/go-sqlite3"
)

type Work struct {
	Id             int      `json:"id"`
	StrId          string   `json:"str_id"`
	Dir            string   `json:"dir"`
	NameRu         string   `json:"name_ru"`
	NameEn         string   `json:"name_en"`
	BaseId         int      `json:"base_id"`
	Year           int      `json:"year"`
	ImgCount       int      `json:"img_count"`
	Descr          string   `json:"descr"`
	BaseRu         string   `json:"base_ru"`
	BaseEn         string   `json:"base_en"`
	Width          int      `json:"width"`
	Height         int      `json:"height"`
	Type           int      `json:"type"`
	Images         []string `json:"images"`
	MaterialsIds   []int    `json:"materials_ids"`
	RemovedIndices []int    `json:"removed_indices,omitempty"`
}
