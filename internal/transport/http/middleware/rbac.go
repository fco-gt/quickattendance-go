package middleware

import (
	"autoattendance-go/internal/domain"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		currentRole := userRole.(domain.Role)

		isAllowed := slices.Contains(allowedRoles, currentRole)

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		c.Next()
	}
}
