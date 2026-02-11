package errs

import "net/http"

// AppError 定义应用统一错误对象。
type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
}

// Error 实现 error 接口，便于 errors.As/Is 统一处理。
func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

// newAppError 构造统一错误实例。
func newAppError(status int, code, message string) *AppError {
	return &AppError{Code: code, Message: message, HTTPStatus: status}
}

var (
	ErrBadRequest          = newAppError(http.StatusBadRequest, "BAD_REQUEST", "invalid request")
	ErrUnauthorized        = newAppError(http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
	ErrForbidden           = newAppError(http.StatusForbidden, "FORBIDDEN", "forbidden")
	ErrNotFound            = newAppError(http.StatusNotFound, "NOT_FOUND", "resource not found")
	ErrConflict            = newAppError(http.StatusConflict, "CONFLICT", "resource conflict")
	ErrUnprocessableEntity = newAppError(http.StatusUnprocessableEntity, "UNPROCESSABLE_ENTITY", "unprocessable entity")
	ErrInternal            = newAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
)

// WithMessage 复制一个新错误并替换 message。
func WithMessage(src *AppError, message string) *AppError {
	if src == nil {
		return ErrInternal
	}
	return &AppError{Code: src.Code, Message: message, HTTPStatus: src.HTTPStatus}
}
