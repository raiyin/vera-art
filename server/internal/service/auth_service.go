package service

import (
	"context"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/port"
	"github.com/raiyin/artserver/pkg/email"
	"github.com/raiyin/artserver/pkg/jwt"
	"github.com/raiyin/artserver/pkg/validator"
)

// AuthService implements port.AuthService.
type AuthService struct {
	userRepo    port.UserRepository
	jwtManager  *jwt.Manager
	emailSender email.EmailSender
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo port.UserRepository, jwtManager *jwt.Manager, emailSender email.EmailSender) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		jwtManager:  jwtManager,
		emailSender: emailSender,
	}
}

// Register registers a new user.
func (s *AuthService) Register(ctx context.Context, emailAddr, password, username string) (*domain.User, error) {
	if err := validator.ValidateEmail(emailAddr); err != nil {
		slog.Warn("AuthService.Register: invalid email",
			"email", emailAddr,
			"error", err,
		)
		return nil, err
	}
	if err := validator.ValidatePassword(password); err != nil {
		slog.Warn("AuthService.Register: invalid password",
			"username", username,
			"error", err,
		)
		return nil, err
	}
	if err := validator.ValidateUsername(username); err != nil {
		slog.Warn("AuthService.Register: invalid username",
			"username", username,
			"error", err,
		)
		return nil, err
	}

	// Check if user already exists
	existing, _ := s.userRepo.GetByEmail(ctx, emailAddr)
	if existing != nil {
		slog.Warn("AuthService.Register: email already exists",
			"email", emailAddr,
			"username", username,
		)
		return nil, domain.ErrDuplicate
	}
	existingByName, _ := s.userRepo.GetByUsername(ctx, username)
	if existingByName != nil {
		slog.Warn("AuthService.Register: username already exists",
			"username", username,
			"email", emailAddr,
		)
		return nil, domain.ErrDuplicate
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("AuthService.Register: failed to hash password",
			"username", username,
			"error", err,
		)
		return nil, domain.ErrInternal
	}

	verificationToken, err := email.GenerateVerificationToken()
	if err != nil {
		slog.Error("AuthService.Register: failed to generate verification token",
			"username", username,
			"error", err,
		)
		return nil, domain.ErrInternal
	}

	now := time.Now()
	user := &domain.User{
		Username:           username,
		Email:              emailAddr,
		PasswordHash:       string(hashedPassword),
		Name:               username,
		Role:               "user",
		EmailVerified:      false,
		VerificationToken:  verificationToken,
		VerificationSentAt: &now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		slog.Error("AuthService.Register: failed to create user",
			"username", username,
			"email", emailAddr,
			"error", err,
		)
		return nil, err
	}

	slog.Info("AuthService.Register: user registered",
		"user_id", user.ID,
		"username", username,
		"email", emailAddr,
	)

	// Send verification email (non-blocking)
	go func() {
		if err := s.emailSender.SendVerificationEmail(user.Email, verificationToken); err != nil {
			slog.Error("AuthService.Register: failed to send verification email",
				"user_id", user.ID,
				"email", user.Email,
				"error", err,
			)
		}
	}()

	return user, nil
}

// Login authenticates a user and returns tokens.
func (s *AuthService) Login(ctx context.Context, username, password string) (string, string, *domain.User, error) {
	if err := validator.ValidateUsername(username); err != nil {
		slog.Warn("AuthService.Login: invalid username",
			"username", username,
			"error", err,
		)
		return "", "", nil, err
	}

	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		slog.Warn("AuthService.Login: user not found",
			"username", username,
			"error", err,
		)
		return "", "", nil, domain.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		slog.Warn("AuthService.Login: invalid password",
			"username", username,
			"user_id", user.ID,
		)
		return "", "", nil, domain.ErrUnauthorized
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user)
	if err != nil {
		slog.Error("AuthService.Login: failed to generate access token",
			"user_id", user.ID,
			"username", username,
			"error", err,
		)
		return "", "", nil, domain.ErrInternal
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user)
	if err != nil {
		slog.Error("AuthService.Login: failed to generate refresh token",
			"user_id", user.ID,
			"username", username,
			"error", err,
		)
		return "", "", nil, domain.ErrInternal
	}

	slog.Info("AuthService.Login: user logged in",
		"user_id", user.ID,
		"username", username,
	)

	return accessToken, refreshToken, user, nil
}

// RefreshToken refreshes an access token using a refresh token.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		slog.Warn("AuthService.RefreshToken: invalid refresh token",
			"error", err,
		)
		return "", "", domain.ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		slog.Warn("AuthService.RefreshToken: user not found",
			"user_id", claims.UserID,
			"error", err,
		)
		return "", "", domain.ErrUnauthorized
	}

	newAccessToken, err := s.jwtManager.GenerateAccessToken(user)
	if err != nil {
		slog.Error("AuthService.RefreshToken: failed to generate access token",
			"user_id", user.ID,
			"error", err,
		)
		return "", "", domain.ErrInternal
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user)
	if err != nil {
		slog.Error("AuthService.RefreshToken: failed to generate refresh token",
			"user_id", user.ID,
			"error", err,
		)
		return "", "", domain.ErrInternal
	}

	slog.Info("AuthService.RefreshToken: tokens refreshed",
		"user_id", user.ID,
	)

	return newAccessToken, newRefreshToken, nil
}

// VerifyEmail verifies a user's email with a token.
func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	user, err := s.userRepo.GetByVerificationToken(ctx, token)
	if err != nil {
		slog.Warn("AuthService.VerifyEmail: invalid verification token",
			"error", err,
		)
		return domain.ErrInvalidToken
	}

	if user.EmailVerified {
		slog.Warn("AuthService.VerifyEmail: email already verified",
			"user_id", user.ID,
			"email", user.Email,
		)
		return nil
	}

	user.EmailVerified = true
	user.VerificationToken = ""
	if err := s.userRepo.Update(ctx, user); err != nil {
		slog.Error("AuthService.VerifyEmail: failed to update user",
			"user_id", user.ID,
			"email", user.Email,
			"error", err,
		)
		return err
	}

	slog.Info("AuthService.VerifyEmail: email verified",
		"user_id", user.ID,
		"email", user.Email,
	)
	return nil
}

// ResendVerification resends the verification email.
func (s *AuthService) ResendVerification(ctx context.Context, emailAddr string) error {
	user, err := s.userRepo.GetByEmail(ctx, emailAddr)
	if err != nil {
		slog.Warn("AuthService.ResendVerification: user not found",
			"email", emailAddr,
			"error", err,
		)
		return domain.ErrNotFound
	}

	if user.EmailVerified {
		slog.Warn("AuthService.ResendVerification: email already verified",
			"user_id", user.ID,
			"email", emailAddr,
		)
		return nil
	}

	verificationToken, err := email.GenerateVerificationToken()
	if err != nil {
		slog.Error("AuthService.ResendVerification: failed to generate verification token",
			"user_id", user.ID,
			"email", emailAddr,
			"error", err,
		)
		return domain.ErrInternal
	}

	now := time.Now()
	user.VerificationToken = verificationToken
	user.VerificationSentAt = &now

	if err := s.userRepo.Update(ctx, user); err != nil {
		slog.Error("AuthService.ResendVerification: failed to update user",
			"user_id", user.ID,
			"email", emailAddr,
			"error", err,
		)
		return err
	}

	slog.Info("AuthService.ResendVerification: verification email resent",
		"user_id", user.ID,
		"email", emailAddr,
	)

	go func() {
		if err := s.emailSender.SendVerificationEmail(user.Email, verificationToken); err != nil {
			slog.Error("AuthService.ResendVerification: failed to send verification email",
				"user_id", user.ID,
				"email", user.Email,
				"error", err,
			)
		}
	}()

	return nil
}
