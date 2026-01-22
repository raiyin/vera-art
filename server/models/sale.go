package models

import (
	_ "github.com/mattn/go-sqlite3"
)

type Sale struct {
	Id           int    `json:"id"`
	Dir          string `json:"dir"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Year         int    `json:"year"`
	Price        int    `json:"price"`
	NameRu       string `json:"name_ru"`
	NameEn       string `json:"name_en"`
	BaseId       int    `json:"base_id"`
	StrId        string `json:"str_id"`
	ImgCount     int    `json:"img_count"`
	Descr        string `json:"descr"`
	MaterialsIds []int  `json:"materials_ids"`
	RemovedIndices []int `json:"removed_indices,omitempty"`
}

type SaleWithBase struct {
	Id       int    `json:"id"`
	Dir      string `json:"dir"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Year     int    `json:"year"`
	Price    int    `json:"price"`
	NameRu   string `json:"name_ru"`
	NameEn   string `json:"name_en"`
	BaseId   int    `json:"base_id"`
	StrId    string `json:"str_id"`
	ImgCount int    `json:"img_count"`
	Descr    string `json:"descr"`
	BaseRu   string `json:"base_ru"`
	BaseEn   string `json:"base_en"`
}
