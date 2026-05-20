package main

import (
	"database/sql"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/raiyin/artserver/dtos"
	"github.com/raiyin/artserver/models"
	"golang.org/x/crypto/bcrypt"
)

// validatePasswordComplexity checks if password meets complexity requirements
func validatePasswordComplexity(password string) (bool, string) {
	if len(password) < 8 {
		return false, "password must be at least 8 characters long"
	}

	// Check for at least one uppercase letter
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false, "password must contain at least one uppercase letter"
	}

	// Check for at least one lowercase letter
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return false, "password must contain at least one lowercase letter"
	}

	// Check for at least one digit
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return false, "password must contain at least one digit"
	}

	// Check for at least one special character
	if !regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~]`).MatchString(password) {
		return false, "password must contain at least one special character"
	}

	// Check for common weak passwords (optional)
	weakPasswords := []string{"password", "12345678", "qwerty", "admin", "letmein"}
	lowerPassword := strings.ToLower(password)
	for _, weak := range weakPasswords {
		if strings.Contains(lowerPassword, weak) {
			return false, "password is too common or weak"
		}
	}

	return true, ""
}

// generateAccessToken creates a short-lived access token (15 minutes)
func generateAccessToken(username string, userID int, role string) (string, time.Time, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &models.Claims{
		Username:  username,
		UserID:    userID,
		Role:      role,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

// logSecurityEvent logs security-related events with structured format
func logSecurityEvent(eventType, username, ip, details string) {
	log.Printf("[SECURITY] type=%s user=%s ip=%s details=%s",
		eventType, username, ip, details)
}

// generateRefreshToken creates a longer-lived refresh token (7 days)
func generateRefreshToken(username string, userID int, role string) (string, time.Time, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) // 7 days
	claims := &models.Claims{
		Username:  username,
		UserID:    userID,
		Role:      role,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

func Register(c *gin.Context) {
	var registerUserRequest dtos.RegisterUserRequest
	if err := c.ShouldBindJSON(&registerUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate password complexity
	if valid, msg := validatePasswordComplexity(registerUserRequest.Password); !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
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
		logSecurityEvent("register_failed", registerUserRequest.Username, c.ClientIP(), "database error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	// Log successful registration
	logSecurityEvent("register_success", registerUserRequest.Username, c.ClientIP(), "user created")
	c.JSON(http.StatusOK, gin.H{"message": "user created successfully"})
}

func Login(c *gin.Context) {
	var loginUserRequest dtos.LoginUserRequest
	if err := c.ShouldBindJSON(&loginUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	query := "SELECT id, username, pass_hash, role, email, full_name, created_at, updated_at FROM users WHERE username = ?"
	row := db.QueryRow(query, loginUserRequest.Username)
	err := row.Scan(&user.Id, &user.Username, &user.PassHash, &user.Role, &user.Email, &user.FullName, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			// Log failed login attempt (user not found)
			logSecurityEvent("login_failed", loginUserRequest.Username, c.ClientIP(), "user not found")
			// Return same error as invalid password to avoid user enumeration
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		} else {
			logSecurityEvent("login_error", loginUserRequest.Username, c.ClientIP(), "database error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
			return
		}
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(loginUserRequest.Password)); err != nil {
		// Log failed login attempt (wrong password)
		logSecurityEvent("login_failed", loginUserRequest.Username, c.ClientIP(), "invalid password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Generate access token (15 minutes)
	accessToken, accessExp, err := generateAccessToken(user.Username, user.Id, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate access token"})
		return
	}

	// Generate refresh token (7 days)
	refreshToken, refreshExp, err := generateRefreshToken(user.Username, user.Id, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate refresh token"})
		return
	}

	// Log successful login
	logSecurityEvent("login_success", loginUserRequest.Username, c.ClientIP(), "tokens issued")

	c.JSON(http.StatusOK, gin.H{
		"access_token":    accessToken,
		"access_expires":  accessExp,
		"refresh_token":   refreshToken,
		"refresh_expires": refreshExp,
		"token_type":      "Bearer",
	})
}

// Refresh generates a new access token using a valid refresh token
func Refresh(c *gin.Context) {
	// Get refresh token from Authorization header or request body
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		// Try to get from request body
		var request struct {
			RefreshToken string `json:"refresh_token"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token required"})
			return
		}
		tokenString = request.RefreshToken
	} else {
		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}
	}

	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh token required"})
		return
	}

	// Parse and validate the refresh token
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		logSecurityEvent("refresh_failed", "", c.ClientIP(), "invalid token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	// Verify it's a refresh token
	if claims.TokenType != "refresh" {
		logSecurityEvent("refresh_failed", claims.Username, c.ClientIP(), "not a refresh token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not a refresh token"})
		return
	}

	// Generate new access token using claims from refresh token
	accessToken, accessExp, err := generateAccessToken(claims.Username, claims.UserID, claims.Role)
	if err != nil {
		logSecurityEvent("refresh_error", claims.Username, c.ClientIP(), "could not generate access token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate access token"})
		return
	}

	// Log successful refresh
	logSecurityEvent("refresh_success", claims.Username, c.ClientIP(), "new access token issued")

	c.JSON(http.StatusOK, gin.H{
		"access_token":   accessToken,
		"access_expires": accessExp,
		"token_type":     "Bearer",
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

		// Check for Bearer prefix
		const bearerPrefix = "Bearer "
		if len(tokenString) <= len(bearerPrefix) || !strings.HasPrefix(tokenString, bearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}
		tokenString = tokenString[len(bearerPrefix):]

		claims := &models.Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Verify it's an access token
		if claims.TokenType != "access" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "not an access token"})
			c.Abort()
			return
		}

		// Set claims and user_id in context for handlers
		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
