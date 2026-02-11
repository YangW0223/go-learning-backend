// package week04 定义 Week04 的分层示例代码。
package week04

// import 分组引入 service 测试所需标准库。
import (
	// context 用于构造测试调用时需要的上下文。
	"context"
	// errors 用于断言返回错误类型。
	"errors"
	// testing 提供 Go 单元测试框架。
	"testing"
)

// spyTodoStore 是测试替身，用于观测 service 是否调用了 store。
type spyTodoStore struct {
	// called 记录 MarkDone 被调用次数。
	called int
	// todo 是替身返回给 service 的模拟结果。
	todo Todo
	// err 是替身返回给 service 的模拟错误。
	err error
}

// MarkDone 实现 TodoStore 接口并记录调用行为。
func (s *spyTodoStore) MarkDone(_ context.Context, _ string) (Todo, error) {
	// 每次调用先累加计数，便于断言是否被触发。
	s.called++
	// 返回预先注入的结果与错误。
	return s.todo, s.err
}

// TestTodoServiceMarkDone_Success 验证 service 正常返回已完成的 todo。
func TestTodoServiceMarkDone_Success(t *testing.T) {
	// 使用真实内存 store 构造一条未完成任务。
	svc := NewTodoService(NewInMemoryTodoStore([]Todo{
		// id=1 的任务初始为未完成。
		{ID: "1", Title: "task", Done: false},
	}))

	// 调用目标方法，传入合法 id。
	got, err := svc.MarkDone(context.Background(), "1")
	// 成功场景不应返回错误。
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 返回结果应被标记为 done=true。
	if !got.Done {
		t.Fatalf("expected done=true, got %+v", got)
	}
}

// TestTodoServiceMarkDone_InvalidID 验证非法 id 会在 service 层被拦截。
func TestTodoServiceMarkDone_InvalidID(t *testing.T) {
	// 构造可观测的 store 替身。
	store := &spyTodoStore{}
	// 使用替身初始化 service。
	svc := NewTodoService(store)

	// 传入非法 id=abc 执行调用。
	_, err := svc.MarkDone(context.Background(), "abc")
	// 预期返回 ErrInvalidTodoID。
	if !errors.Is(err, ErrInvalidTodoID) {
		t.Fatalf("expected ErrInvalidTodoID, got %v", err)
	}
	// 非法输入应在 service 层终止，不应触发 store。
	if store.called != 0 {
		t.Fatalf("store should not be called when id is invalid")
	}
}

// TestTodoServiceMarkDone_NotFound 验证资源不存在时返回 ErrTodoNotFound。
func TestTodoServiceMarkDone_NotFound(t *testing.T) {
	// 使用空数据 store，确保任意合法 id 都查不到。
	svc := NewTodoService(NewInMemoryTodoStore(nil))

	// 对不存在的 id=1 发起调用。
	_, err := svc.MarkDone(context.Background(), "1")
	// 预期返回未找到错误。
	if !errors.Is(err, ErrTodoNotFound) {
		t.Fatalf("expected ErrTodoNotFound, got %v", err)
	}
}

// TestTodoServiceMarkDone_StoreError 验证未知底层错误会被包装返回。
func TestTodoServiceMarkDone_StoreError(t *testing.T) {
	// 构造一个返回未知错误的 store 替身。
	store := &spyTodoStore{
		// 模拟数据库不可用等基础设施错误。
		err: errors.New("db unavailable"),
	}
	// 使用该替身初始化 service。
	svc := NewTodoService(store)

	// 对合法 id 发起调用，触发 store 错误分支。
	_, err := svc.MarkDone(context.Background(), "1")
	// 该场景必须返回错误，不能吞掉异常。
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	// 未知错误不应被误映射成已知业务错误。
	if errors.Is(err, ErrTodoNotFound) || errors.Is(err, ErrInvalidTodoID) {
		t.Fatalf("unexpected business error mapping: %v", err)
	}
}
