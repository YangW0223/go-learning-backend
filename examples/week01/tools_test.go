package week01

import (
	"errors"
	"reflect"
	"testing"
)

// TestUniqueStrings 验证去重逻辑和顺序保持逻辑。
func TestUniqueStrings(t *testing.T) {
	input := []string{"go", "js", "go", "rust", "js"}
	got := UniqueStrings(input)
	want := []string{"go", "js", "rust"}

	// 切片比较使用 reflect.DeepEqual。
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

// TestGroupByFirstLetter 验证分组逻辑，并覆盖空字符串边界。
func TestGroupByFirstLetter(t *testing.T) {
	input := []string{"apple", "ant", "banana", "boat", ""}
	got := GroupByFirstLetter(input)

	want := map[string][]string{
		"a": []string{"apple", "ant"},
		"b": []string{"banana", "boat"},
		"":  []string{""},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

// TestFirstPositive 验证“正常路径”:
// 输入包含正数时，应返回第一个正数且 err=nil。
func TestFirstPositive(t *testing.T) {
	got, err := FirstPositive([]int{-3, -1, 2, 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

// TestFirstPositive_NotFound 验证“异常路径”:
// 没有正数时应返回哨兵错误 ErrNoPositive。
func TestFirstPositive_NotFound(t *testing.T) {
	_, err := FirstPositive([]int{-3, -1, 0})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrNoPositive) {
		t.Fatalf("unexpected error: %v", err)
	}
}
