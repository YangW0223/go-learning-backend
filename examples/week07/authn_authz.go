package week07

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	// ErrUserExists 表示用户已存在。
	ErrUserExists = errors.New("user exists")
	// ErrInvalidCredential 表示用户名或密码错误。
	ErrInvalidCredential = errors.New("invalid credential")
	// ErrInvalidToken 表示 token 非法。
	ErrInvalidToken = errors.New("invalid token")
	// ErrTokenExpired 表示 token 已过期。
	ErrTokenExpired = errors.New("token expired")
	// ErrUnauthorized 表示未认证。
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden 表示已认证但无权限。
	ErrForbidden = errors.New("forbidden")
)

// User 表示用户实体。
type User struct {
	ID           string
	Username     string
	PasswordHash string
	Role         string
}

// Claims 表示 token 载荷。
type Claims struct {
	Sub string `json:"sub"`
	Rol string `json:"rol"`
	Exp int64  `json:"exp"`
}

// AuthService 提供注册、登录、鉴权能力。
type AuthService struct {
	mu       sync.Mutex
	nextID   int
	users    map[string]User
	secret   []byte
	nowFn    func() time.Time
	tokenTTL time.Duration
}

// NewAuthService 创建认证服务。
func NewAuthService(secret string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		nextID:   1,
		users:    make(map[string]User),
		secret:   []byte(secret),
		nowFn:    time.Now,
		tokenTTL: tokenTTL,
	}
}

// Register 注册用户并保存密码摘要。
func (s *AuthService) Register(username, password, role string) error {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	role = strings.TrimSpace(role)
	if username == "" || password == "" || role == "" {
		return ErrInvalidCredential
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[username]; exists {
		return ErrUserExists
	}

	user := User{
		ID:           fmt.Sprintf("u-%d", s.nextID),
		Username:     username,
		PasswordHash: s.hashPassword(password),
		Role:         role,
	}
	s.nextID++
	s.users[username] = user
	return nil
}

// Login 校验凭证并签发 token。
func (s *AuthService) Login(username, password string) (string, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	s.mu.Lock()
	user, ok := s.users[username]
	s.mu.Unlock()
	if !ok {
		return "", ErrInvalidCredential
	}
	if s.hashPassword(password) != user.PasswordHash {
		return "", ErrInvalidCredential
	}

	claims := Claims{
		Sub: user.ID,
		Rol: user.Role,
		Exp: s.nowFn().Add(s.tokenTTL).Unix(),
	}
	return s.signClaims(claims)
}

// ValidateToken 校验 token 并返回 claims。
func (s *AuthService) ValidateToken(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Claims{}, ErrInvalidToken
	}

	headerPayload := parts[0] + "." + parts[1]
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}
	if !hmac.Equal(sig, s.hmac(headerPayload)) {
		return Claims{}, ErrInvalidToken
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}
	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return Claims{}, ErrInvalidToken
	}
	if s.nowFn().Unix() > claims.Exp {
		return Claims{}, ErrTokenExpired
	}
	return claims, nil
}

// RequireRole 校验 token 和角色权限。
func (s *AuthService) RequireRole(token, role string) (Claims, error) {
	claims, err := s.ValidateToken(token)
	if err != nil {
		return Claims{}, ErrUnauthorized
	}
	if claims.Rol != role {
		return Claims{}, ErrForbidden
	}
	return claims, nil
}

// AuthMiddleware 提供最小 HTTP 鉴权中间件。
func AuthMiddleware(auth *AuthService, requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := strings.TrimSpace(r.Header.Get("Authorization"))
		if !strings.HasPrefix(header, "Bearer ") {
			http.Error(w, ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		_, err := auth.RequireRole(token, requiredRole)
		if err != nil {
			if errors.Is(err, ErrForbidden) {
				http.Error(w, ErrForbidden.Error(), http.StatusForbidden)
				return
			}
			http.Error(w, ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func (s *AuthService) hashPassword(password string) string {
	// 示例使用 sha256 + secret，真实生产建议使用 bcrypt/argon2。
	sum := sha256.Sum256(append(s.secret, []byte(password)...))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func (s *AuthService) signClaims(claims Claims) (string, error) {
	headerBytes, _ := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	header := base64.RawURLEncoding.EncodeToString(headerBytes)
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	headerPayload := header + "." + payload
	signature := base64.RawURLEncoding.EncodeToString(s.hmac(headerPayload))
	return headerPayload + "." + signature, nil
}

func (s *AuthService) hmac(content string) []byte {
	h := hmac.New(sha256.New, s.secret)
	_, _ = h.Write([]byte(content))
	return h.Sum(nil)
}
