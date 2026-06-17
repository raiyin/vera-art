package domain

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrValidation        = errors.New("validation error")
	ErrConflict          = errors.New("conflict")
	ErrDuplicate         = errors.New("duplicate entry")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInternal          = errors.New("internal error")
	ErrEmailNotVerified  = errors.New("email not verified")
	ErrInvalidToken      = errors.New("invalid token")
	ErrTokenExpired      = errors.New("token expired")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrPaymentFailed     = errors.New("payment failed")
	ErrPromoCodeInvalid  = errors.New("promo code invalid or expired")
	ErrPromoCodeUsed     = errors.New("promo code already used")
	ErrFileTooLarge      = errors.New("file too large")
	ErrInvalidFileType   = errors.New("invalid file type")
	ErrRateLimited       = errors.New("rate limited")
)
