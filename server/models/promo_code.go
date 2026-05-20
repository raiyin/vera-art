package models

import (
	"database/sql"
	"time"
)

type PromoCode struct {
	Id            int           `json:"id"`
	Code          string        `json:"code"`
	DiscountType  string        `json:"discount_type"` // "percentage", "fixed"
	DiscountValue int           `json:"discount_value"`
	MaxUses       sql.NullInt64 `json:"max_uses,omitempty"` // NULL = без ограничений
	UsedCount     int           `json:"used_count"`
	ValidFrom     sql.NullTime  `json:"valid_from,omitempty"`
	ValidUntil    sql.NullTime  `json:"valid_until,omitempty"`
	IsActive      bool          `json:"is_active"`
	CreatedAt     time.Time     `json:"created_at"`
}

// IsValid проверяет, можно ли использовать промокод
func (p *PromoCode) IsValid() bool {
	if !p.IsActive {
		return false
	}
	if p.MaxUses.Valid && p.UsedCount >= int(p.MaxUses.Int64) {
		return false
	}
	now := time.Now()
	if p.ValidFrom.Valid && now.Before(p.ValidFrom.Time) {
		return false
	}
	if p.ValidUntil.Valid && now.After(p.ValidUntil.Time) {
		return false
	}
	return true
}

// CalculateDiscount вычисляет сумму скидки для указанной цены
func (p *PromoCode) CalculateDiscount(price int) int {
	if p.DiscountType == "percentage" {
		return price * p.DiscountValue / 100
	}
	// fixed
	if p.DiscountValue > price {
		return price
	}
	return p.DiscountValue
}
