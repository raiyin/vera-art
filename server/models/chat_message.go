package models

import (
	"database/sql"
	"time"
)

type ChatMessage struct {
	Id             int            `json:"id"`
	ThreadId       int            `json:"thread_id"`
	SenderId       int            `json:"sender_id"`
	MessageType    string         `json:"message_type"` // "text", "image", "file"
	Content        string         `json:"content"`
	AttachmentUrl  sql.NullString `json:"attachment_url,omitempty"`
	AttachmentSize sql.NullInt64  `json:"attachment_size,omitempty"`
	IsRead         bool           `json:"is_read"`
	ReadAt         sql.NullTime   `json:"read_at,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`

	// Joined fields
	Thread *ChatThread `json:"thread,omitempty"`
	Sender *User       `json:"sender,omitempty"`
}
