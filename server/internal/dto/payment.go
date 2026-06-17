package dto

import "time"

// CreatePaymentRequest represents a create payment request.
type CreatePaymentRequest struct {
	ProductID int64  `json:"product_id" binding:"required"`
	PromoCode string `json:"promo_code,omitempty"`
}

// PaymentResponse represents a payment in API responses.
type PaymentResponse struct {
	ID              int64     `json:"id"`
	UserID          int64     `json:"user_id"`
	Status          string    `json:"status"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	Description     string    `json:"description,omitempty"`
	PaymentMethod   string    `json:"payment_method,omitempty"`
	ConfirmationURL string    `json:"confirmation_url,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CreatePaymentResponse represents the response after creating a payment.
type CreatePaymentResponse struct {
	PaymentID       int64   `json:"payment_id"`
	Amount          float64 `json:"amount"`
	Currency        string  `json:"currency"`
	ConfirmationURL string  `json:"confirmation_url"`
	Description     string  `json:"description"`
	Status          string  `json:"status"`
}

// PurchaseResponse represents a purchase in API responses.
type PurchaseResponse struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"user_id"`
	ProductID   int64      `json:"product_id"`
	PaymentID   int64      `json:"payment_id"`
	PricePaid   float64    `json:"price_paid"`
	Status      string     `json:"status"`
	AccessStart time.Time  `json:"access_start"`
	AccessEnd   *time.Time `json:"access_end,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// CreatePurchaseRequest represents a create purchase request.
type CreatePurchaseRequest struct {
	PaymentID int64 `json:"payment_id" binding:"required"`
}

// CreatePurchaseResponse represents the response after creating a purchase.
type CreatePurchaseResponse struct {
	PurchaseID  int64      `json:"purchase_id"`
	AccessStart time.Time  `json:"access_start"`
	AccessEnd   *time.Time `json:"access_end,omitempty"`
	Status      string     `json:"status"`
}

// YooKassaWebhook represents a YooKassa webhook notification.
type YooKassaWebhook struct {
	Event  string         `json:"event"`
	Object YooKassaObject `json:"object"`
}

// YooKassaObject represents a YooKassa payment object in webhook.
type YooKassaObject struct {
	ID       string         `json:"id"`
	Status   string         `json:"status"`
	Amount   YooKassaAmount `json:"amount"`
	Metadata map[string]any `json:"metadata"`
}

// YooKassaAmount represents a YooKassa amount.
type YooKassaAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}
