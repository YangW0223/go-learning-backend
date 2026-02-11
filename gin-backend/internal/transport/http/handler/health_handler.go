package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yang/go-learning-backend/gin-backend/internal/response"
)

// Healthz 提供存活检查。
// 语义：进程是否存活，通常用于 Kubernetes liveness probe。
func Healthz(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{"status": "ok"})
}

// Readyz 提供就绪检查。
// 语义：服务是否准备好处理请求，通常用于 readiness probe。
func Readyz(c *gin.Context) {
	response.Success(c, http.StatusOK, gin.H{"status": "ready"})
}
