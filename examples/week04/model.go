// package week04 定义 Week04 的分层示例代码。
package week04

// Todo 表示一个最小可用的待办事项实体。
// 当前结构只保留演示分层所必需的 3 个字段。
type Todo struct {
	// ID 是 Todo 的唯一标识，对外通过 JSON 字段 id 传输。
	ID string `json:"id"`
	// Title 是 Todo 的标题内容，对外通过 JSON 字段 title 传输。
	Title string `json:"title"`
	// Done 表示任务是否完成，对外通过 JSON 字段 done 传输。
	Done bool `json:"done"`
}
