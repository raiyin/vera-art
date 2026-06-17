package validator

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/raiyin/artserver/internal/domain"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail validates an email address.
func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return domain.ErrInvalidInput
	}
	if !emailRegex.MatchString(email) {
		return domain.ErrInvalidInput
	}
	if len(email) > 254 {
		return domain.ErrInvalidInput
	}
	return nil
}

// ValidatePassword checks password complexity.
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return domain.ErrInvalidInput
	}
	if len(password) > 128 {
		return domain.ErrInvalidInput
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsNumber(ch):
			hasNumber = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return domain.ErrInvalidInput
	}

	return nil
}

// ValidateName validates a user's display name.
func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.ErrInvalidInput
	}
	if len(name) > 100 {
		return domain.ErrInvalidInput
	}
	return nil
}

// ValidateSlug validates a URL slug.
func ValidateSlug(slug string) error {
	if slug == "" {
		return domain.ErrInvalidInput
	}
	if len(slug) > 200 {
		return domain.ErrInvalidInput
	}
	matched, _ := regexp.MatchString(`^[a-z0-9\-]+$`, slug)
	if !matched {
		return domain.ErrInvalidInput
	}
	return nil
}

// ValidateFileExtension checks if a file has an allowed extension.
func ValidateFileExtension(filename string, allowedExts []string) error {
	ext := strings.ToLower(filename[strings.LastIndex(filename, ".")+1:])
	for _, allowed := range allowedExts {
		if ext == allowed {
			return nil
		}
	}
	return domain.ErrInvalidFileType
}

// ValidateFileSize checks if a file size is within limits.
func ValidateFileSize(size int64, maxSize int64) error {
	if size > maxSize {
		return domain.ErrFileTooLarge
	}
	return nil
}

// Common allowed extensions.
var (
	ImageExtensions = []string{"jpg", "jpeg", "png", "gif", "webp", "svg"}
	VideoExtensions = []string{"mp4", "webm", "ogg", "mov"}
	DocExtensions   = []string{"pdf", "doc", "docx"}
)
