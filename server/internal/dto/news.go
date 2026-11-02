package dto

type NewsResponse struct {
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

type NewsDataPayload struct {
	ID       string   `json:"id"`
	Datetime string   `json:"datetime"`
	TitleRu  string   `json:"title_ru"`
	TitleEn  string   `json:"title_en"`
	Dir      string   `json:"dir"`
	MainImage  string   `json:"main_image"`
	TextRu   string   `json:"text_ru"`
	TextEn   string   `json:"text_en"`
	Images   []string `json:"images"`
	Videos   []string `json:"videos"`
}
