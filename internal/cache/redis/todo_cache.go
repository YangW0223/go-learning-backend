package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yang/go-learning-backend/internal/model"
)

const todoListCacheKey = "todo:list:v1"

// TodoCache 是 Redis 版本的 Todo 列表缓存实现。
type TodoCache struct {
	client *Client
}

// NewTodoCache 创建 Redis Todo 缓存实例。
func NewTodoCache(client *Client) *TodoCache {
	return &TodoCache{client: client}
}

// Ping 执行 Redis 连通性检查。
func (c *TodoCache) Ping(ctx context.Context) error {
	if c.client == nil {
		return fmt.Errorf("redis client is nil")
	}
	return c.client.Ping(ctx)
}

// GetTodoList 从 Redis 获取 Todo 列表缓存。
func (c *TodoCache) GetTodoList(ctx context.Context) ([]model.Todo, bool, error) {
	if c.client == nil {
		return nil, false, fmt.Errorf("redis client is nil")
	}

	payload, hit, err := c.client.Get(ctx, todoListCacheKey)
	if err != nil {
		return nil, false, err
	}
	if !hit {
		return nil, false, nil
	}

	var todos []model.Todo
	if err := json.Unmarshal([]byte(payload), &todos); err != nil {
		return nil, false, fmt.Errorf("decode todo list cache: %w", err)
	}

	return todos, true, nil
}

// SetTodoList 把 Todo 列表写入 Redis。
func (c *TodoCache) SetTodoList(ctx context.Context, todos []model.Todo, ttl time.Duration) error {
	if c.client == nil {
		return fmt.Errorf("redis client is nil")
	}

	encoded, err := json.Marshal(todos)
	if err != nil {
		return fmt.Errorf("encode todo list cache: %w", err)
	}

	if err := c.client.SetEX(ctx, todoListCacheKey, string(encoded), ttl); err != nil {
		return err
	}
	return nil
}

// DeleteTodoList 删除 Todo 列表缓存键。
func (c *TodoCache) DeleteTodoList(ctx context.Context) error {
	if c.client == nil {
		return fmt.Errorf("redis client is nil")
	}
	return c.client.Del(ctx, todoListCacheKey)
}
