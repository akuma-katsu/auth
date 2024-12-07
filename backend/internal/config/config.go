package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"log"
)

type AppConfig struct {
	DataBaseAddress  string `env:"DB_ADDRESS"`
	DataBasePort     string `env:"DB_PORT"`
	DatabaseName     string `env:"DB_NAME"`
	DataBaseUsername string `env:"DB_USERNAME"`
	DataBasePassword string `env:"DB_PASSWORD"`
	Secret           string `env:"SECRET_KEY"`
	AccessLifeTime   string `env:"ACCESS_LIFETIME"` // часы
}

func GetCfg() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := AppConfig{}
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatal("Error parsing config")
	}

	return &cfg
}
