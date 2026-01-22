package models

import (
	_ "github.com/mattn/go-sqlite3"
)

type News struct {
	Id          int    `json:"id"`
	Datetime    string `json:"datetime"`
	TitleRu     string `json:"title_ru"`
	TitleEn     string `json:"title_en"`
	SubtitleRu  string `json:"subTitle_ru"`
	SubtitleEn  string `json:"subTitle_en"`
	Dir         string `json:"dir"`
	ImgBack     string `json:"img_back"`
	ImgBackfull string `json:"img_backfull"`
	ImagesCount int    `json:"imagescount"`
	VideosCount int    `json:"videoscount"`
	TextRu      string `json:"text_ru"`
	TextEn      string `json:"text_en"`
}
