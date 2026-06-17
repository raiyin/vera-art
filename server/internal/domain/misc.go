package domain

import "time"

type Tag struct {
	ID        int64     `json:"id"`
	NameRu    string    `json:"name_ru"`
	NameEn    string    `json:"name_en"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

type Material struct {
	ID        int64     `json:"id"`
	NameRu    string    `json:"name_ru"`
	NameEn    string    `json:"name_en"`
	CreatedAt time.Time `json:"created_at"`
}

type Base struct {
	ID        int64     `json:"id"`
	NameRu    string    `json:"name_ru"`
	NameEn    string    `json:"name_en"`
	CreatedAt time.Time `json:"created_at"`
}

type UserConsent struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	ConsentType string    `json:"consent_type"`
	Granted     bool      `json:"granted"`
	IPAddress   string    `json:"ip_address,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
