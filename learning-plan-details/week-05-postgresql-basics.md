# Week 05：PostgreSQL 入门

## 本周目标

把 Todo 存储从内存切换到 PostgreSQL，并理解事务边界。

## 详细步骤

1. 设计 `todos` 表结构（id/title/done/created_at/updated_at）。
2. 创建索引（至少考虑 `created_at` 或高频查询字段）。
3. 建立数据库连接配置（环境变量管理）。
4. 实现 `Create/List/MarkDone` 的 DB 版本。
5. 设计分页查询（`limit/offset` 或游标）。
6. 标注事务边界：哪些操作必须在同一事务中。

## 建议实践清单

- 迁移脚本版本化管理。
- SQL 语句参数化，避免拼接。
- 对 DB 错误做分类映射。

## 验收清单

- [ ] `todos` 表可创建并可重复初始化。
- [ ] `Create/List/MarkDone` 走 DB 并可用。
- [ ] 能解释至少 1 处事务边界设计。
- [ ] 分页 SQL 可运行且结果正确。

## 产出物

- 数据库迁移脚本。
- DB store 实现和测试。
- 事务边界说明文档。

## 常见风险与排查

- 连接泄漏：查询后检查 `rows.Close()`。
- N+1 查询：优先单 SQL 获取列表数据。
