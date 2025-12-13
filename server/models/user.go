package models

import (
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	PassHash string `json:"passhash"`
}
