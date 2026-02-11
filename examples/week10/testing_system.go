package week10

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	// ErrInvalidTitle 表示标题不合法。
	ErrInvalidTitle = errors.New("invalid title")
	// ErrTodoNotFound 表示资源不存在。
	ErrTodoNotFound = errors.New("todo not found")
)

// Todo 表示最小业务实体。
type Todo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Service 表示业务服务。
type Service struct {
	mu     sync.Mutex
	nextID int
	rows   map[string]Todo
}

// NewService 创建服务。
func NewService() *Service {
	return &Service{nextID: 1, rows: make(map[string]Todo)}
}

// CreateTodo 创建 todo。
func (s *Service) CreateTodo(title string) (Todo, error) {
	title = strings.TrimSpace(title)
	if len(title) < 3 {
		return Todo{}, ErrInvalidTitle
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	id := strconv.Itoa(s.nextID)
	s.nextID++
	todo := Todo{ID: id, Title: title}
	s.rows[id] = todo
	return todo, nil
}

// GetTodo 获取 todo。
func (s *Service) GetTodo(id string) (Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	todo, ok := s.rows[id]
	if !ok {
		return Todo{}, ErrTodoNotFound
	}
	return todo, nil
}

// NewMux 创建 HTTP 路由。
func NewMux(svc *Service) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/todos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Title string `json:"title"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		todo, err := svc.CreateTodo(req.Title)
		if err != nil {
			if errors.Is(err, ErrInvalidTitle) {
				http.Error(w, ErrInvalidTitle.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusCreated, todo)
	})

	mux.HandleFunc("/api/v1/todos/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/api/v1/todos/")
		if strings.TrimSpace(id) == "" || strings.Contains(id, "/") {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		todo, err := svc.GetTodo(id)
		if err != nil {
			if errors.Is(err, ErrTodoNotFound) {
				http.Error(w, ErrTodoNotFound.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, todo)
	})
	return mux
}

func writeJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]any{"data": data, "error": nil})
}

// DebugString 返回可读字符串，方便 cmd 展示。
func DebugString(todo Todo) string {
	return fmt.Sprintf("todo{id=%s,title=%s}", todo.ID, todo.Title)
}
