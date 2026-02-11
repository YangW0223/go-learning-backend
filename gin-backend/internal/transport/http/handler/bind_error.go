package handler

import (
	"errors"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/yang/go-learning-backend/gin-backend/internal/errs"
)

// responseBindError 将 Gin 绑定错误转换成统一 AppError。
// 这样每个 handler 不需要重复实现错误格式转换逻辑。
func responseBindError(err error) *errs.AppError {
	if err == nil {
		return errs.ErrBadRequest
	}
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		return errs.WithMessage(errs.ErrBadRequest, "request validation failed")
	}
	if errors.Is(err, binding.ErrMissingContentType) {
		return errs.WithMessage(errs.ErrBadRequest, "content-type is required")
	}
	return errs.WithMessage(errs.ErrBadRequest, "invalid request body")
}
