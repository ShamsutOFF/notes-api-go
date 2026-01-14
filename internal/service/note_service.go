package service

import (
	"notes-api/internal/domain"
	"notes-api/internal/repository"
)

// NoteService реализует бизнес-логику для работы с заметками
type NoteService struct {
	repo repository.NoteRepository // Изменено на интерфейс!
}

// NewNoteService создает новый сервис
func NewNoteService(repo repository.NoteRepository) *NoteService { // Изменено на интерфейс!
	return &NoteService{repo: repo}
}

// CreateNote создает новую заметку
func (s *NoteService) CreateNote(req domain.CreateNoteRequest) (*domain.Note, error) {
	// Создаем новую заметку
	note := &domain.Note{
		Title:   req.Title,
		Content: req.Content,
	}

	// Валидируем
	if err := note.Validate(); err != nil {
		return nil, err
	}

	// Сохраняем через репозиторий
	return s.repo.Create(note)
}

// GetAllNotes возвращает все заметки
func (s *NoteService) GetAllNotes() ([]*domain.Note, error) {
	return s.repo.GetAll()
}

// GetNoteByID возвращает заметку по ID
func (s *NoteService) GetNoteByID(id int64) (*domain.Note, error) {
	return s.repo.GetByID(id)
}

// UpdateNote обновляет заметку
func (s *NoteService) UpdateNote(id int64, req domain.UpdateNoteRequest) (*domain.Note, error) {
	// Проверяем существование заметки
	if _, err := s.repo.GetByID(id); err != nil {
		return nil, err
	}

	// Создаем обновленную заметку
	note := &domain.Note{
		Title:   req.Title,
		Content: req.Content,
	}

	// Валидируем
	if err := note.Validate(); err != nil {
		return nil, err
	}

	// Обновляем через репозиторий
	return s.repo.Update(id, note)
}

// DeleteNote удаляет заметку
func (s *NoteService) DeleteNote(id int64) error {
	return s.repo.Delete(id)
}
