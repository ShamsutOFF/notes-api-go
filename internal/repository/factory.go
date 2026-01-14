package repository

import (
	"fmt"
	"os"
)

// Config содержит конфигурацию репозитория
type Config struct {
	Type string // "json" или "postgres"
	DSN  string // Для postgres: connection string
	File string // Для json: путь к файлу
}

// NewRepository создает репозиторий на основе конфигурации
func NewRepository(cfg Config) (NoteRepository, error) {
	switch cfg.Type {
	case "json":
		return newJSONRepository(cfg)
	case "postgres":
		return newPostgresRepository(cfg)
	default:
		return nil, fmt.Errorf("unsupported repository type: %s", cfg.Type)
	}
}

// newJSONRepository создает JSON репозиторий
func newJSONRepository(cfg Config) (*JSONRepository, error) {
	if cfg.File == "" {
		cfg.File = "storage/notes.json"
	}
	return NewJSONRepository(cfg.File)
}

// newPostgresRepository создает PostgreSQL репозиторий
func newPostgresRepository(cfg Config) (*PostgresRepository, error) {
	if cfg.DSN == "" {
		// Дефолтные значения для Docker Compose
		cfg.DSN = "host=postgres user=postgres password=postgres dbname=notesdb port=5432 sslmode=disable"
	}
	return NewPostgresRepository(cfg.DSN)
}

// ConfigFromEnv создает конфигурацию из переменных окружения
func ConfigFromEnv() Config {
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "json" // По умолчанию JSON
	}

	return Config{
		Type: storageType,
		DSN:  os.Getenv("DATABASE_URL"),
		File: os.Getenv("STORAGE_FILE"),
	}
}
