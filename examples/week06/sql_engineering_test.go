package week06

import (
	"context"
	"errors"
	"testing"
)

// TestChooseDataAccessTool 验证选型规则。
func TestChooseDataAccessTool(t *testing.T) {
	cases := []struct {
		name string
		in   ToolDecisionInput
		want string
	}{
		{name: "type safety first", in: ToolDecisionInput{NeedTypeSafety: true}, want: "sqlc"},
		{name: "fast prototype", in: ToolDecisionInput{NeedFastPrototype: true}, want: "gorm"},
		{name: "prefer raw sql", in: ToolDecisionInput{PreferRawSQL: true}, want: "sqlc"},
	}

	for _, tc := range cases {
		got := ChooseDataAccessTool(tc.in)
		if got != tc.want {
			t.Fatalf("%s: want %s got %s", tc.name, tc.want, got)
		}
	}
}

// TestSqlcQueriesCreateAndList 验证强类型查询接口。
func TestSqlcQueriesCreateAndList(t *testing.T) {
	q := NewSqlcQueries()
	_, _ = q.CreateUser(context.Background(), CreateUserParams{Name: "alice"})
	_, _ = q.CreateUser(context.Background(), CreateUserParams{Name: "bob"})

	got, err := q.ListUsersByPrefix(context.Background(), ListUsersByPrefixParams{Prefix: "a", Limit: 5})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(got) != 1 || got[0].Name != "alice" {
		t.Fatalf("unexpected list result: %+v", got)
	}
}

// TestSqlcQueriesInvalidName 验证参数非法分支。
func TestSqlcQueriesInvalidName(t *testing.T) {
	q := NewSqlcQueries()
	_, err := q.CreateUser(context.Background(), CreateUserParams{Name: "   "})
	if !errors.Is(err, ErrInvalidName) {
		t.Fatalf("expected ErrInvalidName, got %v", err)
	}
}
