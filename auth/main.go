package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/leedrum/centerlized-emails-management/auth/db"
	"github.com/leedrum/centerlized-emails-management/auth/routes"
)

func main() {
	connectDB(&db.PostgresStorage{})
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading it")
	}
	gin.SetMode(os.Getenv("GIN_MODE"))

	r := gin.Default()
	routes.SetupRoutes(r)

	log.Println("Auth Service listening on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func connectDB(dbs db.DBStorage) {
	dbs.Init()
}
