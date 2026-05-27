package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAdminMiddleware_NoRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mw := NewAdminMiddleware()
	r := gin.New()
	r.GET("/admin", mw.RequireSuperAdmin(), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestAdminMiddleware_UserRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mw := NewAdminMiddleware()
	r := gin.New()
	r.GET("/admin", mw.RequireSuperAdmin(), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/admin", nil)
	req.Header.Set("X-Role", "user")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestAdminMiddleware_SuperAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mw := NewAdminMiddleware()
	r := gin.New()
	r.GET("/admin", func(c *gin.Context) {
		c.Set("role", "super_admin")
	}, mw.RequireSuperAdmin(), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAdminMiddleware_UserRoleSet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mw := NewAdminMiddleware()
	r := gin.New()
	r.GET("/admin", func(c *gin.Context) {
		c.Set("role", "user")
	}, mw.RequireSuperAdmin(), func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/admin", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 403 {
		t.Errorf("expected 403, got %d", w.Code)
	}
}
