package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leedrum/centerlized-emails-management/auth/handlers"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r.POST("/login", handlers.Login)
	r.POST("/register", handlers.Register)
}
