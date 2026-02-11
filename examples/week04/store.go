// package week04 定义 Week04 的分层示例代码。
package week04

// import 分组引入 store 层所需的标准库能力。
import (
	// context 用于在数据访问时透传取消信号与超时控制。
	"context"
	// sync 提供互斥锁，保证 map 在并发访问时安全。
	"sync"
)

// TodoStore 定义 service 层依赖的数据访问抽象。
// service 层只面向接口编程，不直接依赖具体存储实现。
type TodoStore interface {
	// MarkDone 根据 id 将目标 Todo 标记为完成，并返回更新后的实体。
	MarkDone(ctx context.Context, id string) (Todo, error)
}

// InMemoryTodoStore 是内存版 TodoStore 实现，用于教学与本地演示。
// 它通过 map 保存数据，并通过互斥锁保护并发读写。
type InMemoryTodoStore struct {
	// mu 保护 todos，避免并发读写 map 引发数据竞争。
	mu sync.Mutex
	// todos 以 id 为键保存当前内存中的 todo 快照。
	todos map[string]Todo
}

// NewInMemoryTodoStore 根据 seed 初始化一个内存 store。
func NewInMemoryTodoStore(seed []Todo) *InMemoryTodoStore {
	// 按 seed 长度预分配 map 容量，减少扩容次数。
	todos := make(map[string]Todo, len(seed))
	// 遍历初始化数据，逐条写入 map。
	for _, todo := range seed {
		// 使用 todo.ID 作为唯一键保存对应实体。
		todos[todo.ID] = todo
	}
	// 返回可用的内存 store 实例。
	return &InMemoryTodoStore{
		// 将初始化好的 map 注入实例字段。
		todos: todos,
	}
}

// MarkDone 将指定 id 的 Todo 标记为已完成。
func (s *InMemoryTodoStore) MarkDone(ctx context.Context, id string) (Todo, error) {
	// 在进入临界区前先检查 context 是否已经被取消。
	select {
	// 若收到取消信号，直接返回 context 的错误。
	case <-ctx.Done():
		return Todo{}, ctx.Err()
	// default 分支表示 context 仍可用，继续执行后续逻辑。
	default:
	}

	// 加锁进入临界区，开始安全访问共享 map。
	s.mu.Lock()
	// 函数返回前自动解锁，确保任何路径都不会遗漏解锁。
	defer s.mu.Unlock()

	// 读取目标 id 对应的 todo，并判断是否存在。
	todo, ok := s.todos[id]
	// 若不存在则返回统一的业务错误 ErrTodoNotFound。
	if !ok {
		return Todo{}, ErrTodoNotFound
	}

	// 将目标 todo 的完成状态改为 true。
	todo.Done = true
	// 把更新后的 todo 回写到 map 中。
	s.todos[id] = todo
	// 返回更新后的 todo 及 nil 错误。
	return todo, nil
}
