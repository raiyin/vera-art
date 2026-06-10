package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// PromoCodeResponse is a JSON-friendly representation of PromoCode
// that properly serializes sql.Null* types as plain JSON values.
type PromoCodeResponse struct {
	Id            int        `json:"id"`
	Code          string     `json:"code"`
	DiscountType  string     `json:"discount_type"`
	DiscountValue int        `json:"discount_value"`
	MaxUses       *int64     `json:"max_uses"`
	UsedCount     int        `json:"used_count"`
	ValidFrom     *time.Time `json:"valid_from"`
	ValidUntil    *time.Time `json:"valid_until"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
}

// ToResponse converts a PromoCode to a JSON-safe PromoCodeResponse.
func (p *PromoCode) ToResponse() PromoCodeResponse {
	r := PromoCodeResponse{
		Id:            p.Id,
		Code:          p.Code,
		DiscountType:  p.DiscountType,
		DiscountValue: p.DiscountValue,
		UsedCount:     p.UsedCount,
		IsActive:      p.IsActive,
		CreatedAt:     p.CreatedAt,
	}
	if p.MaxUses.Valid {
		r.MaxUses = &p.MaxUses.Int64
	}
	if p.ValidFrom.Valid {
		r.ValidFrom = &p.ValidFrom.Time
	}
	if p.ValidUntil.Valid {
		r.ValidUntil = &p.ValidUntil.Time
	}
	return r
}

// MarshalJSON implements json.Marshaler for PromoCode.
func (p *PromoCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.ToResponse())
}

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
