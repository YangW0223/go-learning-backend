package main

import (
	"errors"
	"fmt"

	"github.com/yang/go-learning-backend/examples/week02"
)

func main() {
	// 1) 演示合法删除路径解析。
	validPath := "/api/v1/todos/20260210112233.123456789"
	id, err := week02.ParseDeleteTodoPath(validPath)
	if err != nil {
		fmt.Println("Parse valid path failed:", err)
		return
	}
	fmt.Println("Parsed ID:", id)

	// 2) 演示构造统一成功响应 JSON。
	successJSON, err := week02.BuildSuccessJSON(week02.DeleteResult{
		ID:      id,
		Deleted: true,
	})
	if err != nil {
		fmt.Println("Build success JSON failed:", err)
		return
	}
	fmt.Println("Success JSON:", string(successJSON))

	// 3) 演示非法路径时的错误分类。
	invalidPath := "/api/v1/todos/abc"
	_, err = week02.ParseDeleteTodoPath(invalidPath)
	if errors.Is(err, week02.ErrInvalidTodoID) {
		errorJSON, _ := week02.BuildErrorJSON("invalid todo id")
		fmt.Println("Error JSON:", string(errorJSON))
	}
}
