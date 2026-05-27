package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort         string
	DatabaseURL        string
	JWTSecret          string
	JWTExpiration      time.Duration
	RefreshExpiration  time.Duration
	EncryptionKey      string
	AllowedOrigins     []string
	Environment        string
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: no .env file found: %v", err)
	}

	jwtExp := getEnv("JWT_ACCESS_EXPIRATION", "30m")
	jwtDuration, err := time.ParseDuration(jwtExp)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_ACCESS_EXPIRATION: %w", err)
	}

	refreshExp := getEnv("JWT_REFRESH_EXPIRATION", "720h")
	refreshDuration, err := time.ParseDuration(refreshExp)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_REFRESH_EXPIRATION: %w", err)
	}

	port := getEnv("SERVER_PORT", "8080")
	if _, err := strconv.Atoi(port); err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %s", port)
	}

	dbURL := getEnv("DATABASE_URL", "")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	encKey := getEnv("ENCRYPTION_KEY", "")
	if encKey == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY is required")
	}

	origins := getEnv("ALLOWED_ORIGINS", "http://localhost:5173")

	return &Config{
		ServerPort:        port,
		DatabaseURL:       dbURL,
		JWTSecret:         jwtSecret,
		JWTExpiration:     jwtDuration,
		RefreshExpiration: refreshDuration,
		EncryptionKey:     encKey,
		AllowedOrigins:    strings.Split(origins, ","),
		Environment:       getEnv("ENVIRONMENT", "development"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
