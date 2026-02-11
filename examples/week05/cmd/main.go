// 详细注释: package main
package main

// 详细注释: import (
import (
	// 详细注释: "context"
	"context"
	// 详细注释: "fmt"
	"fmt"

	// 详细注释: "github.com/yang/go-learning-backend/examples/week05"
	"github.com/yang/go-learning-backend/examples/week05"
	// 详细注释: )
)

// 详细注释: func main() {
func main() {
	// 详细注释: db := week05.NewInMemoryPostgres()
	db := week05.NewInMemoryPostgres()

	// 详细注释: _ = db.WithTx(context.Background(), func(tx *week05.Tx) error {
	_ = db.WithTx(context.Background(), func(tx *week05.Tx) error {
		// 详细注释: _, _ = tx.CreateTodo("design todos table")
		_, _ = tx.CreateTodo("design todos table")
		// 详细注释: _, _ = tx.CreateTodo("implement mark done")
		_, _ = tx.CreateTodo("implement mark done")
		// 详细注释: _, _ = tx.MarkDone("1")
		_, _ = tx.MarkDone("1")
		// 详细注释: return nil
		return nil
		// 详细注释: })
	})

	// 详细注释: items, err := db.ListTodos(1, 10)
	items, err := db.ListTodos(1, 10)
	// 详细注释: if err != nil {
	if err != nil {
		// 详细注释: fmt.Println("list error:", err)
		fmt.Println("list error:", err)
		// 详细注释: return
		return
		// 详细注释: }
	}
	// 详细注释: for _, item := range items {
	for _, item := range items {
		// 详细注释: fmt.Printf("id=%s title=%q done=%v\n", item.ID, item.Title, item.Done)
		fmt.Printf("id=%s title=%q done=%v\n", item.ID, item.Title, item.Done)
		// 详细注释: }
	}
	// 详细注释: }
}
