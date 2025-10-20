package db

import (
	"fmt"
	"log"

	"github.com/leedrum/centerlized-emails-management/auth/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresStorage struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	DSN      string
}

func (p *PostgresStorage) Init() {
	p.LoadEnv()
	p.GenerateDSN()

	var err error
	DB, err = gorm.Open(postgres.Open(p.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed to migrate DB:", err)
	}
}

func (p *PostgresStorage) LoadEnv() {
	p.Host = getEnv("DB_HOST", "localhost")
	p.User = getEnv("DB_USER", "root")
	p.Password = getEnv("DB_PASSWORD", "123456")
	p.DBName = getEnv("DB_NAME", "cem-auth")
	p.Port = getEnv("DB_PORT", "5432")
}

func (p *PostgresStorage) GenerateDSN() {
	p.DSN = fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", p.Host, p.User, p.Password, p.DBName, p.Port)
}
