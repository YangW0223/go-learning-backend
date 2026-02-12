# Week 1（Day 1 ~ Day 7）

### Day 1：环境与启动

- 今日目标：本地完整启动 API + Postgres + Redis。
- 操作清单：
  - 执行 `make docker-up`。
  - 查看容器状态：`docker compose --env-file docker/.env.example -f docker/docker-compose.yaml ps`。
  - 查看 API 日志：`make docker-logs`。
- 完成标准：
  - 3 个服务均为 `Up`。
  - 能在日志里看到服务启动和路由注册信息。

### Day 2：认证接口打通

- 今日目标：完成注册和登录，拿到 JWT。
- 操作清单：
  - 调用 `POST /api/v1/auth/register`。
  - 调用 `POST /api/v1/auth/login`，保存返回 token。
  - 用错误参数各测 1 次（如空邮箱、错误密码）。
- 完成标准：
  - 成功拿到 token。
  - 能区分成功和失败响应差异（状态码、错误信息）。

### Day 3：受保护接口打通

- 今日目标：理解鉴权中间件行为。
- 操作清单：
  - 使用 token 调用 `GET /api/v1/me` 和 `GET /api/v1/todos`。
  - 不带 token 调用同样接口并记录返回。
  - 阅读 `internal/transport/http/middleware` 相关鉴权代码。
- 完成标准：
  - 能解释为什么无 token 返回 401。
  - 能说明 token 在哪里被解析并写入上下文。

### Day 4：Todo 业务链路

- 今日目标：跑通 Todo 的增删改查最小闭环。
- 操作清单：
  - 创建 Todo：`POST /api/v1/todos`。
  - 查询 Todo：`GET /api/v1/todos`。
  - 更新 Todo：`PATCH /api/v1/todos/:id`。
  - 删除 Todo：`DELETE /api/v1/todos/:id`。
- 完成标准：
  - 4 个操作都可成功调用。
  - 能定位每个接口对应的 handler/service/repository 文件。

### Day 5：代码走读与调用链图

- 今日目标：把请求流转写成你自己的图和说明。
- 操作清单：
  - 重点阅读：`internal/transport/http/router`、`internal/transport/http/handler`、`internal/service`。
  - 选择“登录”或“创建 Todo”画调用链图。
  - 写一段 150~300 字说明（输入、关键处理、输出）。
- 完成标准：
  - 能在不看代码情况下复述完整调用链。

### Day 6：最小改动实践

- 今日目标：做一个低风险小改动并验证。
- 操作清单：
  - 选择小改动（如错误提示文案优化、参数校验补充）。
  - 修改后执行 `go test ./...`。
  - 重启服务并手动验证 1 个相关接口。
- 完成标准：
  - 测试通过，接口行为符合预期。
  - 能解释改动影响了哪些层。

### Day 7：周复盘与下周准备

- 今日目标：形成本周复盘并确定 Week 2 改造点。
- 操作清单：
  - 记录本周学到的 5 个后端概念。
  - 记录 3 个踩坑与解决方法。
  - 选定 Week 2 字段改造目标（`due_at` 或 `priority`）。
- 完成标准：
  - 输出一份可执行的下周任务草案（字段、接口、测试点）。

