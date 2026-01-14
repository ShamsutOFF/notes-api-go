package handler

import (
	"strconv"

	"notes-api/internal/domain"
	"notes-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

// NoteHandler обрабатывает HTTP запросы для заметок
type NoteHandler struct {
	service *service.NoteService
}

// NewNoteHandler создает новый обработчик
func NewNoteHandler(service *service.NoteService) *NoteHandler {
	return &NoteHandler{service: service}
}

// CreateNote обрабатывает создание заметки
func (h *NoteHandler) CreateNote(c *fiber.Ctx) error {
	var req domain.CreateNoteRequest

	// Парсим JSON тело запроса
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Создаем заметку через сервис
	note, err := h.service.CreateNote(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Возвращаем ответ
	return c.Status(fiber.StatusCreated).JSON(note)
}

// GetAllNotes обрабатывает получение всех заметок с пагинацией
func (h *NoteHandler) GetAllNotes(c *fiber.Ctx) error {
	// Получаем параметры пагинации из query string
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	// Валидация параметров
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Ограничиваем максимальный лимит
	}

	// Вычисляем offset
	offset := (page - 1) * limit

	// Получаем заметки через сервис с пагинацией
	notes, total, err := h.service.GetAllNotes(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	// Рассчитываем метаданные пагинации
	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit // ceil деление
	}

	// Возвращаем ответ с пагинацией
	return c.JSON(fiber.Map{
		"data": notes,
		"meta": fiber.Map{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": totalPages,
			"hasNext":    page < totalPages,
			"hasPrev":    page > 1,
		},
	})
}

// GetNoteByID обрабатывает получение заметки по ID
func (h *NoteHandler) GetNoteByID(c *fiber.Ctx) error {
	idStr := c.Params("id")

	// Парсим ID
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid note ID",
		})
	}

	// Получаем заметку через сервис
	note, err := h.service.GetNoteByID(id)
	if err != nil {
		if err.Error() == "note not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "note not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	// Возвращаем ответ
	return c.JSON(note)
}

// UpdateNote обрабатывает обновление заметки
func (h *NoteHandler) UpdateNote(c *fiber.Ctx) error {
	idStr := c.Params("id")

	// Парсим ID
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid note ID",
		})
	}

	var req domain.UpdateNoteRequest

	// Парсим JSON тело запроса
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Обновляем заметку через сервис
	note, err := h.service.UpdateNote(id, req)
	if err != nil {
		if err.Error() == "note not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "note not found",
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Возвращаем ответ
	return c.JSON(note)
}

// DeleteNote обрабатывает удаление заметки
func (h *NoteHandler) DeleteNote(c *fiber.Ctx) error {
	idStr := c.Params("id")

	// Парсим ID
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid note ID",
		})
	}

	// Удаляем заметку через сервис
	if err := h.service.DeleteNote(id); err != nil {
		if err.Error() == "note not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "note not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	// Возвращаем пустой ответ с кодом 200
	return c.SendStatus(fiber.StatusOK)
}
