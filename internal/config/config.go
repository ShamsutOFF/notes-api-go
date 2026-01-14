package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string

	StorageType string // "json" or "postgres"

	// JSON storage
	StorageFile string

	// PostgreSQL
	DatabaseURL string
}

func Load() *Config {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		StorageType: getEnv("STORAGE_TYPE", "json"),
		StorageFile: getEnv("STORAGE_FILE", "storage/notes.json"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
