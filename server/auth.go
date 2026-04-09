package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/raiyin/artserver/dtos"
	"github.com/raiyin/artserver/models"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var registerUserRequest dtos.RegisterUserRequest
	if err := c.ShouldBindJSON(&registerUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	var count int64
	query := "select count(*) from users where username = ?"
	err := db.QueryRow(query, registerUserRequest.Username).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not check if user exists"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerUserRequest.Password), bcrypt.DefaultCost)
	var stringHashedPassword = string(hashedPassword)

	// fmt.Println("%X ", hashedPassword)
	// for _, num := range hashedPassword {
	// 	fmt.Printf("%X ", num)
	// }
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	query = "insert into users (username, pass_hash) " +
		"values (?, ?)"
	_, err = db.Exec(query, registerUserRequest.Username, stringHashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user created successfully"})
}

func Login(c *gin.Context) {
	var loginUserRequest dtos.LoginUserRequest
	if err := c.ShouldBindJSON(&loginUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id int
	var username string
	var passhash string

	query := "select * from users where username = ?"
	rows := db.QueryRow(query, loginUserRequest.Username)
	err := rows.Scan(&id, &username, &passhash)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "undefinded servererror"})
			return
		}
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(passhash), []byte(loginUserRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Create JWT token
	expirationTime := time.Now().Add(365 * 24 * time.Hour)
	claims := &models.Claims{
		Username: loginUserRequest.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	fmt.Printf("jwt key is %v", jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   tokenString,
		"expires": expirationTime,
	})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		tokenString = tokenString[len("Bearer "):]
		claims := &models.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
