package dto

import "time"

// ChatThreadResponse represents a chat thread in API responses.
type ChatThreadResponse struct {
	ID            int64         `json:"id"`
	PurchaseID    int64         `json:"purchase_id"`
	UserID        int64         `json:"user_id"`
	AdminID       *int64        `json:"admin_id,omitempty"`
	LastMessageAt time.Time     `json:"last_message_at"`
	IsResolved    bool          `json:"is_resolved"`
	CreatedAt     time.Time     `json:"created_at"`
	User          *ChatUser     `json:"user,omitempty"`
	Purchase      *ChatPurchase `json:"purchase,omitempty"`
}

// ChatUser represents a user in chat responses.
type ChatUser struct {
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

// ChatPurchase represents a purchase in chat responses.
type ChatPurchase struct {
	ProductID   int64  `json:"product_id"`
	ProductName string `json:"product_name"`
}

// CreateChatThreadRequest represents a create chat thread request.
type CreateChatThreadRequest struct {
	PurchaseID int64 `json:"purchase_id" binding:"required"`
}

// CreateChatThreadResponse represents the response after creating a chat thread.
type CreateChatThreadResponse struct {
	ID         int64  `json:"id"`
	PurchaseID int64  `json:"purchase_id"`
	Message    string `json:"message"`
}

// ChatMessageResponse represents a chat message in API responses.
type ChatMessageResponse struct {
	ID             int64      `json:"id"`
	ThreadID       int64      `json:"thread_id"`
	SenderID       int64      `json:"sender_id"`
	MessageType    string     `json:"message_type"`
	Content        string     `json:"content"`
	AttachmentURL  *string    `json:"attachment_url,omitempty"`
	AttachmentSize *int64     `json:"attachment_size,omitempty"`
	IsRead         bool       `json:"is_read"`
	ReadAt         *time.Time `json:"read_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	Sender         *ChatUser  `json:"sender,omitempty"`
}

// SendChatMessageRequest represents a send chat message request.
type SendChatMessageRequest struct {
	Content       string `json:"content" binding:"required"`
	MessageType   string `json:"message_type" binding:"required,oneof=text image file"`
	AttachmentURL string `json:"attachment_url,omitempty"`
}

// SendChatMessageResponse represents the response after sending a chat message.
type SendChatMessageResponse struct {
	ID        int64  `json:"id"`
	ThreadID  int64  `json:"thread_id"`
	CreatedAt string `json:"created_at"`
}

// AdminSendChatMessageRequest represents an admin send chat message request.
type AdminSendChatMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

// PollMessagesResponse represents the response for long-polling messages.
type PollMessagesResponse struct {
	Messages []ChatMessageResponse `json:"messages"`
	HasMore  bool                  `json:"has_more"`
}

// ResolveThreadResponse represents the response for resolving/reopening a thread.
type ResolveThreadResponse struct {
	Message string `json:"message"`
}
