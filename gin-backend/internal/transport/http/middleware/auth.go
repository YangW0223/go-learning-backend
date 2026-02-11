package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yang/go-learning-backend/gin-backend/internal/errs"
	"github.com/yang/go-learning-backend/gin-backend/internal/response"
)

// TokenParser 约束 token 解析能力。
type TokenParser interface {
	ParseToken(token string) (ClaimsView, error)
}

// ClaimsView 避免 middleware 直接依赖 auth 包。
type ClaimsView struct {
	UserID string
	Email  string
	Role   string
}

// AuthJWT 校验 Authorization: Bearer <token>。
func AuthJWT(parser TokenParser) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取 Authorization 请求头，要求 Bearer 方案。
		authz := strings.TrimSpace(c.GetHeader("Authorization"))
		if authz == "" || !strings.HasPrefix(strings.ToLower(authz), "bearer ") {
			response.AbortWithError(c, errs.WithMessage(errs.ErrUnauthorized, "missing bearer token"))
			return
		}

		// 去掉 "Bearer " 前缀并清理空白字符。
		token := strings.TrimSpace(authz[7:])
		if token == "" {
			response.AbortWithError(c, errs.WithMessage(errs.ErrUnauthorized, "missing bearer token"))
			return
		}

		// 交给 service 解析 token；失败统一映射成 401。
		claims, err := parser.ParseToken(token)
		if err != nil {
			response.AbortWithError(c, errs.WithMessage(errs.ErrUnauthorized, "invalid token"))
			return
		}

		// 将登录态信息写入上下文，供后续 handler 使用。
		c.Set(ctxKeyUserID, claims.UserID)
		c.Set(ctxKeyUserEmail, claims.Email)
		c.Set(ctxKeyUserRole, claims.Role)
		c.Next()
	}
}
