package model

import "time"

// Todo 是待办领域模型。
// UserID 用于实现“数据按用户隔离”的访问控制基础。
type Todo struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
