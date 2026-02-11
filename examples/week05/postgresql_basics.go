// 详细注释: package week05
package week05

// 详细注释: import (
import (
	// 详细注释: "context"
	"context"
	// 详细注释: "errors"
	"errors"
	// 详细注释: "fmt"
	"fmt"
	// 详细注释: "sort"
	"sort"
	// 详细注释: "strconv"
	"strconv"
	// 详细注释: "strings"
	"strings"
	// 详细注释: "sync"
	"sync"
	// 详细注释: "time"
	"time"
	// 详细注释: )
)

// 详细注释: var (
var (
	// ErrInvalidPage 表示分页参数不合法。
	// 详细注释: ErrInvalidPage = errors.New("invalid pagination")
	ErrInvalidPage = errors.New("invalid pagination")
	// ErrTodoNotFound 表示目标 todo 不存在。
	// 详细注释: ErrTodoNotFound = errors.New("todo not found")
	ErrTodoNotFound = errors.New("todo not found")

// 详细注释: )
)

// TodoRecord 模拟数据库中的 todos 表记录。
// 详细注释: type TodoRecord struct {
type TodoRecord struct {
	// 详细注释: ID        string
	ID string
	// 详细注释: Title     string
	Title string
	// 详细注释: Done      bool
	Done bool
	// 详细注释: CreatedAt time.Time
	CreatedAt time.Time
	// 详细注释: UpdatedAt time.Time
	UpdatedAt time.Time
	// 详细注释: }
}

// InMemoryPostgres 模拟 Postgres 存储。
// 这个类型用于演示“事务边界 + 分页查询”核心概念。
// 详细注释: type InMemoryPostgres struct {
type InMemoryPostgres struct {
	// 详细注释: mu     sync.Mutex
	mu sync.Mutex
	// 详细注释: nextID int
	nextID int
	// 详细注释: rows   map[string]TodoRecord
	rows map[string]TodoRecord
	// 详细注释: }
}

// NewInMemoryPostgres 创建示例数据库实例。
// 详细注释: func NewInMemoryPostgres() *InMemoryPostgres {
func NewInMemoryPostgres() *InMemoryPostgres {
	// 详细注释: return &InMemoryPostgres{
	return &InMemoryPostgres{
		// 详细注释: nextID: 1,
		nextID: 1,
		// 详细注释: rows:   make(map[string]TodoRecord),
		rows: make(map[string]TodoRecord),
		// 详细注释: }
	}
	// 详细注释: }
}

// WithTx 以“复制-提交”方式模拟事务。
// fn 返回 nil 时提交，否则回滚。
// 详细注释: func (db *InMemoryPostgres) WithTx(ctx context.Context, fn func(tx *Tx) error) error {
func (db *InMemoryPostgres) WithTx(ctx context.Context, fn func(tx *Tx) error) error {
	// 详细注释: select {
	select {
	// 详细注释: case <-ctx.Done():
	case <-ctx.Done():
		// 详细注释: return ctx.Err()
		return ctx.Err()
		// 详细注释: default:
	default:
		// 详细注释: }
	}

	// 详细注释: db.mu.Lock()
	db.mu.Lock()
	// 详细注释: defer db.mu.Unlock()
	defer db.mu.Unlock()

	// 详细注释: shadowRows := make(map[string]TodoRecord, len(db.rows))
	shadowRows := make(map[string]TodoRecord, len(db.rows))
	// 详细注释: for k, v := range db.rows {
	for k, v := range db.rows {
		// 详细注释: shadowRows[k] = v
		shadowRows[k] = v
		// 详细注释: }
	}

	// 详细注释: tx := &Tx{
	tx := &Tx{
		// 详细注释: now:    time.Now().UTC,
		now: time.Now().UTC,
		// 详细注释: nextID: db.nextID,
		nextID: db.nextID,
		// 详细注释: rows:   shadowRows,
		rows: shadowRows,
		// 详细注释: }
	}

	// 详细注释: if err := fn(tx); err != nil {
	if err := fn(tx); err != nil {
		// 详细注释: return err
		return err
		// 详细注释: }
	}

	// 详细注释: db.nextID = tx.nextID
	db.nextID = tx.nextID
	// 详细注释: db.rows = tx.rows
	db.rows = tx.rows
	// 详细注释: return nil
	return nil
	// 详细注释: }
}

// ListTodos 返回分页结果，按 ID 倒序模拟 created_at desc。
// 详细注释: func (db *InMemoryPostgres) ListTodos(page, size int) ([]TodoRecord, error) {
func (db *InMemoryPostgres) ListTodos(page, size int) ([]TodoRecord, error) {
	// 详细注释: if page < 1 || size < 1 || size > 100 {
	if page < 1 || size < 1 || size > 100 {
		// 详细注释: return nil, ErrInvalidPage
		return nil, ErrInvalidPage
		// 详细注释: }
	}

	// 详细注释: db.mu.Lock()
	db.mu.Lock()
	// 详细注释: defer db.mu.Unlock()
	defer db.mu.Unlock()

	// 详细注释: ids := make([]int, 0, len(db.rows))
	ids := make([]int, 0, len(db.rows))
	// 详细注释: for id := range db.rows {
	for id := range db.rows {
		// 详细注释: n, _ := strconv.Atoi(id)
		n, _ := strconv.Atoi(id)
		// 详细注释: ids = append(ids, n)
		ids = append(ids, n)
		// 详细注释: }
	}
	// 详细注释: sort.Sort(sort.Reverse(sort.IntSlice(ids)))
	sort.Sort(sort.Reverse(sort.IntSlice(ids)))

	// 详细注释: start := (page - 1) * size
	start := (page - 1) * size
	// 详细注释: if start >= len(ids) {
	if start >= len(ids) {
		// 详细注释: return []TodoRecord{}, nil
		return []TodoRecord{}, nil
		// 详细注释: }
	}
	// 详细注释: end := start + size
	end := start + size
	// 详细注释: if end > len(ids) {
	if end > len(ids) {
		// 详细注释: end = len(ids)
		end = len(ids)
		// 详细注释: }
	}

	// 详细注释: items := make([]TodoRecord, 0, end-start)
	items := make([]TodoRecord, 0, end-start)
	// 详细注释: for _, idNum := range ids[start:end] {
	for _, idNum := range ids[start:end] {
		// 详细注释: id := strconv.Itoa(idNum)
		id := strconv.Itoa(idNum)
		// 详细注释: items = append(items, db.rows[id])
		items = append(items, db.rows[id])
		// 详细注释: }
	}
	// 详细注释: return items, nil
	return items, nil
	// 详细注释: }
}

// Tx 表示事务上下文。
// 详细注释: type Tx struct {
type Tx struct {
	// 详细注释: now    func() time.Time
	now func() time.Time
	// 详细注释: nextID int
	nextID int
	// 详细注释: rows   map[string]TodoRecord
	rows map[string]TodoRecord
	// 详细注释: }
}

// CreateTodo 在事务中插入一条 todo。
// 详细注释: func (tx *Tx) CreateTodo(title string) (TodoRecord, error) {
func (tx *Tx) CreateTodo(title string) (TodoRecord, error) {
	// 详细注释: title = strings.TrimSpace(title)
	title = strings.TrimSpace(title)
	// 详细注释: if title == "" {
	if title == "" {
		// 详细注释: return TodoRecord{}, fmt.Errorf("title required")
		return TodoRecord{}, fmt.Errorf("title required")
		// 详细注释: }
	}

	// 详细注释: id := strconv.Itoa(tx.nextID)
	id := strconv.Itoa(tx.nextID)
	// 详细注释: tx.nextID++
	tx.nextID++
	// 详细注释: now := tx.now()
	now := tx.now()
	// 详细注释: todo := TodoRecord{
	todo := TodoRecord{
		// 详细注释: ID:        id,
		ID: id,
		// 详细注释: Title:     title,
		Title: title,
		// 详细注释: Done:      false,
		Done: false,
		// 详细注释: CreatedAt: now,
		CreatedAt: now,
		// 详细注释: UpdatedAt: now,
		UpdatedAt: now,
		// 详细注释: }
	}
	// 详细注释: tx.rows[id] = todo
	tx.rows[id] = todo
	// 详细注释: return todo, nil
	return todo, nil
	// 详细注释: }
}

// MarkDone 在事务中将 todo 标记为完成。
// 详细注释: func (tx *Tx) MarkDone(id string) (TodoRecord, error) {
func (tx *Tx) MarkDone(id string) (TodoRecord, error) {
	// 详细注释: todo, ok := tx.rows[id]
	todo, ok := tx.rows[id]
	// 详细注释: if !ok {
	if !ok {
		// 详细注释: return TodoRecord{}, ErrTodoNotFound
		return TodoRecord{}, ErrTodoNotFound
		// 详细注释: }
	}
	// 详细注释: todo.Done = true
	todo.Done = true
	// 详细注释: todo.UpdatedAt = tx.now()
	todo.UpdatedAt = tx.now()
	// 详细注释: tx.rows[id] = todo
	tx.rows[id] = todo
	// 详细注释: return todo, nil
	return todo, nil
	// 详细注释: }
}
