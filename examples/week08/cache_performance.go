package week08

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

// Todo 是缓存示例中的最小实体。
type Todo struct {
	ID    string
	Title string
}

// Source 定义缓存回源接口。
type Source interface {
	List(ctx context.Context, userID string, page, size int) ([]Todo, error)
}

// SlowSource 模拟慢速 DB 查询。
type SlowSource struct {
	Delay time.Duration
	Data  map[string][]Todo
}

// List 在固定延迟后返回数据。
func (s *SlowSource) List(ctx context.Context, userID string, page, size int) ([]Todo, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(s.Delay):
	}

	items := s.Data[userID]
	start := (page - 1) * size
	if start >= len(items) {
		return []Todo{}, nil
	}
	end := start + size
	if end > len(items) {
		end = len(items)
	}
	out := make([]Todo, end-start)
	copy(out, items[start:end])
	return out, nil
}

// TTLCache 是最小内存缓存实现。
type TTLCache struct {
	mu    sync.Mutex
	nowFn func() time.Time
	rows  map[string]cacheEntry
}

type cacheEntry struct {
	items     []Todo
	expiresAt time.Time
}

// NewTTLCache 创建缓存。
func NewTTLCache() *TTLCache {
	return &TTLCache{
		nowFn: time.Now,
		rows:  make(map[string]cacheEntry),
	}
}

// Get 读取缓存。
func (c *TTLCache) Get(key string) ([]Todo, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.rows[key]
	if !ok {
		return nil, false
	}
	if c.nowFn().After(entry.expiresAt) {
		delete(c.rows, key)
		return nil, false
	}
	items := make([]Todo, len(entry.items))
	copy(items, entry.items)
	return items, true
}

// Set 写入缓存。
func (c *TTLCache) Set(key string, items []Todo, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	copyItems := make([]Todo, len(items))
	copy(copyItems, items)
	c.rows[key] = cacheEntry{items: copyItems, expiresAt: c.nowFn().Add(ttl)}
}

// InvalidatePrefix 失效某个前缀下的所有键。
func (c *TTLCache) InvalidatePrefix(prefix string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k := range c.rows {
		if strings.HasPrefix(k, prefix) {
			delete(c.rows, k)
		}
	}
}

// CachedTodoService 实现 cache aside 模式。
type CachedTodoService struct {
	source Source
	cache  *TTLCache
	ttl    time.Duration

	mu     sync.Mutex
	hits   int
	misses int
}

// NewCachedTodoService 创建服务。
func NewCachedTodoService(source Source, cache *TTLCache, ttl time.Duration) *CachedTodoService {
	return &CachedTodoService{source: source, cache: cache, ttl: ttl}
}

// List 先查缓存，miss 再回源并回填。
func (s *CachedTodoService) List(ctx context.Context, userID string, page, size int) ([]Todo, error) {
	key := fmt.Sprintf("todo:list:%s:%d:%d", userID, page, size)
	if cached, ok := s.cache.Get(key); ok {
		s.mu.Lock()
		s.hits++
		s.mu.Unlock()
		return cached, nil
	}

	s.mu.Lock()
	s.misses++
	s.mu.Unlock()

	items, err := s.source.List(ctx, userID, page, size)
	if err != nil {
		return nil, err
	}
	s.cache.Set(key, items, s.ttl)
	return items, nil
}

// InvalidateUser 在写操作后按用户维度失效缓存。
func (s *CachedTodoService) InvalidateUser(userID string) {
	s.cache.InvalidatePrefix("todo:list:" + userID + ":")
}

// Stats 返回命中与未命中统计。
func (s *CachedTodoService) Stats() (hits int, misses int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.hits, s.misses
}

// MeasureLatency 统计一次函数执行耗时。
func MeasureLatency(fn func() error) (time.Duration, error) {
	start := time.Now()
	err := fn()
	return time.Since(start), err
}
