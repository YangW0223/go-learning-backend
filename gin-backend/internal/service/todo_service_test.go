package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/yang/go-learning-backend/gin-backend/internal/model"
	"github.com/yang/go-learning-backend/gin-backend/internal/repository"
)

type fakeTodoRepo struct {
	todos         map[string]model.Todo
	listByUserErr error
}

func newFakeTodoRepo() *fakeTodoRepo {
	return &fakeTodoRepo{todos: map[string]model.Todo{}}
}

func (r *fakeTodoRepo) Create(_ context.Context, todo model.Todo) (model.Todo, error) {
	r.todos[todo.ID] = todo
	return todo, nil
}

func (r *fakeTodoRepo) ListByUserID(_ context.Context, userID string) ([]model.Todo, error) {
	if r.listByUserErr != nil {
		return nil, r.listByUserErr
	}
	out := make([]model.Todo, 0)
	for _, item := range r.todos {
		if item.UserID == userID {
			out = append(out, item)
		}
	}
	return out, nil
}

func (r *fakeTodoRepo) GetByID(_ context.Context, id string) (model.Todo, error) {
	item, ok := r.todos[id]
	if !ok {
		return model.Todo{}, repository.ErrNotFound
	}
	return item, nil
}

func (r *fakeTodoRepo) Update(_ context.Context, todo model.Todo) (model.Todo, error) {
	if _, ok := r.todos[todo.ID]; !ok {
		return model.Todo{}, repository.ErrNotFound
	}
	r.todos[todo.ID] = todo
	return todo, nil
}

func (r *fakeTodoRepo) Delete(_ context.Context, id string) error {
	if _, ok := r.todos[id]; !ok {
		return repository.ErrNotFound
	}
	delete(r.todos, id)
	return nil
}

type fakeTodoCache struct {
	items map[string][]model.Todo
}

func newFakeTodoCache() *fakeTodoCache {
	return &fakeTodoCache{items: map[string][]model.Todo{}}
}

func (c *fakeTodoCache) GetList(_ context.Context, userID string) ([]model.Todo, bool, error) {
	v, ok := c.items[userID]
	if !ok {
		return nil, false, nil
	}
	return v, true, nil
}

func (c *fakeTodoCache) SetList(_ context.Context, userID string, todos []model.Todo, _ time.Duration) error {
	c.items[userID] = todos
	return nil
}

func (c *fakeTodoCache) DeleteList(_ context.Context, userID string) error {
	delete(c.items, userID)
	return nil
}

func (c *fakeTodoCache) Ping(_ context.Context) error {
	return nil
}

func TestTodoService_CreateAndList(t *testing.T) {
	repo := newFakeTodoRepo()
	cache := newFakeTodoCache()
	svc := NewTodoService(repo, cache, 30*time.Second)

	created, err := svc.Create(context.Background(), "u1", "learn gin")
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if created.ID == "" || created.Title != "learn gin" || created.UserID != "u1" {
		t.Fatalf("unexpected created todo: %+v", created)
	}

	items, err := svc.List(context.Background(), "u1")
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(items))
	}
}

func TestTodoService_UpdateNotFound(t *testing.T) {
	repo := newFakeTodoRepo()
	cache := newFakeTodoCache()
	svc := NewTodoService(repo, cache, 30*time.Second)

	title := "x"
	_, err := svc.Update(context.Background(), "u1", "missing", UpdateTodoInput{Title: &title})
	if err == nil {
		t.Fatal("expected not found error")
	}
}

func TestTodoService_ListInternalError(t *testing.T) {
	repo := newFakeTodoRepo()
	repo.listByUserErr = errors.New("db down")
	cache := newFakeTodoCache()
	svc := NewTodoService(repo, cache, 30*time.Second)

	_, err := svc.List(context.Background(), "u1")
	if err == nil {
		t.Fatal("expected internal error")
	}
}
