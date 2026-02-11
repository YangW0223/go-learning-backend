package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateEntry struct {
	count       int
	windowStart time.Time
}

// RateLimit 提供基于 IP 的简易限流。
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	if limit <= 0 {
		limit = 100
	}
	if window <= 0 {
		window = time.Minute
	}

	store := map[string]*rateEntry{}
	var mu sync.Mutex

	return func(c *gin.Context) {
		// 按客户端 IP 做固定窗口限流。
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		entry, ok := store[ip]
		if !ok {
			entry = &rateEntry{count: 0, windowStart: now}
			store[ip] = entry
		}
		if now.Sub(entry.windowStart) >= window {
			entry.count = 0
			entry.windowStart = now
		}
		entry.count++
		current := entry.count
		mu.Unlock()

		// 超过阈值时直接返回 429，不再进入后续 handler。
		if current > limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"request_id": c.GetString(ctxKeyRequestID),
				"error": gin.H{
					"code":    "RATE_LIMITED",
					"message": "too many requests",
				},
			})
			return
		}
		c.Next()
	}
}
