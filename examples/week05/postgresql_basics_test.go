package week05

import (
	"context"
	"errors"
	"testing"
)

// TestWithTxCommit 验证事务成功时会提交。
func TestWithTxCommit(t *testing.T) {
	db := NewInMemoryPostgres()

	err := db.WithTx(context.Background(), func(tx *Tx) error {
		_, createErr := tx.CreateTodo("learn tx")
		return createErr
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	items, err := db.ListTodos(1, 10)
	if err != nil {
		t.Fatalf("list err: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
}

// TestWithTxRollback 验证事务失败时会回滚。
func TestWithTxRollback(t *testing.T) {
	db := NewInMemoryPostgres()

	err := db.WithTx(context.Background(), func(tx *Tx) error {
		_, _ = tx.CreateTodo("will rollback")
		return errors.New("boom")
	})
	if err == nil {
		t.Fatalf("expected rollback error")
	}

	items, listErr := db.ListTodos(1, 10)
	if listErr != nil {
		t.Fatalf("list err: %v", listErr)
	}
	if len(items) != 0 {
		t.Fatalf("expected empty after rollback, got %d", len(items))
	}
}

// TestListTodosInvalidPage 验证非法分页参数。
func TestListTodosInvalidPage(t *testing.T) {
	db := NewInMemoryPostgres()
	_, err := db.ListTodos(0, 10)
	if !errors.Is(err, ErrInvalidPage) {
		t.Fatalf("expected ErrInvalidPage, got %v", err)
	}
}

// TestMarkDoneNotFound 验证资源不存在分支。
func TestMarkDoneNotFound(t *testing.T) {
	db := NewInMemoryPostgres()
	err := db.WithTx(context.Background(), func(tx *Tx) error {
		_, markErr := tx.MarkDone("999")
		return markErr
	})
	if !errors.Is(err, ErrTodoNotFound) {
		t.Fatalf("expected ErrTodoNotFound, got %v", err)
	}
}
