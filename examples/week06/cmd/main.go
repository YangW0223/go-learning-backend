package main

import (
	"context"
	"fmt"

	"github.com/yang/go-learning-backend/examples/week06"
)

func main() {
	decision := week06.ChooseDataAccessTool(week06.ToolDecisionInput{
		NeedTypeSafety: true,
		PreferRawSQL:   true,
	})
	fmt.Println("chosen tool:", decision)

	q := week06.NewSqlcQueries()
	_, _ = q.CreateUser(context.Background(), week06.CreateUserParams{Name: "alice"})
	_, _ = q.CreateUser(context.Background(), week06.CreateUserParams{Name: "allen"})
	_, _ = q.CreateUser(context.Background(), week06.CreateUserParams{Name: "bob"})

	users, _ := q.ListUsersByPrefix(context.Background(), week06.ListUsersByPrefixParams{Prefix: "al", Limit: 10})
	fmt.Println("prefix=al users:", users)
}
