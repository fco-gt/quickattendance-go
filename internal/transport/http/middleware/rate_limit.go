package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

// getVisitorLimiter retrieves or creates a limiter for a given IP with specific rate and burst
func getVisitorLimiter(ip string, r rate.Limit, b int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(r, b)
		visitors[ip] = limiter
	}

	return limiter
}

// RateLimitByIP applies a rate limit per client IP
func RateLimitByIP(r rate.Limit, b int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getVisitorLimiter(ip, r, b)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Límite de peticiones excedido para tu IP. Intenta más tarde.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
