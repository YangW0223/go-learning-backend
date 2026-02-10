package app

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yang/go-learning-backend/internal/handler"
	"github.com/yang/go-learning-backend/internal/store/memory"
)

// TestNewRouter_Ping 验证 /ping 路由是否已注册且返回符合预期。
// 这是一个端到端风格的路由单测：
// 请求 -> 路由 -> handler -> 响应。
func TestNewRouter_Ping(t *testing.T) {
	// 先构造路由依赖，再创建完整 router。
	todoStore := memory.NewTodoStore()
	todoHandler := handler.NewTodoHandler(todoStore)
	router := NewRouter(todoHandler)

	// 构造一个 GET /ping 请求，不需要请求体。
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	// Recorder 用来接收 handler 写入的响应。
	rec := httptest.NewRecorder()

	// 真正执行本次 HTTP 请求。
	router.ServeHTTP(rec, req)

	// 断言状态码。
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// 断言返回体 JSON 内容。
	var got map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got["message"] != "pong" {
		t.Fatalf(`expected message "pong", got %q`, got["message"])
	}
}
