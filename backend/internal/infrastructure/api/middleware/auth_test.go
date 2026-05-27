package middleware

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appAuth "github.com/myinquisitor/backend/internal/application/auth"
)

type mockTokenService struct {
	userID uuid.UUID
	role   string
	err    error
}

func (m *mockTokenService) GenerateAccessToken(_ uuid.UUID, _ string) (string, error) {
	return "token", nil
}

func (m *mockTokenService) GenerateRefreshToken(_ uuid.UUID) (string, error) {
	return "refresh", nil
}

func (m *mockTokenService) ValidateToken(_ string) (*appAuth.Claims, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &appAuth.Claims{UserID: m.userID, Role: m.role}, nil
}

func TestAuthMiddleware_NoHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mw := NewAuthMiddleware(&mockTokenService{userID: uuid.New(), role: "user"})
	r := gin.New()
	r.GET("/protected", mw.Authenticate(), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidScheme(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mw := NewAuthMiddleware(&mockTokenService{userID: uuid.New(), role: "user"})
	r := gin.New()
	r.GET("/protected", mw.Authenticate(), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Basic token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	uid := uuid.New()
	mw := NewAuthMiddleware(&mockTokenService{userID: uid, role: "super_admin"})
	r := gin.New()
	r.GET("/protected", mw.Authenticate(), func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		role, _ := c.Get("role")
		c.JSON(200, gin.H{"user_id": userID, "role": role})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), uid.String()) {
		t.Errorf("expected response to contain user_id")
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mw := NewAuthMiddleware(&mockTokenService{err: errors.New("invalid token")})
	r := gin.New()
	r.GET("/protected", mw.Authenticate(), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalid")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("expected 401, got %d", w.Code)
	}
}
