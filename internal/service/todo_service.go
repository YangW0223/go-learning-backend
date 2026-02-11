package service

import (
	"context"
	"time"

	"github.com/yang/go-learning-backend/internal/cache"
	"github.com/yang/go-learning-backend/internal/model"
	"github.com/yang/go-learning-backend/internal/store"
)

const defaultListCacheTTL = 30 * time.Second

// TodoService 定义 Todo 业务编排层契约。
type TodoService interface {
	Create(ctx context.Context, title string) (model.Todo, error)
	List(ctx context.Context) ([]model.Todo, error)
	MarkDone(ctx context.Context, id string) (model.Todo, error)
	Delete(ctx context.Context, id string) error
}

type todoService struct {
	store        store.TodoStore
	cache        cache.TodoCache
	listCacheTTL time.Duration
}

// NewTodoService 创建默认业务服务。
func NewTodoService(todoStore store.TodoStore, todoCache cache.TodoCache, listCacheTTL time.Duration) TodoService {
	if todoStore == nil {
		panic("todo store is nil")
	}
	if todoCache == nil {
		todoCache = cache.NewNoopTodoCache()
	}
	if listCacheTTL <= 0 {
		listCacheTTL = defaultListCacheTTL
	}

	return &todoService{
		store:        todoStore,
		cache:        todoCache,
		listCacheTTL: listCacheTTL,
	}
}

func (s *todoService) Create(ctx context.Context, title string) (model.Todo, error) {
	todo, err := s.store.Create(title)
	if err != nil {
		return model.Todo{}, err
	}

	_ = s.cache.DeleteTodoList(ctx)
	return todo, nil
}

func (s *todoService) List(ctx context.Context) ([]model.Todo, error) {
	cached, hit, err := s.cache.GetTodoList(ctx)
	if err == nil && hit {
		return cached, nil
	}

	todos, err := s.store.List()
	if err != nil {
		return nil, err
	}

	_ = s.cache.SetTodoList(ctx, todos, s.listCacheTTL)
	return todos, nil
}

func (s *todoService) MarkDone(ctx context.Context, id string) (model.Todo, error) {
	todo, err := s.store.MarkDone(id)
	if err != nil {
		return model.Todo{}, err
	}

	_ = s.cache.DeleteTodoList(ctx)
	return todo, nil
}

func (s *todoService) Delete(ctx context.Context, id string) error {
	if err := s.store.Delete(id); err != nil {
		return err
	}

	_ = s.cache.DeleteTodoList(ctx)
	return nil
}
