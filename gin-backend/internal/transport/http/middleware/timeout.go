package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// Timeout 在请求级别注入超时上下文。
func Timeout(timeout time.Duration) gin.HandlerFunc {
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	return func(c *gin.Context) {
		// 为每个请求创建独立超时上下文，向下游 DB/Redis 调用透传。
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
