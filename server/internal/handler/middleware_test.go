package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/pkg/apperror"
	"github.com/raiyin/artserver/pkg/jwt"
)

func newJWTManager(t *testing.T) *jwt.Manager {
	t.Helper()
	return jwt.NewManager([]byte("test-secret-for-middleware"), 15*time.Minute, 7*24*time.Hour)
}

func generateToken(t *testing.T, m *jwt.Manager, user *domain.User) string {
	t.Helper()
	token, err := m.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}
	return token
}

func setupMiddlewareTest(handlerFunc gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/test", handlerFunc)
	return r
}

type testResponse struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func newTestUser() *domain.User {
	return &domain.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
	}
}

func newTestAdmin() *domain.User {
	return &domain.User{
		ID:       2,
		Username: "admin",
		Email:    "admin@example.com",
		Role:     "admin",
	}
}

func extractError(t *testing.T, body []byte) apperror.APIError {
	t.Helper()
	var apiErr apperror.APIError
	if err := json.Unmarshal(body, &apiErr); err != nil {
		t.Fatalf("unmarshal error response: %v", err)
	}
	return apiErr
}

func TestAuthMiddlewareNoHeader(t *testing.T) {
	m := newJWTManager(t)
	r := setupMiddlewareTest(AuthMiddleware(m))

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	apiErr := extractError(t, w.Body.Bytes())
	if apiErr.Code != "UNAUTHORIZED" {
		t.Errorf("expected code UNAUTHORIZED, got %s", apiErr.Code)
	}
}

func TestAuthMiddlewareInvalidFormat(t *testing.T) {
	m := newJWTManager(t)
	r := setupMiddlewareTest(AuthMiddleware(m))

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat token123")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	apiErr := extractError(t, w.Body.Bytes())
	if apiErr.Code != "INVALID_TOKEN" {
		t.Errorf("expected code INVALID_TOKEN, got %s", apiErr.Code)
	}
}

func TestAuthMiddlewareInvalidToken(t *testing.T) {
	m := newJWTManager(t)
	r := setupMiddlewareTest(AuthMiddleware(m))

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token-value")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}

	apiErr := extractError(t, w.Body.Bytes())
	if apiErr.Code != "INVALID_TOKEN" {
		t.Errorf("expected code INVALID_TOKEN, got %s", apiErr.Code)
	}
}

func TestAuthMiddlewareExpiredToken(t *testing.T) {
	m := jwt.NewManager([]byte("test-secret-for-middleware"), -1*time.Minute, 7*24*time.Hour)
	token := generateToken(t, m, newTestUser())
	r := setupMiddlewareTest(AuthMiddleware(m))

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddlewareValidToken(t *testing.T) {
	m := newJWTManager(t)
	token := generateToken(t, m, newTestUser())

	r := gin.New()
	r.Use(AuthMiddleware(m))
	r.GET("/test", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		username, _ := c.Get("username")
		role, _ := c.Get("role")
		email, _ := c.Get("email")

		if userID.(int64) != 1 {
			t.Errorf("user_id = %d, want 1", userID)
		}
		if username.(string) != "testuser" {
			t.Errorf("username = %q, want %q", username, "testuser")
		}
		if role.(string) != "user" {
			t.Errorf("role = %q, want %q", role, "user")
		}
		if email.(string) != "test@example.com" {
			t.Errorf("email = %q, want %q", email, "test@example.com")
		}
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAuthMiddlewareAdminRoleInContext(t *testing.T) {
	m := newJWTManager(t)
	token := generateToken(t, m, newTestAdmin())

	r := gin.New()
	r.Use(AuthMiddleware(m))
	r.GET("/test", func(c *gin.Context) {
		role, _ := c.Get("role")
		if role.(string) != "admin" {
			t.Errorf("role = %q, want %q", role, "admin")
		}
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAdminMiddlewareMissingRole(t *testing.T) {
	r := gin.New()
	r.Use(AdminMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestAdminMiddlewareUserRole(t *testing.T) {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("role", "user")
		c.Next()
	})
	r.Use(AdminMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}

	apiErr := extractError(t, w.Body.Bytes())
	if apiErr.Code != "FORBIDDEN" {
		t.Errorf("expected code FORBIDDEN, got %s", apiErr.Code)
	}
}

func TestAdminMiddlewareAdminRole(t *testing.T) {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("role", "admin")
		c.Next()
	})
	r.Use(AdminMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAuthMiddlewareBearerPrefixHandling(t *testing.T) {
	m := newJWTManager(t)
	token := generateToken(t, m, newTestUser())

	t.Run("lowercase bearer", func(t *testing.T) {
		r := setupMiddlewareTest(AuthMiddleware(m))
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "bearer "+token)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401 for lowercase bearer, got %d", w.Code)
		}
	})

	t.Run("bearer with extra spaces", func(t *testing.T) {
		r := setupMiddlewareTest(AuthMiddleware(m))
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer  "+token)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected 401 for double space, got %d", w.Code)
		}
	})
}

func TestAuthMiddlewareOptionalSkipsWhenNoHeader(t *testing.T) {
	m := newJWTManager(t)
	r := gin.New()
	r.Use(AuthMiddlewareOptional(m))
	r.GET("/test", func(c *gin.Context) {
		_, exists := c.Get("user_id")
		if exists {
			t.Error("user_id should not be set for unauthenticated requests")
		}
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAuthMiddlewareOptionalValidToken(t *testing.T) {
	m := newJWTManager(t)
	token := generateToken(t, m, newTestUser())

	r := gin.New()
	r.Use(AuthMiddlewareOptional(m))
	r.GET("/test", func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			t.Fatal("user_id should be set")
		}
		if userID.(int64) != 1 {
			t.Errorf("user_id = %d, want 1", userID)
		}
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAuthMiddlewareOptionalInvalidToken(t *testing.T) {
	m := newJWTManager(t)
	r := gin.New()
	r.Use(AuthMiddlewareOptional(m))
	r.GET("/test", func(c *gin.Context) {
		_, exists := c.Get("user_id")
		if exists {
			t.Error("user_id should not be set for invalid tokens")
		}
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
