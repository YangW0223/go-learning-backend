package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/yang/go-learning-backend/gin-backend/internal/model"
	"github.com/yang/go-learning-backend/gin-backend/internal/repository"
)

// UserRepository 是 Postgres 用户仓储实现。
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository 创建 Postgres 用户仓储。
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create 写入用户记录。
func (r *UserRepository) Create(ctx context.Context, user model.User) (model.User, error) {
	const query = `
INSERT INTO users (id, email, password_hash, role, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, email, password_hash, role, created_at;
`
	var out model.User
	err := r.db.QueryRowContext(ctx, query, user.ID, user.Email, user.PasswordHash, user.Role, user.CreatedAt).Scan(
		&out.ID,
		&out.Email,
		&out.PasswordHash,
		&out.Role,
		&out.CreatedAt,
	)
	if err != nil {
		// 唯一索引冲突映射成仓储层通用冲突错误。
		if isUniqueViolation(err) {
			return model.User{}, repository.ErrConflict
		}
		return model.User{}, fmt.Errorf("insert user: %w", err)
	}
	return out, nil
}

// GetByEmail 按邮箱查询用户。
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	const query = `SELECT id, email, password_hash, role, created_at FROM users WHERE email = $1;`
	var out model.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&out.ID,
		&out.Email,
		&out.PasswordHash,
		&out.Role,
		&out.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, repository.ErrNotFound
		}
		return model.User{}, fmt.Errorf("get user by email: %w", err)
	}
	return out, nil
}

// GetByID 按用户 ID 查询用户。
func (r *UserRepository) GetByID(ctx context.Context, id string) (model.User, error) {
	const query = `SELECT id, email, password_hash, role, created_at FROM users WHERE id = $1;`
	var out model.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&out.ID,
		&out.Email,
		&out.PasswordHash,
		&out.Role,
		&out.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, repository.ErrNotFound
		}
		return model.User{}, fmt.Errorf("get user by id: %w", err)
	}
	return out, nil
}

// isUniqueViolation 使用字符串匹配识别唯一约束冲突。
// 在生产场景可进一步基于 pq.Error.Code 做更准确判定。
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate key") || strings.Contains(msg, "unique constraint")
}
