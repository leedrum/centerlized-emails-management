package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		log.Printf("[GIN] %s %s %d %s",
			c.Request.Method,
			c.Request.RequestURI,
			c.Writer.Status(),
			duration,
		)
	}
}
