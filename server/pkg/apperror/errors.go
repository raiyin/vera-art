package apperror

import (
	"errors"
	"net/http"
)

// APIError represents a structured API error response.
type APIError struct {
	Status  int         `json:"-"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return e.Message
}

// New creates a new APIError.
func New(status int, code, message string) *APIError {
	return &APIError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

// NewWithDetails creates a new APIError with additional details.
func NewWithDetails(status int, code, message string, details interface{}) *APIError {
	return &APIError{
		Status:  status,
		Code:    code,
		Message: message,
		Details: details,
	}
}

// FromError maps a domain error to an APIError.
func FromError(err error) *APIError {
	if err == nil {
		return nil
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}

	switch {
	case errors.Is(err, ErrNotFound):
		return New(http.StatusNotFound, "NOT_FOUND", err.Error())
	case errors.Is(err, ErrUnauthorized):
		return New(http.StatusUnauthorized, "UNAUTHORIZED", err.Error())
	case errors.Is(err, ErrForbidden):
		return New(http.StatusForbidden, "FORBIDDEN", err.Error())
	case errors.Is(err, ErrValidation):
		return New(http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	case errors.Is(err, ErrConflict):
		return New(http.StatusConflict, "CONFLICT", err.Error())
	case errors.Is(err, ErrDuplicate):
		return New(http.StatusConflict, "DUPLICATE", err.Error())
	case errors.Is(err, ErrInvalidInput):
		return New(http.StatusBadRequest, "INVALID_INPUT", err.Error())
	case errors.Is(err, ErrEmailNotVerified):
		return New(http.StatusForbidden, "EMAIL_NOT_VERIFIED", err.Error())
	case errors.Is(err, ErrInvalidToken):
		return New(http.StatusUnauthorized, "INVALID_TOKEN", err.Error())
	case errors.Is(err, ErrTokenExpired):
		return New(http.StatusUnauthorized, "TOKEN_EXPIRED", err.Error())
	case errors.Is(err, ErrFileTooLarge):
		return New(http.StatusBadRequest, "FILE_TOO_LARGE", err.Error())
	case errors.Is(err, ErrInvalidFileType):
		return New(http.StatusBadRequest, "INVALID_FILE_TYPE", err.Error())
	case errors.Is(err, ErrRateLimited):
		return New(http.StatusTooManyRequests, "RATE_LIMITED", err.Error())
	case errors.Is(err, ErrPromoCodeInvalid):
		return New(http.StatusBadRequest, "PROMO_CODE_INVALID", err.Error())
	case errors.Is(err, ErrPromoCodeUsed):
		return New(http.StatusBadRequest, "PROMO_CODE_USED", err.Error())
	case errors.Is(err, ErrPaymentFailed):
		return New(http.StatusPaymentRequired, "PAYMENT_FAILED", err.Error())
	default:
		return New(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
	}
}

// Predefined domain errors for use in pkg/apperror.
// These are aliases so the pkg/apperror package can be used standalone.
var (
	ErrNotFound         = errors.New("not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrValidation       = errors.New("validation error")
	ErrConflict         = errors.New("conflict")
	ErrDuplicate        = errors.New("duplicate entry")
	ErrInvalidInput     = errors.New("invalid input")
	ErrEmailNotVerified = errors.New("email not verified")
	ErrInvalidToken     = errors.New("invalid token")
	ErrTokenExpired     = errors.New("token expired")
	ErrFileTooLarge     = errors.New("file too large")
	ErrInvalidFileType  = errors.New("invalid file type")
	ErrRateLimited      = errors.New("rate limited")
	ErrPromoCodeInvalid = errors.New("promo code invalid or expired")
	ErrPromoCodeUsed    = errors.New("promo code already used")
	ErrPaymentFailed    = errors.New("payment failed")
)
