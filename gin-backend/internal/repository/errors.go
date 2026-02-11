package repository

import "errors"

var (
	// ErrNotFound 用于标识资源不存在。
	// 上层通常映射成 HTTP 404。
	ErrNotFound = errors.New("not found")
	// ErrConflict 用于标识唯一索引冲突等资源冲突。
	// 上层通常映射成 HTTP 409。
	ErrConflict = errors.New("conflict")
)
