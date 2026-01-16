package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerConfig
	DatabaseConfig
	JwtConfig
}

type ServerConfig struct {
	ServerPort string
	ServerEnv  string
}

type DatabaseConfig struct {
	DatabaseURL string
}

type JwtConfig struct {
	SecretKey string
	Expiry    time.Duration
}

func Load() *Config {
	if godotenv.Load() != nil {
		log.Println("no .env to load")
	}
	return &Config{
		ServerConfig: ServerConfig{
			ServerPort: os.Getenv("SERVER_PORT"),
			ServerEnv:  os.Getenv("SERVER_ENV"),
		},
		DatabaseConfig: DatabaseConfig{
			DatabaseURL: os.Getenv("DATABASE_URL"),
		},
		JwtConfig: JwtConfig{
			SecretKey: os.Getenv("SECRET_KEY"),
			Expiry:    24 * time.Hour,
		},
	}
}
