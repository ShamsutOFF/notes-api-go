package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config содержит все настройки приложения
type Config struct {
	Port       string
	Repository struct {
		Type string
		DSN  string
		File string
	}
}

// Load загружает конфигурацию из .env файла и переменных окружения
func Load() *Config {
	// Пытаемся загрузить .env файл, но не падаем если его нет
	_ = godotenv.Load()

	cfg := &Config{}

	// Server config
	cfg.Port = getEnv("PORT", "8080")

	// Repository config
	cfg.Repository.Type = getEnv("STORAGE_TYPE", "json")
	cfg.Repository.DSN = os.Getenv("DATABASE_URL")
	cfg.Repository.File = getEnv("STORAGE_FILE", "storage/notes.json")

	return cfg
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
