package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"notes-api/internal/domain"
)

var (
	ErrNoteNotFound = errors.New("note not found")
)

// JSONRepository реализует хранение заметок в JSON файле
type JSONRepository struct {
	filename string
	mu       sync.RWMutex
	notes    map[int64]*domain.Note
	nextID   int64
}

// NewJSONRepository создает новый JSON репозиторий
func NewJSONRepository(filename string) (*JSONRepository, error) {
	repo := &JSONRepository{
		filename: filename,
		notes:    make(map[int64]*domain.Note),
		nextID:   1,
	}

	// Загружаем данные из файла при старте
	if err := repo.loadFromFile(); err != nil {
		return nil, fmt.Errorf("failed to load data from file: %w", err)
	}

	// Находим максимальный ID для генерации новых
	for id := range repo.notes {
		if id >= repo.nextID {
			repo.nextID = id + 1
		}
	}

	return repo, nil
}

// loadFromFile загружает данные из JSON файла
func (r *JSONRepository) loadFromFile() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Проверяем существует ли файл
	if _, err := os.Stat(r.filename); os.IsNotExist(err) {
		// Если файла нет, создаем пустой
		return r.saveToFile()
	}

	// Читаем файл
	data, err := os.ReadFile(r.filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Парсим JSON
	if len(data) > 0 {
		var notes []*domain.Note
		if err := json.Unmarshal(data, &notes); err != nil {
			return fmt.Errorf("failed to parse JSON: %w", err)
		}

		// Загружаем в map
		for _, note := range notes {
			r.notes[note.ID] = note
		}
	}

	return nil
}

// saveToFile сохраняет данные в JSON файл
func (r *JSONRepository) saveToFile() error {
	// Преобразуем map в slice для сохранения
	notes := make([]*domain.Note, 0, len(r.notes))
	for _, note := range r.notes {
		notes = append(notes, note)
	}

	// Сериализуем в JSON с отступами
	data, err := json.MarshalIndent(notes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Записываем в файл
	if err := os.WriteFile(r.filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Create создает новую заметку
func (r *JSONRepository) Create(note *domain.Note) (*domain.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Устанавливаем ID и временные метки
	note.ID = r.nextID
	now := time.Now()
	note.CreatedAt = now
	note.UpdatedAt = now

	// Сохраняем в map
	r.notes[note.ID] = note
	r.nextID++

	// Сохраняем в файл
	if err := r.saveToFile(); err != nil {
		delete(r.notes, note.ID) // Откатываем изменение в случае ошибки
		return nil, err
	}

	return note, nil
}

// GetAll возвращает заметки с пагинацией
func (r *JSONRepository) GetAll(limit, offset int) ([]*domain.Note, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Общее количество записей
	total := len(r.notes)

	// Получаем все заметки в нужном порядке
	allNotes := make([]*domain.Note, 0, len(r.notes))
	for _, note := range r.notes {
		allNotes = append(allNotes, note)
	}

	// Сортируем по created_at DESC (новые первыми)
	sort.Slice(allNotes, func(i, j int) bool {
		return allNotes[i].CreatedAt.After(allNotes[j].CreatedAt)
	})

	// Применяем пагинацию
	start := offset
	end := offset + limit
	if start > len(allNotes) {
		start = len(allNotes)
	}
	if end > len(allNotes) {
		end = len(allNotes)
	}

	return allNotes[start:end], total, nil
}

// GetByID возвращает заметку по ID
func (r *JSONRepository) GetByID(id int64) (*domain.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	note, exists := r.notes[id]
	if !exists {
		return nil, ErrNoteNotFound
	}

	return note, nil
}

// Update обновляет заметку
func (r *JSONRepository) Update(id int64, note *domain.Note) (*domain.Note, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Проверяем существование заметки
	existingNote, exists := r.notes[id]
	if !exists {
		return nil, ErrNoteNotFound
	}

	// Обновляем поля
	existingNote.Title = note.Title
	existingNote.Content = note.Content
	existingNote.UpdatedAt = time.Now()

	// Сохраняем в файл
	if err := r.saveToFile(); err != nil {
		return nil, err
	}

	return existingNote, nil
}

// Delete удаляет заметку
func (r *JSONRepository) Delete(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Проверяем существование заметки
	if _, exists := r.notes[id]; !exists {
		return ErrNoteNotFound
	}

	// Удаляем из map
	delete(r.notes, id)

	// Сохраняем в файл
	if err := r.saveToFile(); err != nil {
		return err
	}

	return nil
}
