package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/yang/go-learning-backend/gin-backend/internal/errs"
	"github.com/yang/go-learning-backend/gin-backend/internal/response"
)

// writeServiceError 把 service 层错误映射成统一 HTTP 错误响应。
// 约定：service 返回 *errs.AppError 时可直接透传；否则兜底 500。
func writeServiceError(c *gin.Context, err error) {
	var appErr *errs.AppError
	if errors.As(err, &appErr) {
		response.Error(c, appErr)
		return
	}
	response.Error(c, errs.ErrInternal)
}
