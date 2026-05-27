package auth

import "github.com/google/uuid"

type TokenService interface {
	GenerateAccessToken(userID uuid.UUID, role string) (string, error)
	GenerateRefreshToken(userID uuid.UUID) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type Claims struct {
	UserID uuid.UUID
	Role   string
}

type PasswordService interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}
