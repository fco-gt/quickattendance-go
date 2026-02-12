package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"quickattendance-go/internal/domain"
	"slices"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		slog.Debug("Checking role", "user_role", userRole, "exists", exists)

		if !exists {
			slog.Warn("Access denied: role not found in context")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		currentRole, ok := userRole.(domain.Role)
		if !ok {
			slog.Error("Access denied: role in context is not of type domain.Role", "actual_type", fmt.Sprintf("%T", userRole))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		isAllowed := slices.Contains(allowedRoles, currentRole)

		if !isAllowed {
			slog.Warn("Access denied: insufficient permissions",
				"current_role", currentRole,
				"allowed_roles", allowedRoles,
			)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		c.Next()
	}
}
