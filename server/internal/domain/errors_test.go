package domain

import (
	"errors"
	"testing"
)

func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		msg  string
	}{
		{"ErrNotFound", ErrNotFound, "not found"},
		{"ErrUnauthorized", ErrUnauthorized, "unauthorized"},
		{"ErrForbidden", ErrForbidden, "forbidden"},
		{"ErrValidation", ErrValidation, "validation error"},
		{"ErrConflict", ErrConflict, "conflict"},
		{"ErrDuplicate", ErrDuplicate, "duplicate entry"},
		{"ErrInvalidInput", ErrInvalidInput, "invalid input"},
		{"ErrInternal", ErrInternal, "internal error"},
		{"ErrEmailNotVerified", ErrEmailNotVerified, "email not verified"},
		{"ErrInvalidToken", ErrInvalidToken, "invalid token"},
		{"ErrTokenExpired", ErrTokenExpired, "token expired"},
		{"ErrInsufficientStock", ErrInsufficientStock, "insufficient stock"},
		{"ErrPaymentFailed", ErrPaymentFailed, "payment failed"},
		{"ErrPromoCodeInvalid", ErrPromoCodeInvalid, "promo code invalid or expired"},
		{"ErrPromoCodeUsed", ErrPromoCodeUsed, "promo code already used"},
		{"ErrFileTooLarge", ErrFileTooLarge, "file too large"},
		{"ErrInvalidFileType", ErrInvalidFileType, "invalid file type"},
		{"ErrRateLimited", ErrRateLimited, "rate limited"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Error("error is nil")
			}
			if tt.err.Error() != tt.msg {
				t.Errorf("got message %q, want %q", tt.err.Error(), tt.msg)
			}
		})
	}
}

func TestErrorMatching(t *testing.T) {
	t.Run("ErrNotFound matches wrapped error", func(t *testing.T) {
		wrapped := errors.Join(ErrNotFound, errors.New("some context"))
		if !errors.Is(wrapped, ErrNotFound) {
			t.Error("expected wrapped error to match ErrNotFound")
		}
	})

	t.Run("distinct errors are not equal", func(t *testing.T) {
		if errors.Is(ErrNotFound, ErrUnauthorized) {
			t.Error("ErrNotFound should not match ErrUnauthorized")
		}
	})

	t.Run("ErrInvalidInput used in validation", func(t *testing.T) {
		err := ErrInvalidInput
		if !errors.Is(err, ErrInvalidInput) {
			t.Error("ErrInvalidInput should match itself")
		}
	})
}
