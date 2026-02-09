package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/yang/go-learning-backend/internal/store"
)

type TodoHandler struct {
	store store.TodoStore
}

func NewTodoHandler(todoStore store.TodoStore) *TodoHandler {
	return &TodoHandler{store: todoStore}
}

type createTodoRequest struct {
	Title string `json:"title"`
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}

	todo, err := h.store.Create(req.Title)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create todo"})
		return
	}

	writeJSON(w, http.StatusCreated, todo)
}

func (h *TodoHandler) List(w http.ResponseWriter, _ *http.Request) {
	todos, err := h.store.List()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list todos"})
		return
	}

	writeJSON(w, http.StatusOK, todos)
}

func (h *TodoHandler) MarkDone(w http.ResponseWriter, r *http.Request, id string) {
	todo, err := h.store.MarkDone(id)
	if err != nil {
		if errors.Is(err, store.ErrTodoNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "todo not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update todo"})
		return
	}

	writeJSON(w, http.StatusOK, todo)
}
