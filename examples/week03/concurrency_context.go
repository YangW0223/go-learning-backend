package week03

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	// ErrRequestTimeout 表示请求在 deadline 之前没有完成。
	// Week03 中用它来模拟“列表接口超时”时的返回错误。
	ErrRequestTimeout = errors.New("request timeout")
)

// SafeCounter 展示 mutex 的典型用法：
// 多个 goroutine 并发写共享状态时，用互斥锁保护临界区。
type SafeCounter struct {
	mu    sync.Mutex
	value int
}

// Add 在加锁状态下更新计数，避免并发写导致数据竞争。
func (c *SafeCounter) Add(n int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += n
}

// Value 返回当前计数值。
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// SumWithMutex 并发累加输入切片，使用 mutex 保护共享变量。
func SumWithMutex(nums []int) int {
	var wg sync.WaitGroup
	counter := &SafeCounter{}

	for _, n := range nums {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			counter.Add(v)
		}(n)
	}

	wg.Wait()
	return counter.Value()
}

// SumWithChannel 并发累加输入切片，使用 channel 汇聚结果。
// 这是“通过通信共享内存”的方式，与 mutex 思路对比学习。
func SumWithChannel(nums []int) int {
	resultCh := make(chan int, len(nums))

	for _, n := range nums {
		go func(v int) {
			resultCh <- v
		}(n)
	}

	total := 0
	for i := 0; i < len(nums); i++ {
		total += <-resultCh
	}
	return total
}

// FirstValue 使用 select 从两个输入 channel 中获取“先到达的值”。
// 这用于演示 select 的核心语义：多个通信操作就绪时，先处理先就绪分支。
func FirstValue(ctx context.Context, ch1 <-chan int, ch2 <-chan int) (int, error) {
	select {
	case v := <-ch1:
		return v, nil
	case v := <-ch2:
		return v, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

// SimulatedListService 模拟一个“列表查询”服务。
// Delay 用来控制慢请求场景，便于复现实验。
type SimulatedListService struct {
	Delay time.Duration
}

// List 在 Delay 后返回固定数据；若 context 先超时/取消则提前返回。
func (s SimulatedListService) List(ctx context.Context) ([]string, error) {
	select {
	case <-time.After(s.Delay):
		return []string{"todo-1", "todo-2", "todo-3"}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// ListWithTimeout 为列表调用增加 deadline 控制。
//
// 这个函数演示 Week03 的关键模式：
// 1) 使用 context.WithTimeout 包裹请求
// 2) 始终 defer cancel() 释放资源
// 3) 将 context 超时错误映射为可读业务错误
func ListWithTimeout(
	parent context.Context,
	timeout time.Duration,
	listFn func(context.Context) ([]string, error),
) ([]string, error) {
	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	items, err := listFn(ctx)
	if err != nil {
		// 这里用 ctx.Err() 判断是否是本次 timeout 引发，
		// 便于把底层错误统一映射为 ErrRequestTimeout。
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return nil, fmt.Errorf("%w: deadline=%s", ErrRequestTimeout, timeout)
		}
		return nil, err
	}
	return items, nil
}
