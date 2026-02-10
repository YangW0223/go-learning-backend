package week01

// UniqueStrings 对字符串切片去重，并保持“首次出现顺序”不变。
//
// 例子:
// 输入: ["go", "js", "go", "rust"]
// 输出: ["go", "js", "rust"]
//
// seen 的 value 使用 struct{} 是 Go 里常见的“集合”写法，
// struct{} 不占额外存储空间，只关心 key 是否存在。
func UniqueStrings(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	result := make([]string, 0, len(items))

	for _, item := range items {
		// 如果 item 已出现过，直接跳过。
		if _, ok := seen[item]; ok {
			continue
		}
		// 第一次出现时，记录到 seen 并加入结果切片。
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}

// GroupByFirstLetter 按字符串首字节分组。
//
// 例子:
// 输入: ["apple", "ant", "banana"]
// 输出: {"a": ["apple", "ant"], "b": ["banana"]}
//
// 注意：这里按 byte 分组（item[0]），
// 对 ASCII 英文场景足够；若要完整支持中文等多字节字符，
// 需要改为按 rune 处理。
func GroupByFirstLetter(items []string) map[string][]string {
	result := make(map[string][]string)
	for _, item := range items {
		// 空字符串单独归到 key=""，避免 item[0] 越界。
		if item == "" {
			result[""] = append(result[""], item)
			continue
		}
		key := string(item[0])
		result[key] = append(result[key], item)
	}
	return result
}
