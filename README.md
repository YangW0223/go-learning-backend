# Go Learning Backend

一个面向前端开发者的 Go 后端学习项目。目标是从 0 到 1 完成一个可运行、可测试、可扩展的 API 服务。

## 当前能力

- HTTP 服务（标准库）
- 健康检查: `GET /healthz`
- Todo API（内存存储 + 可选 Redis 列表缓存）
  - `POST /api/v1/todos`
  - `GET /api/v1/todos`
  - `PATCH /api/v1/todos/{id}/done`
  - `DELETE /api/v1/todos/{id}`
- 配置中心（环境变量）
- 单元测试
- Docker 运行配置（`api + redis`）

## 快速开始

### 本地运行

```bash
make run
```

默认监听 `:8080`。

### 启用 Redis 缓存（本地）

```bash
export REDIS_ENABLED=true
export REDIS_ADDR=localhost:6379
make run
```

## Docker 运行

```bash
make docker-up
```

关闭并清理：

```bash
make docker-down
```

查看日志：

```bash
make docker-logs
```

## 关键环境变量

- `PORT`：服务端口，默认 `8080`
- `REDIS_ENABLED`：是否启用 Redis 缓存，默认 `false`
- `REDIS_ADDR`：Redis 地址，默认 `localhost:6379`
- `REDIS_PASSWORD`：Redis 密码，默认空
- `REDIS_DB`：Redis DB 库编号，默认 `0`
- `REDIS_CACHE_TTL_SECONDS`：Todo 列表缓存 TTL（秒），默认 `30`
- `REDIS_DIAL_TIMEOUT_MS`：Redis 连接超时（毫秒），默认 `1000`
- `REDIS_IO_TIMEOUT_MS`：Redis 读写超时（毫秒），默认 `1000`

## 示例请求

```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H 'Content-Type: application/json' \
  -d '{"title":"learn go basics"}'

curl http://localhost:8080/api/v1/todos

curl -X PATCH http://localhost:8080/api/v1/todos/<todo_id>/done

curl -X DELETE http://localhost:8080/api/v1/todos/<todo_id>
```

## 学习路线

12 周学习计划见 `learning-plan-details/LEARNING_PLAN_12_WEEKS.md`。
