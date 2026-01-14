package main

import (
	"log"
	"os"

	"notes-api/internal/handler"
	"notes-api/internal/repository"
	"notes-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Получаем тип хранилища из переменных окружения
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "json" // По умолчанию JSON
	}

	var repo repository.NoteRepository
	var err error

	switch storageType {
	case "json":
		// Конфигурация JSON хранилища
		filename := os.Getenv("STORAGE_FILE")
		if filename == "" {
			filename = "storage/notes.json"
		}

		repo, err = repository.NewJSONRepository(filename)
		if err != nil {
			log.Fatalf("Failed to create JSON repository: %v", err)
		}
		log.Printf("Using JSON storage: %s", filename)

	case "postgres":
		// Конфигурация PostgreSQL
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			dsn = "host=postgres user=postgres password=postgres dbname=notesdb port=5432 sslmode=disable"
		}

		repo, err = repository.NewPostgresRepository(dsn)
		if err != nil {
			log.Fatalf("Failed to create PostgreSQL repository: %v", err)
		}
		log.Printf("Using PostgreSQL storage")

	default:
		log.Fatalf("Unsupported storage type: %s", storageType)
	}

	// Создаем сервис
	noteService := service.NewNoteService(repo)

	// Создаем обработчики
	noteHandler := handler.NewNoteHandler(noteService)

	// Создаем Fiber приложение
	app := fiber.New()

	// Middleware для логирования
	app.Use(logger.New())

	// API маршруты
	api := app.Group("/api")
	api.Post("/notes", noteHandler.CreateNote)
	api.Get("/notes", noteHandler.GetAllNotes)
	api.Get("/notes/:id", noteHandler.GetNoteByID)
	api.Put("/notes/:id", noteHandler.UpdateNote)
	api.Delete("/notes/:id", noteHandler.DeleteNote)

	// Настраиваем порт
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
