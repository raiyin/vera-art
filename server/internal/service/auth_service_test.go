package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

// mockUserRepo implements port.UserRepository for auth tests.
type mockUserRepo struct {
	users  []domain.User
	err    error
	getBy  func(email, username string) (*domain.User, error)
}

func (m *mockUserRepo) Create(_ context.Context, user *domain.User) error {
	if m.err != nil {
		return m.err
	}
	user.ID = int64(len(m.users) + 1)
	m.users = append(m.users, *user)
	return nil
}

func (m *mockUserRepo) GetByID(_ context.Context, id int64) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, u := range m.users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserRepo) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.getBy != nil {
		return m.getBy(email, "")
	}
	for _, u := range m.users {
		if u.Email == email {
			return &u, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserRepo) GetByUsername(_ context.Context, username string) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.getBy != nil {
		return m.getBy("", username)
	}
	for _, u := range m.users {
		if u.Username == username {
			return &u, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserRepo) GetByVerificationToken(_ context.Context, token string) (*domain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, u := range m.users {
		if u.VerificationToken == token {
			return &u, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockUserRepo) List(_ context.Context, _ domain.UserFilter) ([]domain.User, int, error) {
	return m.users, len(m.users), nil
}

func (m *mockUserRepo) Update(_ context.Context, user *domain.User) error {
	if m.err != nil {
		return m.err
	}
	for i, u := range m.users {
		if u.ID == user.ID {
			m.users[i] = *user
			return nil
		}
	}
	return domain.ErrNotFound
}

func (m *mockUserRepo) Delete(_ context.Context, _ int64) error {
	return m.err
}

// mockEmailSender implements email.EmailSender.
type mockEmailSender struct {
	err error
}

func (m *mockEmailSender) SendVerificationEmail(to, token string) error {
	return m.err
}

func newAuthService(userRepo *mockUserRepo) *AuthService {
	m := jwt.NewManager([]byte("test-secret-for-auth"), 15*time.Minute, 7*24*time.Hour)
	return NewAuthService(userRepo, m, &mockEmailSender{})
}

func validUser() (string, string, string) {
	return "test@example.com", "Str0ng!Pass1", "testuser"
}

func hashPassword(password string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(h)
}

func TestRegister(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := newAuthService(&mockUserRepo{})
		email, password, username := validUser()
		user, err := svc.Register(context.Background(), email, password, username)
		if err != nil {
			t.Fatalf("Register failed: %v", err)
		}
		if user.ID == 0 {
			t.Error("expected non-zero ID")
		}
		if user.Username != username {
			t.Errorf("Username = %q, want %q", user.Username, username)
		}
		if user.Email != email {
			t.Errorf("Email = %q, want %q", user.Email, email)
		}
		if user.Role != "user" {
			t.Errorf("Role = %q, want %q", user.Role, "user")
		}
		if user.PasswordHash == "" {
			t.Error("expected non-empty PasswordHash")
		}
		if user.VerificationToken == "" {
			t.Error("expected non-empty VerificationToken")
		}
		if user.EmailVerified {
			t.Error("EmailVerified should be false")
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		svc := newAuthService(&mockUserRepo{})
		_, err := svc.Register(context.Background(), "bad-email", "Str0ng!Pass1", "testuser")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("invalid password too short", func(t *testing.T) {
		svc := newAuthService(&mockUserRepo{})
		_, err := svc.Register(context.Background(), "test@example.com", "Sh0rt!", "testuser")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("invalid password no special char", func(t *testing.T) {
		svc := newAuthService(&mockUserRepo{})
		_, err := svc.Register(context.Background(), "test@example.com", "Password1", "testuser")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("invalid username too short", func(t *testing.T) {
		svc := newAuthService(&mockUserRepo{})
		_, err := svc.Register(context.Background(), "test@example.com", "Str0ng!Pass1", "ab")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("invalid username with special chars", func(t *testing.T) {
		svc := newAuthService(&mockUserRepo{})
		_, err := svc.Register(context.Background(), "test@example.com", "Str0ng!Pass1", "user name!")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("duplicate email", func(t *testing.T) {
		repo := &mockUserRepo{}
		svc := newAuthService(repo)
		email, password, username := validUser()
		if _, err := svc.Register(context.Background(), email, password, username); err != nil {
			t.Fatalf("first register failed: %v", err)
		}
		_, err := svc.Register(context.Background(), email, "Diff3r!ntP@ss", "otheruser")
		if !errors.Is(err, domain.ErrDuplicate) {
			t.Errorf("expected ErrDuplicate, got %v", err)
		}
	})

	t.Run("duplicate username", func(t *testing.T) {
		repo := &mockUserRepo{}
		svc := newAuthService(repo)
		email, password, username := validUser()
		if _, err := svc.Register(context.Background(), email, password, username); err != nil {
			t.Fatalf("first register failed: %v", err)
		}
		_, err := svc.Register(context.Background(), "other@example.com", "Diff3r!ntP@ss", username)
		if !errors.Is(err, domain.ErrDuplicate) {
			t.Errorf("expected ErrDuplicate, got %v", err)
		}
	})

	t.Run("empty email", func(t *testing.T) {
		svc := newAuthService(&mockUserRepo{})
		_, err := svc.Register(context.Background(), "", "Str0ng!Pass1", "testuser")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})
}

func TestLogin(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockUserRepo{}
		svc := newAuthService(repo)

		email, password, username := validUser()
		if _, err := svc.Register(context.Background(), email, password, username); err != nil {
			t.Fatalf("register failed: %v", err)
		}

		accessToken, refreshToken, user, err := svc.Login(context.Background(), username, password)
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}
		if accessToken == "" {
			t.Error("expected non-empty access token")
		}
		if refreshToken == "" {
			t.Error("expected non-empty refresh token")
		}
		if user == nil {
			t.Fatal("expected non-nil user")
		}
		if user.Username != username {
			t.Errorf("Username = %q, want %q", user.Username, username)
		}
	})

	t.Run("invalid username", func(t *testing.T) {
		svc := newAuthService(&mockUserRepo{})
		_, _, _, err := svc.Login(context.Background(), "ab", "Str0ng!Pass1")
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		svc := newAuthService(&mockUserRepo{})
		_, _, _, err := svc.Login(context.Background(), "nonexistent", "Str0ng!Pass1")
		if !errors.Is(err, domain.ErrUnauthorized) {
			t.Errorf("expected ErrUnauthorized, got %v", err)
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		repo := &mockUserRepo{}
		svc := newAuthService(repo)
		email, password, username := validUser()
		if _, err := svc.Register(context.Background(), email, password, username); err != nil {
			t.Fatalf("register failed: %v", err)
		}
		_, _, _, err := svc.Login(context.Background(), "testuser", "WrongP@ss1")
		if !errors.Is(err, domain.ErrUnauthorized) {
			t.Errorf("expected ErrUnauthorized, got %v", err)
		}
	})
}

func TestRefreshToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockUserRepo{}
		svc := newAuthService(repo)

		email, password, username := validUser()
		if _, err := svc.Register(context.Background(), email, password, username); err != nil {
			t.Fatalf("register failed: %v", err)
		}

		_, refreshToken, _, err := svc.Login(context.Background(), "testuser", "Str0ng!Pass1")
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}

		newAccess, newRefresh, err := svc.RefreshToken(context.Background(), refreshToken)
		if err != nil {
			t.Fatalf("RefreshToken failed: %v", err)
		}
		if newAccess == "" {
			t.Error("expected non-empty new access token")
		}
		if newRefresh == "" {
			t.Error("expected non-empty new refresh token")
		}
		if newAccess == refreshToken {
			t.Error("new access token should differ from old refresh token")
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		svc := newAuthService(&mockUserRepo{})
		_, _, err := svc.RefreshToken(context.Background(), "invalid-token")
		if !errors.Is(err, domain.ErrInvalidToken) {
			t.Errorf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("empty token", func(t *testing.T) {
		svc := newAuthService(&mockUserRepo{})
		_, _, err := svc.RefreshToken(context.Background(), "")
		if !errors.Is(err, domain.ErrInvalidToken) {
			t.Errorf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("expired refresh token", func(t *testing.T) {
		m := jwt.NewManager([]byte("test-secret-for-auth"), 15*time.Minute, -1*time.Hour)
		repo := &mockUserRepo{}
		svc := NewAuthService(repo, m, &mockEmailSender{})

		user := &domain.User{
			ID:       1,
			Username: "testuser",
			Email:    "test@example.com",
			PasswordHash: hashPassword("Str0ng!Pass1"),
			Role:    "user",
		}
		repo.users = append(repo.users, *user)

		_, refreshToken, _, err := svc.Login(context.Background(), "testuser", "Str0ng!Pass1")
		if err != nil {
			t.Fatalf("Login failed: %v", err)
		}

		_, _, err = svc.RefreshToken(context.Background(), refreshToken)
		if !errors.Is(err, domain.ErrInvalidToken) {
			t.Errorf("expected ErrInvalidToken for expired refresh, got %v", err)
		}
	})
}

func TestVerifyEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockUserRepo{}
		svc := newAuthService(repo)
		email, password, username := validUser()

		if _, err := svc.Register(context.Background(), email, password, username); err != nil {
			t.Fatalf("register failed: %v", err)
		}

		user := &repo.users[0]
		if err := svc.VerifyEmail(context.Background(), user.VerificationToken); err != nil {
			t.Fatalf("VerifyEmail failed: %v", err)
		}

		updated := &repo.users[0]
		if !updated.EmailVerified {
			t.Error("EmailVerified should be true after verification")
		}
		if updated.VerificationToken != "" {
			t.Error("VerificationToken should be cleared after verification")
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		svc := newAuthService(&mockUserRepo{})
		err := svc.VerifyEmail(context.Background(), "invalid-token")
		if !errors.Is(err, domain.ErrInvalidToken) {
			t.Errorf("expected ErrInvalidToken, got %v", err)
		}
	})

	t.Run("already verified", func(t *testing.T) {
		repo := &mockUserRepo{}
		svc := newAuthService(repo)
		email, password, username := validUser()

		if _, err := svc.Register(context.Background(), email, password, username); err != nil {
			t.Fatalf("register failed: %v", err)
		}

		user := &repo.users[0]
		if err := svc.VerifyEmail(context.Background(), user.VerificationToken); err != nil {
			t.Fatalf("first verify failed: %v", err)
		}

		err := svc.VerifyEmail(context.Background(), user.VerificationToken)
		if err != nil {
			t.Errorf("expected nil for already verified, got %v", err)
		}
	})
}

func TestResendVerification(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockUserRepo{}
		svc := newAuthService(repo)
		email, password, username := validUser()

		if _, err := svc.Register(context.Background(), email, password, username); err != nil {
			t.Fatalf("register failed: %v", err)
		}

		oldToken := repo.users[0].VerificationToken
		if err := svc.ResendVerification(context.Background(), "test@example.com"); err != nil {
			t.Fatalf("ResendVerification failed: %v", err)
		}

		if repo.users[0].VerificationToken == "" {
			t.Error("expected new verification token")
		}
		newToken := repo.users[0].VerificationToken
		if newToken == oldToken {
			t.Error("expected different verification token on resend")
		}
	})

	t.Run("email not found", func(t *testing.T) {
		svc := newAuthService(&mockUserRepo{})
		err := svc.ResendVerification(context.Background(), "nonexistent@example.com")
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("already verified", func(t *testing.T) {
		repo := &mockUserRepo{}
		svc := newAuthService(repo)
		email, password, username := validUser()

		if _, err := svc.Register(context.Background(), email, password, username); err != nil {
			t.Fatalf("register failed: %v", err)
		}

		repo.users[0].EmailVerified = true
		if err := svc.ResendVerification(context.Background(), "test@example.com"); err != nil {
			t.Errorf("expected nil for already verified, got %v", err)
		}
	})
}

func TestRegisterPasswordHashed(t *testing.T) {
	svc := newAuthService(&mockUserRepo{})
	_, err := svc.Register(context.Background(), "test@example.com", "Str0ng!Pass1", "testuser")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
}

func TestLoginTokensAreValidJWTs(t *testing.T) {
	repo := &mockUserRepo{}
	svc := newAuthService(repo)
	m := jwt.NewManager([]byte("test-secret-for-auth"), 15*time.Minute, 7*24*time.Hour)
	email, password, username := validUser()

	if _, err := svc.Register(context.Background(), email, password, username); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	accessToken, refreshToken, _, err := svc.Login(context.Background(), "testuser", "Str0ng!Pass1")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	accessClaims, err := m.ValidateToken(accessToken)
	if err != nil {
		t.Errorf("access token validation failed: %v", err)
	} else {
		if accessClaims.UserID != 1 {
			t.Errorf("access token UserID = %d, want 1", accessClaims.UserID)
		}
		if accessClaims.Role != "user" {
			t.Errorf("access token Role = %q, want %q", accessClaims.Role, "user")
		}
	}

	refreshClaims, err := m.ValidateToken(refreshToken)
	if err != nil {
		t.Errorf("refresh token validation failed: %v", err)
	} else {
		if refreshClaims.UserID != 1 {
			t.Errorf("refresh token UserID = %d, want 1", refreshClaims.UserID)
		}
	}
}

func TestRegisterCreatesUserInRepo(t *testing.T) {
	repo := &mockUserRepo{}
	svc := newAuthService(repo)
	email, password, username := validUser()

	if _, err := svc.Register(context.Background(), email, password, username); err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if len(repo.users) != 1 {
		t.Errorf("expected 1 user in repo, got %d", len(repo.users))
	}

	user := repo.users[0]
	if user.Username != "testuser" {
		t.Errorf("Username = %q, want %q", user.Username, "testuser")
	}
	if user.Email != "test@example.com" {
		t.Errorf("Email = %q, want %q", user.Email, "test@example.com")
	}
	if user.Name != "testuser" {
		t.Errorf("Name = %q, want %q", user.Name, "testuser")
	}
	if user.Role != "user" {
		t.Errorf("Role = %q, want %q", user.Role, "user")
	}
}

func TestRegisterEmptyEmail(t *testing.T) {
	svc := newAuthService(&mockUserRepo{})
	_, err := svc.Register(context.Background(), "", "Str0ng!Pass1", "testuser")
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Errorf("expected ErrInvalidInput, got %v", err)
	}
}
