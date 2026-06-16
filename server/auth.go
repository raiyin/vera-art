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
		log.Printf("[ERROR] could not check if user exists: %v", err)
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

	if err != nil {
		log.Printf("[ERROR] could not hash password for user %s: %v", registerUserRequest.Username, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	// Generate verification token
	verificationToken, tokenExpiry, err := generateVerificationToken()
	if err != nil {
		log.Printf("[ERROR] could not generate verification token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	query = "insert into users (username, pass_hash, email, email_verified, verification_token, verification_token_expires_at, created_at, updated_at) " +
		"values (?, ?, ?, 0, ?, ?, datetime('now'), datetime('now'))"
	_, err = db.Exec(query, registerUserRequest.Username, stringHashedPassword, registerUserRequest.Email, verificationToken, tokenExpiry)
	if err != nil {
		logSecurityEvent("register_failed", registerUserRequest.Username, c.ClientIP(), "database error: "+err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	// Send verification email (non-blocking — log error but don't fail the request)
	if err := SendVerificationEmail(registerUserRequest.Email, verificationToken); err != nil {
		log.Printf("[ERROR] failed to send verification email to %s: %v", registerUserRequest.Email, err)
	}

	// Log successful registration
	logSecurityEvent("register_success", registerUserRequest.Username, c.ClientIP(), "user created (unverified)")
	c.JSON(http.StatusOK, gin.H{"message": "Registration successful! A verification link has been sent to your email."})
}

func Login(c *gin.Context) {
	var loginUserRequest dtos.LoginUserRequest
	if err := c.ShouldBindJSON(&loginUserRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	query := "SELECT id, username, pass_hash, role, email, full_name, email_verified, created_at, updated_at FROM users WHERE username = ?"
	row := db.QueryRow(query, loginUserRequest.Username)
	err := row.Scan(&user.Id, &user.Username, &user.PassHash, &user.Role, &user.Email, &user.FullName, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)

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

	// Check if email is verified
	if !user.EmailVerified {
		logSecurityEvent("login_failed", loginUserRequest.Username, c.ClientIP(), "email not verified")
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Please verify your email before logging in",
			"code":  "email_not_verified",
			"email": user.Email,
		})
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

// VerifyEmail handles email verification via token
func VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token is required"})
		return
	}

	var userID int
	var expiresAt time.Time

	query := "SELECT id, verification_token_expires_at FROM users WHERE verification_token = ? AND email_verified = 0"
	err := db.QueryRow(query, token).Scan(&userID, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired verification token"})
			return
		}
		log.Printf("[ERROR] could not look up verification token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	// Check if token has expired
	if time.Now().After(expiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Verification token has expired. Please request a new one."})
		return
	}

	// Mark user as verified and clear the token
	_, err = db.Exec("UPDATE users SET email_verified = 1, verification_token = NULL, verification_token_expires_at = NULL, updated_at = datetime('now') WHERE id = ?", userID)
	if err != nil {
		log.Printf("[ERROR] could not verify email for user %d: %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not verify email"})
		return
	}

	logSecurityEvent("email_verified", "", c.ClientIP(), "user_id="+string(rune(userID))+" email verified")
	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully! You can now log in."})
}

// ResendVerification resends the verification email
func ResendVerification(c *gin.Context) {
	var req dtos.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Always return success to prevent email enumeration
	// Look up user by email
	var user models.User
	query := "SELECT id, username, email, email_verified, verification_token, verification_token_expires_at FROM users WHERE email = ?"
	err := db.QueryRow(query, req.Email).Scan(&user.Id, &user.Username, &user.Email, &user.EmailVerified, &user.VerificationToken, &user.VerificationTokenExpiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return success even if email not found (anti-enumeration)
			c.JSON(http.StatusOK, gin.H{"message": "If this email is registered, a verification link has been sent."})
			return
		}
		log.Printf("[ERROR] could not look up user by email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}

	// Check if already verified
	if user.EmailVerified {
		c.JSON(http.StatusOK, gin.H{"message": "This email is already verified. You can log in."})
		return
	}

	// Generate new token
	newToken, newExpiry, err := generateVerificationToken()
	if err != nil {
		log.Printf("[ERROR] could not generate verification token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not resend verification email"})
		return
	}

	// Update token in database
	_, err = db.Exec("UPDATE users SET verification_token = ?, verification_token_expires_at = ?, updated_at = datetime('now') WHERE id = ?", newToken, newExpiry, user.Id)
	if err != nil {
		log.Printf("[ERROR] could not update verification token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not resend verification email"})
		return
	}

	// Send verification email
	if err := SendVerificationEmail(user.Email, newToken); err != nil {
		log.Printf("[ERROR] failed to send verification email to %s: %v", user.Email, err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "If this email is registered, a verification link has been sent."})
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

// AuthMiddlewareOptional — опциональная авторизация.
// Если токен предоставлен и валиден, устанавливает user_id в контекст.
// Если токена нет или он невалиден, запрос всё равно проходит дальше.
func AuthMiddlewareOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.Next()
			return
		}

		const bearerPrefix = "Bearer "
		if len(tokenString) <= len(bearerPrefix) || !strings.HasPrefix(tokenString, bearerPrefix) {
			c.Next()
			return
		}
		tokenString = tokenString[len(bearerPrefix):]

		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.Next()
			return
		}

		if claims.TokenType != "access" {
			c.Next()
			return
		}

		c.Set("claims", claims)
		c.Set("user_id", claims.UserID)
		c.Next()
	}
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
