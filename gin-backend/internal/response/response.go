package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yang/go-learning-backend/gin-backend/internal/errs"
)

// ErrorPayload 定义对外错误体。
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Body 是统一响应结构。
type Body struct {
	RequestID string        `json:"request_id,omitempty"`
	Data      any           `json:"data,omitempty"`
	Error     *ErrorPayload `json:"error,omitempty"`
}

// Success 返回成功响应。
func Success(c *gin.Context, status int, data any) {
	// 统一携带 request_id，便于日志和客户端串联问题定位。
	c.JSON(status, Body{
		RequestID: requestID(c),
		Data:      data,
	})
}

// Error 返回错误响应。
func Error(c *gin.Context, appErr *errs.AppError) {
	if appErr == nil {
		appErr = errs.ErrInternal
	}

	// 统一输出错误结构，避免各 handler 自己拼装导致字段不一致。
	c.JSON(appErr.HTTPStatus, Body{
		RequestID: requestID(c),
		Error: &ErrorPayload{
			Code:    appErr.Code,
			Message: appErr.Message,
		},
	})
}

// AbortWithError 输出错误并终止链路。
func AbortWithError(c *gin.Context, appErr *errs.AppError) {
	Error(c, appErr)
	c.Abort()
}

// NoContent 用于删除等无返回体场景。
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// requestID 从 gin context 读取 request_id。
// 注意：request_id 由 RequestID 中间件写入。
func requestID(c *gin.Context) string {
	v, ok := c.Get("request_id")
	if !ok {
		return ""
	}
	id, _ := v.(string)
	return id
}
