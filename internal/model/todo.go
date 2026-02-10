package model

import "time"

// Todo 是项目的核心领域模型。
// 它会被 store 层读写，也会被 handler 层序列化为 JSON 返回给前端。
type Todo struct {
	// ID 是待办项唯一标识。
	ID string `json:"id"`
	// Title 是待办文本内容。
	Title string `json:"title"`
	// Done 表示是否完成。
	Done bool `json:"done"`
	// CreatedAt 记录创建时间（UTC）。
	CreatedAt time.Time `json:"created_at"`
}
