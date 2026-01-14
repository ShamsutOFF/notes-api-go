package main

import (
	"fmt"
	"log"
	"notes-api/internal/domain"
	"notes-api/internal/repository"
)

func main() {
	// Тестируем JSON репозиторий
	fmt.Println("=== Testing JSON Repository ===")
	testJSONRepository()

	fmt.Println("\n=== Testing PostgreSQL Repository ===")
	testPostgresRepository()
}

func testJSONRepository() {
	repo, err := repository.NewJSONRepository("test_storage.json")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		// Очистка тестового файла
		// os.Remove("test_storage.json")
	}()

	note := &domain.Note{
		Title:   "Test Title",
		Content: "Test Content",
	}

	created, err := repo.Create(note)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created note with ID: %d\n", created.ID)

	notes, err := repo.GetAll()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total notes: %d\n", len(notes))
}

func testPostgresRepository() {
	dsn := "host=localhost user=postgres password=postgres dbname=notesdb port=5432 sslmode=disable"
	repo, err := repository.NewPostgresRepository(dsn)
	if err != nil {
		fmt.Printf("PostgreSQL not available: %v\n", err)
		return
	}

	note := &domain.Note{
		Title:   "Postgres Test",
		Content: "Testing PostgreSQL",
	}

	created, err := repo.Create(note)
	if err != nil {
		fmt.Printf("Failed to create note: %v\n", err)
		return
	}
	fmt.Printf("Created note with ID: %d\n", created.ID)

	notes, err := repo.GetAll()
	if err != nil {
		fmt.Printf("Failed to get notes: %v\n", err)
		return
	}
	fmt.Printf("Total notes in PostgreSQL: %d\n", len(notes))
}
