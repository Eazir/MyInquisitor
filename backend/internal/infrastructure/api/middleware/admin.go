package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myinquisitor/backend/internal/infrastructure/api/response"
)

type AdminMiddleware struct{}

func NewAdminMiddleware() *AdminMiddleware {
	return &AdminMiddleware{}
}

func (m *AdminMiddleware) RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "access denied")
			c.Abort()
			return
		}

		if role.(string) != "super_admin" {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "super admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}
