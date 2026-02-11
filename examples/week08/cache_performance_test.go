// 详细注释: package week08
package week08

// 详细注释: import (
import (
	// 详细注释: "context"
	"context"
	// 详细注释: "testing"
	"testing"
	// 详细注释: "time"
	"time"
	// 详细注释: )
)

// TestCacheHitMiss 验证首次 miss、再次 hit。
// 详细注释: func TestCacheHitMiss(t *testing.T) {
func TestCacheHitMiss(t *testing.T) {
	// 详细注释: source := &SlowSource{
	source := &SlowSource{
		// 详细注释: Delay: 1 * time.Millisecond,
		Delay: 1 * time.Millisecond,
		// 详细注释: Data: map[string][]Todo{
		Data: map[string][]Todo{
			// 详细注释: "u1": {{ID: "1", Title: "a"}},
			"u1": {{ID: "1", Title: "a"}},
			// 详细注释: },
		},
		// 详细注释: }
	}
	// 详细注释: service := NewCachedTodoService(source, NewTTLCache(), 2*time.Second)
	service := NewCachedTodoService(source, NewTTLCache(), 2*time.Second)

	// 详细注释: _, err := service.List(context.Background(), "u1", 1, 10)
	_, err := service.List(context.Background(), "u1", 1, 10)
	// 详细注释: if err != nil {
	if err != nil {
		// 详细注释: t.Fatalf("first list err: %v", err)
		t.Fatalf("first list err: %v", err)
		// 详细注释: }
	}
	// 详细注释: _, err = service.List(context.Background(), "u1", 1, 10)
	_, err = service.List(context.Background(), "u1", 1, 10)
	// 详细注释: if err != nil {
	if err != nil {
		// 详细注释: t.Fatalf("second list err: %v", err)
		t.Fatalf("second list err: %v", err)
		// 详细注释: }
	}

	// 详细注释: hit, miss := service.Stats()
	hit, miss := service.Stats()
	// 详细注释: if hit != 1 || miss != 1 {
	if hit != 1 || miss != 1 {
		// 详细注释: t.Fatalf("unexpected stats hit=%d miss=%d", hit, miss)
		t.Fatalf("unexpected stats hit=%d miss=%d", hit, miss)
		// 详细注释: }
	}
	// 详细注释: }
}

// TestInvalidateUser 验证失效后会再次 miss。
// 详细注释: func TestInvalidateUser(t *testing.T) {
func TestInvalidateUser(t *testing.T) {
	// 详细注释: source := &SlowSource{
	source := &SlowSource{
		// 详细注释: Delay: 1 * time.Millisecond,
		Delay: 1 * time.Millisecond,
		// 详细注释: Data: map[string][]Todo{
		Data: map[string][]Todo{
			// 详细注释: "u1": {{ID: "1", Title: "a"}},
			"u1": {{ID: "1", Title: "a"}},
			// 详细注释: },
		},
		// 详细注释: }
	}
	// 详细注释: service := NewCachedTodoService(source, NewTTLCache(), 10*time.Second)
	service := NewCachedTodoService(source, NewTTLCache(), 10*time.Second)

	// 详细注释: _, _ = service.List(context.Background(), "u1", 1, 10)
	_, _ = service.List(context.Background(), "u1", 1, 10)
	// 详细注释: _, _ = service.List(context.Background(), "u1", 1, 10)
	_, _ = service.List(context.Background(), "u1", 1, 10)
	// 详细注释: service.InvalidateUser("u1")
	service.InvalidateUser("u1")
	// 详细注释: _, _ = service.List(context.Background(), "u1", 1, 10)
	_, _ = service.List(context.Background(), "u1", 1, 10)

	// 详细注释: hit, miss := service.Stats()
	hit, miss := service.Stats()
	// 详细注释: if hit != 1 || miss != 2 {
	if hit != 1 || miss != 2 {
		// 详细注释: t.Fatalf("unexpected stats after invalidate hit=%d miss=%d", hit, miss)
		t.Fatalf("unexpected stats after invalidate hit=%d miss=%d", hit, miss)
		// 详细注释: }
	}
	// 详细注释: }
}
