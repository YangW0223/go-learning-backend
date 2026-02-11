# Backend 服务体系改造计划（internal）

更新时间：2026-02-11

## 目标

在不破坏现有 Todo API 行为的前提下，为 `internal` 引入更完整的后端服务基础能力，并新增 Redis 与 Docker 运行配置，形成可本地开发、可容器化运行、可逐步扩展到生产形态的骨架。

## 本次改造范围

1. 分层增强：在现有 `handler -> store` 基础上增加 `service` 层。
2. 配置中心：新增 `internal/config`，统一读取环境变量并提供默认值。
3. Redis 接入：新增缓存抽象与 Redis 实现，优先用于 Todo 列表缓存。
4. 应用启动：`cmd/api/main.go` 改为按配置装配依赖；Redis 可开关。
5. Docker 化：新增 `Dockerfile` 与 `docker-compose`（应用 + Redis）。
6. 测试与文档：补齐核心单元测试并更新 README 的运行方式。

## 设计原则

1. 向后兼容：现有路由、返回结构、错误码语义保持不变。
2. 渐进改造：先把基础骨架搭好，再逐步引入数据库、鉴权、观测性。
3. 可降级：Redis 不可用时可通过配置关闭；关闭后服务仍可运行。
4. 可验证：每一步改动都对应可执行验证命令。

## 启动配置（可直接执行）

### 方式 1：本地启动（不启用 Redis）

```bash
PORT=8080 REDIS_ENABLED=false go run ./cmd/api
```

### 方式 2：本地启动（启用 Redis）

前置条件：本机已有可访问 Redis（默认 `localhost:6379`）。

```bash
PORT=8080 \
REDIS_ENABLED=true \
REDIS_ADDR=localhost:6379 \
REDIS_DB=0 \
REDIS_CACHE_TTL_SECONDS=30 \
REDIS_DIAL_TIMEOUT_MS=1000 \
REDIS_IO_TIMEOUT_MS=1000 \
go run ./cmd/api
```

### 方式 3：Docker Compose 启动（api + redis）

```bash
docker compose --env-file docker/.env.example -f docker/docker-compose.redis.yaml up --build
```

停止并清理：

```bash
docker compose --env-file docker/.env.example -f docker/docker-compose.redis.yaml down -v
```

### 启动配置项说明

1. `PORT`：HTTP 服务监听端口，默认 `8080`。
2. `REDIS_ENABLED`：是否启用 Redis 缓存，默认 `false`。
3. `REDIS_ADDR`：Redis 地址，默认 `localhost:6379`。
4. `REDIS_PASSWORD`：Redis 密码，默认空字符串。
5. `REDIS_DB`：Redis DB 编号，默认 `0`。
6. `REDIS_CACHE_TTL_SECONDS`：Todo 列表缓存 TTL（秒），默认 `30`。
7. `REDIS_DIAL_TIMEOUT_MS`：Redis 连接超时（毫秒），默认 `1000`。
8. `REDIS_IO_TIMEOUT_MS`：Redis 读写超时（毫秒），默认 `1000`。

## 分阶段实施

## 阶段 1：文档与骨架

1. 新增本计划文档，明确改造边界与顺序。
2. 新增目录：
   - `internal/config`
   - `internal/service`
   - `internal/cache`

验收：目录与计划文档存在，结构清晰。

## 阶段 2：配置层与服务层

1. 在 `internal/config` 增加 App 配置加载逻辑：
   - `PORT`
   - `REDIS_ENABLED`
   - `REDIS_ADDR`
   - `REDIS_PASSWORD`
   - `REDIS_DB`
   - `REDIS_CACHE_TTL_SECONDS`
2. 在 `internal/service` 新增 TodoService：
   - 封装 Create/List/MarkDone/Delete
   - 统一处理缓存失效逻辑
3. `handler` 改为依赖 service（保留兼容构造函数）。

验收：`go test ./...` 通过，现有 API 测试行为一致。

## 阶段 3：Redis 缓存接入

1. 在 `internal/cache` 定义缓存接口与 no-op 实现。
2. 实现 Redis Todo 缓存：
   - 列表读取命中缓存
   - 创建/更新/删除后清理缓存键
   - 提供 `Ping` 用于启动阶段连通性校验
3. `cmd/api/main.go` 根据配置决定是否启用 Redis 缓存。

验收：
1. `REDIS_ENABLED=false` 时服务可正常运行。
2. `REDIS_ENABLED=true` 且 Redis 可达时服务正常运行。

## 阶段 4：Docker 运行配置

1. 新增根目录 `Dockerfile`（多阶段构建）。
2. 新增 `docker/docker-compose.redis.yaml`：
   - `api` 服务
   - `redis` 服务（含健康检查）
3. 新增 `docker/.env.example`，提供常用变量示例。
4. 增补 `Makefile` 的 Docker 命令（up/down/logs）。

验收：
1. `docker compose --env-file docker/.env.example -f docker/docker-compose.redis.yaml up --build` 可启动。
2. `curl http://localhost:8080/healthz` 返回 200。

## 阶段 5：文档与回归

1. 更新 `README.md`：
   - 本地运行
   - Docker 运行
   - Redis 开关说明
2. 执行并记录验证：
   - `gofmt -w ...`
   - `go test ./...`

验收：README 可指导首次运行，验证命令有明确结果。

## 非本次范围（后续迭代）

1. PostgreSQL 持久化与 migration。
2. 鉴权与权限（JWT/RBAC）。
3. 结构化日志、metrics、trace 全量接入。
4. CI/CD 与镜像发布流水线。

## 风险与应对

1. 风险：Redis 连接失败导致服务不可用。
   - 应对：通过 `REDIS_ENABLED` 开关控制；启用模式下启动前 `Ping` 快速失败。
2. 风险：分层重构影响现有测试。
   - 应对：保留向后兼容构造函数，优先修复断言。
3. 风险：Docker 配置与本地端口冲突。
   - 应对：`.env` 支持自定义端口映射。
