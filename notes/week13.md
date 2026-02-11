# Week 13 API 治理与稳定性工程每日打卡

## 本周目标

- 建立统一 API 契约（响应结构、错误码、版本约定）。
- 为关键写接口落地幂等机制并补测试。
- 引入限流与优雅停机能力，完成一次稳定性演练。

## 周一（API 契约规范）

- [ ] 统一响应结构（`data` / `error` / `request_id`）。
- [ ] 统一错误码与 HTTP 状态码映射（400/401/403/404/409/429/500）。
- [ ] 输出初版规范文档。

目标文件：

- `docs/week13-api-contract.md`
- `docs/week13-error-codes.md`

## 周二（OpenAPI 文档）

- [ ] 为登录与 Todo CRUD 输出 OpenAPI 文档。
- [ ] 导入 Postman/Apifox 做一次可执行验证。
- [ ] 补充示例请求与响应。

目标文件：

- `docs/week13-openapi.yaml`
- `docs/week13-api-examples.md`

## 周三（幂等机制）

- [ ] 选择 1 个写接口实现幂等（建议创建 Todo）。
- [ ] 设计幂等键存储与过期策略。
- [ ] 编写重复请求回归测试。

目标文件：

- `internal/handler/todo.go`
- `internal/service/todo_service.go`
- `internal/service/todo_service_test.go`

## 周四（限流中间件）

- [ ] 实现按用户或 IP 的限流中间件。
- [ ] 统一 429 响应结构。
- [ ] 加入最小配置项（速率、突发值）。

目标文件：

- `internal/middleware/ratelimit.go`
- `internal/app/router.go`
- `internal/app/router_test.go`

## 周五（超时/重试策略）

- [ ] 为关键依赖调用补超时控制。
- [ ] 为可重试错误加入有限重试。
- [ ] 记录不可重试错误类型。

目标文件：

- `internal/store/*`
- `internal/service/*`
- `docs/week13-retry-timeout-guidelines.md`

## 周六（优雅停机）

- [ ] 接入 `server.Shutdown`，支持优雅退出。
- [ ] 模拟在途请求，验证退出行为。
- [ ] 记录停机流程图与验证命令。

目标文件：

- `cmd/api/main.go`
- `docs/week13-graceful-shutdown.md`

## 周日（演练与复盘）

- [ ] 执行一次最小故障演练（限流触发、依赖超时、服务重启）。
- [ ] 输出“现象-定位-修复-回归”复盘。
- [ ] 对照验收清单逐项确认。

复盘模板：

- 故障现象：
- 定位证据：
- 根因结论：
- 修复动作：
- 预防措施：

## 验收清单

- [ ] 有统一 API 契约和错误码文档。
- [ ] 至少 1 个写接口支持幂等并有测试。
- [ ] 限流可生效并返回 429。
- [ ] 服务支持优雅停机且有验证记录。
- [ ] 完成一次稳定性演练并形成复盘。
