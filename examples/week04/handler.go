// package week04 定义 Week04 的分层示例代码。
package week04

// import 分组引入 handler 层实现所需标准库。
import (
	// encoding/json 用于将响应结构编码为 JSON。
	"encoding/json"
	// errors 用于判断 service 返回错误的具体类型。
	"errors"
	// net/http 提供 HTTP handler、状态码和响应写入能力。
	"net/http"
	// strings 用于解析路径中的 id 片段。
	"strings"
)

// const 分组声明路由解析时使用的固定路径片段。
const (
	// week04TodoPathPrefix 表示 Todo 资源路径前缀。
	week04TodoPathPrefix = "/api/v1/todos/"
	// week04TodoPathSuffix 表示“完成 Todo”动作路径后缀。
	week04TodoPathSuffix = "/done"
)

// TodoHandler 负责 HTTP 协议层编排，不直接实现业务逻辑。
type TodoHandler struct {
	// service 是 handler 调用的业务层依赖。
	service *TodoService
}

// NewTodoHandler 创建并返回 TodoHandler 实例。
func NewTodoHandler(service *TodoService) *TodoHandler {
	// 将 service 注入 handler，完成层间依赖连接。
	return &TodoHandler{service: service}
}

// NewMux 构建并返回 week04 示例用的路由复用入口。
func NewMux(handler *TodoHandler) *http.ServeMux {
	// 创建标准库提供的多路复用器。
	mux := http.NewServeMux()
	// 把 Todo 路径前缀挂载到 TodoHandler。
	mux.Handle("/api/v1/todos/", handler)
	// 返回配置完成的 mux 供 main 与测试共用。
	return mux
}

// ServeHTTP 处理 PATCH /api/v1/todos/{id}/done 请求。
func (h *TodoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 只允许 PATCH 方法，其他方法统一返回 405。
	if r.Method != http.MethodPatch {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	// 解析路径并提取 todo id。
	id, err := parseMarkDonePath(r.URL.Path)
	// 路径不合法时返回 400。
	if err != nil {
		writeError(w, http.StatusBadRequest, ErrInvalidTodoID.Error())
		return
	}

	// 调用 service 执行业务逻辑。
	todo, err := h.service.MarkDone(r.Context(), id)
	// 对 service 返回错误进行 HTTP 状态码映射。
	if err != nil {
		// switch true 模式用于按条件顺序匹配错误类型。
		switch {
		// 非法参数映射为 400。
		case errors.Is(err, ErrInvalidTodoID):
			writeError(w, http.StatusBadRequest, ErrInvalidTodoID.Error())
		// 资源不存在映射为 404。
		case errors.Is(err, ErrTodoNotFound):
			writeError(w, http.StatusNotFound, ErrTodoNotFound.Error())
		// 其他未分类错误映射为 500。
		default:
			writeError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	// 业务成功时返回 200 与统一响应体结构。
	writeSuccess(w, http.StatusOK, todo)
}

// parseMarkDonePath 从完整路径中提取 id，并做结构合法性校验。
func parseMarkDonePath(path string) (string, error) {
	// 若路径前缀不正确，判定为非法 id 请求。
	if !strings.HasPrefix(path, week04TodoPathPrefix) {
		return "", ErrInvalidTodoID
	}
	// 若路径后缀不正确，判定为非法 id 请求。
	if !strings.HasSuffix(path, week04TodoPathSuffix) {
		return "", ErrInvalidTodoID
	}

	// 先移除前缀，保留中间 id + 后缀片段。
	id := strings.TrimPrefix(path, week04TodoPathPrefix)
	// 再移除后缀，得到理论上的 id 片段。
	id = strings.TrimSuffix(id, week04TodoPathSuffix)
	// 去掉两端斜杠，防止 `/1/` 这类形式影响判断。
	id = strings.Trim(id, "/")
	// 空 id 或包含额外路径层级都视为非法。
	if id == "" || strings.Contains(id, "/") {
		return "", ErrInvalidTodoID
	}
	// 返回提取后的 id 与 nil 错误。
	return id, nil
}

// writeSuccess 写入统一成功响应结构。
func writeSuccess(w http.ResponseWriter, status int, data any) {
	// 成功响应中 data 存放业务结果，error 固定为 nil。
	writeJSON(w, status, map[string]any{
		"data":  data,
		"error": nil,
	})
}

// writeError 写入统一失败响应结构。
func writeError(w http.ResponseWriter, status int, message string) {
	// 失败响应中 data 为 nil，error 存放错误描述。
	writeJSON(w, status, map[string]any{
		"data":  nil,
		"error": message,
	})
}

// writeJSON 负责设置响应头、状态码并编码 JSON。
func writeJSON(w http.ResponseWriter, status int, payload map[string]any) {
	// 告知客户端响应体格式为 JSON。
	w.Header().Set("Content-Type", "application/json")
	// 先写入 HTTP 状态码。
	w.WriteHeader(status)
	// 编码并写入响应体；编码失败时降级返回 500 文本错误。
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "json encode failed", http.StatusInternalServerError)
	}
}
