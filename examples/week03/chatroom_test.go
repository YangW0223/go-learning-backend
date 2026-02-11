package week03

import (
	"context"
	"errors"
	"testing"
	"time"
)

// TestChatRoom_BroadcastToAll 验证广播会发送给所有在线用户。
func TestChatRoom_BroadcastToAll(t *testing.T) {
	room := NewChatRoom("week03-room", nil)
	defer room.Close()

	aliceCh, err := room.Join(context.Background(), "alice")
	if err != nil {
		t.Fatalf("alice join error: %v", err)
	}
	bobCh, err := room.Join(context.Background(), "bob")
	if err != nil {
		t.Fatalf("bob join error: %v", err)
	}

	if err := room.Publish(context.Background(), "alice", "hello everyone"); err != nil {
		t.Fatalf("publish error: %v", err)
	}

	aliceMsg, err := WaitMessage(aliceCh, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("alice wait message error: %v", err)
	}
	bobMsg, err := WaitMessage(bobCh, 200*time.Millisecond)
	if err != nil {
		t.Fatalf("bob wait message error: %v", err)
	}

	if aliceMsg.Content != "hello everyone" || bobMsg.Content != "hello everyone" {
		t.Fatalf("unexpected broadcast content: alice=%q bob=%q", aliceMsg.Content, bobMsg.Content)
	}
}

// TestChatRoom_DuplicateUser 验证重复用户名加入会失败。
func TestChatRoom_DuplicateUser(t *testing.T) {
	room := NewChatRoom("week03-room", nil)
	defer room.Close()

	_, err := room.Join(context.Background(), "alice")
	if err != nil {
		t.Fatalf("first join error: %v", err)
	}

	_, err = room.Join(context.Background(), "alice")
	if !errors.Is(err, ErrUserExists) {
		t.Fatalf("expected ErrUserExists, got %v", err)
	}
}

// TestChatRoom_ContextCancelAutoLeave 验证 context 取消后用户会自动离开。
func TestChatRoom_ContextCancelAutoLeave(t *testing.T) {
	room := NewChatRoom("week03-room", nil)
	defer room.Close()

	ctx, cancel := context.WithCancel(context.Background())
	_, err := room.Join(ctx, "alice")
	if err != nil {
		t.Fatalf("join error: %v", err)
	}

	cancel()

	deadline := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(deadline) {
		stats, statsErr := room.Stats(context.Background())
		if statsErr != nil {
			t.Fatalf("stats error: %v", statsErr)
		}
		if stats.UsersCount == 0 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("expected users count to become 0 after context cancel")
}

// TestChatRoom_CloseClosesUserChannel 验证关闭房间时会关闭用户接收通道。
func TestChatRoom_CloseClosesUserChannel(t *testing.T) {
	room := NewChatRoom("week03-room", nil)

	aliceCh, err := room.Join(context.Background(), "alice")
	if err != nil {
		t.Fatalf("join error: %v", err)
	}

	room.Close()

	select {
	case _, ok := <-aliceCh:
		if ok {
			t.Fatal("expected alice channel to be closed")
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("expected alice channel closed quickly after room close")
	}
}
