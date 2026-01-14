package app

import (
	"notes-api/internal/handler"
	"notes-api/internal/repository"
	"notes-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// App представляет основное приложение с внедренными зависимостями
type App struct {
	repo    repository.NoteRepository
	service *service.NoteService
	handler *handler.NoteHandler
	fiber   *fiber.App
}

// New создает новое приложение с внедрением зависимостей
func New(repo repository.NoteRepository) *App {
	// Создаем цепочку зависимостей (Dependency Injection)
	noteService := service.NewNoteService(repo)
	noteHandler := handler.NewNoteHandler(noteService)

	// Создаем Fiber приложение
	app := fiber.New(fiber.Config{
		AppName: "Notes API",
	})

	// Middleware
	app.Use(logger.New())

	// Настраиваем маршруты
	setupRoutes(app, noteHandler)

	return &App{
		repo:    repo,
		service: noteService,
		handler: noteHandler,
		fiber:   app,
	}
}

// setupRoutes настраивает все API маршруты
func setupRoutes(app *fiber.App, handler *handler.NoteHandler) {
	api := app.Group("/api")

	// Notes endpoints
	api.Post("/notes", handler.CreateNote)
	api.Get("/notes", handler.GetAllNotes)
	api.Get("/notes/:id", handler.GetNoteByID)
	api.Put("/notes/:id", handler.UpdateNote)
	api.Delete("/notes/:id", handler.DeleteNote)

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "notes-api",
		})
	})
}

// Run запускает приложение на указанном адресе
func (a *App) Run(addr string) error {
	return a.fiber.Listen(addr)
}

// Shutdown корректно останавливает приложение
func (a *App) Shutdown() error {
	return a.fiber.Shutdown()
}

// Fiber возвращает экземпляр Fiber приложения (для тестов или кастомной конфигурации)
func (a *App) Fiber() *fiber.App {
	return a.fiber
}

// Service возвращает сервис (для тестов)
func (a *App) Service() *service.NoteService {
	return a.service
}
