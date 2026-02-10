# Week 02 错误码处理规范

本文档用于约束 HTTP 接口错误处理，减少“业务错误”和“系统错误”混淆。

## 目标

- 错误码可预测、可测试、可复用。
- 前端或调用方能根据状态码快速判断问题类型。

## 分类规则

1. `400 Bad Request`
- 请求参数不合法（格式错误、缺失、类型不匹配）。
- 例如：`DELETE /api/v1/todos/abc`（id 不符合规则）。

2. `404 Not Found`
- 请求参数合法，但资源不存在。
- 例如：删除一个合法 id，但存储中无对应 todo。

3. `500 Internal Server Error`
- 系统内部错误（存储不可用、未知异常）。
- 不向外暴露底层细节，统一返回可读消息。

## 实施建议

- 在 store 层定义可判定错误（例如 `store.ErrTodoNotFound`）。
- 在 handler 层集中映射 HTTP 状态码，不在路由层散落业务判断。
- 每个接口至少覆盖：成功、参数错误、资源不存在。
- 对关键分支补单测，保证回归稳定。

## Week 02 实际落地

- 400 映射：`internal/handler/todo.go` 中 `Delete` 的 id 校验失败。
- 404 映射：`errors.Is(err, store.ErrTodoNotFound)`。
- 500 映射：删除时遇到非 `ErrTodoNotFound` 错误。

## 测试建议

- 接口级测试验证：200 / 400 / 404。
- 单元测试补充：500 分支（通过 fake store 注入错误）。
