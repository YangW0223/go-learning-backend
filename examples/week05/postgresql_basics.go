package week05

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// ErrInvalidPage 表示分页参数不合法。
	ErrInvalidPage = errors.New("invalid pagination")
	// ErrTodoNotFound 表示目标 todo 不存在。
	ErrTodoNotFound = errors.New("todo not found")
)

// TodoRecord 模拟数据库中的 todos 表记录。
type TodoRecord struct {
	ID        string
	Title     string
	Done      bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// InMemoryPostgres 模拟 Postgres 存储。
// 这个类型用于演示“事务边界 + 分页查询”核心概念。
type InMemoryPostgres struct {
	mu     sync.Mutex
	nextID int
	rows   map[string]TodoRecord
}

// NewInMemoryPostgres 创建示例数据库实例。
func NewInMemoryPostgres() *InMemoryPostgres {
	return &InMemoryPostgres{
		nextID: 1,
		rows:   make(map[string]TodoRecord),
	}
}

// WithTx 以“复制-提交”方式模拟事务。
// fn 返回 nil 时提交，否则回滚。
func (db *InMemoryPostgres) WithTx(ctx context.Context, fn func(tx *Tx) error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	shadowRows := make(map[string]TodoRecord, len(db.rows))
	for k, v := range db.rows {
		shadowRows[k] = v
	}

	tx := &Tx{
		now:    time.Now().UTC,
		nextID: db.nextID,
		rows:   shadowRows,
	}

	if err := fn(tx); err != nil {
		return err
	}

	db.nextID = tx.nextID
	db.rows = tx.rows
	return nil
}

// ListTodos 返回分页结果，按 ID 倒序模拟 created_at desc。
func (db *InMemoryPostgres) ListTodos(page, size int) ([]TodoRecord, error) {
	if page < 1 || size < 1 || size > 100 {
		return nil, ErrInvalidPage
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	ids := make([]int, 0, len(db.rows))
	for id := range db.rows {
		n, _ := strconv.Atoi(id)
		ids = append(ids, n)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(ids)))

	start := (page - 1) * size
	if start >= len(ids) {
		return []TodoRecord{}, nil
	}
	end := start + size
	if end > len(ids) {
		end = len(ids)
	}

	items := make([]TodoRecord, 0, end-start)
	for _, idNum := range ids[start:end] {
		id := strconv.Itoa(idNum)
		items = append(items, db.rows[id])
	}
	return items, nil
}

// Tx 表示事务上下文。
type Tx struct {
	now    func() time.Time
	nextID int
	rows   map[string]TodoRecord
}

// CreateTodo 在事务中插入一条 todo。
func (tx *Tx) CreateTodo(title string) (TodoRecord, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return TodoRecord{}, fmt.Errorf("title required")
	}

	id := strconv.Itoa(tx.nextID)
	tx.nextID++
	now := tx.now()
	todo := TodoRecord{
		ID:        id,
		Title:     title,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	tx.rows[id] = todo
	return todo, nil
}

// MarkDone 在事务中将 todo 标记为完成。
func (tx *Tx) MarkDone(id string) (TodoRecord, error) {
	todo, ok := tx.rows[id]
	if !ok {
		return TodoRecord{}, ErrTodoNotFound
	}
	todo.Done = true
	todo.UpdatedAt = tx.now()
	tx.rows[id] = todo
	return todo, nil
}
