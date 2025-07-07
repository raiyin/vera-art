package models

import (
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/mattn/go-sqlite3"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
