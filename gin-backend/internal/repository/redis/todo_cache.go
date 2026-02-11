package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yang/go-learning-backend/gin-backend/internal/model"
)

// TodoCache 是 Redis Todo 列表缓存实现。
type TodoCache struct {
	client *Client
}

// NewTodoCache 创建 Redis Todo 缓存实例。
func NewTodoCache(client *Client) *TodoCache {
	return &TodoCache{client: client}
}

// Ping 校验 Redis 连通性。
func (c *TodoCache) Ping(ctx context.Context) error {
	if c.client == nil {
		return fmt.Errorf("redis client is nil")
	}
	return c.client.Ping(ctx)
}

// GetList 读取用户 Todo 列表缓存。
func (c *TodoCache) GetList(ctx context.Context, userID string) ([]model.Todo, bool, error) {
	if c.client == nil {
		return nil, false, fmt.Errorf("redis client is nil")
	}
	payload, hit, err := c.client.Get(ctx, listKey(userID))
	if err != nil {
		return nil, false, err
	}
	if !hit {
		return nil, false, nil
	}
	var todos []model.Todo
	if err := json.Unmarshal([]byte(payload), &todos); err != nil {
		return nil, false, fmt.Errorf("unmarshal todo cache: %w", err)
	}
	return todos, true, nil
}

// SetList 写入用户 Todo 列表缓存。
func (c *TodoCache) SetList(ctx context.Context, userID string, todos []model.Todo, ttl time.Duration) error {
	if c.client == nil {
		return fmt.Errorf("redis client is nil")
	}
	encoded, err := json.Marshal(todos)
	if err != nil {
		return fmt.Errorf("marshal todo cache: %w", err)
	}
	return c.client.SetEX(ctx, listKey(userID), string(encoded), ttl)
}

// DeleteList 删除用户 Todo 列表缓存。
func (c *TodoCache) DeleteList(ctx context.Context, userID string) error {
	if c.client == nil {
		return fmt.Errorf("redis client is nil")
	}
	return c.client.Del(ctx, listKey(userID))
}

// listKey 生成按用户隔离的缓存 key。
func listKey(userID string) string {
	return "todo:list:v1:user:" + userID
}
