package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCORSMiddleware_AllowedOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mw := NewCORSMiddleware([]string{"http://example.com"})
	r := gin.New()
	r.Use(mw.Handler())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "http://example.com" {
		t.Errorf("expected Allow-Origin header, got %q", w.Header().Get("Access-Control-Allow-Origin"))
	}
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestCORSMiddleware_DisallowedOrigin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mw := NewCORSMiddleware([]string{"http://example.com"})
	r := gin.New()
	r.Use(mw.Handler())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://evil.com")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "" {
		t.Errorf("expected no Allow-Origin for disallowed origin, got %q", w.Header().Get("Access-Control-Allow-Origin"))
	}
	if w.Code != 200 {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestCORSMiddleware_Preflight(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mw := NewCORSMiddleware([]string{"http://example.com"})
	r := gin.New()
	r.Use(mw.Handler())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 204 {
		t.Errorf("expected 204 for preflight, got %d", w.Code)
	}
}

func TestCORSMiddleware_MultipleOrigins(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mw := NewCORSMiddleware([]string{"http://a.com", "http://b.com"})
	r := gin.New()
	r.Use(mw.Handler())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	for _, origin := range []string{"http://a.com", "http://b.com"} {
		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Origin", origin)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Header().Get("Access-Control-Allow-Origin") != origin {
			t.Errorf("expected Allow-Origin %q, got %q", origin, w.Header().Get("Access-Control-Allow-Origin"))
		}
	}
}
