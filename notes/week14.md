# Week 14 可观测性进阶与异步任务每日打卡

## 本周目标

- 建立日志 + 指标 + Trace 的排障闭环。
- 定义最小 SLI/SLO 与告警阈值。
- 实现一个可重试、可追踪的异步任务。

## 周一（Trace 接入）

- [ ] 接入 OpenTelemetry SDK。
- [ ] 完成基础 Trace 导出配置。
- [ ] 验证能看到最小请求链路。

目标文件：

- `internal/observability/tracing.go`
- `cmd/api/main.go`
- `docs/week14-tracing-setup.md`

## 周二（核心链路打点）

- [ ] 为登录、创建 Todo、查询 Todo、完成 Todo 补 Span。
- [ ] 统一 Span 命名和关键标签（status、error、user_id）。
- [ ] 校验 Trace 与 request_id 能相互关联。

目标文件：

- `internal/handler/*`
- `internal/service/*`

## 周三（SLI/SLO 与告警）

- [ ] 定义最小 SLI（可用性、P95 延迟、错误率）。
- [ ] 设定初版 SLO 目标值与告警阈值。
- [ ] 写告警说明（触发条件、响应动作、升级路径）。

目标文件：

- `docs/week14-sli-slo.md`
- `docs/week14-alerting.md`

## 周四（观测看板）

- [ ] 搭建最小看板：QPS、错误率、P95、缓存命中率。
- [ ] 增加按接口维度过滤。
- [ ] 记录看板截图或导出配置。

目标文件：

- `docs/week14-dashboard.md`

## 周五（异步任务落地）

- [ ] 选择一个异步任务场景（如 Todo 完成后通知）。
- [ ] 完成入队、消费、状态流转设计。
- [ ] 为失败场景补重试和失败原因记录。

目标文件：

- `internal/worker/*`
- `internal/store/*`
- `docs/week14-async-job-design.md`

## 周六（审计日志与安全行为）

- [ ] 为登录成功/失败、权限拒绝、关键写操作补审计日志。
- [ ] 明确审计字段（actor、action、resource、result、timestamp）。
- [ ] 做一次审计日志检索演练。

目标文件：

- `internal/middleware/audit.go`
- `docs/week14-audit-log.md`

## 周日（全链路排障演练）

- [ ] 人为制造 1 个故障（任务失败或依赖抖动）。
- [ ] 用日志 + 指标 + Trace 完成定位。
- [ ] 记录“告警 -> 定位 -> 修复 -> 回归”全过程。

复盘模板：

- 告警触发时间：
- 关键观测证据：
- 根因：
- 修复方案：
- 后续防再发动作：

## 验收清单

- [ ] 核心链路在 Trace 中可完整追踪。
- [ ] 有明确 SLI/SLO 与告警规则文档。
- [ ] 至少 1 个异步任务支持重试与失败追踪。
- [ ] 审计日志覆盖登录与关键写操作。
- [ ] 完成一次全链路故障排查复盘。
