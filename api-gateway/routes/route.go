package routes

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/leedrum/centerlized-emails-management/api-gateway/middlewares"

	"github.com/gin-gonic/gin"
)

// ReverseProxy returns a handler forwarding requests to targetHost
func ReverseProxy(targetHost string) gin.HandlerFunc {
	return func(c *gin.Context) {
		target, err := url.Parse(targetHost)
		if err != nil {
			c.String(http.StatusInternalServerError, "Bad target URL")
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		c.Request.URL.Path = c.Param("proxyPath")

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// SetupRoutes configures all routes and applies middleware
func SetupRoutes(r *gin.Engine) {
	r.Use(middlewares.RateLimit())
	r.Use(middlewares.Logger())

	authGroup := r.Group("/auth")
	authGroup.Any("/*proxyPath", ReverseProxy("http://localhost:9001"))
}
