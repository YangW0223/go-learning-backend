package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yang/go-learning-backend/gin-backend/internal/errs"
	"github.com/yang/go-learning-backend/gin-backend/internal/response"
	"github.com/yang/go-learning-backend/gin-backend/internal/service"
	"github.com/yang/go-learning-backend/gin-backend/internal/transport/http/middleware"
)

// UserHandler 处理用户相关接口。
type UserHandler struct {
	authService service.AuthService
}

// NewUserHandler 创建用户处理器。
func NewUserHandler(authService service.AuthService) *UserHandler {
	return &UserHandler{authService: authService}
}

// Me 返回当前登录用户资料。
func (h *UserHandler) Me(c *gin.Context) {
	// user_id 来自 JWT 中间件解析结果。
	userID := c.GetString(middleware.UserIDKey())
	if userID == "" {
		response.Error(c, errs.WithMessage(errs.ErrUnauthorized, "missing user identity"))
		return
	}
	user, err := h.authService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	response.Success(c, http.StatusOK, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	})
}
