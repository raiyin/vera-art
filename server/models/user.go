package models

import (
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type User struct {
	Id                        int       `json:"id"`
	Username                  string    `json:"username"`
	PassHash                  string    `json:"passhash"`
	Role                      string    `json:"role"`
	Email                     string    `json:"email,omitempty"`
	FullName                  string    `json:"full_name,omitempty"`
	Avatar                    string    `json:"avatar,omitempty"`
	EmailVerified             bool      `json:"email_verified"`
	VerificationToken         string    `json:"-"`
	VerificationTokenExpiresAt time.Time `json:"-"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}
