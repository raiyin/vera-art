package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/dto"
	"github.com/raiyin/artserver/internal/port"
	"github.com/raiyin/artserver/pkg/apperror"
)

// AuthHandler handles authentication HTTP requests.
type AuthHandler struct {
	authService port.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authService port.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles user registration.
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Register: invalid request body",
			"client_ip", c.ClientIP(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Email, req.Password, req.Username)
	if err != nil {
		slog.Warn("Register: registration failed",
			"email", req.Email,
			"username", req.Username,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("User registered successfully",
		"user_id", user.ID,
		"username", user.Username,
		"email", user.Email,
	)
	c.JSON(http.StatusCreated, dto.UserDTO{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		Name:          user.Name,
		Role:          user.Role,
		AvatarPath:    user.AvatarPath,
		EmailVerified: user.EmailVerified,
	})
}

// Login handles user login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Login: invalid request body",
			"client_ip", c.ClientIP(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	accessToken, refreshToken, user, err := h.authService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		slog.Warn("Login: authentication failed",
			"username", req.Username,
			"client_ip", c.ClientIP(),
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("User logged in successfully",
		"user_id", user.ID,
		"username", user.Username,
	)
	c.JSON(http.StatusOK, dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserDTO{
			ID:            user.ID,
			Username:      user.Username,
			Email:         user.Email,
			Name:          user.Name,
			Role:          user.Role,
			AvatarPath:    user.AvatarPath,
			EmailVerified: user.EmailVerified,
		},
	})
}

// RefreshToken handles token refresh.
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("RefreshToken: invalid request body",
			"client_ip", c.ClientIP(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	newAccessToken, newRefreshToken, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		slog.Warn("RefreshToken: token refresh failed",
			"client_ip", c.ClientIP(),
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Token refreshed successfully")
	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

// VerifyEmail handles email verification.
func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	var req dto.VerifyEmailRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		slog.Warn("VerifyEmail: invalid request",
			"client_ip", c.ClientIP(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid token",
		})
		return
	}

	if err := h.authService.VerifyEmail(c.Request.Context(), req.Token); err != nil {
		slog.Warn("VerifyEmail: verification failed",
			"client_ip", c.ClientIP(),
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Email verified successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}

// ResendVerification resends the verification email.
func (h *AuthHandler) ResendVerification(c *gin.Context) {
	var req dto.ResendVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("ResendVerification: invalid request body",
			"client_ip", c.ClientIP(),
			"error", err,
		)
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	if err := h.authService.ResendVerification(c.Request.Context(), req.Email); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			// Don't reveal if email exists
			c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a verification email has been sent"})
			return
		}
		slog.Warn("ResendVerification: failed to resend",
			"email", req.Email,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Verification email resent",
		"email", req.Email,
	)
	c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a verification email has been sent"})
}
