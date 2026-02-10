package store

import (
	"errors"

	"github.com/yang/go-learning-backend/internal/model"
)

// ErrTodoNotFound 是存储层暴露给上层的哨兵错误。
// handler 层会用 errors.Is 判断该错误并映射为 HTTP 404。
var ErrTodoNotFound = errors.New("todo not found")

// TodoStore 定义“待办存储”的最小行为契约。
//
// 设计意义：
// 1) handler 依赖接口而不是具体实现，便于替换存储（内存/数据库）。
// 2) 测试时可用 mock 或 fake 实现，减少耦合。
type TodoStore interface {
	Create(title string) (model.Todo, error)
	List() ([]model.Todo, error)
	MarkDone(id string) (model.Todo, error)
}
