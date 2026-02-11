package middleware

import "github.com/gin-gonic/gin"

// CORS 注入基础跨域响应头。
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 当前策略对学习项目采用宽松配置（*），生产环境建议按域名白名单收敛。
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
