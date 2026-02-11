# Week09 日志与可观测示例说明

## 目标

建立最小可观测能力：

- 结构化日志（JSON）
- `request_id` 贯穿单次请求
- 基础指标（请求数、错误数、耗时桶）

## 代码位置

- `examples/week09/logging_observability.go`
- `examples/week09/logging_observability_test.go`
- `examples/week09/cmd/main.go`

## 关键行为

1. `WithObservability` 包装 handler，记录 method/path/status/latency/request_id。
2. `Metrics.Observe` 累计请求数、错误数和耗时分桶。
3. 通过 `?fail=1` 可制造一次 500 便于排障演练。

## 运行与验证

运行示例：

```bash
go run ./examples/week09/cmd
```

预期输出（日志顺序固定，耗时字段可能有微小变化）：

```text
request1 status: 200
request2 status: 500
metrics: requests=2 errors=1 ...
logs:
{..."request_id":"req-1"...}
{..."request_id":"req-2"..."status":500...}
```

运行测试：

```bash
go test -v ./examples/week09
```
