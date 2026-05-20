package models

import (
	"database/sql"
	"time"
)

type Purchase struct {
	Id           int           `json:"id"`
	UserId       int           `json:"user_id"`
	ProductId    int           `json:"product_id"`
	PurchaseDate time.Time     `json:"purchase_date"`
	AccessStart  time.Time     `json:"access_start"`
	AccessEnd    sql.NullTime  `json:"access_end,omitempty"`
	Status       string        `json:"status"` // "active", "expired", "cancelled"
	PaymentId    sql.NullInt64 `json:"payment_id,omitempty"`
	PromoCodeId  sql.NullInt64 `json:"promo_code_id,omitempty"`
	PricePaid    int           `json:"price_paid"` // фактически уплаченная сумма
	CreatedAt    time.Time     `json:"created_at"`

	// Joined fields
	User      *User      `json:"user,omitempty"`
	Product   *Product   `json:"product,omitempty"`
	Payment   *Payment   `json:"payment,omitempty"`
	PromoCode *PromoCode `json:"promo_code,omitempty"`
}

// IsActive проверяет, активна ли покупка в данный момент
func (p *Purchase) IsActive() bool {
	if p.Status != "active" {
		return false
	}
	if p.AccessEnd.Valid && time.Now().After(p.AccessEnd.Time) {
		return false
	}
	return true
}

// DaysRemaining возвращает количество оставшихся дней доступа (или -1 для бессрочного)
func (p *Purchase) DaysRemaining() int {
	if !p.AccessEnd.Valid {
		return -1 // бессрочный доступ
	}
	remaining := p.AccessEnd.Time.Sub(time.Now())
	days := int(remaining.Hours() / 24)
	if days < 0 {
		return 0
	}
	return days
}
