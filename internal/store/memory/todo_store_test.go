package memory

import (
	"errors"
	"testing"

	"github.com/yang/go-learning-backend/internal/store"
)

// TestTodoStore_CreateAndList 验证最基本的“创建 + 列表查询”流程。
func TestTodoStore_CreateAndList(t *testing.T) {
	s := NewTodoStore()

	created, err := s.Create("learn go")
	if err != nil {
		t.Fatalf("create should not fail: %v", err)
	}

	if created.Title != "learn go" {
		t.Fatalf("unexpected title: %s", created.Title)
	}

	items, err := s.List()
	if err != nil {
		t.Fatalf("list should not fail: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
}

// TestTodoStore_MarkDone 验证按 ID 更新完成状态的流程。
func TestTodoStore_MarkDone(t *testing.T) {
	s := NewTodoStore()

	created, err := s.Create("learn go")
	if err != nil {
		t.Fatalf("create should not fail: %v", err)
	}

	updated, err := s.MarkDone(created.ID)
	if err != nil {
		t.Fatalf("mark done should not fail: %v", err)
	}

	if !updated.Done {
		t.Fatalf("expected done=true")
	}
}

// TestTodoStore_Delete 验证删除成功后，列表中不再包含该数据。
func TestTodoStore_Delete(t *testing.T) {
	s := NewTodoStore()

	created, err := s.Create("to be deleted")
	if err != nil {
		t.Fatalf("create should not fail: %v", err)
	}

	if err := s.Delete(created.ID); err != nil {
		t.Fatalf("delete should not fail: %v", err)
	}

	items, err := s.List()
	if err != nil {
		t.Fatalf("list should not fail: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected 0 item, got %d", len(items))
	}
}

// TestTodoStore_Delete_NotFound 验证删除不存在数据时返回可判定错误。
func TestTodoStore_Delete_NotFound(t *testing.T) {
	s := NewTodoStore()

	err := s.Delete("20000101000000.000000000")
	if !errors.Is(err, store.ErrTodoNotFound) {
		t.Fatalf("expected ErrTodoNotFound, got %v", err)
	}
}
