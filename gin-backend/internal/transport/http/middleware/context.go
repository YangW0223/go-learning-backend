package middleware

const (
	// ctxKeyRequestID 保存请求级 trace id。
	ctxKeyRequestID = "request_id"
	// ctxKeyUserID 保存当前登录用户 ID。
	ctxKeyUserID = "user_id"
	// ctxKeyUserEmail 保存当前登录邮箱。
	ctxKeyUserEmail = "user_email"
	// ctxKeyUserRole 保存当前登录角色。
	ctxKeyUserRole = "user_role"
)

// RequestIDKey 暴露 request id 上下文 key。
func RequestIDKey() string { return ctxKeyRequestID }

// UserIDKey 暴露 user id 上下文 key。
func UserIDKey() string { return ctxKeyUserID }

// UserEmailKey 暴露 user email 上下文 key。
func UserEmailKey() string { return ctxKeyUserEmail }

// UserRoleKey 暴露 user role 上下文 key。
func UserRoleKey() string { return ctxKeyUserRole }
