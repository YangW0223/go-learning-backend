// package week04 定义 Week04 的分层示例代码。
package week04

// import 分组引入 service 层实现所依赖的标准库。
import (
	// context 用于透传请求生命周期与取消信号。
	"context"
	// errors 用于判断错误链中是否包含指定业务错误。
	"errors"
	// fmt 用于包装错误并保留底层错误信息。
	"fmt"
	// regexp 用于预编译并复用 id 校验规则。
	"regexp"
	// strings 用于处理输入字符串中的空白字符。
	"strings"
)

// week04TodoIDPattern 约束 id 格式为 1~9 位数字且首位非 0。
var week04TodoIDPattern = regexp.MustCompile(`^[1-9][0-9]{0,8}$`)

// TodoService 负责承载 Todo 相关业务规则。
type TodoService struct {
	// store 是 service 层依赖的数据访问抽象。
	store TodoStore
}

// NewTodoService 通过依赖注入创建 TodoService。
func NewTodoService(store TodoStore) *TodoService {
	// 返回绑定了具体 store 实现的 service 实例。
	return &TodoService{store: store}
}

// MarkDone 执行“将指定 Todo 标记为完成”的业务流程。
func (s *TodoService) MarkDone(ctx context.Context, id string) (Todo, error) {
	// 先去除输入两端空白，避免无效空格干扰校验。
	id = strings.TrimSpace(id)
	// 若 id 不匹配业务规则，则直接返回非法参数错误。
	if !week04TodoIDPattern.MatchString(id) {
		return Todo{}, ErrInvalidTodoID
	}

	// 调用 store 执行状态更新并获取更新结果。
	todo, err := s.store.MarkDone(ctx, id)
	// 若 store 返回错误，则进行业务错误分类与包装。
	if err != nil {
		// 对“资源不存在”这类已知业务错误进行透传。
		if errors.Is(err, ErrTodoNotFound) {
			return Todo{}, ErrTodoNotFound
		}
		// 对未知错误进行包装，保留原始错误链便于排查。
		return Todo{}, fmt.Errorf("mark done failed: %w", err)
	}

	// 正常路径返回更新后的 todo。
	return todo, nil
}
