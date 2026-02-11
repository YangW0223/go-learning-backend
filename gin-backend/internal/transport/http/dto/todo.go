package dto

// CreateTodoRequest 定义创建 todo 请求。
type CreateTodoRequest struct {
	Title string `json:"title" binding:"required"`
}

// UpdateTodoRequest 定义更新 todo 请求。
// 字段使用指针是为了区分“未传字段”和“字段传空值”两种语义。
type UpdateTodoRequest struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}
