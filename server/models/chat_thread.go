package models

import (
	"database/sql"
	"time"
)

type ChatThread struct {
	Id            int           `json:"id"`
	PurchaseId    int           `json:"purchase_id"`
	UserId        int           `json:"user_id"`
	AdminId       sql.NullInt64 `json:"admin_id,omitempty"`
	LastMessageAt time.Time     `json:"last_message_at"`
	IsResolved    bool          `json:"is_resolved"`
	CreatedAt     time.Time     `json:"created_at"`

	// Joined fields
	Purchase *Purchase     `json:"purchase,omitempty"`
	User     *User         `json:"user,omitempty"`
	Admin    *User         `json:"admin,omitempty"`
	Messages []ChatMessage `json:"messages,omitempty"`
}
