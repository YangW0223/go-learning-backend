package memory

import (
	"sync"
	"time"

	"github.com/yang/go-learning-backend/internal/model"
	"github.com/yang/go-learning-backend/internal/store"
)

// TodoStore keeps todos in memory for quick local development.
type TodoStore struct {
	mu    sync.RWMutex
	todos []model.Todo
}

func NewTodoStore() *TodoStore {
	return &TodoStore{todos: make([]model.Todo, 0)}
}

func (s *TodoStore) Create(title string) (model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo := model.Todo{
		ID:        generateID(),
		Title:     title,
		Done:      false,
		CreatedAt: time.Now().UTC(),
	}
	s.todos = append(s.todos, todo)

	return todo, nil
}

func (s *TodoStore) List() ([]model.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]model.Todo, len(s.todos))
	copy(result, s.todos)
	return result, nil
}

func (s *TodoStore) MarkDone(id string) (model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.todos {
		if s.todos[i].ID == id {
			s.todos[i].Done = true
			return s.todos[i], nil
		}
	}

	return model.Todo{}, store.ErrTodoNotFound
}

func generateID() string {
	return time.Now().UTC().Format("20060102150405.000000000")
}
