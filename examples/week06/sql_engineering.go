// 详细注释: package week06
package week06

// 详细注释: import (
import (
	// 详细注释: "context"
	"context"
	// 详细注释: "errors"
	"errors"
	// 详细注释: "sort"
	"sort"
	// 详细注释: "strings"
	"strings"
	// 详细注释: )
)

// 详细注释: var (
var (
	// ErrInvalidName 表示用户名不合法。
	// 详细注释: ErrInvalidName = errors.New("invalid name")
	ErrInvalidName = errors.New("invalid name")

// 详细注释: )
)

// ToolDecisionInput 表示选型约束。
// 详细注释: type ToolDecisionInput struct {
type ToolDecisionInput struct {
	// 详细注释: NeedTypeSafety    bool
	NeedTypeSafety bool
	// 详细注释: PreferRawSQL      bool
	PreferRawSQL bool
	// 详细注释: NeedFastPrototype bool
	NeedFastPrototype bool
	// 详细注释: }
}

// ChooseDataAccessTool 根据约束返回 sqlc 或 gorm。
// 详细注释: func ChooseDataAccessTool(in ToolDecisionInput) string {
func ChooseDataAccessTool(in ToolDecisionInput) string {
	// 详细注释: if in.NeedFastPrototype && !in.NeedTypeSafety {
	if in.NeedFastPrototype && !in.NeedTypeSafety {
		// 详细注释: return "gorm"
		return "gorm"
		// 详细注释: }
	}
	// 详细注释: if in.NeedTypeSafety || in.PreferRawSQL {
	if in.NeedTypeSafety || in.PreferRawSQL {
		// 详细注释: return "sqlc"
		return "sqlc"
		// 详细注释: }
	}
	// 详细注释: return "sqlc"
	return "sqlc"
	// 详细注释: }
}

// User 是示例用户模型。
// 详细注释: type User struct {
type User struct {
	// 详细注释: ID   int64
	ID int64
	// 详细注释: Name string
	Name string
	// 详细注释: }
}

// CreateUserParams 模拟 sqlc 生成的参数结构。
// 详细注释: type CreateUserParams struct {
type CreateUserParams struct {
	// 详细注释: Name string
	Name string
	// 详细注释: }
}

// ListUsersByPrefixParams 模拟 sqlc 生成的查询参数。
// 详细注释: type ListUsersByPrefixParams struct {
type ListUsersByPrefixParams struct {
	// 详细注释: Prefix string
	Prefix string
	// 详细注释: Limit  int
	Limit int
	// 详细注释: }
}

// SqlcQueries 模拟 sqlc 生成后的类型安全 API。
// 详细注释: type SqlcQueries struct {
type SqlcQueries struct {
	// 详细注释: nextID int64
	nextID int64
	// 详细注释: users  []User
	users []User
	// 详细注释: }
}

// NewSqlcQueries 创建查询对象。
// 详细注释: func NewSqlcQueries() *SqlcQueries {
func NewSqlcQueries() *SqlcQueries {
	// 详细注释: return &SqlcQueries{nextID: 1}
	return &SqlcQueries{nextID: 1}
	// 详细注释: }
}

// CreateUser 使用强类型参数创建用户。
// 详细注释: func (q *SqlcQueries) CreateUser(_ context.Context, p CreateUserParams) (User, error) {
func (q *SqlcQueries) CreateUser(_ context.Context, p CreateUserParams) (User, error) {
	// 详细注释: name := strings.TrimSpace(p.Name)
	name := strings.TrimSpace(p.Name)
	// 详细注释: if name == "" {
	if name == "" {
		// 详细注释: return User{}, ErrInvalidName
		return User{}, ErrInvalidName
		// 详细注释: }
	}
	// 详细注释: user := User{ID: q.nextID, Name: name}
	user := User{ID: q.nextID, Name: name}
	// 详细注释: q.nextID++
	q.nextID++
	// 详细注释: q.users = append(q.users, user)
	q.users = append(q.users, user)
	// 详细注释: return user, nil
	return user, nil
	// 详细注释: }
}

// ListUsersByPrefix 使用强类型参数查询用户。
// 详细注释: func (q *SqlcQueries) ListUsersByPrefix(_ context.Context, p ListUsersByPrefixParams) ([]User, error) {
func (q *SqlcQueries) ListUsersByPrefix(_ context.Context, p ListUsersByPrefixParams) ([]User, error) {
	// 详细注释: prefix := strings.TrimSpace(p.Prefix)
	prefix := strings.TrimSpace(p.Prefix)
	// 详细注释: limit := p.Limit
	limit := p.Limit
	// 详细注释: if limit <= 0 {
	if limit <= 0 {
		// 详细注释: limit = 10
		limit = 10
		// 详细注释: }
	}

	// 详细注释: matched := make([]User, 0)
	matched := make([]User, 0)
	// 详细注释: for _, u := range q.users {
	for _, u := range q.users {
		// 详细注释: if strings.HasPrefix(strings.ToLower(u.Name), strings.ToLower(prefix)) {
		if strings.HasPrefix(strings.ToLower(u.Name), strings.ToLower(prefix)) {
			// 详细注释: matched = append(matched, u)
			matched = append(matched, u)
			// 详细注释: }
		}
		// 详细注释: }
	}
	// 详细注释: sort.Slice(matched, func(i, j int) bool { return matched[i].ID < matched[j].ID })
	sort.Slice(matched, func(i, j int) bool { return matched[i].ID < matched[j].ID })
	// 详细注释: if len(matched) > limit {
	if len(matched) > limit {
		// 详细注释: matched = matched[:limit]
		matched = matched[:limit]
		// 详细注释: }
	}
	// 详细注释: return matched, nil
	return matched, nil
	// 详细注释: }
}
