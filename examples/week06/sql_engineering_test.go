// 详细注释: package week06
package week06

// 详细注释: import (
import (
	// 详细注释: "context"
	"context"
	// 详细注释: "errors"
	"errors"
	// 详细注释: "testing"
	"testing"
	// 详细注释: )
)

// TestChooseDataAccessTool 验证选型规则。
// 详细注释: func TestChooseDataAccessTool(t *testing.T) {
func TestChooseDataAccessTool(t *testing.T) {
	// 详细注释: cases := []struct {
	cases := []struct {
		// 详细注释: name string
		name string
		// 详细注释: in   ToolDecisionInput
		in ToolDecisionInput
		// 详细注释: want string
		want string
		// 详细注释: }{
	}{
		// 详细注释: {name: "type safety first", in: ToolDecisionInput{NeedTypeSafety: true}, want: "sqlc"},
		{name: "type safety first", in: ToolDecisionInput{NeedTypeSafety: true}, want: "sqlc"},
		// 详细注释: {name: "fast prototype", in: ToolDecisionInput{NeedFastPrototype: true}, want: "gorm"},
		{name: "fast prototype", in: ToolDecisionInput{NeedFastPrototype: true}, want: "gorm"},
		// 详细注释: {name: "prefer raw sql", in: ToolDecisionInput{PreferRawSQL: true}, want: "sqlc"},
		{name: "prefer raw sql", in: ToolDecisionInput{PreferRawSQL: true}, want: "sqlc"},
		// 详细注释: }
	}

	// 详细注释: for _, tc := range cases {
	for _, tc := range cases {
		// 详细注释: got := ChooseDataAccessTool(tc.in)
		got := ChooseDataAccessTool(tc.in)
		// 详细注释: if got != tc.want {
		if got != tc.want {
			// 详细注释: t.Fatalf("%s: want %s got %s", tc.name, tc.want, got)
			t.Fatalf("%s: want %s got %s", tc.name, tc.want, got)
			// 详细注释: }
		}
		// 详细注释: }
	}
	// 详细注释: }
}

// TestSqlcQueriesCreateAndList 验证强类型查询接口。
// 详细注释: func TestSqlcQueriesCreateAndList(t *testing.T) {
func TestSqlcQueriesCreateAndList(t *testing.T) {
	// 详细注释: q := NewSqlcQueries()
	q := NewSqlcQueries()
	// 详细注释: _, _ = q.CreateUser(context.Background(), CreateUserParams{Name: "alice"})
	_, _ = q.CreateUser(context.Background(), CreateUserParams{Name: "alice"})
	// 详细注释: _, _ = q.CreateUser(context.Background(), CreateUserParams{Name: "bob"})
	_, _ = q.CreateUser(context.Background(), CreateUserParams{Name: "bob"})

	// 详细注释: got, err := q.ListUsersByPrefix(context.Background(), ListUsersByPrefixParams{Prefix: "a", Limit: 5})
	got, err := q.ListUsersByPrefix(context.Background(), ListUsersByPrefixParams{Prefix: "a", Limit: 5})
	// 详细注释: if err != nil {
	if err != nil {
		// 详细注释: t.Fatalf("unexpected err: %v", err)
		t.Fatalf("unexpected err: %v", err)
		// 详细注释: }
	}
	// 详细注释: if len(got) != 1 || got[0].Name != "alice" {
	if len(got) != 1 || got[0].Name != "alice" {
		// 详细注释: t.Fatalf("unexpected list result: %+v", got)
		t.Fatalf("unexpected list result: %+v", got)
		// 详细注释: }
	}
	// 详细注释: }
}

// TestSqlcQueriesInvalidName 验证参数非法分支。
// 详细注释: func TestSqlcQueriesInvalidName(t *testing.T) {
func TestSqlcQueriesInvalidName(t *testing.T) {
	// 详细注释: q := NewSqlcQueries()
	q := NewSqlcQueries()
	// 详细注释: _, err := q.CreateUser(context.Background(), CreateUserParams{Name: "   "})
	_, err := q.CreateUser(context.Background(), CreateUserParams{Name: "   "})
	// 详细注释: if !errors.Is(err, ErrInvalidName) {
	if !errors.Is(err, ErrInvalidName) {
		// 详细注释: t.Fatalf("expected ErrInvalidName, got %v", err)
		t.Fatalf("expected ErrInvalidName, got %v", err)
		// 详细注释: }
	}
	// 详细注释: }
}
