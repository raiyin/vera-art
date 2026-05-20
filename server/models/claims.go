package models

import (
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
)

type Claims struct {
	Username string `json:"username"`
	UserID   int    `json:"user_id"`
	Role     string `json:"role"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}
