// 详细注释: package week05
package week05

// 详细注释: import (
import (
	// 详细注释: "context"
	"context"
	// 详细注释: "errors"
	"errors"
	// 详细注释: "testing"
	"testing"
	// 详细注释: )
)

// TestWithTxCommit 验证事务成功时会提交。
// 详细注释: func TestWithTxCommit(t *testing.T) {
func TestWithTxCommit(t *testing.T) {
	// 详细注释: db := NewInMemoryPostgres()
	db := NewInMemoryPostgres()

	// 详细注释: err := db.WithTx(context.Background(), func(tx *Tx) error {
	err := db.WithTx(context.Background(), func(tx *Tx) error {
		// 详细注释: _, createErr := tx.CreateTodo("learn tx")
		_, createErr := tx.CreateTodo("learn tx")
		// 详细注释: return createErr
		return createErr
		// 详细注释: })
	})
	// 详细注释: if err != nil {
	if err != nil {
		// 详细注释: t.Fatalf("unexpected err: %v", err)
		t.Fatalf("unexpected err: %v", err)
		// 详细注释: }
	}

	// 详细注释: items, err := db.ListTodos(1, 10)
	items, err := db.ListTodos(1, 10)
	// 详细注释: if err != nil {
	if err != nil {
		// 详细注释: t.Fatalf("list err: %v", err)
		t.Fatalf("list err: %v", err)
		// 详细注释: }
	}
	// 详细注释: if len(items) != 1 {
	if len(items) != 1 {
		// 详细注释: t.Fatalf("expected 1 item, got %d", len(items))
		t.Fatalf("expected 1 item, got %d", len(items))
		// 详细注释: }
	}
	// 详细注释: }
}

// TestWithTxRollback 验证事务失败时会回滚。
// 详细注释: func TestWithTxRollback(t *testing.T) {
func TestWithTxRollback(t *testing.T) {
	// 详细注释: db := NewInMemoryPostgres()
	db := NewInMemoryPostgres()

	// 详细注释: err := db.WithTx(context.Background(), func(tx *Tx) error {
	err := db.WithTx(context.Background(), func(tx *Tx) error {
		// 详细注释: _, _ = tx.CreateTodo("will rollback")
		_, _ = tx.CreateTodo("will rollback")
		// 详细注释: return errors.New("boom")
		return errors.New("boom")
		// 详细注释: })
	})
	// 详细注释: if err == nil {
	if err == nil {
		// 详细注释: t.Fatalf("expected rollback error")
		t.Fatalf("expected rollback error")
		// 详细注释: }
	}

	// 详细注释: items, listErr := db.ListTodos(1, 10)
	items, listErr := db.ListTodos(1, 10)
	// 详细注释: if listErr != nil {
	if listErr != nil {
		// 详细注释: t.Fatalf("list err: %v", listErr)
		t.Fatalf("list err: %v", listErr)
		// 详细注释: }
	}
	// 详细注释: if len(items) != 0 {
	if len(items) != 0 {
		// 详细注释: t.Fatalf("expected empty after rollback, got %d", len(items))
		t.Fatalf("expected empty after rollback, got %d", len(items))
		// 详细注释: }
	}
	// 详细注释: }
}

// TestListTodosInvalidPage 验证非法分页参数。
// 详细注释: func TestListTodosInvalidPage(t *testing.T) {
func TestListTodosInvalidPage(t *testing.T) {
	// 详细注释: db := NewInMemoryPostgres()
	db := NewInMemoryPostgres()
	// 详细注释: _, err := db.ListTodos(0, 10)
	_, err := db.ListTodos(0, 10)
	// 详细注释: if !errors.Is(err, ErrInvalidPage) {
	if !errors.Is(err, ErrInvalidPage) {
		// 详细注释: t.Fatalf("expected ErrInvalidPage, got %v", err)
		t.Fatalf("expected ErrInvalidPage, got %v", err)
		// 详细注释: }
	}
	// 详细注释: }
}

// TestMarkDoneNotFound 验证资源不存在分支。
// 详细注释: func TestMarkDoneNotFound(t *testing.T) {
func TestMarkDoneNotFound(t *testing.T) {
	// 详细注释: db := NewInMemoryPostgres()
	db := NewInMemoryPostgres()
	// 详细注释: err := db.WithTx(context.Background(), func(tx *Tx) error {
	err := db.WithTx(context.Background(), func(tx *Tx) error {
		// 详细注释: _, markErr := tx.MarkDone("999")
		_, markErr := tx.MarkDone("999")
		// 详细注释: return markErr
		return markErr
		// 详细注释: })
	})
	// 详细注释: if !errors.Is(err, ErrTodoNotFound) {
	if !errors.Is(err, ErrTodoNotFound) {
		// 详细注释: t.Fatalf("expected ErrTodoNotFound, got %v", err)
		t.Fatalf("expected ErrTodoNotFound, got %v", err)
		// 详细注释: }
	}
	// 详细注释: }
}
