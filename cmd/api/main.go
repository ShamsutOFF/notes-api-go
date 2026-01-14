package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"notes-api/internal/handler"
	"notes-api/internal/repository"
	"notes-api/internal/service"
)

func main() {
	// Настраиваем путь к файлу хранилища
	filename := "storage/notes.json"
	if os.Getenv("STORAGE_PATH") != "" {
		filename = os.Getenv("STORAGE_PATH")
	}

	// Создаем репозиторий
	repo, err := repository.NewJSONRepository(filename)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
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
	port := ":8080"
	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Storage file: %s", filename)

	if err := app.Listen(port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
