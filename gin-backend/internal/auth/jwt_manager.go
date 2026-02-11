package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims 定义 token 载荷。
type Claims struct {
	UserID string `json:"uid"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTManager 负责 token 签发与解析。
type JWTManager struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

// NewJWTManager 创建 JWT 管理器。
func NewJWTManager(secret, issuer string, ttl time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), issuer: issuer, ttl: ttl}
}

// Generate 生成访问令牌。
func (m *JWTManager) Generate(userID, email, role string) (string, error) {
	// 使用 UTC 时间统一签发时间，避免时区差异带来的调试困扰。
	now := time.Now().UTC()
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}

	// 使用 HS256 对 claims 进行签名。
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}

// Parse 解析并校验访问令牌。
func (m *JWTManager) Parse(tokenString string) (Claims, error) {
	// ParseWithClaims 会校验签名与标准字段（如 exp）。
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
		}
		return m.secret, nil
	})
	if err != nil {
		return Claims{}, fmt.Errorf("parse token: %w", err)
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return Claims{}, fmt.Errorf("invalid token claims")
	}
	// 额外校验 issuer，避免跨系统 token 混用。
	if claims.Issuer != m.issuer {
		return Claims{}, fmt.Errorf("invalid issuer")
	}
	return *claims, nil
}
