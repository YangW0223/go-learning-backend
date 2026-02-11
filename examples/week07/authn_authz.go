// 详细注释: package week07
package week07

// 详细注释: import (
import (
	// 详细注释: "crypto/hmac"
	"crypto/hmac"
	// 详细注释: "crypto/sha256"
	"crypto/sha256"
	// 详细注释: "encoding/base64"
	"encoding/base64"
	// 详细注释: "encoding/json"
	"encoding/json"
	// 详细注释: "errors"
	"errors"
	// 详细注释: "fmt"
	"fmt"
	// 详细注释: "net/http"
	"net/http"
	// 详细注释: "strings"
	"strings"
	// 详细注释: "sync"
	"sync"
	// 详细注释: "time"
	"time"
	// 详细注释: )
)

// 详细注释: var (
var (
	// ErrUserExists 表示用户已存在。
	// 详细注释: ErrUserExists = errors.New("user exists")
	ErrUserExists = errors.New("user exists")
	// ErrInvalidCredential 表示用户名或密码错误。
	// 详细注释: ErrInvalidCredential = errors.New("invalid credential")
	ErrInvalidCredential = errors.New("invalid credential")
	// ErrInvalidToken 表示 token 非法。
	// 详细注释: ErrInvalidToken = errors.New("invalid token")
	ErrInvalidToken = errors.New("invalid token")
	// ErrTokenExpired 表示 token 已过期。
	// 详细注释: ErrTokenExpired = errors.New("token expired")
	ErrTokenExpired = errors.New("token expired")
	// ErrUnauthorized 表示未认证。
	// 详细注释: ErrUnauthorized = errors.New("unauthorized")
	ErrUnauthorized = errors.New("unauthorized")
	// ErrForbidden 表示已认证但无权限。
	// 详细注释: ErrForbidden = errors.New("forbidden")
	ErrForbidden = errors.New("forbidden")

// 详细注释: )
)

// User 表示用户实体。
// 详细注释: type User struct {
type User struct {
	// 详细注释: ID           string
	ID string
	// 详细注释: Username     string
	Username string
	// 详细注释: PasswordHash string
	PasswordHash string
	// 详细注释: Role         string
	Role string
	// 详细注释: }
}

// Claims 表示 token 载荷。
// 详细注释: type Claims struct {
type Claims struct {
	// 详细注释: Sub string `json:"sub"`
	Sub string `json:"sub"`
	// 详细注释: Rol string `json:"rol"`
	Rol string `json:"rol"`
	// 详细注释: Exp int64  `json:"exp"`
	Exp int64 `json:"exp"`
	// 详细注释: }
}

// AuthService 提供注册、登录、鉴权能力。
// 详细注释: type AuthService struct {
type AuthService struct {
	// 详细注释: mu       sync.Mutex
	mu sync.Mutex
	// 详细注释: nextID   int
	nextID int
	// 详细注释: users    map[string]User
	users map[string]User
	// 详细注释: secret   []byte
	secret []byte
	// 详细注释: nowFn    func() time.Time
	nowFn func() time.Time
	// 详细注释: tokenTTL time.Duration
	tokenTTL time.Duration
	// 详细注释: }
}

// NewAuthService 创建认证服务。
// 详细注释: func NewAuthService(secret string, tokenTTL time.Duration) *AuthService {
func NewAuthService(secret string, tokenTTL time.Duration) *AuthService {
	// 详细注释: return &AuthService{
	return &AuthService{
		// 详细注释: nextID:   1,
		nextID: 1,
		// 详细注释: users:    make(map[string]User),
		users: make(map[string]User),
		// 详细注释: secret:   []byte(secret),
		secret: []byte(secret),
		// 详细注释: nowFn:    time.Now,
		nowFn: time.Now,
		// 详细注释: tokenTTL: tokenTTL,
		tokenTTL: tokenTTL,
		// 详细注释: }
	}
	// 详细注释: }
}

// Register 注册用户并保存密码摘要。
// 详细注释: func (s *AuthService) Register(username, password, role string) error {
func (s *AuthService) Register(username, password, role string) error {
	// 详细注释: username = strings.TrimSpace(username)
	username = strings.TrimSpace(username)
	// 详细注释: password = strings.TrimSpace(password)
	password = strings.TrimSpace(password)
	// 详细注释: role = strings.TrimSpace(role)
	role = strings.TrimSpace(role)
	// 详细注释: if username == "" || password == "" || role == "" {
	if username == "" || password == "" || role == "" {
		// 详细注释: return ErrInvalidCredential
		return ErrInvalidCredential
		// 详细注释: }
	}

	// 详细注释: s.mu.Lock()
	s.mu.Lock()
	// 详细注释: defer s.mu.Unlock()
	defer s.mu.Unlock()
	// 详细注释: if _, exists := s.users[username]; exists {
	if _, exists := s.users[username]; exists {
		// 详细注释: return ErrUserExists
		return ErrUserExists
		// 详细注释: }
	}

	// 详细注释: user := User{
	user := User{
		// 详细注释: ID:           fmt.Sprintf("u-%d", s.nextID),
		ID: fmt.Sprintf("u-%d", s.nextID),
		// 详细注释: Username:     username,
		Username: username,
		// 详细注释: PasswordHash: s.hashPassword(password),
		PasswordHash: s.hashPassword(password),
		// 详细注释: Role:         role,
		Role: role,
		// 详细注释: }
	}
	// 详细注释: s.nextID++
	s.nextID++
	// 详细注释: s.users[username] = user
	s.users[username] = user
	// 详细注释: return nil
	return nil
	// 详细注释: }
}

// Login 校验凭证并签发 token。
// 详细注释: func (s *AuthService) Login(username, password string) (string, error) {
func (s *AuthService) Login(username, password string) (string, error) {
	// 详细注释: username = strings.TrimSpace(username)
	username = strings.TrimSpace(username)
	// 详细注释: password = strings.TrimSpace(password)
	password = strings.TrimSpace(password)

	// 详细注释: s.mu.Lock()
	s.mu.Lock()
	// 详细注释: user, ok := s.users[username]
	user, ok := s.users[username]
	// 详细注释: s.mu.Unlock()
	s.mu.Unlock()
	// 详细注释: if !ok {
	if !ok {
		// 详细注释: return "", ErrInvalidCredential
		return "", ErrInvalidCredential
		// 详细注释: }
	}
	// 详细注释: if s.hashPassword(password) != user.PasswordHash {
	if s.hashPassword(password) != user.PasswordHash {
		// 详细注释: return "", ErrInvalidCredential
		return "", ErrInvalidCredential
		// 详细注释: }
	}

	// 详细注释: claims := Claims{
	claims := Claims{
		// 详细注释: Sub: user.ID,
		Sub: user.ID,
		// 详细注释: Rol: user.Role,
		Rol: user.Role,
		// 详细注释: Exp: s.nowFn().Add(s.tokenTTL).Unix(),
		Exp: s.nowFn().Add(s.tokenTTL).Unix(),
		// 详细注释: }
	}
	// 详细注释: return s.signClaims(claims)
	return s.signClaims(claims)
	// 详细注释: }
}

// ValidateToken 校验 token 并返回 claims。
// 详细注释: func (s *AuthService) ValidateToken(token string) (Claims, error) {
func (s *AuthService) ValidateToken(token string) (Claims, error) {
	// 详细注释: parts := strings.Split(token, ".")
	parts := strings.Split(token, ".")
	// 详细注释: if len(parts) != 3 {
	if len(parts) != 3 {
		// 详细注释: return Claims{}, ErrInvalidToken
		return Claims{}, ErrInvalidToken
		// 详细注释: }
	}

	// 详细注释: headerPayload := parts[0] + "." + parts[1]
	headerPayload := parts[0] + "." + parts[1]
	// 详细注释: sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	// 详细注释: if err != nil {
	if err != nil {
		// 详细注释: return Claims{}, ErrInvalidToken
		return Claims{}, ErrInvalidToken
		// 详细注释: }
	}
	// 详细注释: if !hmac.Equal(sig, s.hmac(headerPayload)) {
	if !hmac.Equal(sig, s.hmac(headerPayload)) {
		// 详细注释: return Claims{}, ErrInvalidToken
		return Claims{}, ErrInvalidToken
		// 详细注释: }
	}

	// 详细注释: payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	// 详细注释: if err != nil {
	if err != nil {
		// 详细注释: return Claims{}, ErrInvalidToken
		return Claims{}, ErrInvalidToken
		// 详细注释: }
	}
	// 详细注释: var claims Claims
	var claims Claims
	// 详细注释: if err := json.Unmarshal(payloadBytes, &claims); err != nil {
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		// 详细注释: return Claims{}, ErrInvalidToken
		return Claims{}, ErrInvalidToken
		// 详细注释: }
	}
	// 详细注释: if s.nowFn().Unix() > claims.Exp {
	if s.nowFn().Unix() > claims.Exp {
		// 详细注释: return Claims{}, ErrTokenExpired
		return Claims{}, ErrTokenExpired
		// 详细注释: }
	}
	// 详细注释: return claims, nil
	return claims, nil
	// 详细注释: }
}

// RequireRole 校验 token 和角色权限。
// 详细注释: func (s *AuthService) RequireRole(token, role string) (Claims, error) {
func (s *AuthService) RequireRole(token, role string) (Claims, error) {
	// 详细注释: claims, err := s.ValidateToken(token)
	claims, err := s.ValidateToken(token)
	// 详细注释: if err != nil {
	if err != nil {
		// 详细注释: return Claims{}, ErrUnauthorized
		return Claims{}, ErrUnauthorized
		// 详细注释: }
	}
	// 详细注释: if claims.Rol != role {
	if claims.Rol != role {
		// 详细注释: return Claims{}, ErrForbidden
		return Claims{}, ErrForbidden
		// 详细注释: }
	}
	// 详细注释: return claims, nil
	return claims, nil
	// 详细注释: }
}

// AuthMiddleware 提供最小 HTTP 鉴权中间件。
// 详细注释: func AuthMiddleware(auth *AuthService, requiredRole string, next http.HandlerFunc) http.HandlerFunc {
func AuthMiddleware(auth *AuthService, requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	// 详细注释: return func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// 详细注释: header := strings.TrimSpace(r.Header.Get("Authorization"))
		header := strings.TrimSpace(r.Header.Get("Authorization"))
		// 详细注释: if !strings.HasPrefix(header, "Bearer ") {
		if !strings.HasPrefix(header, "Bearer ") {
			// 详细注释: http.Error(w, ErrUnauthorized.Error(), http.StatusUnauthorized)
			http.Error(w, ErrUnauthorized.Error(), http.StatusUnauthorized)
			// 详细注释: return
			return
			// 详细注释: }
		}
		// 详细注释: token := strings.TrimPrefix(header, "Bearer ")
		token := strings.TrimPrefix(header, "Bearer ")
		// 详细注释: _, err := auth.RequireRole(token, requiredRole)
		_, err := auth.RequireRole(token, requiredRole)
		// 详细注释: if err != nil {
		if err != nil {
			// 详细注释: if errors.Is(err, ErrForbidden) {
			if errors.Is(err, ErrForbidden) {
				// 详细注释: http.Error(w, ErrForbidden.Error(), http.StatusForbidden)
				http.Error(w, ErrForbidden.Error(), http.StatusForbidden)
				// 详细注释: return
				return
				// 详细注释: }
			}
			// 详细注释: http.Error(w, ErrUnauthorized.Error(), http.StatusUnauthorized)
			http.Error(w, ErrUnauthorized.Error(), http.StatusUnauthorized)
			// 详细注释: return
			return
			// 详细注释: }
		}
		// 详细注释: next(w, r)
		next(w, r)
		// 详细注释: }
	}
	// 详细注释: }
}

// 详细注释: func (s *AuthService) hashPassword(password string) string {
func (s *AuthService) hashPassword(password string) string {
	// 示例使用 sha256 + secret，真实生产建议使用 bcrypt/argon2。
	// 详细注释: sum := sha256.Sum256(append(s.secret, []byte(password)...))
	sum := sha256.Sum256(append(s.secret, []byte(password)...))
	// 详细注释: return base64.RawURLEncoding.EncodeToString(sum[:])
	return base64.RawURLEncoding.EncodeToString(sum[:])
	// 详细注释: }
}

// 详细注释: func (s *AuthService) signClaims(claims Claims) (string, error) {
func (s *AuthService) signClaims(claims Claims) (string, error) {
	// 详细注释: headerBytes, _ := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	headerBytes, _ := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	// 详细注释: payloadBytes, err := json.Marshal(claims)
	payloadBytes, err := json.Marshal(claims)
	// 详细注释: if err != nil {
	if err != nil {
		// 详细注释: return "", err
		return "", err
		// 详细注释: }
	}
	// 详细注释: header := base64.RawURLEncoding.EncodeToString(headerBytes)
	header := base64.RawURLEncoding.EncodeToString(headerBytes)
	// 详细注释: payload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	// 详细注释: headerPayload := header + "." + payload
	headerPayload := header + "." + payload
	// 详细注释: signature := base64.RawURLEncoding.EncodeToString(s.hmac(headerPayload))
	signature := base64.RawURLEncoding.EncodeToString(s.hmac(headerPayload))
	// 详细注释: return headerPayload + "." + signature, nil
	return headerPayload + "." + signature, nil
	// 详细注释: }
}

// 详细注释: func (s *AuthService) hmac(content string) []byte {
func (s *AuthService) hmac(content string) []byte {
	// 详细注释: h := hmac.New(sha256.New, s.secret)
	h := hmac.New(sha256.New, s.secret)
	// 详细注释: _, _ = h.Write([]byte(content))
	_, _ = h.Write([]byte(content))
	// 详细注释: return h.Sum(nil)
	return h.Sum(nil)
	// 详细注释: }
}
