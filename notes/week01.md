# Week 01 Go 基础每日打卡

## 本周目标

- 能独立写出函数、结构体、接口、错误处理。
- 理解本项目目录结构和入口逻辑。
- 完成 `GET /ping` 与至少 1 个测试。
- 完成 2 个 Go 小工具函数重写。

## 周一（语法热身）

- [ ] 通读变量、切片、map 基础语法。
- [ ] 手写最小例子并运行。

最小例子：

```go
name := "go"
nums := []int{1, 2, 3}
counts := map[string]int{"go": 1}
counts["go"]++
```

## 周二（struct/interface/error）

- [ ] 写 `struct`、`interface`、`error` 的最小例子。
- [ ] 理解“值接收者/指针接收者”的差异（本周先用值接收者即可）。

最小例子：

```go
type User struct {
	Name string
}

type Greeter interface {
	Greet() string
}

var ErrNoPositive = errors.New("no positive number found")
```

## 周三（JS 小工具重写 1）

- [ ] 完成“去重”函数重写：`UniqueStrings`。
- [ ] 补一个对应测试。

代码位置：

- `examples/week01/tools.go`
- `examples/week01/tools_test.go`

## 周四（JS 小工具重写 2）

- [ ] 完成“分组”函数重写：`GroupByFirstLetter`。
- [ ] 补一个对应测试，覆盖空字符串边界。

代码位置：

- `examples/week01/tools.go`
- `examples/week01/tools_test.go`

## 周五（项目入口阅读）

- [ ] 阅读 `cmd/api/main.go`。
- [ ] 阅读 `internal/app/router.go` 与 `internal/handler/...`。
- [ ] 用自己的话记录 `cmd` 与 `internal` 的职责。

记录：

- `cmd`：程序入口，负责组装依赖并启动服务。
- `internal`：项目内部实现（handler、store、model 等），不对外暴露。

## 周六（接口实践）

- [ ] 新增 `GET /ping`，返回 `{"message":"pong"}`。
- [ ] 至少补 1 个单测验证状态码和返回体。

代码位置：

- `internal/handler/ping.go`
- `internal/app/router.go`
- `internal/app/router_test.go`

## 周日（验收与复盘）

- [ ] 运行 `make test` 或 `go test ./...`。
- [ ] 对照验收清单逐项打勾。
- [ ] 写一段复盘：本周最卡点、下周改进点。

复盘模板：

- 本周最卡点：
- 我如何排查：
- 下周要提前做的事：

## 验收清单

- [ ] 能口头解释 `cmd` 与 `internal` 的作用。
- [ ] `GET /ping` 可访问且返回正确。
- [ ] 至少完成 2 个 Go 小工具函数重写。
- [ ] 本周代码通过 `make test`。

## 已有产出映射

- Go 基础语法示例：`examples/week01/basics.go`
- 小工具与测试：`examples/week01/tools.go`、`examples/week01/tools_test.go`
- `GET /ping` 代码与测试：`internal/handler/ping.go`、`internal/app/router.go`、`internal/app/router_test.go`
