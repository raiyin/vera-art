package dto

// RegisterUserRequest represents a registration request.
type RegisterUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

// LoginUserRequest represents a login request.
type LoginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest represents a token refresh request.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ResendVerificationRequest represents a resend verification request.
type ResendVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// AuthResponse represents an authentication response.
type AuthResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	User         UserDTO `json:"user"`
}

// UserDTO represents a user in API responses.
type UserDTO struct {
	ID            int64  `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Role          string `json:"role"`
	AvatarPath    string `json:"avatar_path,omitempty"`
	EmailVerified bool   `json:"email_verified"`
}

// VerifyEmailRequest represents an email verification request.
type VerifyEmailRequest struct {
	Token string `json:"token" form:"token" binding:"required"`
}
