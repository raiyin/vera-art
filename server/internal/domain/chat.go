package domain

import "time"

type ChatThread struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Subject   string    `json:"subject,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChatMessage struct {
	ID        int64      `json:"id"`
	ThreadID  int64      `json:"thread_id"`
	UserID    int64      `json:"user_id"`
	Content   string     `json:"content"`
	IsAdmin   bool       `json:"is_admin"`
	IsRead    bool       `json:"is_read"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}
