package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"notes-api/internal/domain"
	"notes-api/internal/service"
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

// GetAllNotes обрабатывает получение всех заметок
func (h *NoteHandler) GetAllNotes(c *fiber.Ctx) error {
	// Получаем все заметки через сервис
	notes, err := h.service.GetAllNotes()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	// Возвращаем ответ
	return c.JSON(notes)
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
