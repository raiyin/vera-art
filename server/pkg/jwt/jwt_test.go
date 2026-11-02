package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/raiyin/artserver/internal/domain"
)

func newTestManager(t *testing.T) *Manager {
	t.Helper()
	return NewManager([]byte("test-secret-key-12345"), 15*time.Minute, 7*24*time.Hour)
}

func testUser() *domain.User {
	return &domain.User{
		ID:       42,
		Username: "testuser",
		Email:    "test@example.com",
		Role:     "user",
	}
}

func TestGenerateAccessToken(t *testing.T) {
	m := newTestManager(t)
	user := testUser()

	token, err := m.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	m := newTestManager(t)
	user := testUser()

	token, err := m.GenerateRefreshToken(user)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestAccessTokenRoundTrip(t *testing.T) {
	m := newTestManager(t)
	user := testUser()

	tokenString, err := m.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != 42 {
		t.Errorf("UserID = %d, want 42", claims.UserID)
	}
	if claims.Username != "testuser" {
		t.Errorf("Username = %q, want %q", claims.Username, "testuser")
	}
	if claims.Email != "test@example.com" {
		t.Errorf("Email = %q, want %q", claims.Email, "test@example.com")
	}
	if claims.Role != "user" {
		t.Errorf("Role = %q, want %q", claims.Role, "user")
	}
}

func TestRefreshTokenRoundTrip(t *testing.T) {
	m := newTestManager(t)
	user := testUser()

	tokenString, err := m.GenerateRefreshToken(user)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}

	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != 42 {
		t.Errorf("UserID = %d, want 42", claims.UserID)
	}
}

func TestValidateExpiredToken(t *testing.T) {
	m := NewManager([]byte("test-secret-key-12345"), -1*time.Minute, -1*time.Minute)
	user := testUser()

	tokenString, err := m.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	_, err = m.ValidateToken(tokenString)
	if err != domain.ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken for expired token, got %v", err)
	}
}

func TestValidateMalformedToken(t *testing.T) {
	m := newTestManager(t)

	_, err := m.ValidateToken("not-a-jwt-token")
	if err != domain.ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}

func TestValidateEmptyToken(t *testing.T) {
	m := newTestManager(t)

	_, err := m.ValidateToken("")
	if err != domain.ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken, got %v", err)
	}
}

func TestValidateTokenWithWrongSecret(t *testing.T) {
	m1 := NewManager([]byte("my-secret-key-one-12345"), 15*time.Minute, 7*24*time.Hour)
	m2 := NewManager([]byte("my-secret-key-two-67890"), 15*time.Minute, 7*24*time.Hour)

	tokenString, err := m1.GenerateAccessToken(testUser())
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	_, err = m2.ValidateToken(tokenString)
	if err != domain.ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken for wrong secret, got %v", err)
	}
}

func TestValidateTokenWithWrongSigningMethod(t *testing.T) {
	m := newTestManager(t)
	user := testUser()

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("sign with None method failed: %v", err)
	}

	_, err = m.ValidateToken(tokenString)
	if err != domain.ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken for wrong signing method, got %v", err)
	}
}

func TestAccessTokenExpiry(t *testing.T) {
	m := NewManager([]byte("test-secret-key-12345"), 1*time.Second, 7*24*time.Hour)
	user := testUser()

	tokenString, err := m.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		t.Fatalf("ValidateToken should succeed immediately: %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("UserID = %d, want 42", claims.UserID)
	}

	time.Sleep(1500 * time.Millisecond)

	_, err = m.ValidateToken(tokenString)
	if err != domain.ErrInvalidToken {
		t.Errorf("expected ErrInvalidToken after expiry, got %v", err)
	}
}

func TestTokenAdminRole(t *testing.T) {
	m := newTestManager(t)
	admin := &domain.User{
		ID:       1,
		Username: "admin",
		Email:    "admin@example.com",
		Role:     "admin",
	}

	tokenString, err := m.GenerateAccessToken(admin)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if claims.Role != "admin" {
		t.Errorf("Role = %q, want %q", claims.Role, "admin")
	}
}

func TestTokenIATClaim(t *testing.T) {
	m := newTestManager(t)
	user := testUser()

	before := time.Now().Add(-time.Second)
	tokenString, err := m.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.IssuedAt.Time.Before(before) {
		t.Error("IssuedAt is before the token was generated")
	}
	if claims.IssuedAt.Time.After(time.Now().Add(time.Second)) {
		t.Error("IssuedAt is in the future")
	}
}

func TestAccessAndRefreshTokenDifferences(t *testing.T) {
	m := newTestManager(t)
	user := testUser()

	accessToken, err := m.GenerateAccessToken(user)
	if err != nil {
		t.Fatalf("GenerateAccessToken failed: %v", err)
	}

	refreshToken, err := m.GenerateRefreshToken(user)
	if err != nil {
		t.Fatalf("GenerateRefreshToken failed: %v", err)
	}

	if accessToken == refreshToken {
		t.Error("access and refresh tokens should be different")
	}
}
