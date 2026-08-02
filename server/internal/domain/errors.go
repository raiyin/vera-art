package domain

import (
	"errors"

	"github.com/raiyin/artserver/pkg/apperror"
)

var (
	// Shared sentinels are aliases of the ones in pkg/apperror. Keeping a single
	// canonical value lets apperror.FromError (which uses errors.Is) map a domain
	// error to the correct HTTP status — a domain copy would silently degrade to
	// a generic 500.
	ErrNotFound          = apperror.ErrNotFound
	ErrUnauthorized      = apperror.ErrUnauthorized
	ErrForbidden         = apperror.ErrForbidden
	ErrValidation        = apperror.ErrValidation
	ErrConflict          = apperror.ErrConflict
	ErrDuplicate         = apperror.ErrDuplicate
	ErrInvalidInput      = apperror.ErrInvalidInput
	ErrEmailNotVerified  = apperror.ErrEmailNotVerified
	ErrInvalidToken      = apperror.ErrInvalidToken
	ErrTokenExpired      = apperror.ErrTokenExpired
	ErrPaymentFailed     = apperror.ErrPaymentFailed
	ErrPromoCodeInvalid  = apperror.ErrPromoCodeInvalid
	ErrPromoCodeUsed     = apperror.ErrPromoCodeUsed
	ErrFileTooLarge      = apperror.ErrFileTooLarge
	ErrInvalidFileType   = apperror.ErrInvalidFileType
	ErrRateLimited       = apperror.ErrRateLimited

	// Domain-only errors without a dedicated HTTP mapping in apperror.FromError;
	// they fall through to a generic 500.
	ErrInternal          = errors.New("internal error")
	ErrInsufficientStock = errors.New("insufficient stock")
)
