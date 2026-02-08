package middleware

import (
	"autoattendance-go/pkg/security"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(jwtSvc *security.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			slog.Warn("Unauthorized access attempt: token required")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			slog.Warn("Unauthorized access attempt: invalid token format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token format"})
			return
		}

		claims, err := jwtSvc.Verify(parts[1])
		if err != nil {
			slog.Warn("Unauthorized access attempt: invalid token", "error", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		slog.Info("User authenticated",
			"user_id", claims.UserID,
			"agency_id", claims.AgencyID,
			"role", claims.Role,
		)

		c.Set("user_id", claims.UserID)
		c.Set("agency_id", claims.AgencyID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
