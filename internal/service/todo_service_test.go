package service

import (
	"context"
	"testing"
	"time"

	"github.com/yang/go-learning-backend/internal/model"
	"github.com/yang/go-learning-backend/internal/store"
)

type fakeTodoStore struct {
	todos     []model.Todo
	listCalls int
}

func (s *fakeTodoStore) Create(title string) (model.Todo, error) {
	todo := model.Todo{ID: "new-id", Title: title, Done: false}
	s.todos = append(s.todos, todo)
	return todo, nil
}

func (s *fakeTodoStore) List() ([]model.Todo, error) {
	s.listCalls++
	items := make([]model.Todo, len(s.todos))
	copy(items, s.todos)
	return items, nil
}

func (s *fakeTodoStore) MarkDone(id string) (model.Todo, error) {
	for i := range s.todos {
		if s.todos[i].ID == id {
			s.todos[i].Done = true
			return s.todos[i], nil
		}
	}
	return model.Todo{}, store.ErrTodoNotFound
}

func (s *fakeTodoStore) Delete(id string) error {
	for i := range s.todos {
		if s.todos[i].ID == id {
			s.todos = append(s.todos[:i], s.todos[i+1:]...)
			return nil
		}
	}
	return store.ErrTodoNotFound
}

type fakeTodoCache struct {
	hit         bool
	items       []model.Todo
	getCalls    int
	setCalls    int
	deleteCalls int
}

func (c *fakeTodoCache) GetTodoList(_ context.Context) ([]model.Todo, bool, error) {
	c.getCalls++
	if !c.hit {
		return nil, false, nil
	}
	items := make([]model.Todo, len(c.items))
	copy(items, c.items)
	return items, true, nil
}

func (c *fakeTodoCache) SetTodoList(_ context.Context, todos []model.Todo, _ time.Duration) error {
	c.setCalls++
	c.hit = true
	c.items = make([]model.Todo, len(todos))
	copy(c.items, todos)
	return nil
}

func (c *fakeTodoCache) DeleteTodoList(_ context.Context) error {
	c.deleteCalls++
	c.hit = false
	c.items = nil
	return nil
}

func (c *fakeTodoCache) Ping(_ context.Context) error {
	return nil
}

func TestTodoService_ListHitsCache(t *testing.T) {
	store := &fakeTodoStore{
		todos: []model.Todo{{ID: "store-1", Title: "from-store", Done: false}},
	}
	cache := &fakeTodoCache{
		hit:   true,
		items: []model.Todo{{ID: "cache-1", Title: "from-cache", Done: true}},
	}
	service := NewTodoService(store, cache, 45*time.Second)

	items, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(items) != 1 || items[0].ID != "cache-1" {
		t.Fatalf("expected cached item, got %+v", items)
	}
	if store.listCalls != 0 {
		t.Fatalf("expected store list calls 0, got %d", store.listCalls)
	}
	if cache.setCalls != 0 {
		t.Fatalf("expected cache set calls 0, got %d", cache.setCalls)
	}
}

func TestTodoService_ListMissWritesCache(t *testing.T) {
	store := &fakeTodoStore{
		todos: []model.Todo{{ID: "store-1", Title: "from-store", Done: false}},
	}
	cache := &fakeTodoCache{hit: false}
	service := NewTodoService(store, cache, 45*time.Second)

	items, err := service.List(context.Background())
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(items) != 1 || items[0].ID != "store-1" {
		t.Fatalf("expected store item, got %+v", items)
	}
	if store.listCalls != 1 {
		t.Fatalf("expected store list calls 1, got %d", store.listCalls)
	}
	if cache.setCalls != 1 {
		t.Fatalf("expected cache set calls 1, got %d", cache.setCalls)
	}
}

func TestTodoService_MutationsInvalidateCache(t *testing.T) {
	store := &fakeTodoStore{
		todos: []model.Todo{{ID: "todo-1", Title: "task", Done: false}},
	}
	cache := &fakeTodoCache{hit: true, items: store.todos}
	service := NewTodoService(store, cache, 45*time.Second)

	if _, err := service.Create(context.Background(), "new task"); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if _, err := service.MarkDone(context.Background(), "todo-1"); err != nil {
		t.Fatalf("MarkDone returned error: %v", err)
	}
	if err := service.Delete(context.Background(), "todo-1"); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	if cache.deleteCalls != 3 {
		t.Fatalf("expected cache delete calls 3, got %d", cache.deleteCalls)
	}
}
