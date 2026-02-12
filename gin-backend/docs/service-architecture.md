# Gin Backend 服务组织说明

本文说明这个项目如何由目录中的代码拼成一个可运行服务，重点是：
- 启动时如何装配依赖
- 请求进入后如何在各层流转
- 前端联调最关心的协议约定

## 1. 分层与职责

| 层 | 目录 | 职责 |
| --- | --- | --- |
| 进程入口 | `cmd/server` | 进程生命周期：启动、监听信号、优雅停机 |
| 装配层 | `internal/bootstrap` | 依赖注入：配置、DB、缓存、service、handler、router |
| 协议层 | `internal/transport/http/router` + `internal/transport/http/handler` + `internal/transport/http/dto` | 路由注册、请求绑定、响应输出 |
| 业务层 | `internal/service` | 业务规则、权限判断、缓存失效策略 |
| 数据层 | `internal/repository/postgres` + `internal/repository/redis` | Postgres 持久化、Redis 缓存 |
| 横切能力 | `internal/config` + `internal/auth` + `internal/response` + `internal/transport/http/middleware` | 配置、JWT、统一响应、中间件能力 |

## 2. 启动装配流程

入口是 `cmd/server/main.go`，核心装配在 `internal/bootstrap/app.go`：

1. `config.Load()` 读取环境变量并校验。
2. 初始化 logger 和 metrics。
3. 连接 Postgres，并在启动阶段 `EnsureSchema` 保证核心表存在。
4. 创建 repository（user/todo）。
5. 创建 auth service（密码哈希 + JWT）。
6. 根据 `REDIS_ENABLED` 装配真实 Redis 缓存或 no-op 缓存。
7. 创建 todo service。
8. 创建各 handler。
9. 构建 Gin router（含全局中间件和路由组）。
10. 生成 `http.Server` 并启动监听。

最终 `main` 负责处理 `SIGINT/SIGTERM`，执行优雅停机并释放资源。

## 3. 一次请求如何流转

以 `POST /api/v1/todos` 为例：

1. `router` 注册到 `/api/v1` 受保护路由组。
2. `AuthJWT` 中间件校验 `Authorization: Bearer <token>`，并把 `user_id` 写入上下文。
3. `TodoHandler.Create` 绑定 JSON，调用 `TodoService.Create`。
4. `TodoService.Create` 做参数校验，生成业务对象，调用 `TodoRepository.Create` 入库。
5. 写操作成功后删除该用户的 todo 列表缓存。
6. `response.Success` 返回统一响应结构。

一句话链路：`router -> middleware -> handler -> service -> repository -> response`。

## 4. 前端联调约定

### 4.1 认证方式

- 登录接口返回 `access_token`。
- 受保护接口统一使用：`Authorization: Bearer <token>`。

### 4.2 统一响应结构

成功：

```json
{
  "request_id": "trace-id",
  "data": {}
}
```

失败：

```json
{
  "request_id": "trace-id",
  "error": {
    "code": "BAD_REQUEST",
    "message": "invalid request"
  }
}
```

### 4.3 常见错误码

- `BAD_REQUEST` (400)
- `UNAUTHORIZED` (401)
- `FORBIDDEN` (403)
- `NOT_FOUND` (404)
- `CONFLICT` (409)
- `INTERNAL_ERROR` (500)
- `RATE_LIMITED` (429，限流中间件返回)

## 5. 运行方式

本地：

1. 准备 Postgres/Redis
2. 设置环境变量（见 `README.md`）
3. `make migrate-up`
4. `make run`

Docker：

1. `make docker-up` 启动 API + Postgres + Redis
2. `make docker-down` 停止并清理
