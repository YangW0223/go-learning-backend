// 详细注释: package main
package main

// 详细注释: import (
import (
	// 详细注释: "context"
	"context"
	// 详细注释: "fmt"
	"fmt"

	// 详细注释: "github.com/yang/go-learning-backend/examples/week06"
	"github.com/yang/go-learning-backend/examples/week06"
	// 详细注释: )
)

// 详细注释: func main() {
func main() {
	// 详细注释: decision := week06.ChooseDataAccessTool(week06.ToolDecisionInput{
	decision := week06.ChooseDataAccessTool(week06.ToolDecisionInput{
		// 详细注释: NeedTypeSafety: true,
		NeedTypeSafety: true,
		// 详细注释: PreferRawSQL:   true,
		PreferRawSQL: true,
		// 详细注释: })
	})
	// 详细注释: fmt.Println("chosen tool:", decision)
	fmt.Println("chosen tool:", decision)

	// 详细注释: q := week06.NewSqlcQueries()
	q := week06.NewSqlcQueries()
	// 详细注释: _, _ = q.CreateUser(context.Background(), week06.CreateUserParams{Name: "alice"})
	_, _ = q.CreateUser(context.Background(), week06.CreateUserParams{Name: "alice"})
	// 详细注释: _, _ = q.CreateUser(context.Background(), week06.CreateUserParams{Name: "allen"})
	_, _ = q.CreateUser(context.Background(), week06.CreateUserParams{Name: "allen"})
	// 详细注释: _, _ = q.CreateUser(context.Background(), week06.CreateUserParams{Name: "bob"})
	_, _ = q.CreateUser(context.Background(), week06.CreateUserParams{Name: "bob"})

	// 详细注释: users, _ := q.ListUsersByPrefix(context.Background(), week06.ListUsersByPrefixParams{Prefix: "al", Limit: 10})
	users, _ := q.ListUsersByPrefix(context.Background(), week06.ListUsersByPrefixParams{Prefix: "al", Limit: 10})
	// 详细注释: fmt.Println("prefix=al users:", users)
	fmt.Println("prefix=al users:", users)
	// 详细注释: }
}
