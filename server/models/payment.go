package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Payment struct {
	Id            int            `json:"id"`
	UserId        int            `json:"user_id"`
	ExternalId    sql.NullString `json:"external_id,omitempty"` // ID платежа в ЮKassa
	Status        string         `json:"status"`                // "pending", "waiting_for_capture", "succeeded", "canceled", "refunded"
	Amount        int            `json:"amount"`                // в копейках
	Currency      string         `json:"currency"`              // "RUB"
	Description   sql.NullString `json:"description,omitempty"`
	PaymentMethod sql.NullString `json:"payment_method,omitempty"`
	Metadata      map[string]any `json:"metadata,omitempty"` // дополнительные данные
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`

	// Joined fields
	User *User `json:"user,omitempty"`
}

// MetadataJSON implements sql.Scanner and driver.Valuer for JSON map
type MetadataJSON map[string]any

func (m *MetadataJSON) Scan(value any) error {
	if value == nil {
		*m = nil
		return nil
	}
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return nil
	}
	return json.Unmarshal(data, m)
}

func (m MetadataJSON) Value() (any, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}
