package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/leedrum/centerlized-emails-management/api-gateway/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading it")
	}
	gin.SetMode(os.Getenv("GIN_MODE"))

	r := gin.Default()
	routes.SetupRoutes(r)

	log.Println("API Gateway listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
