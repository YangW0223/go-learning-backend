package store

import (
	"errors"

	"github.com/yang/go-learning-backend/internal/model"
)

var ErrTodoNotFound = errors.New("todo not found")

// TodoStore defines the storage behavior for todos.
type TodoStore interface {
	Create(title string) (model.Todo, error)
	List() ([]model.Todo, error)
	MarkDone(id string) (model.Todo, error)
}
