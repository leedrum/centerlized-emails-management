package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	KafkaBroker    []string
	KafkaTopic     string
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	return &Config{
		MinIOEndpoint:  viper.GetString("MINIO_ENDPOINT"),
		MinIOAccessKey: viper.GetString("MINIO_ACCESS_KEY"),
		MinIOSecretKey: viper.GetString("MINIO_SECRET_KEY"),
		KafkaBroker:    viper.New().GetStringSlice("KAFKA_BROKERS"),
		KafkaTopic:     viper.GetString("KAFKA_TOPICS"),
	}
}
