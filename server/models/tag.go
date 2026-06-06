package models

import "time"

// Tag represents a tag for products (many-to-many via product_tags)
type Tag struct {
	Id        int       `json:"id"`
	NameRu    string    `json:"name_ru"`
	NameEn    string    `json:"name_en"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

// ProductTag is the junction record linking products to tags
type ProductTag struct {
	Id        int       `json:"id"`
	ProductId int       `json:"product_id"`
	TagId     int       `json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`

	// Joined fields
	Tag *Tag `json:"tag,omitempty"`
}

// MasterClassResponse — ответ для публичного API мастер-классов.
// Собирается из products + product_tags + tags.
type MasterClassResponse struct {
	Id                 int      `json:"id"`
	TitleRu            string   `json:"title_ru"`
	TitleEn            string   `json:"title_en"`
	DescriptionRu      string   `json:"description_ru"`
	DescriptionEn      string   `json:"description_en"`
	ShortDescriptionRu string   `json:"short_description_ru"`
	ShortDescriptionEn string   `json:"short_description_en"`
	Price              int      `json:"price"` // in kopecks
	IsFree             bool     `json:"is_free"`
	ThumbnailUrl       string   `json:"thumbnail_url"`
	VideoUrl           string   `json:"video_url"`
	DurationMinutes    int      `json:"duration_minutes"`
	Difficulty         string   `json:"difficulty"`
	CategoryId         int      `json:"category_id"`
	Tags               []TagDTO `json:"tags"`
	IsFeatured         bool     `json:"is_featured"`
	ViewCount          int      `json:"view_count"`
}

// TagDTO — упрощённое представление тега для клиента
type TagDTO struct {
	NameRu string `json:"name_ru"`
	NameEn string `json:"name_en"`
	Slug   string `json:"slug"`
}
