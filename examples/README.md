# Examples 运行说明

本文档说明如何运行 `examples` 下的示例代码。

## 前置条件

- 已安装 Go（建议使用与项目一致的版本）。
- 当前终端目录在项目根目录（包含 `go.mod` 的目录）。

## 目录说明

当前示例位于：

- `examples/week01`
- `examples/week01/cmd`（`go run` 入口）

说明：

- `examples/week01` 是普通包（`package week01`），用于放可复用函数与测试。
- `examples/week01/cmd` 是 `main` 程序入口，用于直接演示输出结果。

## 运行方式

1. 直接运行示例入口（推荐先跑这个）：

```bash
go run ./examples/week01/cmd
```

示例输出（你的 map 键顺序可能不同，这是 Go 的正常行为）：

```text
Greet: hello, Go Learner
CountWords: map[go:3 js:1 rust:1]
UniqueStrings: [go js rust]
GroupByFirstLetter: map[:[] a:[ant apple] b:[bear boat]]
FirstPositive(normal): 2
FirstPositive(error): no positive number found
```

2. 运行 week01 全部示例测试：

```bash
go test -v ./examples/week01
```

3. 只运行某一个测试（示例：去重函数）：

```bash
go test -v ./examples/week01 -run TestUniqueStrings
```

4. 运行仓库所有测试（包含 examples）：

```bash
go test ./...
```

## 文档规则

1. 新增 `examples/*` 示例代码时，必须同步更新本文件中的“目录说明”和“运行方式”。
2. 如果新增了 `go run` 入口，必须写明完整可执行命令（从项目根目录开始执行）。
3. 如果命令、路径或包名发生变化，必须在同一次改动中更新本文件，避免文档过期。

## 常见问题

1. 报错 `go.mod file not found`
- 原因：不在项目根目录执行命令。
- 解决：先 `cd` 到仓库根目录，再运行上述命令。

2. `go run ./examples/week01` 报错
- 原因：`examples/week01` 是普通包，不是 `main` 包。
- 解决：使用 `go run ./examples/week01/cmd`。
