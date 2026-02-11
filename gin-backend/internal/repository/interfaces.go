package repository

import (
	"context"
	"time"

	"github.com/yang/go-learning-backend/gin-backend/internal/model"
)

// UserRepository 定义用户存储契约。
// service 层仅依赖该接口，便于替换 Postgres 实现或注入 mock。
type UserRepository interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	GetByID(ctx context.Context, id string) (model.User, error)
}

// TodoRepository 定义 Todo 存储契约。
// 该接口描述 Todo 的持久化能力，不包含缓存语义。
type TodoRepository interface {
	Create(ctx context.Context, todo model.Todo) (model.Todo, error)
	ListByUserID(ctx context.Context, userID string) ([]model.Todo, error)
	GetByID(ctx context.Context, id string) (model.Todo, error)
	Update(ctx context.Context, todo model.Todo) (model.Todo, error)
	Delete(ctx context.Context, id string) error
}

// TodoCache 定义 Todo 列表缓存契约。
// 约定：GetList 返回 (nil,false,nil) 表示缓存未命中。
type TodoCache interface {
	GetList(ctx context.Context, userID string) ([]model.Todo, bool, error)
	SetList(ctx context.Context, userID string, todos []model.Todo, ttl time.Duration) error
	DeleteList(ctx context.Context, userID string) error
	Ping(ctx context.Context) error
}
