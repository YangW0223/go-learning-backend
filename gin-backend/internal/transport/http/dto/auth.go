package dto

// RegisterRequest 定义注册请求。
// 这里只定义协议层字段，不包含任何业务行为。
type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRequest 定义登录请求。
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse 定义认证成功返回。
// User 使用 any 是为了兼容不同用户视图（可按需要替换为明确结构体）。
type AuthResponse struct {
	AccessToken string `json:"access_token"`
	User        any    `json:"user"`
}
