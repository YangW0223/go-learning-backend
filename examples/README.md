# Examples 运行说明

本文档说明如何运行 `examples` 下每周示例代码。

## 前置条件

- 已安装 Go（建议使用与项目一致的版本）。
- 终端位于项目根目录（包含 `go.mod` 的目录）。

## 目录说明

- `examples/week01`：Week01 示例函数与测试（普通包）。
- `examples/week01/cmd`：Week01 独立运行入口（`go run`）。
- `examples/week02`：Week02 示例函数与测试（普通包）。
- `examples/week02/cmd`：Week02 独立运行入口（`go run`）。
- `examples/week03`：Week03 并发与 context 示例函数与测试（普通包）。
- `examples/week03/cmd`：Week03 独立运行入口（`go run`）。

## 运行方式

1. 运行 Week01 示例入口：

```bash
go run ./examples/week01/cmd
```

示例输出（map 键顺序可能不同）：

```text
Greet: hello, Go Learner
CountWords: map[go:3 js:1 rust:1]
UniqueStrings: [go js rust]
GroupByFirstLetter: map[:[] a:[ant apple] b:[bear boat]]
FirstPositive(normal): 2
FirstPositive(error): no positive number found
```

2. 运行 Week02 示例入口：

```bash
go run ./examples/week02/cmd
```

示例输出：

```text
Parsed ID: 20260210112233.123456789
Success JSON: {"data":{"id":"20260210112233.123456789","deleted":true},"error":null}
Error JSON: {"data":null,"error":"invalid todo id"}
```

3. 运行 Week01 示例测试：

```bash
go test -v ./examples/week01
```

4. 运行 Week02 示例测试：

```bash
go test -v ./examples/week02
```

5. 运行仓库所有测试（包含 examples）：

```bash
go test ./...
```

6. 运行 Week03 示例入口：

```bash
go run ./examples/week03/cmd
```

示例输出：

```text
SumWithMutex: 15
SumWithChannel: 15
FirstValue: 7
ListWithTimeout(success): [todo-1 todo-2 todo-3]
ListWithTimeout(timeout): request timeout: deadline=40ms
ChatRoom alice received: hello bob
ChatRoom bob received: hello bob
ChatRoom stats: users=1 messages=1
```

7. 运行 Week03 示例测试：

```bash
go test -v ./examples/week03
```

## 文档规则

1. 新增 `examples/weekXX` 时，必须同步更新本文件：
- 目录说明
- `go run` 命令
- `go test` 命令
- 示例输出

2. 如果路径、包名或命令变化，必须同次改动更新本文件，避免文档过期。

## 常见问题

1. 报错 `go.mod file not found`
- 原因：不在项目根目录执行命令。
- 解决：先 `cd` 到仓库根目录，再运行命令。

2. 报错 `package ... is not a main package`
- 原因：直接运行了普通包（如 `go run ./examples/week01`）。
- 解决：运行 `cmd` 入口，例如 `go run ./examples/week01/cmd`。
