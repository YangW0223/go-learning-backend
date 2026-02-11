package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/yang/go-learning-backend/internal/cache"
	"github.com/yang/go-learning-backend/internal/service"
	"github.com/yang/go-learning-backend/internal/store"
)

var todoIDPattern = regexp.MustCompile(`^\d{14}\.\d{9}$`)

// TodoHandler 是 HTTP 层的待办事项处理器。
// 它只负责:
// 1) 解析/校验请求
// 2) 调用 store 执行业务动作
// 3) 把结果映射成 HTTP 状态码和 JSON 响应
type TodoHandler struct {
	service service.TodoService
}

// NewTodoHandler 用依赖注入方式创建 handler，便于测试替换 store。
func NewTodoHandler(todoStore store.TodoStore) *TodoHandler {
	todoService := service.NewTodoService(todoStore, cache.NewNoopTodoCache(), 0)
	return &TodoHandler{service: todoService}
}

// NewTodoHandlerWithService 允许上层注入完整业务服务（例如带 Redis 缓存）。
func NewTodoHandlerWithService(todoService service.TodoService) *TodoHandler {
	if todoService == nil {
		panic("todo service is nil")
	}
	return &TodoHandler{service: todoService}
}

// createTodoRequest 对应 POST /api/v1/todos 的 JSON 请求体。
// 期望格式: {"title":"some text"}
type createTodoRequest struct {
	Title string `json:"title"`
}

// Create 处理创建 todo 请求。
// 流程:
// 1) 解析 JSON
// 2) 清洗并校验 title
// 3) 调用 store 创建
// 4) 返回 201 + 创建结果
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTodoRequest
	// JSON 格式错误直接返回 400，避免进入后续业务逻辑。
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}

	// 去掉用户输入前后空格，避免 "   " 被当成有效标题。
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}

	// 调用存储层执行创建。这里把存储层错误统一映射为 500。
	todo, err := h.service.Create(r.Context(), req.Title)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create todo"})
		return
	}

	// 创建成功返回 201（Created）。
	writeJSON(w, http.StatusCreated, todo)
}

// List 返回所有 todo。
// 成功: 200 + JSON 数组
// 失败: 500
func (h *TodoHandler) List(w http.ResponseWriter, r *http.Request) {
	todos, err := h.service.List(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to list todos"})
		return
	}

	writeJSON(w, http.StatusOK, todos)
}

// MarkDone 按 id 把某个 todo 标记为完成。
// 路由层已经负责从 URL 里提取 id，这里只关注业务执行与错误映射。
func (h *TodoHandler) MarkDone(w http.ResponseWriter, r *http.Request, id string) {
	todo, err := h.service.MarkDone(r.Context(), id)
	if err != nil {
		// 业务可识别错误 -> 404，更符合 REST 语义。
		if errors.Is(err, store.ErrTodoNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "todo not found"})
			return
		}

		// 其他未知错误统一按 500 返回。
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update todo"})
		return
	}

	// 更新成功返回 200 + 更新后的 todo。
	writeJSON(w, http.StatusOK, todo)
}

// Delete 按 id 删除某个 todo。
// 返回规则：
// 1) id 非法 -> 400
// 2) todo 不存在 -> 404
// 3) 其他存储错误 -> 500
// 4) 删除成功 -> 200
func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	id = strings.TrimSpace(id)
	if !isValidTodoID(id) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid todo id"})
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		if errors.Is(err, store.ErrTodoNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "todo not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to delete todo"})
		return
	}

	// Week 02 采用统一响应结构: data + error。
	writeJSON(w, http.StatusOK, map[string]any{
		"data": map[string]any{
			"id":      id,
			"deleted": true,
		},
		"error": nil,
	})
}

// isValidTodoID 校验 todo id 格式是否符合当前内存存储实现。
// 当前规则与 generateID 输出一致: YYYYMMDDhhmmss.nanoseconds。
func isValidTodoID(id string) bool {
	return todoIDPattern.MatchString(id)
}
