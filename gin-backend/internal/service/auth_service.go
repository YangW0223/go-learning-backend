package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yang/go-learning-backend/gin-backend/internal/auth"
	"github.com/yang/go-learning-backend/gin-backend/internal/errs"
	"github.com/yang/go-learning-backend/gin-backend/internal/model"
	"github.com/yang/go-learning-backend/gin-backend/internal/repository"
)

// AuthService 定义鉴权服务契约。
type AuthService interface {
	Register(ctx context.Context, email, password string) (model.User, error)
	Login(ctx context.Context, email, password string) (string, model.User, error)
	ParseToken(token string) (auth.Claims, error)
	GetProfile(ctx context.Context, userID string) (model.User, error)
}

type authService struct {
	users           repository.UserRepository
	jwt             *auth.JWTManager
	passwordMinSize int
}

// NewAuthService 创建鉴权服务。
// 参数校验策略：
// 1) 核心依赖为空时直接 panic，避免运行时隐式失败；
// 2) 密码最小长度设置下限，避免弱配置。
func NewAuthService(users repository.UserRepository, jwt *auth.JWTManager, passwordMinSize int) AuthService {
	if users == nil || jwt == nil {
		panic("auth service dependencies cannot be nil")
	}
	if passwordMinSize < 6 {
		passwordMinSize = 6
	}
	return &authService{users: users, jwt: jwt, passwordMinSize: passwordMinSize}
}

// Register 负责新用户注册流程：
// 1) 标准化输入；
// 2) 校验邮箱与密码；
// 3) 生成密码哈希；
// 4) 持久化用户并映射冲突错误。
func (s *authService) Register(ctx context.Context, email, password string) (model.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	password = strings.TrimSpace(password)
	if email == "" || !strings.Contains(email, "@") {
		return model.User{}, errs.WithMessage(errs.ErrBadRequest, "invalid email")
	}
	if len(password) < s.passwordMinSize {
		return model.User{}, errs.WithMessage(errs.ErrBadRequest, "password too short")
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return model.User{}, errs.WithMessage(errs.ErrInternal, "failed to hash password")
	}

	created, err := s.users.Create(ctx, model.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: hash,
		Role:         "user",
		CreatedAt:    time.Now().UTC(),
	})
	if err != nil {
		// 对外暴露明确业务语义：邮箱冲突 -> 409。
		if errors.Is(err, repository.ErrConflict) {
			return model.User{}, errs.WithMessage(errs.ErrConflict, "email already exists")
		}
		return model.User{}, errs.WithMessage(errs.ErrInternal, "failed to create user")
	}
	return created, nil
}

// Login 负责用户名密码登录流程：
// 1) 查询用户；
// 2) 校验密码；
// 3) 生成 JWT token。
func (s *authService) Login(ctx context.Context, email, password string) (string, model.User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return "", model.User{}, errs.WithMessage(errs.ErrBadRequest, "email and password are required")
	}

	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		// 出于安全考虑，“用户不存在”和“密码错误”统一返回 invalid credentials。
		if errors.Is(err, repository.ErrNotFound) {
			return "", model.User{}, errs.WithMessage(errs.ErrUnauthorized, "invalid credentials")
		}
		return "", model.User{}, errs.WithMessage(errs.ErrInternal, "failed to query user")
	}
	if err := auth.ComparePassword(user.PasswordHash, password); err != nil {
		return "", model.User{}, errs.WithMessage(errs.ErrUnauthorized, "invalid credentials")
	}

	token, err := s.jwt.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		return "", model.User{}, errs.WithMessage(errs.ErrInternal, "failed to generate token")
	}
	return token, user, nil
}

// ParseToken 调用 JWT 管理器解析 token，并将错误映射成统一 401 语义。
func (s *authService) ParseToken(token string) (auth.Claims, error) {
	claims, err := s.jwt.Parse(token)
	if err != nil {
		return auth.Claims{}, errs.WithMessage(errs.ErrUnauthorized, "invalid token")
	}
	return claims, nil
}

// GetProfile 根据 userID 查询用户资料。
func (s *authService) GetProfile(ctx context.Context, userID string) (model.User, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.User{}, errs.WithMessage(errs.ErrNotFound, "user not found")
		}
		return model.User{}, errs.WithMessage(errs.ErrInternal, "failed to query user")
	}
	return user, nil
}
