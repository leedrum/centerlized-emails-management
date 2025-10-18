package middlewares

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Limit each IP to 5 requests per second with burst of 10
const (
	rateLimitRPS   = 5
	rateLimitBurst = 10
)

var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getVisitorLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(rateLimitRPS), rateLimitBurst)
		visitors[ip] = limiter
	}

	return limiter
}

// RateLimitMiddleware limits requests per IP
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getVisitorLimiter(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}
