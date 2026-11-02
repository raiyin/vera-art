package domain

type News struct {
	ID       string   `json:"id"`
	DateTime string   `json:"datetime"`
	TitleRu  string   `json:"title_ru"`
	TitleEn  string   `json:"title_en"`
	Dir      string   `json:"dir"`
	MainImage  string   `json:"main_image"`
	TextRu   string   `json:"text_ru"`
	TextEn   string   `json:"text_en"`
	Images   []string `json:"images,omitempty"`
	Videos   []string `json:"videos,omitempty"`
}

type NewsFilter struct {
	Status string
	Query  string
	Page   int
	Limit  int
}
