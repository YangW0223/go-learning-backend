package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/yang/go-learning-backend/gin-backend/internal/model"
	"github.com/yang/go-learning-backend/gin-backend/internal/repository"
)

// TodoRepository 是 Postgres Todo 仓储实现。
type TodoRepository struct {
	db *sql.DB
}

// NewTodoRepository 创建 Postgres Todo 仓储。
func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

// Create 插入 Todo 并返回完整记录。
func (r *TodoRepository) Create(ctx context.Context, todo model.Todo) (model.Todo, error) {
	const query = `
INSERT INTO todos (id, user_id, title, done, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, user_id, title, done, created_at, updated_at;
`
	var out model.Todo
	err := r.db.QueryRowContext(ctx, query, todo.ID, todo.UserID, todo.Title, todo.Done, todo.CreatedAt, todo.UpdatedAt).Scan(
		&out.ID,
		&out.UserID,
		&out.Title,
		&out.Done,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		return model.Todo{}, fmt.Errorf("insert todo: %w", err)
	}
	return out, nil
}

// ListByUserID 查询用户的全部 Todo（按创建时间倒序）。
func (r *TodoRepository) ListByUserID(ctx context.Context, userID string) ([]model.Todo, error) {
	const query = `
SELECT id, user_id, title, done, created_at, updated_at
FROM todos
WHERE user_id = $1
ORDER BY created_at DESC;
`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list todos: %w", err)
	}
	defer rows.Close()

	out := make([]model.Todo, 0)
	for rows.Next() {
		var item model.Todo
		if err := rows.Scan(&item.ID, &item.UserID, &item.Title, &item.Done, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan todo row: %w", err)
		}
		out = append(out, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate todo rows: %w", err)
	}
	return out, nil
}

// GetByID 查询单个 Todo。
func (r *TodoRepository) GetByID(ctx context.Context, id string) (model.Todo, error) {
	const query = `SELECT id, user_id, title, done, created_at, updated_at FROM todos WHERE id = $1;`
	var out model.Todo
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&out.ID,
		&out.UserID,
		&out.Title,
		&out.Done,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Todo{}, repository.ErrNotFound
		}
		return model.Todo{}, fmt.Errorf("get todo by id: %w", err)
	}
	return out, nil
}

// Update 更新 Todo 的标题、完成状态与更新时间。
func (r *TodoRepository) Update(ctx context.Context, todo model.Todo) (model.Todo, error) {
	const query = `
UPDATE todos
SET title = $2, done = $3, updated_at = $4
WHERE id = $1
RETURNING id, user_id, title, done, created_at, updated_at;
`
	var out model.Todo
	err := r.db.QueryRowContext(ctx, query, todo.ID, todo.Title, todo.Done, todo.UpdatedAt).Scan(
		&out.ID,
		&out.UserID,
		&out.Title,
		&out.Done,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Todo{}, repository.ErrNotFound
		}
		return model.Todo{}, fmt.Errorf("update todo: %w", err)
	}
	return out, nil
}

// Delete 删除指定 Todo。
func (r *TodoRepository) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM todos WHERE id = $1;`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete todo: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected for delete todo: %w", err)
	}
	if affected == 0 {
		// 没有命中记录时返回统一 not found 错误。
		return repository.ErrNotFound
	}
	return nil
}
