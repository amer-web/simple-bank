package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

var Source Config

type Config struct {
	Port          string        `mapstructure:"SERVER_PORT"`
	DBHost        string        `mapstructure:"DB_HOST"`
	DBPort        int           `mapstructure:"DB_PORT"`
	DBUser        string        `mapstructure:"DB_USER"`
	DBPassword    string        `mapstructure:"DB_PASSWORD"`
	DBName        string        `mapstructure:"DB_NAME"`
	DRIVER        string        `mapstructure:"DRIVER"`
	TOKENDRIVER   string        `mapstructure:"TOKEN_DRIVER"`
	TOKENKEY      string        `mapstructure:"TOKEN_KEY"`
	TOKENDURATION time.Duration `mapstructure:"TOKEN_DURATION"`
}

func LoadConfig() error {

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("unable to decode into struct: %w", err)
	}
	if err := viper.Unmarshal(&Source); err != nil {
		return fmt.Errorf("unable to decode into struct: %w", err)
	}
	return nil
}
