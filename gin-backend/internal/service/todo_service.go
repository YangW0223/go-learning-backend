package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yang/go-learning-backend/gin-backend/internal/errs"
	"github.com/yang/go-learning-backend/gin-backend/internal/model"
	"github.com/yang/go-learning-backend/gin-backend/internal/repository"
)

// UpdateTodoInput 定义可更新字段。
type UpdateTodoInput struct {
	Title *string
	Done  *bool
}

// TodoService 定义 Todo 业务契约。
type TodoService interface {
	Create(ctx context.Context, userID, title string) (model.Todo, error)
	List(ctx context.Context, userID string) ([]model.Todo, error)
	Update(ctx context.Context, userID, id string, in UpdateTodoInput) (model.Todo, error)
	Delete(ctx context.Context, userID, id string) error
}

type todoService struct {
	repo     repository.TodoRepository
	cache    repository.TodoCache
	cacheTTL time.Duration
}

// NewTodoService 创建 Todo 业务服务。
// 设计点：
// 1) repo/cache 为空直接 panic，尽早暴露装配错误；
// 2) cacheTTL 非法时回退默认值，避免写入 0 秒缓存。
func NewTodoService(repo repository.TodoRepository, cache repository.TodoCache, cacheTTL time.Duration) TodoService {
	if repo == nil || cache == nil {
		panic("todo service dependencies cannot be nil")
	}
	if cacheTTL <= 0 {
		cacheTTL = 30 * time.Second
	}
	return &todoService{repo: repo, cache: cache, cacheTTL: cacheTTL}
}

// Create 创建 Todo 并触发缓存失效。
func (s *todoService) Create(ctx context.Context, userID, title string) (model.Todo, error) {
	title = strings.TrimSpace(title)
	if userID == "" {
		return model.Todo{}, errs.WithMessage(errs.ErrUnauthorized, "missing user identity")
	}
	if title == "" {
		return model.Todo{}, errs.WithMessage(errs.ErrBadRequest, "title is required")
	}
	now := time.Now().UTC()
	created, err := s.repo.Create(ctx, model.Todo{
		ID:        uuid.NewString(),
		UserID:    userID,
		Title:     title,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return model.Todo{}, errs.WithMessage(errs.ErrInternal, "failed to create todo")
	}
	// 写操作后删除列表缓存，避免读取脏数据。
	_ = s.cache.DeleteList(ctx, userID)
	return created, nil
}

// List 先读缓存，未命中再回源数据库，并回填缓存。
func (s *todoService) List(ctx context.Context, userID string) ([]model.Todo, error) {
	if userID == "" {
		return nil, errs.WithMessage(errs.ErrUnauthorized, "missing user identity")
	}
	cached, hit, err := s.cache.GetList(ctx, userID)
	if err == nil && hit {
		return cached, nil
	}
	items, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, errs.WithMessage(errs.ErrInternal, "failed to list todos")
	}
	_ = s.cache.SetList(ctx, userID, items, s.cacheTTL)
	return items, nil
}

// Update 更新 Todo（支持更新标题和完成状态）。
// 关键规则：
// 1) 必须是资源拥有者；
// 2) 至少更新一个字段；
// 3) 更新成功后失效缓存。
func (s *todoService) Update(ctx context.Context, userID, id string, in UpdateTodoInput) (model.Todo, error) {
	if userID == "" {
		return model.Todo{}, errs.WithMessage(errs.ErrUnauthorized, "missing user identity")
	}
	if strings.TrimSpace(id) == "" {
		return model.Todo{}, errs.WithMessage(errs.ErrBadRequest, "todo id is required")
	}
	if in.Title == nil && in.Done == nil {
		return model.Todo{}, errs.WithMessage(errs.ErrBadRequest, "at least one field is required")
	}

	todo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Todo{}, errs.WithMessage(errs.ErrNotFound, "todo not found")
		}
		return model.Todo{}, errs.WithMessage(errs.ErrInternal, "failed to query todo")
	}
	if todo.UserID != userID {
		return model.Todo{}, errs.WithMessage(errs.ErrForbidden, "cannot modify this todo")
	}

	if in.Title != nil {
		trimmed := strings.TrimSpace(*in.Title)
		if trimmed == "" {
			return model.Todo{}, errs.WithMessage(errs.ErrBadRequest, "title is required")
		}
		todo.Title = trimmed
	}
	if in.Done != nil {
		todo.Done = *in.Done
	}
	todo.UpdatedAt = time.Now().UTC()

	updated, err := s.repo.Update(ctx, todo)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Todo{}, errs.WithMessage(errs.ErrNotFound, "todo not found")
		}
		return model.Todo{}, errs.WithMessage(errs.ErrInternal, "failed to update todo")
	}
	_ = s.cache.DeleteList(ctx, userID)
	return updated, nil
}

// Delete 删除 Todo。
// 删除前先读取资源做所有权校验，防止跨用户删除。
func (s *todoService) Delete(ctx context.Context, userID, id string) error {
	if userID == "" {
		return errs.WithMessage(errs.ErrUnauthorized, "missing user identity")
	}
	if strings.TrimSpace(id) == "" {
		return errs.WithMessage(errs.ErrBadRequest, "todo id is required")
	}

	todo, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return errs.WithMessage(errs.ErrNotFound, "todo not found")
		}
		return errs.WithMessage(errs.ErrInternal, "failed to query todo")
	}
	if todo.UserID != userID {
		return errs.WithMessage(errs.ErrForbidden, "cannot delete this todo")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return errs.WithMessage(errs.ErrNotFound, "todo not found")
		}
		return errs.WithMessage(errs.ErrInternal, "failed to delete todo")
	}
	// 删除后失效缓存，保证列表读取一致性。
	_ = s.cache.DeleteList(ctx, userID)
	return nil
}
