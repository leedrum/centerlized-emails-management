package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	GmailClientID     string
	GmailClientSecret string
	GmailTokenPath    string
	S3Bucket          string
	KafkaBroker       []string
	KafkaTopic        string
	Tenants           string
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	return &Config{
		GmailClientID:     viper.GetString("GMAIL_CLIENT_ID"),
		GmailClientSecret: viper.GetString("GMAIL_CLIENT_SECRET"),
		GmailTokenPath:    viper.GetString("GMAIL_TOKEN_PATH"),
		S3Bucket:          viper.GetString("S3_BUCKET"),
		KafkaBroker:       viper.New().GetStringSlice("KAFKA_BROKER"),
		KafkaTopic:        viper.GetString("KAFKA_TOPIC"),
	}
}
