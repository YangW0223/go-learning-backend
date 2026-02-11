package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery 捕获 panic，避免服务崩溃。
func Recovery(logger *log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				// 输出 panic 与堆栈，便于定位根因。
				if logger != nil {
					logger.Printf("level=error msg=panic recovered err=%v stack=%q", rec, string(debug.Stack()))
				}

				// 使用统一结构响应 500，避免返回默认 HTML 错误页。
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"request_id": c.GetString(ctxKeyRequestID),
					"error": gin.H{
						"code":    "INTERNAL_ERROR",
						"message": "internal server error",
					},
				})
			}
		}()
		c.Next()
	}
}
