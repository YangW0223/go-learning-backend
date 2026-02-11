package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/yang/go-learning-backend/examples/week03"
)

func main() {
	// 1) 对比 mutex 与 channel 的并发求和结果。
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println("SumWithMutex:", week03.SumWithMutex(nums))
	fmt.Println("SumWithChannel:", week03.SumWithChannel(nums))

	// 2) 演示 select: 谁先到达就先返回。
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 7
	first, err := week03.FirstValue(context.Background(), ch1, ch2)
	if err != nil {
		fmt.Println("FirstValue error:", err)
	} else {
		fmt.Println("FirstValue:", first)
	}

	// 3) 演示列表接口超时控制（成功场景）。
	fastService := week03.SimulatedListService{Delay: 30 * time.Millisecond}
	items, err := week03.ListWithTimeout(
		context.Background(),
		120*time.Millisecond,
		fastService.List,
	)
	if err != nil {
		fmt.Println("ListWithTimeout(success) error:", err)
	} else {
		fmt.Println("ListWithTimeout(success):", items)
	}

	// 4) 演示慢请求触发超时（失败场景）。
	slowService := week03.SimulatedListService{Delay: 180 * time.Millisecond}
	_, err = week03.ListWithTimeout(
		context.Background(),
		40*time.Millisecond,
		slowService.List,
	)
	if errors.Is(err, week03.ErrRequestTimeout) {
		fmt.Println("ListWithTimeout(timeout):", err)
	}

	// 5) 演示聊天室：加入 -> 广播 -> 接收 -> 取消退出。
	room := week03.NewChatRoom("demo-room", nil)
	defer room.Close()

	aliceCtx, aliceCancel := context.WithCancel(context.Background())
	defer aliceCancel()

	aliceCh, err := room.Join(aliceCtx, "alice")
	if err != nil {
		fmt.Println("ChatRoom join alice error:", err)
		return
	}
	bobCh, err := room.Join(context.Background(), "bob")
	if err != nil {
		fmt.Println("ChatRoom join bob error:", err)
		return
	}

	if err := room.Publish(context.Background(), "alice", "hello bob"); err != nil {
		fmt.Println("ChatRoom publish error:", err)
		return
	}

	aliceMsg, err := week03.WaitMessage(aliceCh, 200*time.Millisecond)
	if err != nil {
		fmt.Println("ChatRoom alice receive error:", err)
		return
	}
	bobMsg, err := week03.WaitMessage(bobCh, 200*time.Millisecond)
	if err != nil {
		fmt.Println("ChatRoom bob receive error:", err)
		return
	}
	fmt.Println("ChatRoom alice received:", aliceMsg.Content)
	fmt.Println("ChatRoom bob received:", bobMsg.Content)

	// 取消 alice 会话，触发自动离开。
	aliceCancel()
	time.Sleep(60 * time.Millisecond)

	stats, err := room.Stats(context.Background())
	if err != nil {
		fmt.Println("ChatRoom stats error:", err)
		return
	}
	fmt.Printf("ChatRoom stats: users=%d messages=%d\n", stats.UsersCount, stats.MessagesCount)
}
