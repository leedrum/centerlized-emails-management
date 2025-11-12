package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leedrum/centerlized-emails-management/tenant/config"
	"github.com/leedrum/centerlized-emails-management/tenant/handlers"
)

func SetupRoutes(r *gin.Engine, cfg config.Config) {
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	r.POST("/create", func(c *gin.Context) {
		handlers.Create(c, cfg)
	})
	r.POST("/delete", func(c *gin.Context) {
		handlers.Delete(c, cfg)
	})
}
