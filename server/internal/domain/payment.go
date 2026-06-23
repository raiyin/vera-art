package domain

import "time"

type Payment struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	ProductID     int64     `json:"product_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	PaymentMethod string    `json:"payment_method,omitempty"`
	YooKassaID    string    `json:"yookassa_id,omitempty"`
	PromoCode     string    `json:"promo_code,omitempty"`
	Discount      float64   `json:"discount,omitempty"`
	Metadata      string    `json:"metadata,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Purchase struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	ProductID   int64      `json:"product_id"`
	PaymentID   int64      `json:"payment_id"`
	PricePaid   int64      `json:"price_paid"` // в копейках
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	AccessStart time.Time  `json:"access_start"`
	AccessEnd   *time.Time `json:"access_end,omitempty"`
}

type PromoCode struct {
	ID              int64     `json:"id"`
	Code            string    `json:"code"`
	DiscountPercent float64   `json:"discount_percent"`
	MaxUses         int       `json:"max_uses"`
	CurrentUses     int       `json:"current_uses"`
	ExpiresAt       time.Time `json:"expires_at"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
}

type Review struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	Rating    int       `json:"rating"`
	Text      string    `json:"text,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
