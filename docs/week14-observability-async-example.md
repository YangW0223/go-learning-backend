# Week14 可观测与异步任务示例说明

## 目标

实现“Trace + 审计日志 + 异步重试 + 死信”最小闭环：

- 任务处理写入 Span
- 失败任务记录错误原因
- 超过重试次数进入 Dead Letter
- 审计日志记录成功/失败结果

## 代码位置

- `examples/week14/observability_async.go`
- `examples/week14/observability_async_test.go`
- `examples/week14/cmd/main.go`

## 关键行为

1. `AsyncWorker.ProcessBatch` 支持最多 `MaxRetries+1` 次尝试。
2. 临时错误可重试成功；永久错误写入死信列表。
3. `Tracer` 记录每个 job 的 span，失败 span 含错误信息。
4. `AuditLogger` 输出 `success/failed` 事件。

## 运行与验证

运行示例：

```bash
go run ./examples/week14/cmd
```

预期输出：

```text
results: [{job-1 true 2 } {job-2 false 3 permanent error}]
dead letters: [{job-2 false 3 permanent error}]
spans: [...]
audit events: [...]
```

运行测试：

```bash
go test -v ./examples/week14
```
