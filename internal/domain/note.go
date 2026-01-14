package domain

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Note представляет структуру заметки
type Note struct {
	ID        int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Title     string         `json:"title" gorm:"not null"`
	Content   string         `json:"content" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
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
