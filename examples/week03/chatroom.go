package week03

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	// ErrRoomClosed 表示聊天室已经关闭，不能再接收新操作。
	ErrRoomClosed = errors.New("chat room is closed")
	// ErrUserExists 表示用户名在当前房间内已被占用。
	ErrUserExists = errors.New("user already exists")
	// ErrInvalidUser 表示用户名为空或仅空白字符。
	ErrInvalidUser = errors.New("invalid user")
	// ErrInvalidMessage 表示消息内容为空或仅空白字符。
	ErrInvalidMessage = errors.New("invalid message")
)

// ChatMessage 是聊天室广播的最小消息模型。
type ChatMessage struct {
	From      string    `json:"from"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// ChatRoomStats 用于观察当前房间状态，便于测试和排查。
type ChatRoomStats struct {
	UsersCount    int
	MessagesCount int
}

type joinRequest struct {
	user string
	resp chan joinResponse
}

type joinResponse struct {
	ch  chan ChatMessage
	err error
}

// ChatRoom 使用“单事件循环 + channel”管理并发状态。
//
// 设计要点：
// 1) 用户列表和广播都在同一个 goroutine 内处理，减少锁竞争。
// 2) 外部通过 Join/Leave/Publish/Stats 与房间通信，不直接访问内部 map。
// 3) 对慢消费者采用“丢弃消息并记录日志”策略，避免拖垮全局广播。
type ChatRoom struct {
	name string

	joinCh    chan joinRequest
	leaveCh   chan string
	publishCh chan ChatMessage
	statsCh   chan chan ChatRoomStats

	doneCh chan struct{}
	once   sync.Once

	logf func(format string, args ...any)
}

// NewChatRoom 创建并启动聊天室事件循环。
// buffer 用于控制每个用户消息队列长度，避免瞬时广播阻塞。
func NewChatRoom(name string, logf func(format string, args ...any)) *ChatRoom {
	room := &ChatRoom{
		name:      strings.TrimSpace(name),
		joinCh:    make(chan joinRequest),
		leaveCh:   make(chan string, 16),
		publishCh: make(chan ChatMessage, 64),
		statsCh:   make(chan chan ChatRoomStats),
		doneCh:    make(chan struct{}),
		logf:      logf,
	}
	go room.run(16)
	return room
}

// run 是聊天室核心事件循环，串行处理所有状态变更。
func (r *ChatRoom) run(userBuffer int) {
	users := make(map[string]chan ChatMessage)
	messagesCount := 0

	for {
		select {
		case req := <-r.joinCh:
			if _, exists := users[req.user]; exists {
				req.resp <- joinResponse{err: ErrUserExists}
				continue
			}
			userCh := make(chan ChatMessage, userBuffer)
			users[req.user] = userCh
			req.resp <- joinResponse{ch: userCh, err: nil}
			r.debugf("room=%s user_join user=%s", r.name, req.user)

		case user := <-r.leaveCh:
			if userCh, ok := users[user]; ok {
				close(userCh)
				delete(users, user)
				r.debugf("room=%s user_leave user=%s", r.name, user)
			}

		case msg := <-r.publishCh:
			messagesCount++
			for user, userCh := range users {
				select {
				case userCh <- msg:
				default:
					// 慢消费者策略：丢弃消息并记录日志，避免全房间阻塞。
					r.debugf("room=%s drop_message user=%s from=%s", r.name, user, msg.From)
				}
			}

		case resp := <-r.statsCh:
			resp <- ChatRoomStats{
				UsersCount:    len(users),
				MessagesCount: messagesCount,
			}

		case <-r.doneCh:
			// 房间关闭时统一关闭所有用户队列，防止读端永远阻塞。
			for user, userCh := range users {
				close(userCh)
				delete(users, user)
			}
			return
		}
	}
}

// Join 将用户加入聊天室并返回接收消息的 channel。
// 当传入 context 被取消时，会自动触发 Leave，避免 goroutine 泄漏。
func (r *ChatRoom) Join(ctx context.Context, user string) (<-chan ChatMessage, error) {
	user = strings.TrimSpace(user)
	if user == "" {
		return nil, ErrInvalidUser
	}

	req := joinRequest{
		user: user,
		resp: make(chan joinResponse, 1),
	}

	select {
	case <-r.doneCh:
		return nil, ErrRoomClosed
	case <-ctx.Done():
		return nil, ctx.Err()
	case r.joinCh <- req:
	}

	select {
	case <-r.doneCh:
		return nil, ErrRoomClosed
	case <-ctx.Done():
		return nil, ctx.Err()
	case resp := <-req.resp:
		if resp.err != nil {
			return nil, resp.err
		}

		// 仅当 ctx 可取消时才启动自动退出协程，避免无意义常驻 goroutine。
		if ctx.Done() != nil {
			go func(username string, done <-chan struct{}) {
				select {
				case <-done:
					_ = r.Leave(username)
				case <-r.doneCh:
				}
			}(user, ctx.Done())
		}

		return resp.ch, nil
	}
}

// Leave 将用户移出聊天室。重复离开是幂等的。
func (r *ChatRoom) Leave(user string) error {
	user = strings.TrimSpace(user)
	if user == "" {
		return ErrInvalidUser
	}

	select {
	case <-r.doneCh:
		return ErrRoomClosed
	case r.leaveCh <- user:
		return nil
	}
}

// Publish 发送一条消息到聊天室。
func (r *ChatRoom) Publish(ctx context.Context, from, content string) error {
	from = strings.TrimSpace(from)
	content = strings.TrimSpace(content)
	if from == "" {
		return ErrInvalidUser
	}
	if content == "" {
		return ErrInvalidMessage
	}

	msg := ChatMessage{
		From:      from,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}

	select {
	case <-r.doneCh:
		return ErrRoomClosed
	case <-ctx.Done():
		return ctx.Err()
	case r.publishCh <- msg:
		return nil
	}
}

// Stats 返回聊天室当前统计信息。
func (r *ChatRoom) Stats(ctx context.Context) (ChatRoomStats, error) {
	resp := make(chan ChatRoomStats, 1)
	select {
	case <-r.doneCh:
		return ChatRoomStats{}, ErrRoomClosed
	case <-ctx.Done():
		return ChatRoomStats{}, ctx.Err()
	case r.statsCh <- resp:
	}

	select {
	case <-r.doneCh:
		return ChatRoomStats{}, ErrRoomClosed
	case <-ctx.Done():
		return ChatRoomStats{}, ctx.Err()
	case stats := <-resp:
		return stats, nil
	}
}

// Close 关闭聊天室并释放资源，调用多次安全。
func (r *ChatRoom) Close() {
	r.once.Do(func() {
		close(r.doneCh)
	})
}

func (r *ChatRoom) debugf(format string, args ...any) {
	if r.logf != nil {
		r.logf(format, args...)
	}
}

// WaitMessage 在指定超时时间内等待一条消息。
// 这是给示例和测试复用的辅助函数。
func WaitMessage(ch <-chan ChatMessage, timeout time.Duration) (ChatMessage, error) {
	select {
	case msg, ok := <-ch:
		if !ok {
			return ChatMessage{}, fmt.Errorf("channel closed")
		}
		return msg, nil
	case <-time.After(timeout):
		return ChatMessage{}, fmt.Errorf("wait message timeout: %s", timeout)
	}
}
