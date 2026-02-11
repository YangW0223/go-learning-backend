package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

// RequestID 为请求注入 request id。
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 优先使用上游透传的 X-Request-ID，便于跨服务链路追踪。
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			// 若上游未提供则本地生成，保证每个请求都有 request id。
			id = randomID()
		}
		// 同时写入 context 与响应头，方便日志和客户端侧排查。
		c.Set(ctxKeyRequestID, id)
		c.Writer.Header().Set("X-Request-ID", id)
		c.Next()
	}
}

// randomID 生成短随机 ID，足够用于请求级追踪场景。
func randomID() string {
	buf := make([]byte, 12)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}
