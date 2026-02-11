# Gin Backend

基于 Gin 的独立后端项目，包含配置中心、分层架构、JWT 鉴权、Postgres 存储、Redis 缓存、中间件体系、基础观测与 Docker 运行配置。

## 目录

- `cmd/server`：服务入口
- `internal/config`：环境变量配置
- `internal/transport/http`：路由、handler、中间件、DTO
- `internal/service`：业务编排层
- `internal/repository/postgres`：Postgres 仓储实现
- `internal/repository/redis`：Redis 缓存实现
- `internal/auth`：密码与 JWT
- `internal/observability`：日志与指标
- `migrations`：SQL migration
- `scripts/migrate`：migration 执行脚本

## 快速开始

### 1. 本地依赖

需要本机可用：

- PostgreSQL
- Redis
- Go 1.23+

### 2. 环境变量（最小集）

```bash
export APP_ENV=dev
export HTTP_PORT=8081
export PG_DSN='postgres://postgres:postgres@localhost:5432/gin_backend?sslmode=disable'
export REDIS_ENABLED=true
export REDIS_ADDR=localhost:6379
export JWT_SECRET='change-me-please'
```

### 3. 初始化表结构

```bash
make migrate-up
```

### 4. 启动服务

```bash
make run
```

## Docker 启动

```bash
make docker-up
```

停止并清理：

```bash
make docker-down
```

## 核心接口

- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `GET /api/v1/me`
- `POST /api/v1/todos`
- `GET /api/v1/todos`
- `PATCH /api/v1/todos/:id`
- `DELETE /api/v1/todos/:id`
- `GET /healthz`
- `GET /readyz`
- `GET /metrics`

## curl 验证

### 1. 注册

```bash
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"u1@example.com","password":"password123"}'
```

### 2. 登录

```bash
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"u1@example.com","password":"password123"}'
```

### 3. 使用 token 创建 Todo

```bash
TOKEN='<replace-token>'
curl -X POST http://localhost:8081/api/v1/todos \
  -H "Authorization: Bearer ${TOKEN}" \
  -H 'Content-Type: application/json' \
  -d '{"title":"learn gin architecture"}'
```

### 4. 无 token 访问失败示例（401）

```bash
curl http://localhost:8081/api/v1/todos
```

## 常用命令

```bash
make fmt
make test
make tidy
make migrate-up
make migrate-down
```

## 默认配置项

- `APP_NAME=gin-backend`
- `APP_ENV=dev`
- `HTTP_PORT=8081`
- `HTTP_REQUEST_TIMEOUT_MS=3000`
- `PG_DSN=postgres://postgres:postgres@localhost:5432/gin_backend?sslmode=disable`
- `REDIS_ENABLED=true`
- `REDIS_ADDR=localhost:6379`
- `REDIS_CACHE_TTL_SECONDS=30`
- `JWT_SECRET=change-me`
- `TOKEN_TTL_MINUTES=120`
