package main

import (
	"log"

	"github.com/leedrum/centerlized-emails-management/api-gateway/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	routes.SetupRoutes(r)

	log.Println("API Gateway listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
