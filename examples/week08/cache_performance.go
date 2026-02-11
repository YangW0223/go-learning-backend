// 详细注释: package week08
package week08

// 详细注释: import (
import (
	// 详细注释: "context"
	"context"
	// 详细注释: "fmt"
	"fmt"
	// 详细注释: "strings"
	"strings"
	// 详细注释: "sync"
	"sync"
	// 详细注释: "time"
	"time"
	// 详细注释: )
)

// Todo 是缓存示例中的最小实体。
// 详细注释: type Todo struct {
type Todo struct {
	// 详细注释: ID    string
	ID string
	// 详细注释: Title string
	Title string
	// 详细注释: }
}

// Source 定义缓存回源接口。
// 详细注释: type Source interface {
type Source interface {
	// 详细注释: List(ctx context.Context, userID string, page, size int) ([]Todo, error)
	List(ctx context.Context, userID string, page, size int) ([]Todo, error)
	// 详细注释: }
}

// SlowSource 模拟慢速 DB 查询。
// 详细注释: type SlowSource struct {
type SlowSource struct {
	// 详细注释: Delay time.Duration
	Delay time.Duration
	// 详细注释: Data  map[string][]Todo
	Data map[string][]Todo
	// 详细注释: }
}

// List 在固定延迟后返回数据。
// 详细注释: func (s *SlowSource) List(ctx context.Context, userID string, page, size int) ([]Todo, error) {
func (s *SlowSource) List(ctx context.Context, userID string, page, size int) ([]Todo, error) {
	// 详细注释: select {
	select {
	// 详细注释: case <-ctx.Done():
	case <-ctx.Done():
		// 详细注释: return nil, ctx.Err()
		return nil, ctx.Err()
		// 详细注释: case <-time.After(s.Delay):
	case <-time.After(s.Delay):
		// 详细注释: }
	}

	// 详细注释: items := s.Data[userID]
	items := s.Data[userID]
	// 详细注释: start := (page - 1) * size
	start := (page - 1) * size
	// 详细注释: if start >= len(items) {
	if start >= len(items) {
		// 详细注释: return []Todo{}, nil
		return []Todo{}, nil
		// 详细注释: }
	}
	// 详细注释: end := start + size
	end := start + size
	// 详细注释: if end > len(items) {
	if end > len(items) {
		// 详细注释: end = len(items)
		end = len(items)
		// 详细注释: }
	}
	// 详细注释: out := make([]Todo, end-start)
	out := make([]Todo, end-start)
	// 详细注释: copy(out, items[start:end])
	copy(out, items[start:end])
	// 详细注释: return out, nil
	return out, nil
	// 详细注释: }
}

// TTLCache 是最小内存缓存实现。
// 详细注释: type TTLCache struct {
type TTLCache struct {
	// 详细注释: mu    sync.Mutex
	mu sync.Mutex
	// 详细注释: nowFn func() time.Time
	nowFn func() time.Time
	// 详细注释: rows  map[string]cacheEntry
	rows map[string]cacheEntry
	// 详细注释: }
}

// 详细注释: type cacheEntry struct {
type cacheEntry struct {
	// 详细注释: items     []Todo
	items []Todo
	// 详细注释: expiresAt time.Time
	expiresAt time.Time
	// 详细注释: }
}

// NewTTLCache 创建缓存。
// 详细注释: func NewTTLCache() *TTLCache {
func NewTTLCache() *TTLCache {
	// 详细注释: return &TTLCache{
	return &TTLCache{
		// 详细注释: nowFn: time.Now,
		nowFn: time.Now,
		// 详细注释: rows:  make(map[string]cacheEntry),
		rows: make(map[string]cacheEntry),
		// 详细注释: }
	}
	// 详细注释: }
}

// Get 读取缓存。
// 详细注释: func (c *TTLCache) Get(key string) ([]Todo, bool) {
func (c *TTLCache) Get(key string) ([]Todo, bool) {
	// 详细注释: c.mu.Lock()
	c.mu.Lock()
	// 详细注释: defer c.mu.Unlock()
	defer c.mu.Unlock()
	// 详细注释: entry, ok := c.rows[key]
	entry, ok := c.rows[key]
	// 详细注释: if !ok {
	if !ok {
		// 详细注释: return nil, false
		return nil, false
		// 详细注释: }
	}
	// 详细注释: if c.nowFn().After(entry.expiresAt) {
	if c.nowFn().After(entry.expiresAt) {
		// 详细注释: delete(c.rows, key)
		delete(c.rows, key)
		// 详细注释: return nil, false
		return nil, false
		// 详细注释: }
	}
	// 详细注释: items := make([]Todo, len(entry.items))
	items := make([]Todo, len(entry.items))
	// 详细注释: copy(items, entry.items)
	copy(items, entry.items)
	// 详细注释: return items, true
	return items, true
	// 详细注释: }
}

// Set 写入缓存。
// 详细注释: func (c *TTLCache) Set(key string, items []Todo, ttl time.Duration) {
func (c *TTLCache) Set(key string, items []Todo, ttl time.Duration) {
	// 详细注释: c.mu.Lock()
	c.mu.Lock()
	// 详细注释: defer c.mu.Unlock()
	defer c.mu.Unlock()
	// 详细注释: copyItems := make([]Todo, len(items))
	copyItems := make([]Todo, len(items))
	// 详细注释: copy(copyItems, items)
	copy(copyItems, items)
	// 详细注释: c.rows[key] = cacheEntry{items: copyItems, expiresAt: c.nowFn().Add(ttl)}
	c.rows[key] = cacheEntry{items: copyItems, expiresAt: c.nowFn().Add(ttl)}
	// 详细注释: }
}

// InvalidatePrefix 失效某个前缀下的所有键。
// 详细注释: func (c *TTLCache) InvalidatePrefix(prefix string) {
func (c *TTLCache) InvalidatePrefix(prefix string) {
	// 详细注释: c.mu.Lock()
	c.mu.Lock()
	// 详细注释: defer c.mu.Unlock()
	defer c.mu.Unlock()
	// 详细注释: for k := range c.rows {
	for k := range c.rows {
		// 详细注释: if strings.HasPrefix(k, prefix) {
		if strings.HasPrefix(k, prefix) {
			// 详细注释: delete(c.rows, k)
			delete(c.rows, k)
			// 详细注释: }
		}
		// 详细注释: }
	}
	// 详细注释: }
}

// CachedTodoService 实现 cache aside 模式。
// 详细注释: type CachedTodoService struct {
type CachedTodoService struct {
	// 详细注释: source Source
	source Source
	// 详细注释: cache  *TTLCache
	cache *TTLCache
	// 详细注释: ttl    time.Duration
	ttl time.Duration

	// 详细注释: mu     sync.Mutex
	mu sync.Mutex
	// 详细注释: hits   int
	hits int
	// 详细注释: misses int
	misses int
	// 详细注释: }
}

// NewCachedTodoService 创建服务。
// 详细注释: func NewCachedTodoService(source Source, cache *TTLCache, ttl time.Duration) *CachedTodoService {
func NewCachedTodoService(source Source, cache *TTLCache, ttl time.Duration) *CachedTodoService {
	// 详细注释: return &CachedTodoService{source: source, cache: cache, ttl: ttl}
	return &CachedTodoService{source: source, cache: cache, ttl: ttl}
	// 详细注释: }
}

// List 先查缓存，miss 再回源并回填。
// 详细注释: func (s *CachedTodoService) List(ctx context.Context, userID string, page, size int) ([]Todo, error) {
func (s *CachedTodoService) List(ctx context.Context, userID string, page, size int) ([]Todo, error) {
	// 详细注释: key := fmt.Sprintf("todo:list:%s:%d:%d", userID, page, size)
	key := fmt.Sprintf("todo:list:%s:%d:%d", userID, page, size)
	// 详细注释: if cached, ok := s.cache.Get(key); ok {
	if cached, ok := s.cache.Get(key); ok {
		// 详细注释: s.mu.Lock()
		s.mu.Lock()
		// 详细注释: s.hits++
		s.hits++
		// 详细注释: s.mu.Unlock()
		s.mu.Unlock()
		// 详细注释: return cached, nil
		return cached, nil
		// 详细注释: }
	}

	// 详细注释: s.mu.Lock()
	s.mu.Lock()
	// 详细注释: s.misses++
	s.misses++
	// 详细注释: s.mu.Unlock()
	s.mu.Unlock()

	// 详细注释: items, err := s.source.List(ctx, userID, page, size)
	items, err := s.source.List(ctx, userID, page, size)
	// 详细注释: if err != nil {
	if err != nil {
		// 详细注释: return nil, err
		return nil, err
		// 详细注释: }
	}
	// 详细注释: s.cache.Set(key, items, s.ttl)
	s.cache.Set(key, items, s.ttl)
	// 详细注释: return items, nil
	return items, nil
	// 详细注释: }
}

// InvalidateUser 在写操作后按用户维度失效缓存。
// 详细注释: func (s *CachedTodoService) InvalidateUser(userID string) {
func (s *CachedTodoService) InvalidateUser(userID string) {
	// 详细注释: s.cache.InvalidatePrefix("todo:list:" + userID + ":")
	s.cache.InvalidatePrefix("todo:list:" + userID + ":")
	// 详细注释: }
}

// Stats 返回命中与未命中统计。
// 详细注释: func (s *CachedTodoService) Stats() (hits int, misses int) {
func (s *CachedTodoService) Stats() (hits int, misses int) {
	// 详细注释: s.mu.Lock()
	s.mu.Lock()
	// 详细注释: defer s.mu.Unlock()
	defer s.mu.Unlock()
	// 详细注释: return s.hits, s.misses
	return s.hits, s.misses
	// 详细注释: }
}

// MeasureLatency 统计一次函数执行耗时。
// 详细注释: func MeasureLatency(fn func() error) (time.Duration, error) {
func MeasureLatency(fn func() error) (time.Duration, error) {
	// 详细注释: start := time.Now()
	start := time.Now()
	// 详细注释: err := fn()
	err := fn()
	// 详细注释: return time.Since(start), err
	return time.Since(start), err
	// 详细注释: }
}
