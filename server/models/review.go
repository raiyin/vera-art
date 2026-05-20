package models

import (
	"time"
)

type Review struct {
	Id         int       `json:"id"`
	UserId     int       `json:"user_id"`
	ProductId  int       `json:"product_id"`
	PurchaseId int       `json:"purchase_id"`
	Rating     int       `json:"rating"` // 1-5
	TitleRu    string    `json:"title_ru,omitempty"`
	TitleEn    string    `json:"title_en,omitempty"`
	CommentRu  string    `json:"comment_ru,omitempty"`
	CommentEn  string    `json:"comment_en,omitempty"`
	IsApproved bool      `json:"is_approved"` // модерация админом
	IsVisible  bool      `json:"is_visible"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Joined fields
	User     *User     `json:"user,omitempty"`
	Product  *Product  `json:"product,omitempty"`
	Purchase *Purchase `json:"purchase,omitempty"`
}
