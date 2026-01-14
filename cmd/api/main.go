package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"notes-api/internal/app"
	"notes-api/internal/config"
	"notes-api/internal/repository"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Создаем репозиторий на основе конфигурации
	repoCfg := repository.Config{
		Type: cfg.Repository.Type,
		DSN:  cfg.Repository.DSN,
		File: cfg.Repository.File,
	}

	repo, err := repository.NewRepository(repoCfg)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}

	log.Printf("Using %s storage", cfg.Repository.Type)
	if cfg.Repository.Type == "json" {
		log.Printf("Storage file: %s", cfg.Repository.File)
	}

	// Создаем приложение с внедренной зависимостью
	application := app.New(repo)

	// Настраиваем graceful shutdown
	setupGracefulShutdown(application)

	// Запускаем приложение
	log.Printf("Server starting on :%s", cfg.Port)
	if err := application.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// setupGracefulShutdown настраивает корректное завершение работы
func setupGracefulShutdown(app *app.App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")

		if err := app.Shutdown(); err != nil {
			log.Fatalf("Error shutting down: %v", err)
		}

		log.Println("Server stopped")
		os.Exit(0)
	}()
}
