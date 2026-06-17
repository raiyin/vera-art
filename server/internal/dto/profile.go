package dto

import "time"

// UpdateProfileRequest represents a profile update request.
type UpdateProfileRequest struct {
	Name *string `json:"name,omitempty"`
}

// ProfileResponse represents a user profile in API responses.
type ProfileResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UploadAvatarResponse represents the response after uploading an avatar.
type UploadAvatarResponse struct {
	Message   string `json:"message"`
	AvatarURL string `json:"avatar_url"`
}

// DeleteAvatarResponse represents the response after deleting an avatar.
type DeleteAvatarResponse struct {
	Message string `json:"message"`
}
