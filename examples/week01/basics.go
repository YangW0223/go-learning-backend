package week01

import (
	"errors"
	"strings"
)

// ErrNoPositive 是一个哨兵错误（sentinel error）。
// 当 FirstPositive 在输入切片中找不到任何正数时，会返回该错误。
// 测试里可以用 errors.Is(err, ErrNoPositive) 做稳定判断。
var ErrNoPositive = errors.New("no positive number found")

// User 用于演示最小 struct 定义。
type User struct {
	Name string
}

// Greeter 用于演示最小 interface 定义。
// 任何实现了 Greet() string 的类型都满足该接口。
type Greeter interface {
	Greet() string
}

// Greet 是 User 对 Greeter 接口的实现。
// strings.TrimSpace 用于避免名字前后空格影响输出。
func (u User) Greet() string {
	return "hello, " + strings.TrimSpace(u.Name)
}

// CountWords 统计每个单词出现次数。
// 输入:  ["go", "js", "go"]
// 输出:  {"go": 2, "js": 1}
//
// make(map[string]int, len(words)) 的第二个参数是初始容量，
// 不是长度，目的是减少 map 扩容次数。
func CountWords(words []string) map[string]int {
	result := make(map[string]int, len(words))
	for _, word := range words {
		result[word]++
	}
	return result
}

// FirstPositive 返回切片中第一个大于 0 的数字。
// 若不存在正数，则返回 0 和 ErrNoPositive。
//
// 注意：返回 0 不代表“找到了 0”，
// 是否成功应以 error 是否为 nil 为准。
func FirstPositive(nums []int) (int, error) {
	for _, n := range nums {
		if n > 0 {
			return n, nil
		}
	}
	return 0, ErrNoPositive
}
