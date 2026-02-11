package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yang/go-learning-backend/gin-backend/internal/response"
	"github.com/yang/go-learning-backend/gin-backend/internal/service"
	"github.com/yang/go-learning-backend/gin-backend/internal/transport/http/dto"
)

// AuthHandler 负责认证接口。
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler 创建认证处理器。
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register 处理用户注册请求。
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	// 绑定并校验 JSON 请求体。
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, responseBindError(err))
		return
	}

	// 调用 service 执行业务逻辑。
	user, err := h.authService.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	// 对外返回必要字段，避免暴露敏感信息（如密码 hash）。
	response.Success(c, http.StatusCreated, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
	})
}

// Login 处理登录请求并返回 access token。
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, responseBindError(err))
		return
	}

	// service 返回 token 与用户信息，handler 只负责协议层映射。
	token, user, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	response.Success(c, http.StatusOK, dto.AuthResponse{
		AccessToken: token,
		User: gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"role":       user.Role,
			"created_at": user.CreatedAt,
		},
	})
}
