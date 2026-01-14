package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"test/internal/domain"
	"test/internal/service"
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
func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateNoteRequest

	// Парсим JSON тело запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Создаем заметку через сервис
	note, err := h.service.CreateNote(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Возвращаем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

// GetAllNotes обрабатывает получение всех заметок
func (h *NoteHandler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	// Получаем все заметки через сервис
	notes, err := h.service.GetAllNotes()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Возвращаем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notes)
}

// GetNoteByID обрабатывает получение заметки по ID
func (h *NoteHandler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Парсим ID
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid note ID")
		return
	}

	// Получаем заметку через сервис
	note, err := h.service.GetNoteByID(id)
	if err != nil {
		if err.Error() == "note not found" {
			writeError(w, http.StatusNotFound, "note not found")
		} else {
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	// Возвращаем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(note)
}

// UpdateNote обрабатывает обновление заметки
func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Парсим ID
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid note ID")
		return
	}

	var req domain.UpdateNoteRequest

	// Парсим JSON тело запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Обновляем заметку через сервис
	note, err := h.service.UpdateNote(id, req)
	if err != nil {
		if err.Error() == "note not found" {
			writeError(w, http.StatusNotFound, "note not found")
		} else {
			writeError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	// Возвращаем ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(note)
}

// DeleteNote обрабатывает удаление заметки
func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Парсим ID
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid note ID")
		return
	}

	// Удаляем заметку через сервис
	if err := h.service.DeleteNote(id); err != nil {
		if err.Error() == "note not found" {
			writeError(w, http.StatusNotFound, "note not found")
		} else {
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	// Возвращаем пустой ответ с кодом 200
	w.WriteHeader(http.StatusOK)
}

// writeError записывает ошибку в ответ
func writeError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
