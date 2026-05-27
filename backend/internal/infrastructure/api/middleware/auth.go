package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	appAuth "github.com/myinquisitor/backend/internal/application/auth"
	"github.com/myinquisitor/backend/internal/infrastructure/api/response"
)

type AuthMiddleware struct {
	token appAuth.TokenService
}

func NewAuthMiddleware(token appAuth.TokenService) *AuthMiddleware {
	return &AuthMiddleware{token: token}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "AUTH_REQUIRED", "authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, http.StatusUnauthorized, "INVALID_TOKEN", "authorization header must be Bearer token")
			c.Abort()
			return
		}

		claims, err := m.token.ValidateToken(parts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "INVALID_TOKEN", "token is invalid or expired")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
