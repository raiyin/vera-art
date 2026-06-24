package domain

import "time"

type User struct {
	ID                 int64      `json:"id"`
	Username           string     `json:"username"`
	Email              string     `json:"email"`
	PasswordHash       string     `json:"-"`
	Name               string     `json:"name"`
	Role               string     `json:"role"`
	AvatarPath         string     `json:"avatar_path,omitempty"`
	EmailVerified      bool       `json:"email_verified"`
	VerificationToken  string     `json:"-"`
	VerificationSentAt *time.Time `json:"verification_sent_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type UserFilter struct {
	Role      string
	Query     string
	Page      int
	Limit     int
	SortBy    string
	SortOrder string
}
