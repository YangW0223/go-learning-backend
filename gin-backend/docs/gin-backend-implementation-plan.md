# Gin 完整后端项目实施文档（根目录独立运行）

更新时间：2026-02-11

## 目标

在当前仓库根目录下新增一个独立的 Gin 后端项目，目录建议为 `gin-backend/`，满足：

1. 可独立运行（不依赖现有 `cmd/api` 与 `internal`）。
2. 具备完整后端体系（配置、路由、业务、存储、缓存、鉴权、日志、观测、测试、容器化）。
3. 具备最小可生产化能力（优雅停机、健康检查、错误治理、基础限流与安全中间件）。

## 交付边界

本次将先搭建一套可运行且可扩展的“完整骨架 + 核心业务”版本，默认业务以 `Todo` + `Auth` 为主。

1. 强制交付：代码、测试、文档、验证脚本。
2. 可选增强（后续迭代）：消息队列、复杂权限模型、分布式追踪后端（如 Jaeger/Tempo）等。

## 目录规划（独立项目）

```text
gin-backend/
  cmd/server/main.go
  internal/
    bootstrap/
    config/
    transport/http/
      router/
      handler/
      middleware/
      dto/
    service/
    repository/
      postgres/
      redis/
    model/
    auth/
    errs/
    observability/
    testkit/
  migrations/
  scripts/
  docker/
    docker-compose.yaml
    .env.example
  Dockerfile
  Makefile
  README.md
```

## 后端体系标准

## 1. 配置体系

1. 支持环境变量加载与默认值。
2. 配置分层：`app`、`http`、`postgres`、`redis`、`auth`、`log`。
3. 启动时进行配置校验，非法配置快速失败。

## 2. 分层架构

1. `handler`：仅做参数绑定、校验、响应映射。
2. `service`：业务编排与事务边界。
3. `repository`：数据访问实现（Postgres）与缓存实现（Redis）。
4. `model/dto`：领域模型与传输对象分离。

## 3. 数据与缓存

1. Postgres 作为主存储（含 migration）。
2. Redis 作为缓存层（读穿透、写后失效）。
3. 提供 repository 接口，支持单测 mock。

## 4. 认证与授权

1. JWT 登录态（access token）。
2. 中间件注入用户上下文。
3. 基础 RBAC（至少 `admin` / `user` 两级）或资源所有者校验。

## 5. API 治理

1. 统一响应结构（`data`/`error`/`request_id`）。
2. 统一错误码映射（至少覆盖 400/401/403/404/409/422/500）。
3. API 版本化（`/api/v1`）。
4. 参数校验统一化（Gin binding + 自定义校验错误转换）。

## 6. 中间件体系

1. Recovery（防 panic 崩溃）。
2. Request ID。
3. 访问日志（结构化日志）。
4. 超时控制。
5. CORS。
6. 基础限流（按 IP 或 token）。

## 7. 可观测性

1. 健康检查：`/healthz`、`/readyz`。
2. Prometheus 指标：请求量、耗时、状态码。
3. 日志分级：`debug/info/warn/error`。

## 8. 工程质量

1. 单元测试：service、handler、repository（接口 mock）。
2. 集成测试：带 Postgres/Redis 容器或测试实例。
3. 覆盖成功路径 + 参数非法 + 资源不存在 + 权限失败。
4. `go test ./...` 可稳定通过。

## 9. 交付与运行

1. Dockerfile（多阶段构建）。
2. docker compose（api + postgres + redis）。
3. Makefile 提供 `run/test/lint/migrate-up/migrate-down/docker-up/docker-down`。

## 初始业务接口（第一版）

1. `POST /api/v1/auth/register`
2. `POST /api/v1/auth/login`
3. `GET /api/v1/me`
4. `POST /api/v1/todos`
5. `GET /api/v1/todos`
6. `PATCH /api/v1/todos/:id`
7. `DELETE /api/v1/todos/:id`

## 实施节奏（按天）

### Day 1：脚手架与配置

1. 初始化 `gin-backend` 独立模块与目录。
2. 接入配置加载、Gin 路由、基础中间件。
3. 提供 `healthz/readyz`。

验收：`go run ./gin-backend/cmd/server` 可启动并返回健康检查 200。

### Day 2：数据库与仓储层

1. 接入 Postgres 连接池与 migration。
2. 完成 Todo repository + service + handler。
3. 补齐基础 CRUD 单测。

验收：Todo API 可用，`go test ./...` 通过。

### Day 3：Redis 与鉴权

1. 接入 Redis 缓存，完成列表缓存与失效逻辑。
2. 接入 JWT 鉴权与用户上下文中间件。
3. 增加受保护路由与权限校验。

验收：带 token 调用通过，无 token 返回 401。

### Day 4：治理与观测

1. 统一错误码与响应结构。
2. 增加请求日志、request id、prometheus 指标。
3. 补齐失败路径测试与集成测试。

验收：指标可访问，日志包含 request id。

### Day 5：容器化与文档

1. 完成 Dockerfile + docker-compose。
2. 完成 README（启动、配置、示例 curl、故障排查）。
3. 完成最终回归验证。

验收：`docker compose up --build` 后 API 可完整访问。

## 启动配置（计划值）

1. `GIN_MODE=release`
2. `HTTP_PORT=8081`
3. `PG_DSN=postgres://postgres:postgres@localhost:5432/gin_backend?sslmode=disable`
4. `REDIS_ADDR=localhost:6379`
5. `JWT_SECRET=change-me`
6. `TOKEN_TTL_MINUTES=120`

## 质量门禁

1. `gofmt -w` 与 `go test ./...` 必须通过。
2. 至少提供成功与失败各一组 `curl` 示例。
3. 文档与代码同次提交，避免脱节。

## 下一步执行说明

本文件确认后，我将按 Day 1 开始在根目录创建 `gin-backend/` 独立工程，先完成可启动骨架，再逐步补齐数据库、缓存、鉴权和 Docker。
