package db

import (
	"os"

	"gorm.io/gorm"
)

var DB *gorm.DB

type DBStorage interface {
	Init()
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
