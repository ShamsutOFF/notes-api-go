package repository

import (
	"fmt"

	"notes-api/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(dsn string) (*PostgresRepository, error) {
	// Подключаемся к базе данных
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Автомиграция - создаст таблицу если её нет
	if err := db.AutoMigrate(&domain.Note{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %v", err)
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) Create(note *domain.Note) (*domain.Note, error) {
	result := r.db.Create(note)
	if result.Error != nil {
		return nil, result.Error
	}
	return note, nil
}

func (r *PostgresRepository) GetAll(limit, offset int) ([]*domain.Note, int, error) {
	var notes []*domain.Note
	var total int64

	// Сначала получаем общее количество записей
	if err := r.db.Model(&domain.Note{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Затем получаем данные с пагинацией
	result := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&notes)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return notes, int(total), nil
}

func (r *PostgresRepository) GetByID(id int64) (*domain.Note, error) {
	var note domain.Note
	result := r.db.First(&note, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNoteNotFound
		}
		return nil, result.Error
	}
	return &note, nil
}

func (r *PostgresRepository) Update(id int64, note *domain.Note) (*domain.Note, error) {
	// Сначала проверяем существование заметки
	var existingNote domain.Note
	result := r.db.First(&existingNote, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNoteNotFound
		}
		return nil, result.Error
	}

	// Обновляем только необходимые поля
	updates := map[string]interface{}{
		"title":      note.Title,
		"content":    note.Content,
		"updated_at": gorm.Expr("NOW()"),
	}

	result = r.db.Model(&domain.Note{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return nil, result.Error
	}

	// Получаем обновленную запись
	var updatedNote domain.Note
	r.db.First(&updatedNote, id)
	return &updatedNote, nil
}

func (r *PostgresRepository) Delete(id int64) error {
	result := r.db.Delete(&domain.Note{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNoteNotFound
	}

	return nil
}
