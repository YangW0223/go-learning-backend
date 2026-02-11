# Week08 缓存与性能示例说明

## 目标

演示 Cache-Aside 模式及命中率对延迟的影响：

- 首次请求 miss 回源
- 再次请求 hit 直接返回
- 写后失效策略（按用户前缀失效）

## 代码位置

- `examples/week08/cache_performance.go`
- `examples/week08/cache_performance_test.go`
- `examples/week08/cmd/main.go`

## 关键行为

1. `CachedTodoService.List`：`cache -> source -> 回填 cache`。
2. `InvalidateUser`：模拟写操作后失效。
3. `Stats`：统计 hit/miss。

## 运行与验证

运行示例：

```bash
go run ./examples/week08/cmd
```

预期输出（耗时会因机器差异变化）：

```text
first call latency=...
second call latency=...
cache stats hit=1 miss=1
```

运行测试：

```bash
go test -v ./examples/week08
```
