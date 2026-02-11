# Examples 运行说明

本文档说明如何运行 `examples` 下每周示例代码。

## 前置条件

- 已安装 Go（建议使用与项目一致的版本）。
- 终端位于项目根目录（包含 `go.mod` 的目录）。

## 目录说明

- `examples/week01`：Week01 Go 基础示例与测试。
- `examples/week01/cmd`：Week01 独立运行入口。
- `examples/week02`：Week02 HTTP/JSON 示例与测试。
- `examples/week02/cmd`：Week02 独立运行入口。
- `examples/week03`：Week03 并发与 context 示例与测试。
- `examples/week03/cmd`：Week03 独立运行入口。
- `examples/week04`：Week04 分层重构示例（handler/service/store，含逐行注释版代码）。
- `examples/week04/cmd`：Week04 独立运行入口。
- `examples/week05`：Week05 PostgreSQL 基础概念示例（事务边界/分页）。
- `examples/week05/cmd`：Week05 独立运行入口。
- `examples/week06`：Week06 SQL 工程化示例（sqlc vs gorm 选型思路）。
- `examples/week06/cmd`：Week06 独立运行入口。
- `examples/week07`：Week07 认证鉴权示例（注册/登录/token/中间件）。
- `examples/week07/cmd`：Week07 独立运行入口。
- `examples/week08`：Week08 缓存与性能示例（cache aside + hit/miss）。
- `examples/week08/cmd`：Week08 独立运行入口。
- `examples/week09`：Week09 日志与可观测示例（request_id + metrics）。
- `examples/week09/cmd`：Week09 独立运行入口。
- `examples/week10`：Week10 测试体系示例（service/handler/E2E 风格）。
- `examples/week10/cmd`：Week10 独立运行入口。
- `examples/week11`：Week11 交付部署示例（env 配置/health/readiness）。
- `examples/week11/cmd`：Week11 独立运行入口。
- `examples/week12`：Week12 总结与进阶示例（能力矩阵与下阶段计划）。
- `examples/week12/cmd`：Week12 独立运行入口。
- `examples/week13`：Week13 API 治理与稳定性示例（统一响应/幂等/限流）。
- `examples/week13/cmd`：Week13 独立运行入口。
- `examples/week14`：Week14 可观测与异步任务示例（trace/retry/dead-letter/audit）。
- `examples/week14/cmd`：Week14 独立运行入口。

## 运行方式

1. Week01

```bash
go run ./examples/week01/cmd
```

2. Week02

```bash
go run ./examples/week02/cmd
```

3. Week03

```bash
go run ./examples/week03/cmd
```

4. Week04

```bash
go run ./examples/week04/cmd
```

5. Week04 HTTP 服务模式（用于 curl 验证）

```bash
go run ./examples/week04/cmd -mode server -addr :18084
```

更多请求示例见：`docs/week04-curl-examples.md`

6. Week05

```bash
go run ./examples/week05/cmd
```

示例输出：

```text
id=2 title="implement mark done" done=false
id=1 title="design todos table" done=true
```

7. Week06

```bash
go run ./examples/week06/cmd
```

示例输出：

```text
chosen tool: sqlc
prefix=al users: [{1 alice} {2 allen}]
```

8. Week07

```bash
go run ./examples/week07/cmd
```

示例输出：

```text
issued token: <token>
authorized => status=200 body=access granted
missing token => status=401 body=unauthorized
```

9. Week08

```bash
go run ./examples/week08/cmd
```

示例输出：

```text
first call latency=...
second call latency=...
cache stats hit=1 miss=1
```

10. Week09

```bash
go run ./examples/week09/cmd
```

示例输出：

```text
request1 status: 200
request2 status: 500
metrics: requests=2 errors=1 ...
```

11. Week10

```bash
go run ./examples/week10/cmd
```

示例输出：

```text
POST /api/v1/todos => 201 ...
GET /api/v1/todos/1 => 200 ...
```

12. Week11

```bash
go run ./examples/week11/cmd
```

示例输出：

```text
loaded config: port=8080 env=dev version=v1.1.0
GET /healthz => 200 ok
GET /readyz => 200 ready
```

13. Week12

```bash
go run ./examples/week12/cmd
```

示例输出：

```text
score: 66
narrative: project=todo-service capabilities=...
```

14. Week13

```bash
go run ./examples/week13/cmd
```

示例输出：

```text
first => status=201 ...
replay => status=201 replay=1 ...
```

15. Week14

```bash
go run ./examples/week14/cmd
```

示例输出：

```text
results: [...]
dead letters: [...]
spans: [...]
audit events: [...]
```

## 测试命令

1. 分周执行测试

```bash
go test -v ./examples/week01
go test -v ./examples/week02
go test -v ./examples/week03
go test -v ./examples/week04
go test -v ./examples/week05
go test -v ./examples/week06
go test -v ./examples/week07
go test -v ./examples/week08
go test -v ./examples/week09
go test -v ./examples/week10
go test -v ./examples/week11
go test -v ./examples/week12
go test -v ./examples/week13
go test -v ./examples/week14
```

2. 运行仓库所有测试（包含 examples）

```bash
go test ./...
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
