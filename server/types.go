package main

type Sale struct {
	Id     int    `json:"id"`
	Dir    string `json:"dir"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Year   int    `json:"year"`
	Price  int    `json:"price"`
	NameRu string `json:"name_ru"`
	NameEn string `json:"name_en"`
	BaseId string `json:"base_id"`
	StrId  string `json:"str_id"`
}

type Painting struct {
	Id     int    `json:"id"`
	Dir    string `json:"dir"`
	Width  string `json:"width"`
	Height string `json:"height"`
	Year   int    `json:"year"`
	NameRu string `json:"name_ru"`
	NameEn string `json:"name_en"`
	BaseId string `json:"base_id"`
	StrId  string `json:"str_id"`
}

type Threed struct {
	Id     int    `json:"id"`
	StrId  string `json:"str_id"`
	Dir    string `json:"dir"`
	NameRu string `json:"name_ru"`
	NameEn string `json:"name_en"`
	BaseId string `json:"base_id"`
	Year   int    `json:"year"`
}

type Illustration struct {
	Id     int    `json:"id"`
	StrId  string `json:"str_id"`
	Dir    string `json:"dir"`
	NameRu string `json:"name_ru"`
	NameEn string `json:"name_en"`
	BaseId string `json:"base_id"`
	Year   int    `json:"year"`
}
