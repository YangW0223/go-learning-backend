package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yang/go-learning-backend/internal/handler"
	"github.com/yang/go-learning-backend/internal/store/memory"
)

func newTestRouter() http.Handler {
	todoStore := memory.NewTodoStore()
	todoHandler := handler.NewTodoHandler(todoStore)
	return NewRouter(todoHandler)
}

func createTodoForTest(t *testing.T, router http.Handler, title string) string {
	t.Helper()

	body := bytes.NewBufferString(`{"title":"` + title + `"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/todos", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("create todo expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var created struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("decode created todo: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected created todo id, got empty")
	}
	return created.ID
}

// TestNewRouter_Ping 验证 /ping 路由是否已注册且返回符合预期。
// 这是一个端到端风格的路由单测：
// 请求 -> 路由 -> handler -> 响应。
func TestNewRouter_Ping(t *testing.T) {
	router := newTestRouter()

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

// TestNewRouter_DeleteTodo_Success 验证删除成功路径。
func TestNewRouter_DeleteTodo_Success(t *testing.T) {
	router := newTestRouter()
	id := createTodoForTest(t, router, "learn http json")

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/todos/"+id, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got struct {
		Data struct {
			ID      string `json:"id"`
			Deleted bool   `json:"deleted"`
		} `json:"data"`
		Error any `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Data.ID != id || !got.Data.Deleted {
		t.Fatalf("unexpected delete response: %+v", got)
	}

	// 删除后再次查询列表，应为空。
	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/todos", nil)
	listRec := httptest.NewRecorder()
	router.ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("list expected status %d, got %d", http.StatusOK, listRec.Code)
	}
	var items []map[string]any
	if err := json.NewDecoder(listRec.Body).Decode(&items); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected list length 0, got %d", len(items))
	}
}

// TestNewRouter_DeleteTodo_NotFound 验证“id 合法但数据不存在”返回 404。
func TestNewRouter_DeleteTodo_NotFound(t *testing.T) {
	router := newTestRouter()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/todos/20000101000000.000000000", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

// TestNewRouter_DeleteTodo_InvalidID 验证“id 非法”返回 400。
func TestNewRouter_DeleteTodo_InvalidID(t *testing.T) {
	router := newTestRouter()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/todos/abc", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var got map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got["error"] != "invalid todo id" {
		t.Fatalf("unexpected error message: %v", got["error"])
	}
}
