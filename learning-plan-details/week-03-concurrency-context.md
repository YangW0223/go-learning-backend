# Week 03：并发与 context

## 本周目标

理解 goroutine、channel、mutex、context，并在列表接口中加入超时控制。

## 详细步骤

1. 用最小示例分别练习 goroutine/channel/select。
2. 对比两种并发同步方式：mutex 与 channel。
3. 为列表接口加 `context.WithTimeout`（如 200ms/500ms，按你的实现定）。
4. 模拟慢请求并验证超时返回行为。
5. 在日志中记录超时触发点（便于后续排查）。
6. 总结“什么时候用 mutex，什么时候用 channel”。

## 建议实践清单

- 避免 goroutine 泄漏，确保协程有退出路径。
- 在超时场景返回明确错误码和信息。
- 对共享状态读写必须有同步机制。

## 验收清单

- [ ] 列表接口具备 deadline 控制。
- [ ] 能清楚解释 mutex 和 channel 的使用边界。
- [ ] 超时场景有测试或可复现实验。

## 产出物

- 超时控制实现代码。
- 并发模型总结笔记。

## 常见风险与排查

- 忘记 `cancel()`：始终在创建 timeout 后 `defer cancel()`。
- 无缓冲 channel 阻塞：确认发送和接收是否成对。
