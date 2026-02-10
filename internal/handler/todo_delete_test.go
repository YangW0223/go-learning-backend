package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yang/go-learning-backend/internal/model"
)

// fakeTodoStoreForDeleteError 仅用于触发 Delete 的 500 分支测试。
type fakeTodoStoreForDeleteError struct {
	deleteErr error
}

func (f *fakeTodoStoreForDeleteError) Create(title string) (model.Todo, error) {
	return model.Todo{}, nil
}

func (f *fakeTodoStoreForDeleteError) List() ([]model.Todo, error) {
	return nil, nil
}

func (f *fakeTodoStoreForDeleteError) MarkDone(id string) (model.Todo, error) {
	return model.Todo{}, nil
}

func (f *fakeTodoStoreForDeleteError) Delete(id string) error {
	return f.deleteErr
}

// TestTodoHandler_Delete_InternalError 验证未知存储错误会映射为 500。
func TestTodoHandler_Delete_InternalError(t *testing.T) {
	h := NewTodoHandler(&fakeTodoStoreForDeleteError{
		deleteErr: errors.New("db temporary unavailable"),
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/todos/20260210121212.123456789", nil)

	h.Delete(rec, req, "20260210121212.123456789")

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}

	var got map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got["error"] != "failed to delete todo" {
		t.Fatalf("unexpected error message: %v", got["error"])
	}
}
