package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUrl     string
	JWTSecret string
}

func Load() *Config {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("⚠ WARNING: DATABASE_URL is empty")
	}

	return &Config{
		DBUrl:     dbURL,
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
