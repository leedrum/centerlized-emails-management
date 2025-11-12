package main

import (
	"log"

	"github.com/leedrum/centerlized-emails-management/tenant/config"
	"github.com/leedrum/centerlized-emails-management/tenant/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/tursodatabase/turso-go"
)

func main() {

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Println("Can not load the env file")
	}
	gin.SetMode(cfg.GinMode)

	r := gin.Default()
	routes.SetupRoutes(r, cfg)

	log.Println("API Gateway listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
