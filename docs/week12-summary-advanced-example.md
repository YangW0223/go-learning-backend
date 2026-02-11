# Week12 总结与进阶示例说明

## 目标

沉淀作品集能力与下一阶段学习计划：

- 计算当前能力覆盖率
- 生成项目叙述（便于面试/演示）
- 生成下一阶段 4 周计划

## 代码位置

- `examples/week12/summary_advanced.go`
- `examples/week12/summary_advanced_test.go`
- `examples/week12/cmd/main.go`

## 关键行为

1. `Portfolio.MarkCompleted` 累加能力项。
2. `Score(total)` 计算完成率。
3. `NextStagePlan(topic)` 返回按周拆分计划。

## 运行与验证

运行示例：

```bash
go run ./examples/week12/cmd
```

预期输出：

```text
score: 66
narrative: project=todo-service capabilities=...
next stage (mq): [week1: ... week4: ...]
```

运行测试：

```bash
go test -v ./examples/week12
```
