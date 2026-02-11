package cache

import (
	"context"
	"time"

	"github.com/yang/go-learning-backend/internal/model"
)

// TodoCache 定义 Todo 列表缓存的最小能力。
type TodoCache interface {
	GetTodoList(ctx context.Context) ([]model.Todo, bool, error)
	SetTodoList(ctx context.Context, todos []model.Todo, ttl time.Duration) error
	DeleteTodoList(ctx context.Context) error
	Ping(ctx context.Context) error
}

// NoopTodoCache 在未启用 Redis 时提供透明降级。
type NoopTodoCache struct{}

func NewNoopTodoCache() *NoopTodoCache {
	return &NoopTodoCache{}
}

func (c *NoopTodoCache) GetTodoList(_ context.Context) ([]model.Todo, bool, error) {
	return nil, false, nil
}

func (c *NoopTodoCache) SetTodoList(_ context.Context, _ []model.Todo, _ time.Duration) error {
	return nil
}

func (c *NoopTodoCache) DeleteTodoList(_ context.Context) error {
	return nil
}

func (c *NoopTodoCache) Ping(_ context.Context) error {
	return nil
}
