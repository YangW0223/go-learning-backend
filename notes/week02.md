# Week 02 HTTP 与 JSON 每日打卡

## 本周目标

- 掌握 `net/http` 请求处理流程：路由匹配 -> 参数解析 -> 业务逻辑 -> JSON 响应。
- 完成 `DELETE /api/v1/todos/{id}`。
- 明确并实现错误分类：400 / 404 / 500。
- 输出删除接口行为文档与错误码规范文档。

## 周一（请求链路梳理）

- [ ] 阅读并画出请求流转图：`internal/app/router.go` -> `internal/handler/todo.go` -> `internal/store/...`。
- [ ] 记录每层职责与边界。

重点文件：

- `internal/app/router.go`
- `internal/handler/todo.go`
- `internal/store/todo_store.go`

## 周二（错误分类设计）

- [ ] 明确 400 / 404 / 500 判定规则。
- [ ] 定义删除接口的输入合法性规则（id 格式）。
- [ ] 记录错误码与错误消息对照关系。

## 周三（删除接口实现）

- [ ] 在路由层接入 `DELETE /api/v1/todos/{id}`。
- [ ] 在 handler 层实现参数校验与错误映射。
- [ ] 在 store 层补齐删除能力。

完成后检查：

- [ ] 删除成功返回 200。
- [ ] id 非法返回 400。
- [ ] id 不存在返回 404。

## 周四（接口测试）

- [ ] 补齐至少 3 个接口测试：成功删除、ID 不存在、ID 非法。
- [ ] 额外补 500 分支单测（存储异常）。

目标文件：

- `internal/app/router_test.go`
- `internal/handler/todo_delete_test.go`

## 周五（文档产出）

- [ ] 写删除接口行为说明（输入、输出、错误码）。
- [ ] 写错误码处理规范文档。
- [ ] 补 `curl` 手动验证命令。

目标文件：

- `docs/week02-delete-api.md`
- `docs/week02-error-code-guidelines.md`
- `docs/week02-curl-examples.md`

## 周六（示例与运行入口）

- [ ] 新增 `examples/week02` 示例代码与测试。
- [ ] 新增独立入口 `examples/week02/cmd/main.go`。
- [ ] 更新 `examples/README.md` 的目录、命令、示例输出。

## 周日（验收与复盘）

- [ ] 运行 `go test ./...`。
- [ ] 对照验收清单逐项确认。
- [ ] 完成复盘记录。

复盘模板：

- 本周最卡点：
- 我如何定位问题：
- 下周要提前避免的问题：

## 验收清单

- [ ] `DELETE /api/v1/todos/{id}` 可用。
- [ ] 400 / 404 / 500 分支行为明确。
- [ ] 至少 3 个接口测试稳定通过。
- [ ] 文档已补齐删除接口描述。

## 本周产出映射

- 删除接口代码：`internal/app/router.go`、`internal/handler/todo.go`、`internal/store/*`
- 删除接口测试：`internal/app/router_test.go`、`internal/handler/todo_delete_test.go`
- 文档：`docs/week02-*.md`
- 示例：`examples/week02/*`
