package memory

import (
	"sync"
	"time"

	"github.com/yang/go-learning-backend/internal/model"
	"github.com/yang/go-learning-backend/internal/store"
)

// TodoStore 是基于内存的存储实现，适合本地学习和快速调试。
//
// 说明：
// 1) 数据只存在进程内，服务重启后会丢失。
// 2) 使用读写锁保证并发安全（多个请求同时读写时不会数据竞争）。
type TodoStore struct {
	// mu 保护 todos 切片的并发访问。
	mu sync.RWMutex
	// todos 保存当前进程内的所有待办项。
	todos []model.Todo
}

// NewTodoStore 创建一个空的内存存储。
func NewTodoStore() *TodoStore {
	return &TodoStore{todos: make([]model.Todo, 0)}
}

// Create 新增一个 todo。
// 它会生成 ID、写入创建时间，并默认 Done=false。
func (s *TodoStore) Create(title string) (model.Todo, error) {
	// 写操作需要独占锁，避免并发 append 造成数据竞争。
	s.mu.Lock()
	defer s.mu.Unlock()

	todo := model.Todo{
		ID:        generateID(),
		Title:     title,
		Done:      false,
		CreatedAt: time.Now().UTC(),
	}
	s.todos = append(s.todos, todo)

	return todo, nil
}

// List 返回 todo 快照。
//
// 为什么要 copy：
// 直接返回 s.todos 会把内部切片暴露给调用方，调用方可能误改底层数据。
// copy 后返回的是独立切片，避免外部破坏内部状态。
func (s *TodoStore) List() ([]model.Todo, error) {
	// 读操作使用读锁，允许多个并发读取。
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]model.Todo, len(s.todos))
	copy(result, s.todos)
	return result, nil
}

// MarkDone 按 id 标记 todo 为完成。
// 找到则返回更新后的 todo；找不到返回 store.ErrTodoNotFound。
func (s *TodoStore) MarkDone(id string) (model.Todo, error) {
	// 写操作（更新 Done 字段）需要独占锁。
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.todos {
		if s.todos[i].ID == id {
			s.todos[i].Done = true
			return s.todos[i], nil
		}
	}

	return model.Todo{}, store.ErrTodoNotFound
}

// generateID 生成一个基于 UTC 时间戳的字符串 ID。
//
// 格式: YYYYMMDDhhmmss.nanoseconds
// 在学习项目中足够使用；生产系统更常见做法是 UUID 或雪花算法。
func generateID() string {
	return time.Now().UTC().Format("20060102150405.000000000")
}
