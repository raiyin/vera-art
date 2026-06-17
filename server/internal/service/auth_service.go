package service

import (
	"context"
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
func (s *AuthService) Register(ctx context.Context, emailAddr, password, name string) (*domain.User, error) {
	if err := validator.ValidateEmail(emailAddr); err != nil {
		return nil, err
	}
	if err := validator.ValidatePassword(password); err != nil {
		return nil, err
	}
	if err := validator.ValidateName(name); err != nil {
		return nil, err
	}

	// Check if user already exists
	existing, _ := s.userRepo.GetByEmail(ctx, emailAddr)
	if existing != nil {
		return nil, domain.ErrDuplicate
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, domain.ErrInternal
	}

	verificationToken, err := email.GenerateVerificationToken()
	if err != nil {
		return nil, domain.ErrInternal
	}

	now := time.Now()
	user := &domain.User{
		Email:              emailAddr,
		PasswordHash:       string(hashedPassword),
		Name:               name,
		Role:               "user",
		EmailVerified:      false,
		VerificationToken:  verificationToken,
		VerificationSentAt: &now,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Send verification email (non-blocking)
	go func() {
		_ = s.emailSender.SendVerificationEmail(user.Email, verificationToken)
	}()

	return user, nil
}

// Login authenticates a user and returns tokens.
func (s *AuthService) Login(ctx context.Context, emailAddr, password string) (string, string, *domain.User, error) {
	if err := validator.ValidateEmail(emailAddr); err != nil {
		return "", "", nil, err
	}

	user, err := s.userRepo.GetByEmail(ctx, emailAddr)
	if err != nil {
		return "", "", nil, domain.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", nil, domain.ErrUnauthorized
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(user)
	if err != nil {
		return "", "", nil, domain.ErrInternal
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user)
	if err != nil {
		return "", "", nil, domain.ErrInternal
	}

	return accessToken, refreshToken, user, nil
}

// RefreshToken refreshes an access token using a refresh token.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return "", "", domain.ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return "", "", domain.ErrUnauthorized
	}

	newAccessToken, err := s.jwtManager.GenerateAccessToken(user)
	if err != nil {
		return "", "", domain.ErrInternal
	}

	newRefreshToken, err := s.jwtManager.GenerateRefreshToken(user)
	if err != nil {
		return "", "", domain.ErrInternal
	}

	return newAccessToken, newRefreshToken, nil
}

// VerifyEmail verifies a user's email with a token.
func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
	user, err := s.userRepo.GetByVerificationToken(ctx, token)
	if err != nil {
		return domain.ErrInvalidToken
	}

	if user.EmailVerified {
		return nil
	}

	user.EmailVerified = true
	user.VerificationToken = ""
	return s.userRepo.Update(ctx, user)
}

// ResendVerification resends the verification email.
func (s *AuthService) ResendVerification(ctx context.Context, emailAddr string) error {
	user, err := s.userRepo.GetByEmail(ctx, emailAddr)
	if err != nil {
		return domain.ErrNotFound
	}

	if user.EmailVerified {
		return nil
	}

	verificationToken, err := email.GenerateVerificationToken()
	if err != nil {
		return domain.ErrInternal
	}

	now := time.Now()
	user.VerificationToken = verificationToken
	user.VerificationSentAt = &now

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	go func() {
		_ = s.emailSender.SendVerificationEmail(user.Email, verificationToken)
	}()

	return nil
}
