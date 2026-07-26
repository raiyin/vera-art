package domain

type Work struct {
	ID          int64   `json:"id"`
	StrID       string  `json:"str_id"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	Year        int     `json:"year"`
	NameRu      string  `json:"name_ru"`
	NameEn      string  `json:"name_en"`
	BaseID      int64   `json:"base_id"`
	DescrRu     string  `json:"descr_ru,omitempty"`
	DescrEn     string  `json:"descr_en,omitempty"`
	WorkPath    string  `json:"work_path"`
	Images      string  `json:"images,omitempty"`
	MaterialIDs []int64 `json:"material_ids,omitempty"`
}

type WorkFilter struct {
	Query string
	Page  int
	Limit int
}
