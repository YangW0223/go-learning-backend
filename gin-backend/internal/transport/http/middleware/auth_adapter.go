package middleware

import "github.com/yang/go-learning-backend/gin-backend/internal/service"

// AuthParserAdapter 把 service.AuthService 适配成 middleware.TokenParser。
type AuthParserAdapter struct {
	Auth service.AuthService
}

// ParseToken 把 service 层 claims 转换成 middleware 可识别的轻量结构。
func (a AuthParserAdapter) ParseToken(token string) (ClaimsView, error) {
	claims, err := a.Auth.ParseToken(token)
	if err != nil {
		return ClaimsView{}, err
	}
	return ClaimsView{
		UserID: claims.UserID,
		Email:  claims.Email,
		Role:   claims.Role,
	}, nil
}
