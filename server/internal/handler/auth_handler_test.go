package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/dto"
	"github.com/raiyin/artserver/pkg/apperror"
)

// mockAuthService implements port.AuthService.
type mockAuthService struct {
	register          func(ctx context.Context, email, password, username string) (*domain.User, error)
	login             func(ctx context.Context, username, password string) (string, string, *domain.User, error)
	refreshToken      func(ctx context.Context, refreshToken string) (string, string, error)
	verifyEmail       func(ctx context.Context, token string) error
	resendVerification func(ctx context.Context, email string) error
}

func (m *mockAuthService) Register(ctx context.Context, email, password, username string) (*domain.User, error) {
	if m.register != nil {
		return m.register(ctx, email, password, username)
	}
	return &domain.User{ID: 1, Username: username, Email: email, Role: "user"}, nil
}

func (m *mockAuthService) Login(ctx context.Context, username, password string) (string, string, *domain.User, error) {
	if m.login != nil {
		return m.login(ctx, username, password)
	}
	return "access-token", "refresh-token", &domain.User{ID: 1, Username: username, Role: "user"}, nil
}

func (m *mockAuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	if m.refreshToken != nil {
		return m.refreshToken(ctx, refreshToken)
	}
	return "new-access", "new-refresh", nil
}

func (m *mockAuthService) VerifyEmail(ctx context.Context, token string) error {
	if m.verifyEmail != nil {
		return m.verifyEmail(ctx, token)
	}
	return nil
}

func (m *mockAuthService) ResendVerification(ctx context.Context, email string) error {
	if m.resendVerification != nil {
		return m.resendVerification(ctx, email)
	}
	return nil
}

func setupAuthRouter(h *AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.POST("/refresh", h.RefreshToken)
	r.GET("/verify-email", h.VerifyEmail)
	r.POST("/resend-verification", h.ResendVerification)
	return r
}

func authHandlerWithService(svc *mockAuthService) *AuthHandler {
	return NewAuthHandler(svc)
}

func TestRegister(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockAuthService{
			register: func(_ context.Context, email, password, username string) (*domain.User, error) {
				return &domain.User{ID: 1, Username: username, Email: email, Role: "user", Name: username}, nil
			},
		}
		h := authHandlerWithService(svc)
		r := setupAuthRouter(h)

		body := `{"username":"newuser","email":"new@example.com","password":"ValidP@ss1"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d. Body: %s", w.Code, w.Body.String())
		}

		var resp dto.UserDTO
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if resp.Username != "newuser" {
			t.Errorf("Username = %q, want %q", resp.Username, "newuser")
		}
		if resp.Email != "new@example.com" {
			t.Errorf("Email = %q, want %q", resp.Email, "new@example.com")
		}
		if resp.Role != "user" {
			t.Errorf("Role = %q, want %q", resp.Role, "user")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		svc := &mockAuthService{}
		h := authHandlerWithService(svc)
		r := setupAuthRouter(h)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/register", strings.NewReader("not-json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		svc := &mockAuthService{}
		h := authHandlerWithService(svc)
		r := setupAuthRouter(h)

		body := `{"username":"u"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("duplicate user", func(t *testing.T) {
		svc := &mockAuthService{
			register: func(_ context.Context, _, _, _ string) (*domain.User, error) {
				return nil, apperror.ErrDuplicate
			},
		}
		h := authHandlerWithService(svc)
		r := setupAuthRouter(h)

		body := `{"username":"existing","email":"existing@example.com","password":"ValidP@ss1"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusConflict {
			t.Errorf("expected 409, got %d", w.Code)
		}
	})
}

func TestLogin(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockAuthService{
			login: func(_ context.Context, username, password string) (string, string, *domain.User, error) {
				return "access-token-123", "refresh-token-456",
					&domain.User{ID: 1, Username: username, Email: "test@example.com", Role: "user", Name: username},
					nil
			},
		}
		h := authHandlerWithService(svc)
		r := setupAuthRouter(h)

		body := `{"username":"testuser","password":"ValidP@ss1"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var resp dto.AuthResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if resp.AccessToken != "access-token-123" {
			t.Errorf("AccessToken = %q, want %q", resp.AccessToken, "access-token-123")
		}
		if resp.RefreshToken != "refresh-token-456" {
			t.Errorf("RefreshToken = %q, want %q", resp.RefreshToken, "refresh-token-456")
		}
		if resp.User.Username != "testuser" {
			t.Errorf("Username = %q, want %q", resp.User.Username, "testuser")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		svc := &mockAuthService{}
		h := authHandlerWithService(svc)
		r := setupAuthRouter(h)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/login", strings.NewReader("bad"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		svc := &mockAuthService{
			login: func(_ context.Context, _, _ string) (string, string, *domain.User, error) {
				return "", "", nil, apperror.ErrUnauthorized
			},
		}
		h := authHandlerWithService(svc)
		r := setupAuthRouter(h)

		body := `{"username":"testuser","password":"WrongP@ss1"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", w.Code)
		}
	})

	t.Run("missing fields", func(t *testing.T) {
		svc := &mockAuthService{}
		h := authHandlerWithService(svc)
		r := setupAuthRouter(h)

		body := `{}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestRefreshToken(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockAuthService{
			refreshToken: func(_ context.Context, _ string) (string, string, error) {
				return "new-access-token", "new-refresh-token", nil
			},
		}
		h := authHandlerWithService(svc)
		r := setupAuthRouter(h)

		body := `{"refresh_token":"valid-refresh-token"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/refresh", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var resp struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if resp.AccessToken != "new-access-token" {
			t.Errorf("AccessToken = %q, want %q", resp.AccessToken, "new-access-token")
		}
		if resp.RefreshToken != "new-refresh-token" {
			t.Errorf("RefreshToken = %q, want %q", resp.RefreshToken, "new-refresh-token")
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		svc := &mockAuthService{}
		h := authHandlerWithService(svc)
		r := setupAuthRouter(h)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/refresh", strings.NewReader("{}"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("invalid refresh token", func(t *testing.T) {
		svc := &mockAuthService{
			refreshToken: func(_ context.Context, _ string) (string, string, error) {
				return "", "", apperror.ErrInvalidToken
			},
		}
		h := authHandlerWithService(svc)
		r := setupAuthRouter(h)

		body := `{"refresh_token":"bad-token"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/refresh", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", w.Code)
		}
	})
}

func TestVerifyEmail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockAuthService{
			verifyEmail: func(_ context.Context, token string) error {
				return nil
			},
		}
		h := authHandlerWithService(svc)
		r := setupAuthRouter(h)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/verify-email?token=valid-token", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var resp map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if resp["message"] != "Email verified successfully" {
			t.Errorf("message = %q, want %q", resp["message"], "Email verified successfully")
		}
	})

	t.Run("missing token", func(t *testing.T) {
		svc := &mockAuthService{}
		h := authHandlerWithService(svc)
		r := setupAuthRouter(h)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/verify-email", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("invalid token", func(t *testing.T) {
		svc := &mockAuthService{
			verifyEmail: func(_ context.Context, token string) error {
				return apperror.ErrInvalidToken
			},
		}
		h := authHandlerWithService(svc)
		r := setupAuthRouter(h)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/verify-email?token=bad-token", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401, got %d", w.Code)
		}
	})
}

func TestResendVerification(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockAuthService{
			resendVerification: func(_ context.Context, email string) error {
				return nil
			},
		}
		h := authHandlerWithService(svc)
		r := setupAuthRouter(h)

		body := `{"email":"test@example.com"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/resend-verification", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("always returns 200 even when email not found", func(t *testing.T) {
		svc := &mockAuthService{
			resendVerification: func(_ context.Context, email string) error {
				return domain.ErrNotFound
			},
		}
		h := authHandlerWithService(svc)
		r := setupAuthRouter(h)

		body := `{"email":"nonexistent@example.com"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/resend-verification", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200 (privacy), got %d", w.Code)
		}
	})

	t.Run("invalid email format", func(t *testing.T) {
		svc := &mockAuthService{}
		h := authHandlerWithService(svc)
		r := setupAuthRouter(h)

		body := `{"email":"not-an-email"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/resend-verification", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		svc := &mockAuthService{}
		h := authHandlerWithService(svc)
		r := setupAuthRouter(h)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/resend-verification", strings.NewReader("bad"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("already verified returns 200", func(t *testing.T) {
		svc := &mockAuthService{
			resendVerification: func(_ context.Context, email string) error {
				return nil
			},
		}
		h := authHandlerWithService(svc)
		r := setupAuthRouter(h)

		body := `{"email":"verified@example.com"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/resend-verification", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})
}

func TestAuthHandlerImplementsPort(t *testing.T) {
	svc := &mockAuthService{}
	h := NewAuthHandler(svc)
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}

func TestRegisterReturnsUserDTO(t *testing.T) {
	svc := &mockAuthService{
		register: func(_ context.Context, email, password, username string) (*domain.User, error) {
			return &domain.User{
				ID:            10,
				Username:      username,
				Email:         email,
				Name:          "Display Name",
				Role:          "user",
				AvatarPath:    "avatars/10.jpg",
				EmailVerified: false,
			}, nil
		},
	}
	h := authHandlerWithService(svc)
	r := setupAuthRouter(h)

	body := `{"username":"johndoe","email":"john@example.com","password":"Str0ng!Pass"}`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	var resp dto.UserDTO
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.ID != 10 {
		t.Errorf("ID = %d, want 10", resp.ID)
	}
	if resp.Name != "Display Name" {
		t.Errorf("Name = %q, want %q", resp.Name, "Display Name")
	}
	if resp.AvatarPath != "avatars/10.jpg" {
		t.Errorf("AvatarPath = %q, want %q", resp.AvatarPath, "avatars/10.jpg")
	}
	if resp.EmailVerified {
		t.Error("EmailVerified should be false")
	}
}
