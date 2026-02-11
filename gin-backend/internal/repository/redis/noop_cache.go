package redis

import (
	"context"
	"time"

	"github.com/yang/go-learning-backend/gin-backend/internal/model"
)

// NoopTodoCache 在未启用 Redis 时提供降级实现。
type NoopTodoCache struct{}

// NewNoopTodoCache 创建 no-op 缓存实例。
func NewNoopTodoCache() *NoopTodoCache {
	return &NoopTodoCache{}
}

func (c *NoopTodoCache) Ping(_ context.Context) error {
	return nil
}

func (c *NoopTodoCache) GetList(_ context.Context, _ string) ([]model.Todo, bool, error) {
	return nil, false, nil
}

func (c *NoopTodoCache) SetList(_ context.Context, _ string, _ []model.Todo, _ time.Duration) error {
	return nil
}

func (c *NoopTodoCache) DeleteList(_ context.Context, _ string) error {
	return nil
}
