package main

import (
	"context"
	"fmt"

	"github.com/yang/go-learning-backend/examples/week05"
)

func main() {
	db := week05.NewInMemoryPostgres()

	_ = db.WithTx(context.Background(), func(tx *week05.Tx) error {
		_, _ = tx.CreateTodo("design todos table")
		_, _ = tx.CreateTodo("implement mark done")
		_, _ = tx.MarkDone("1")
		return nil
	})

	items, err := db.ListTodos(1, 10)
	if err != nil {
		fmt.Println("list error:", err)
		return
	}
	for _, item := range items {
		fmt.Printf("id=%s title=%q done=%v\n", item.ID, item.Title, item.Done)
	}
}
