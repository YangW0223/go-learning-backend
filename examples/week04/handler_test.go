// package week04 定义 Week04 的分层示例代码。
package week04

// import 分组引入 handler 测试所需标准库。
import (
	// context 用于满足 TodoStore 接口方法签名。
	"context"
	// encoding/json 用于解析 handler 返回的 JSON 响应体。
	"encoding/json"
	// errors 用于构造模拟错误。
	"errors"
	// net/http 提供 HTTP 方法与状态码常量。
	"net/http"
	// net/http/httptest 用于构造请求与记录响应。
	"net/http/httptest"
	// testing 提供测试框架。
	"testing"
)

// failingTodoStore 是始终返回错误的 store 替身。
type failingTodoStore struct {
	// err 表示替身每次调用都会返回的错误。
	err error
}

// MarkDone 实现 TodoStore 接口，并固定返回预设错误。
func (s *failingTodoStore) MarkDone(_ context.Context, _ string) (Todo, error) {
	// 返回空 Todo 和预设错误，用于触发 500 分支测试。
	return Todo{}, s.err
}

// buildWeek04Mux 封装测试中重复的依赖组装逻辑。
func buildWeek04Mux(store TodoStore) *http.ServeMux {
	// 先构建 service 层。
	service := NewTodoService(store)
	// 再构建 handler 层。
	handler := NewTodoHandler(service)
	// 最后返回包含路由的 mux。
	return NewMux(handler)
}

// TestTodoHandlerMarkDone_Success 验证正常请求返回 200 且 done=true。
func TestTodoHandlerMarkDone_Success(t *testing.T) {
	// 准备包含一条可命中 todo 的路由。
	mux := buildWeek04Mux(NewInMemoryTodoStore([]Todo{
		// id=1 的任务初始为未完成。
		{ID: "1", Title: "layered", Done: false},
	}))

	// 构造 PATCH 成功请求。
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/todos/1/done", nil)
	// 创建响应记录器。
	rec := httptest.NewRecorder()
	// 执行请求。
	mux.ServeHTTP(rec, req)

	// 断言状态码为 200。
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	// 用 map 解析通用响应结构。
	var got map[string]any
	// 反序列化响应体 JSON。
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal response failed: %v", err)
	}
	// 读取 data 字段并断言其为对象。
	data, ok := got["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %#v", got["data"])
	}
	// 读取 done 字段并断言其为 true。
	done, ok := data["done"].(bool)
	if !ok || !done {
		t.Fatalf("expected done=true, got %#v", data["done"])
	}
}

// TestTodoHandlerMarkDone_InvalidID 验证非法路径参数返回 400。
func TestTodoHandlerMarkDone_InvalidID(t *testing.T) {
	// 准备空 store，重点验证路径参数校验。
	mux := buildWeek04Mux(NewInMemoryTodoStore(nil))

	// 构造 id 非数字的请求。
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/todos/abc/done", nil)
	// 创建响应记录器。
	rec := httptest.NewRecorder()
	// 执行请求。
	mux.ServeHTTP(rec, req)

	// 断言状态码为 400。
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
}

// TestTodoHandlerMarkDone_NotFound 验证资源不存在返回 404。
func TestTodoHandlerMarkDone_NotFound(t *testing.T) {
	// 准备空 store，确保 id=1 查不到。
	mux := buildWeek04Mux(NewInMemoryTodoStore(nil))

	// 构造合法但不存在的 id 请求。
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/todos/1/done", nil)
	// 创建响应记录器。
	rec := httptest.NewRecorder()
	// 执行请求。
	mux.ServeHTTP(rec, req)

	// 断言状态码为 404。
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

// TestTodoHandlerMarkDone_InternalError 验证未知错误映射为 500。
func TestTodoHandlerMarkDone_InternalError(t *testing.T) {
	// 准备一个固定返回未知错误的 store 替身。
	mux := buildWeek04Mux(&failingTodoStore{err: errors.New("db down")})

	// 构造合法请求触发 service/store 调用链。
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/todos/1/done", nil)
	// 创建响应记录器。
	rec := httptest.NewRecorder()
	// 执行请求。
	mux.ServeHTTP(rec, req)

	// 断言状态码为 500。
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}

// TestTodoHandlerMarkDone_MethodNotAllowed 验证非 PATCH 请求返回 405。
func TestTodoHandlerMarkDone_MethodNotAllowed(t *testing.T) {
	// 准备基础路由。
	mux := buildWeek04Mux(NewInMemoryTodoStore(nil))

	// 使用 GET 方法请求同一路径。
	req := httptest.NewRequest(http.MethodGet, "/api/v1/todos/1/done", nil)
	// 创建响应记录器。
	rec := httptest.NewRecorder()
	// 执行请求。
	mux.ServeHTTP(rec, req)

	// 断言状态码为 405。
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rec.Code)
	}
}
