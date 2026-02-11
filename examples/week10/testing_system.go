// 详细注释: package week10
package week10

// 详细注释: import (
import (
	// 详细注释: "encoding/json"
	"encoding/json"
	// 详细注释: "errors"
	"errors"
	// 详细注释: "fmt"
	"fmt"
	// 详细注释: "net/http"
	"net/http"
	// 详细注释: "strconv"
	"strconv"
	// 详细注释: "strings"
	"strings"
	// 详细注释: "sync"
	"sync"
	// 详细注释: )
)

// 详细注释: var (
var (
	// ErrInvalidTitle 表示标题不合法。
	// 详细注释: ErrInvalidTitle = errors.New("invalid title")
	ErrInvalidTitle = errors.New("invalid title")
	// ErrTodoNotFound 表示资源不存在。
	// 详细注释: ErrTodoNotFound = errors.New("todo not found")
	ErrTodoNotFound = errors.New("todo not found")

// 详细注释: )
)

// Todo 表示最小业务实体。
// 详细注释: type Todo struct {
type Todo struct {
	// 详细注释: ID    string `json:"id"`
	ID string `json:"id"`
	// 详细注释: Title string `json:"title"`
	Title string `json:"title"`
	// 详细注释: }
}

// Service 表示业务服务。
// 详细注释: type Service struct {
type Service struct {
	// 详细注释: mu     sync.Mutex
	mu sync.Mutex
	// 详细注释: nextID int
	nextID int
	// 详细注释: rows   map[string]Todo
	rows map[string]Todo
	// 详细注释: }
}

// NewService 创建服务。
// 详细注释: func NewService() *Service {
func NewService() *Service {
	// 详细注释: return &Service{nextID: 1, rows: make(map[string]Todo)}
	return &Service{nextID: 1, rows: make(map[string]Todo)}
	// 详细注释: }
}

// CreateTodo 创建 todo。
// 详细注释: func (s *Service) CreateTodo(title string) (Todo, error) {
func (s *Service) CreateTodo(title string) (Todo, error) {
	// 详细注释: title = strings.TrimSpace(title)
	title = strings.TrimSpace(title)
	// 详细注释: if len(title) < 3 {
	if len(title) < 3 {
		// 详细注释: return Todo{}, ErrInvalidTitle
		return Todo{}, ErrInvalidTitle
		// 详细注释: }
	}

	// 详细注释: s.mu.Lock()
	s.mu.Lock()
	// 详细注释: defer s.mu.Unlock()
	defer s.mu.Unlock()
	// 详细注释: id := strconv.Itoa(s.nextID)
	id := strconv.Itoa(s.nextID)
	// 详细注释: s.nextID++
	s.nextID++
	// 详细注释: todo := Todo{ID: id, Title: title}
	todo := Todo{ID: id, Title: title}
	// 详细注释: s.rows[id] = todo
	s.rows[id] = todo
	// 详细注释: return todo, nil
	return todo, nil
	// 详细注释: }
}

// GetTodo 获取 todo。
// 详细注释: func (s *Service) GetTodo(id string) (Todo, error) {
func (s *Service) GetTodo(id string) (Todo, error) {
	// 详细注释: s.mu.Lock()
	s.mu.Lock()
	// 详细注释: defer s.mu.Unlock()
	defer s.mu.Unlock()
	// 详细注释: todo, ok := s.rows[id]
	todo, ok := s.rows[id]
	// 详细注释: if !ok {
	if !ok {
		// 详细注释: return Todo{}, ErrTodoNotFound
		return Todo{}, ErrTodoNotFound
		// 详细注释: }
	}
	// 详细注释: return todo, nil
	return todo, nil
	// 详细注释: }
}

// NewMux 创建 HTTP 路由。
// 详细注释: func NewMux(svc *Service) *http.ServeMux {
func NewMux(svc *Service) *http.ServeMux {
	// 详细注释: mux := http.NewServeMux()
	mux := http.NewServeMux()
	// 详细注释: mux.HandleFunc("/api/v1/todos", func(w http.ResponseWriter, r *http.Request) {
	mux.HandleFunc("/api/v1/todos", func(w http.ResponseWriter, r *http.Request) {
		// 详细注释: if r.Method != http.MethodPost {
		if r.Method != http.MethodPost {
			// 详细注释: http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			// 详细注释: return
			return
			// 详细注释: }
		}

		// 详细注释: var req struct {
		var req struct {
			// 详细注释: Title string `json:"title"`
			Title string `json:"title"`
			// 详细注释: }
		}
		// 详细注释: if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			// 详细注释: http.Error(w, "bad json", http.StatusBadRequest)
			http.Error(w, "bad json", http.StatusBadRequest)
			// 详细注释: return
			return
			// 详细注释: }
		}
		// 详细注释: todo, err := svc.CreateTodo(req.Title)
		todo, err := svc.CreateTodo(req.Title)
		// 详细注释: if err != nil {
		if err != nil {
			// 详细注释: if errors.Is(err, ErrInvalidTitle) {
			if errors.Is(err, ErrInvalidTitle) {
				// 详细注释: http.Error(w, ErrInvalidTitle.Error(), http.StatusBadRequest)
				http.Error(w, ErrInvalidTitle.Error(), http.StatusBadRequest)
				// 详细注释: return
				return
				// 详细注释: }
			}
			// 详细注释: http.Error(w, "internal error", http.StatusInternalServerError)
			http.Error(w, "internal error", http.StatusInternalServerError)
			// 详细注释: return
			return
			// 详细注释: }
		}
		// 详细注释: writeJSON(w, http.StatusCreated, todo)
		writeJSON(w, http.StatusCreated, todo)
		// 详细注释: })
	})

	// 详细注释: mux.HandleFunc("/api/v1/todos/", func(w http.ResponseWriter, r *http.Request) {
	mux.HandleFunc("/api/v1/todos/", func(w http.ResponseWriter, r *http.Request) {
		// 详细注释: if r.Method != http.MethodGet {
		if r.Method != http.MethodGet {
			// 详细注释: http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			// 详细注释: return
			return
			// 详细注释: }
		}
		// 详细注释: id := strings.TrimPrefix(r.URL.Path, "/api/v1/todos/")
		id := strings.TrimPrefix(r.URL.Path, "/api/v1/todos/")
		// 详细注释: if strings.TrimSpace(id) == "" || strings.Contains(id, "/") {
		if strings.TrimSpace(id) == "" || strings.Contains(id, "/") {
			// 详细注释: http.Error(w, "invalid id", http.StatusBadRequest)
			http.Error(w, "invalid id", http.StatusBadRequest)
			// 详细注释: return
			return
			// 详细注释: }
		}
		// 详细注释: todo, err := svc.GetTodo(id)
		todo, err := svc.GetTodo(id)
		// 详细注释: if err != nil {
		if err != nil {
			// 详细注释: if errors.Is(err, ErrTodoNotFound) {
			if errors.Is(err, ErrTodoNotFound) {
				// 详细注释: http.Error(w, ErrTodoNotFound.Error(), http.StatusNotFound)
				http.Error(w, ErrTodoNotFound.Error(), http.StatusNotFound)
				// 详细注释: return
				return
				// 详细注释: }
			}
			// 详细注释: http.Error(w, "internal error", http.StatusInternalServerError)
			http.Error(w, "internal error", http.StatusInternalServerError)
			// 详细注释: return
			return
			// 详细注释: }
		}
		// 详细注释: writeJSON(w, http.StatusOK, todo)
		writeJSON(w, http.StatusOK, todo)
		// 详细注释: })
	})
	// 详细注释: return mux
	return mux
	// 详细注释: }
}

// 详细注释: func writeJSON(w http.ResponseWriter, code int, data any) {
func writeJSON(w http.ResponseWriter, code int, data any) {
	// 详细注释: w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json")
	// 详细注释: w.WriteHeader(code)
	w.WriteHeader(code)
	// 详细注释: _ = json.NewEncoder(w).Encode(map[string]any{"data": data, "error": nil})
	_ = json.NewEncoder(w).Encode(map[string]any{"data": data, "error": nil})
	// 详细注释: }
}

// DebugString 返回可读字符串，方便 cmd 展示。
// 详细注释: func DebugString(todo Todo) string {
func DebugString(todo Todo) string {
	// 详细注释: return fmt.Sprintf("todo{id=%s,title=%s}", todo.ID, todo.Title)
	return fmt.Sprintf("todo{id=%s,title=%s}", todo.ID, todo.Title)
	// 详细注释: }
}
