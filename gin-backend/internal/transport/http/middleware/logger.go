package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yang/go-learning-backend/gin-backend/internal/observability"
)

// Logger 记录访问日志并上报请求指标。
func Logger(logger *log.Logger, metrics *observability.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在进入 handler 前记录起始时间。
		start := time.Now()
		c.Next()
		// c.Next 返回后即可拿到最终状态码和耗时。
		duration := time.Since(start)
		status := c.Writer.Status()
		path := c.FullPath()
		if path == "" {
			// 未命中路由时 FullPath 为空，回退到原始路径。
			path = c.Request.URL.Path
		}
		if metrics != nil {
			// 上报聚合指标，便于 Prometheus 抓取。
			metrics.ObserveRequest(c.Request.Method, path, status, duration)
		}
		if logger != nil {
			// 输出结构化日志字段，便于检索分析。
			logger.Printf(
				"level=info method=%s path=%s status=%d latency_ms=%d request_id=%s client_ip=%s",
				c.Request.Method,
				path,
				status,
				duration.Milliseconds(),
				c.GetString(ctxKeyRequestID),
				c.ClientIP(),
			)
		}
	}
}
