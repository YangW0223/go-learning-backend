// package week04 定义 Week04 的分层示例代码。
package week04

// errors 包用于创建静态业务错误变量。
import "errors"

// var 分组声明 Week04 中会复用的业务错误常量。
var (
	// ErrInvalidTodoID 表示传入的 todo id 不符合业务规则。
	ErrInvalidTodoID = errors.New("invalid todo id")
	// ErrTodoNotFound 表示在存储层没有找到对应 id 的 todo。
	ErrTodoNotFound = errors.New("todo not found")
)
