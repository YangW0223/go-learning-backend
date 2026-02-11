package week06

import (
	"context"
	"errors"
	"sort"
	"strings"
)

var (
	// ErrInvalidName 表示用户名不合法。
	ErrInvalidName = errors.New("invalid name")
)

// ToolDecisionInput 表示选型约束。
type ToolDecisionInput struct {
	NeedTypeSafety    bool
	PreferRawSQL      bool
	NeedFastPrototype bool
}

// ChooseDataAccessTool 根据约束返回 sqlc 或 gorm。
func ChooseDataAccessTool(in ToolDecisionInput) string {
	if in.NeedFastPrototype && !in.NeedTypeSafety {
		return "gorm"
	}
	if in.NeedTypeSafety || in.PreferRawSQL {
		return "sqlc"
	}
	return "sqlc"
}

// User 是示例用户模型。
type User struct {
	ID   int64
	Name string
}

// CreateUserParams 模拟 sqlc 生成的参数结构。
type CreateUserParams struct {
	Name string
}

// ListUsersByPrefixParams 模拟 sqlc 生成的查询参数。
type ListUsersByPrefixParams struct {
	Prefix string
	Limit  int
}

// SqlcQueries 模拟 sqlc 生成后的类型安全 API。
type SqlcQueries struct {
	nextID int64
	users  []User
}

// NewSqlcQueries 创建查询对象。
func NewSqlcQueries() *SqlcQueries {
	return &SqlcQueries{nextID: 1}
}

// CreateUser 使用强类型参数创建用户。
func (q *SqlcQueries) CreateUser(_ context.Context, p CreateUserParams) (User, error) {
	name := strings.TrimSpace(p.Name)
	if name == "" {
		return User{}, ErrInvalidName
	}
	user := User{ID: q.nextID, Name: name}
	q.nextID++
	q.users = append(q.users, user)
	return user, nil
}

// ListUsersByPrefix 使用强类型参数查询用户。
func (q *SqlcQueries) ListUsersByPrefix(_ context.Context, p ListUsersByPrefixParams) ([]User, error) {
	prefix := strings.TrimSpace(p.Prefix)
	limit := p.Limit
	if limit <= 0 {
		limit = 10
	}

	matched := make([]User, 0)
	for _, u := range q.users {
		if strings.HasPrefix(strings.ToLower(u.Name), strings.ToLower(prefix)) {
			matched = append(matched, u)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].ID < matched[j].ID })
	if len(matched) > limit {
		matched = matched[:limit]
	}
	return matched, nil
}
