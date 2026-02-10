package main

import (
	"errors"
	"fmt"

	"github.com/yang/go-learning-backend/examples/week01"
)

func main() {
	// 1) struct + interface 示例。
	// User 实现了 Greet()，因此满足 Greeter 接口。
	user := week01.User{Name: "  Go Learner  "}
	fmt.Println("Greet:", user.Greet())

	// 2) map 统计示例。
	words := []string{"go", "js", "go", "rust", "go"}
	fmt.Println("CountWords:", week01.CountWords(words))

	// 3) 去重示例（JS 工具函数重写）。
	items := []string{"go", "js", "go", "rust", "js"}
	fmt.Println("UniqueStrings:", week01.UniqueStrings(items))

	// 4) 分组示例（JS 工具函数重写）。
	animals := []string{"ant", "apple", "bear", "boat", ""}
	fmt.Println("GroupByFirstLetter:", week01.GroupByFirstLetter(animals))

	// 5) 错误处理示例（正常路径）。
	n, err := week01.FirstPositive([]int{-3, -1, 2, 5})
	if err != nil {
		fmt.Println("FirstPositive(normal): unexpected error:", err)
	} else {
		fmt.Println("FirstPositive(normal):", n)
	}

	// 6) 错误处理示例（异常路径）。
	_, err = week01.FirstPositive([]int{-3, -1, 0})
	if errors.Is(err, week01.ErrNoPositive) {
		fmt.Println("FirstPositive(error):", week01.ErrNoPositive)
	}
}
