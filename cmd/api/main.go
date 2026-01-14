package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"test/internal/handler"
	"test/internal/repository"
	"test/internal/service"
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

	// Настраиваем маршрутизатор
	r := mux.NewRouter()

	// API маршруты
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/notes", noteHandler.CreateNote).Methods("POST")
	api.HandleFunc("/notes", noteHandler.GetAllNotes).Methods("GET")
	api.HandleFunc("/notes/{id}", noteHandler.GetNoteByID).Methods("GET")
	api.HandleFunc("/notes/{id}", noteHandler.UpdateNote).Methods("PUT")
	api.HandleFunc("/notes/{id}", noteHandler.DeleteNote).Methods("DELETE")

	// Middleware для логирования
	r.Use(loggingMiddleware)

	// Настраиваем сервер
	port := ":8080"
	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Storage file: %s", filename)

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// loggingMiddleware логирует HTTP запросы
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
