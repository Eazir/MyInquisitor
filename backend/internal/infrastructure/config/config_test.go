package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	os.Clearenv()
	os.Setenv("DATABASE_URL", "postgres://localhost/test")
	os.Setenv("JWT_SECRET", "supersecretkeythatisatleast32charslong!!")
	os.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.ServerPort != "8080" {
		t.Errorf("expected 8080, got %s", cfg.ServerPort)
	}
	if cfg.JWTExpiration.String() != "30m0s" {
		t.Errorf("expected 30m0s, got %s", cfg.JWTExpiration)
	}
	if cfg.RefreshExpiration.String() != "720h0m0s" {
		t.Errorf("expected 720h0m0s, got %s", cfg.RefreshExpiration)
	}
	if len(cfg.AllowedOrigins) != 1 || cfg.AllowedOrigins[0] != "http://localhost:5173" {
		t.Errorf("expected default origins, got %v", cfg.AllowedOrigins)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	os.Clearenv()
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DATABASE_URL", "postgres://user:pass@host:5432/db")
	os.Setenv("JWT_SECRET", "supersecretkeythatisatleast32charslong!!")
	os.Setenv("JWT_ACCESS_EXPIRATION", "15m")
	os.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	os.Setenv("ALLOWED_ORIGINS", "http://a.com,http://b.com")
	os.Setenv("ENVIRONMENT", "production")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.ServerPort != "9090" {
		t.Errorf("expected 9090, got %s", cfg.ServerPort)
	}
	if cfg.JWTExpiration.String() != "15m0s" {
		t.Errorf("expected 15m0s, got %s", cfg.JWTExpiration)
	}
	if cfg.Environment != "production" {
		t.Errorf("expected production, got %s", cfg.Environment)
	}
	if len(cfg.AllowedOrigins) != 2 {
		t.Errorf("expected 2 origins, got %d", len(cfg.AllowedOrigins))
	}
}

func TestLoad_MissingDatabaseURL(t *testing.T) {
	os.Clearenv()
	os.Setenv("JWT_SECRET", "supersecretkeythatisatleast32charslong!!")
	os.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing DATABASE_URL")
	}
}

func TestLoad_MissingJWTSecret(t *testing.T) {
	os.Clearenv()
	os.Setenv("DATABASE_URL", "postgres://localhost/test")
	os.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing JWT_SECRET")
	}
}

func TestLoad_InvalidJWTExpiration(t *testing.T) {
	os.Clearenv()
	os.Setenv("DATABASE_URL", "postgres://localhost/test")
	os.Setenv("JWT_SECRET", "supersecretkeythatisatleast32charslong!!")
	os.Setenv("ENCRYPTION_KEY", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	os.Setenv("JWT_ACCESS_EXPIRATION", "invalid")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid JWT expiration")
	}
}
