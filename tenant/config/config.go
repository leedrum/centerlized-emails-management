package config

import "github.com/spf13/viper"

type Config struct {
	DBSource      string
	KafkaRetryMax int
	Brokers       []string
	GinMode       string
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
