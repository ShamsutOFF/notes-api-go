package domain

import (
	"errors"
	"time"
)

// Note представляет структуру заметки
type Note struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Validate проверяет корректность данных заметки
func (n *Note) Validate() error {
	if n.Title == "" {
		return errors.New("title cannot be empty")
	}
	if n.Content == "" {
		return errors.New("content cannot be empty")
	}
	return nil
}

// CreateNoteRequest представляет запрос на создание заметки
type CreateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// UpdateNoteRequest представляет запрос на обновление заметки
type UpdateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
