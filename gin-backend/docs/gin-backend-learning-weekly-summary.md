# Gin Backend 学习执行总览

本页作为索引入口，按“每周一个文件”组织学习内容。

## 周文件索引

- Week 1：`docs/learning-weeks/week-01.md`
- Week 2：`docs/learning-weeks/week-02.md`
- Week 3：`docs/learning-weeks/week-03.md`
- Week 4：`docs/learning-weeks/week-04.md`
- Week 5：`docs/learning-weeks/week-05.md`
- Week 6：`docs/learning-weeks/week-06.md`
- Week 7：`docs/learning-weeks/week-07.md`
- Week 8：`docs/learning-weeks/week-08.md`

## 8 周总览

| 周次 | 主题 | 周目标 | 关键交付物 |
| --- | --- | --- | --- |
| Week 1 | 全链路跑通 | 理解请求流转与基础调用 | 调通核心接口 + 调用链图 |
| Week 2 | 数据库与迁移 | 掌握 schema 演进 | 新字段 migration + CRUD 改造 |
| Week 3 | 认证与授权 | 吃透 JWT 与中间件 | 受保护接口 + 401/403 场景验证 |
| Week 4 | Redis 缓存 | 建立缓存命中和失效策略 | Todo 列表缓存 + 失效逻辑 |
| Week 5 | 错误与配置 | 统一错误语义和配置 | 错误码规范 + 配置清单 |
| Week 6 | 测试体系 | 建立可重构安全网 | service/handler 测试补齐 |
| Week 7 | 可观测性 | 掌握基础排障流程 | 故障演练与排障清单 |
| Week 8 | 综合交付 | 独立完成端到端功能 | 一次完整功能交付（含文档/测试） |

## 打卡清单（可勾选）

### Week 1

- [ ] Day 1：环境与启动
- [ ] Day 2：认证接口打通
- [ ] Day 3：受保护接口打通
- [ ] Day 4：Todo 业务链路
- [ ] Day 5：代码走读与调用链图
- [ ] Day 6：最小改动实践
- [ ] Day 7：周复盘与下周准备

### Week 2

- [ ] Day 1：读懂现有数据模型
- [ ] Day 2：编写并验证 migration
- [ ] Day 3：改 repository
- [ ] Day 4：改 service + handler + DTO
- [ ] Day 5：联调与回归
- [ ] Day 6：问题修复日
- [ ] Day 7：复盘与归档

### Week 3

- [ ] Day 1：通读认证链路
- [ ] Day 2：设计受保护接口
- [ ] Day 3：实现接口主体
- [ ] Day 4：补失败场景
- [ ] Day 5：自动化验证
- [ ] Day 6：安全加固
- [ ] Day 7：复盘与沉淀

### Week 4

- [ ] Day 1：缓存方案设计
- [ ] Day 2：实现读取与回源
- [ ] Day 3：实现失效策略
- [ ] Day 4：压力与边界验证
- [ ] Day 5：测试与文档
- [ ] Day 6：优化清理
- [ ] Day 7：复盘与迁移

### Week 5

- [ ] Day 1：错误现状盘点
- [ ] Day 2：统一错误模型
- [ ] Day 3：配置治理
- [ ] Day 4：校验与容错
- [ ] Day 5：回归测试
- [ ] Day 6：文档与规范固化
- [ ] Day 7：周复盘

### Week 6

- [ ] Day 1：测试策略设计
- [ ] Day 2：service happy path
- [ ] Day 3：service error path
- [ ] Day 4：handler 接口测试
- [ ] Day 5：测试整合与稳定性
- [ ] Day 6：补漏与优化
- [ ] Day 7：测试周复盘

### Week 7

- [ ] Day 1：可观测入口盘点
- [ ] Day 2：故障演练设计
- [ ] Day 3：执行故障演练
- [ ] Day 4：排障流程抽象
- [ ] Day 5：改进可观测性
- [ ] Day 6：复盘与演练二次验证
- [ ] Day 7：周总结

### Week 8

- [ ] Day 1：选题与设计
- [ ] Day 2：数据层改造
- [ ] Day 3：核心功能实现
- [ ] Day 4：配套能力补齐
- [ ] Day 5：测试与文档
- [ ] Day 6：PR 整理
- [ ] Day 7：最终复盘

## 日常打卡模板

```md
### Day X 打卡（YYYY-MM-DD）
- 今日目标：
- 完成内容：
- 遇到问题：
- 如何解决：
- 明日计划：
```
