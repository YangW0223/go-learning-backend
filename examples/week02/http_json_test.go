package week02

import (
	"encoding/json"
	"errors"
	"testing"
)

// TestParseDeleteTodoPath_Valid 验证合法路径可正确提取 id。
func TestParseDeleteTodoPath_Valid(t *testing.T) {
	id, err := ParseDeleteTodoPath("/api/v1/todos/20260210112233.123456789")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "20260210112233.123456789" {
		t.Fatalf("unexpected id: %s", id)
	}
}

// TestParseDeleteTodoPath_InvalidID 验证非法 id 会返回 ErrInvalidTodoID。
func TestParseDeleteTodoPath_InvalidID(t *testing.T) {
	_, err := ParseDeleteTodoPath("/api/v1/todos/abc")
	if !errors.Is(err, ErrInvalidTodoID) {
		t.Fatalf("expected ErrInvalidTodoID, got %v", err)
	}
}

// TestParseDeleteTodoPath_InvalidPrefix 验证错误前缀会返回 ErrInvalidDeletePath。
func TestParseDeleteTodoPath_InvalidPrefix(t *testing.T) {
	_, err := ParseDeleteTodoPath("/api/v1/tasks/20260210112233.123456789")
	if !errors.Is(err, ErrInvalidDeletePath) {
		t.Fatalf("expected ErrInvalidDeletePath, got %v", err)
	}
}

// TestBuildSuccessJSON 验证统一成功响应结构。
func TestBuildSuccessJSON(t *testing.T) {
	b, err := BuildSuccessJSON(DeleteResult{
		ID:      "20260210112233.123456789",
		Deleted: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got["error"] != nil {
		t.Fatalf("expected error=nil, got %v", got["error"])
	}
}

// TestBuildErrorJSON 验证统一错误响应结构。
func TestBuildErrorJSON(t *testing.T) {
	b, err := BuildErrorJSON("invalid todo id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got["error"] != "invalid todo id" {
		t.Fatalf("unexpected error message: %v", got["error"])
	}
}
