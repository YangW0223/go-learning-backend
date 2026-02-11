package week03

import (
	"context"
	"errors"
	"testing"
	"time"
)

// TestSumWithMutexAndChannel 验证两种并发同步方式得到一致结果。
func TestSumWithMutexAndChannel(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5}

	gotMutex := SumWithMutex(nums)
	gotChannel := SumWithChannel(nums)

	if gotMutex != 15 {
		t.Fatalf("SumWithMutex expected 15, got %d", gotMutex)
	}
	if gotChannel != 15 {
		t.Fatalf("SumWithChannel expected 15, got %d", gotChannel)
	}
}

// TestFirstValue 验证 select 能拿到先到达的 channel 值。
func TestFirstValue(t *testing.T) {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch2 <- 99

	got, err := FirstValue(context.Background(), ch1, ch2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 99 {
		t.Fatalf("expected 99, got %d", got)
	}
}

// TestListWithTimeout_Success 验证“服务在超时前返回”场景。
func TestListWithTimeout_Success(t *testing.T) {
	service := SimulatedListService{Delay: 20 * time.Millisecond}

	items, err := ListWithTimeout(
		context.Background(),
		100*time.Millisecond,
		service.List,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}
}

// TestListWithTimeout_Timeout 验证“慢请求触发超时”场景。
func TestListWithTimeout_Timeout(t *testing.T) {
	service := SimulatedListService{Delay: 120 * time.Millisecond}

	_, err := ListWithTimeout(
		context.Background(),
		30*time.Millisecond,
		service.List,
	)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !errors.Is(err, ErrRequestTimeout) {
		t.Fatalf("expected ErrRequestTimeout, got %v", err)
	}
}
