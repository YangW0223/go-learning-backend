package week02

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

var (
	// ErrInvalidDeletePath 表示删除接口路径格式错误。
	ErrInvalidDeletePath = errors.New("invalid delete path")
	// ErrInvalidTodoID 表示路径中的 id 不满足当前规则。
	ErrInvalidTodoID = errors.New("invalid todo id")
)

var week02TodoIDPattern = regexp.MustCompile(`^\d{14}\.\d{9}$`)

// DeleteResult 表示删除接口成功时返回的最小数据结构。
type DeleteResult struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

// ParseDeleteTodoPath 从路径中提取并校验 todo id。
// 期望格式: /api/v1/todos/{id}
func ParseDeleteTodoPath(path string) (string, error) {
	const prefix = "/api/v1/todos/"
	if !strings.HasPrefix(path, prefix) {
		return "", ErrInvalidDeletePath
	}

	id := strings.TrimPrefix(path, prefix)
	id = strings.Trim(id, "/")
	if id == "" || strings.Contains(id, "/") {
		return "", ErrInvalidTodoID
	}
	if !week02TodoIDPattern.MatchString(id) {
		return "", ErrInvalidTodoID
	}
	return id, nil
}

// BuildSuccessJSON 构造统一成功响应：{"data":..., "error":null}。
func BuildSuccessJSON(data any) ([]byte, error) {
	return json.Marshal(map[string]any{
		"data":  data,
		"error": nil,
	})
}

// BuildErrorJSON 构造统一错误响应：{"data":null, "error":"..."}。
func BuildErrorJSON(message string) ([]byte, error) {
	return json.Marshal(map[string]any{
		"data":  nil,
		"error": message,
	})
}
