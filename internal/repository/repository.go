package repository

import (
	"notes-api/internal/domain"
)

// NoteRepository определяет интерфейс для работы с заметками
type NoteRepository interface {
	Create(note *domain.Note) (*domain.Note, error)
	GetAll() ([]*domain.Note, error)
	GetByID(id int64) (*domain.Note, error)
	Update(id int64, note *domain.Note) (*domain.Note, error)
	Delete(id int64) error
}
