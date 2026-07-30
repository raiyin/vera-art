package dto

type NewsResponse struct {
	ID          string   `json:"id"`
	DateTime    string   `json:"datetime"`
	TitleRu     string   `json:"title_ru"`
	TitleEn     string   `json:"title_en"`
	SubTitleRu  string   `json:"sub_title_ru,omitempty"`
	SubTitleEn  string   `json:"sub_title_en,omitempty"`
	Dir         string   `json:"dir"`
	ImgBack     string   `json:"img_back"`
	ImgBackfull string   `json:"img_backfull"`
	TextRu      string   `json:"text_ru"`
	TextEn      string   `json:"text_en"`
	Images      []string `json:"images,omitempty"`
	Videos      []string `json:"videos,omitempty"`
}

type NewsDataPayload struct {
	ID         string   `json:"id"`
	Datetime   string   `json:"datetime"`
	TitleRu    string   `json:"title_ru"`
	TitleEn    string   `json:"title_en"`
	SubTitleRu string   `json:"sub_title_ru"`
	SubTitleEn string   `json:"sub_title_en"`
	Dir        string   `json:"dir"`
	ImgBack    string   `json:"img_back"`
	ImgBackfull string  `json:"img_backfull"`
	TextRu     string   `json:"text_ru"`
	TextEn     string   `json:"text_en"`
	Images     []string `json:"images"`
	Videos     []string `json:"videos"`
}
