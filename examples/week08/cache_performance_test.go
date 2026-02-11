package week08

import (
	"context"
	"testing"
	"time"
)

// TestCacheHitMiss 验证首次 miss、再次 hit。
func TestCacheHitMiss(t *testing.T) {
	source := &SlowSource{
		Delay: 1 * time.Millisecond,
		Data: map[string][]Todo{
			"u1": {{ID: "1", Title: "a"}},
		},
	}
	service := NewCachedTodoService(source, NewTTLCache(), 2*time.Second)

	_, err := service.List(context.Background(), "u1", 1, 10)
	if err != nil {
		t.Fatalf("first list err: %v", err)
	}
	_, err = service.List(context.Background(), "u1", 1, 10)
	if err != nil {
		t.Fatalf("second list err: %v", err)
	}

	hit, miss := service.Stats()
	if hit != 1 || miss != 1 {
		t.Fatalf("unexpected stats hit=%d miss=%d", hit, miss)
	}
}

// TestInvalidateUser 验证失效后会再次 miss。
func TestInvalidateUser(t *testing.T) {
	source := &SlowSource{
		Delay: 1 * time.Millisecond,
		Data: map[string][]Todo{
			"u1": {{ID: "1", Title: "a"}},
		},
	}
	service := NewCachedTodoService(source, NewTTLCache(), 10*time.Second)

	_, _ = service.List(context.Background(), "u1", 1, 10)
	_, _ = service.List(context.Background(), "u1", 1, 10)
	service.InvalidateUser("u1")
	_, _ = service.List(context.Background(), "u1", 1, 10)

	hit, miss := service.Stats()
	if hit != 1 || miss != 2 {
		t.Fatalf("unexpected stats after invalidate hit=%d miss=%d", hit, miss)
	}
}
