package main

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
}

type Painting struct {
	Id       int    `json:"id"`
	Dir      string `json:"dir"`
	Width    string `json:"width"`
	Height   string `json:"height"`
	Year     int    `json:"year"`
	NameRu   string `json:"name_ru"`
	NameEn   string `json:"name_en"`
	BaseId   string `json:"base_id"`
	StrId    string `json:"str_id"`
	ImgCount int    `json:"img_count"`
}

type Threed struct {
	Id       int    `json:"id"`
	StrId    string `json:"str_id"`
	Dir      string `json:"dir"`
	NameRu   string `json:"name_ru"`
	NameEn   string `json:"name_en"`
	BaseId   string `json:"base_id"`
	Year     int    `json:"year"`
	ImgCount int    `json:"img_count"`
}

type Illustration struct {
	Id       int    `json:"id"`
	StrId    string `json:"str_id"`
	Dir      string `json:"dir"`
	NameRu   string `json:"name_ru"`
	NameEn   string `json:"name_en"`
	BaseId   string `json:"base_id"`
	Year     int    `json:"year"`
	ImgCount int    `json:"img_count"`
}

type News struct {
	Id          int    `json:"id"`
	Datetime    string `json:"datetime"`
	TitleRu     string `json:"title_ru"`
	TitleEn     string `json:"title_en"`
	SubtitleRu  string `json:"subtitle_ru"`
	SubtitleEn  string `json:"subtitle_en"`
	Dir         string `json:"dir"`
	ImgBack     string `json:"img_back"`
	ImgBackfull string `json:"img_backfull"`
	ImagesCount int    `json:"imagescount"`
	VideosCount int    `json:"videoscount"`
	TextRu      string `json:"text_ru"`
	TextEn      string `json:"text_en"`
}
