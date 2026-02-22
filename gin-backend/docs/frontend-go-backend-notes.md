# 前端转 Go 后端学习笔记（基于 gin-backend）

> 目标：用前端视角理解这个仓库的后端体系。  
> 规则：每次讲解先沉淀到本文件，再继续下一讲。

## 第 1 讲：进程入口与生命周期（`cmd/server/main.go`）

### 1. 一句话结论
`cmd/server/main.go` 是后端进程的生命周期控制器，负责启动、监听退出信号、优雅停机和资源释放。

### 2. 请求链路定位
完整业务链路是：

`router -> middleware -> handler -> service -> repository`

`main.go` 不直接处理业务请求，它负责把这条链路“组装并跑起来”。

### 3. 关键代码定位（按启动顺序）
1. 创建根上下文：`cmd/server/main.go:21`
2. 装配应用依赖：`cmd/server/main.go:24` -> `internal/bootstrap/app.go:32`
3. 启动 HTTP 服务：`cmd/server/main.go:39`
4. 监听退出事件（异常/信号）：`cmd/server/main.go:56`
5. 优雅停机（带超时）：`cmd/server/main.go:64`、`cmd/server/main.go:68`
6. 资源释放（DB close）：`cmd/server/main.go:31`、`internal/bootstrap/app.go:117`

### 4. 前端类比理解
- 类比 `main.tsx`：负责应用入口。
- 类比“页面卸载清理”：进程退出前做清理，但这里清理的是 DB 连接和 HTTP server。
- 类比请求拦截链初始化：真正的接口处理在 router/middleware/handler 中，`main` 只负责把系统托管起来。

### 5. 前端可感知影响
1. 服务退出时不会粗暴断连接，减少联调时随机失败。
2. 若启动阶段依赖异常（例如 DB 不可达），服务会直接失败退出，避免“半可用”状态。
3. 停机有超时上限，不会无限卡住发布或重启流程。

### 6. 最小验证动作
```bash
make run
curl -i http://localhost:8081/healthz
```

然后在运行服务的终端按 `Ctrl+C`，观察日志中的信号接收与停机过程（对应优雅停机分支）。

### 7. 本讲涉及文件
- `cmd/server/main.go`
- `internal/bootstrap/app.go`
- `internal/transport/http/router/router.go`

## 第 2 讲：依赖装配（`internal/bootstrap/app.go`）

### 1. 一句话结论
`bootstrap.New(ctx)` 是后端的“应用组装工厂”，把配置、数据库、缓存、业务服务、路由一次性装配成可运行的 HTTP 服务。

### 2. 请求链路定位
请求链路仍然是：

`router -> middleware -> handler -> service -> repository`

`bootstrap` 的职责是把这条链路上每个组件实例化并连接起来。

### 3. 装配顺序（按代码真实顺序）
1. 加载配置：`internal/bootstrap/app.go:34` -> `internal/config/config.go:92`
2. 初始化日志和指标：`internal/bootstrap/app.go:40`、`internal/bootstrap/app.go:41`
3. 连接 Postgres + 确保表结构：`internal/bootstrap/app.go:44`、`internal/bootstrap/app.go:48`  
   对应实现：`internal/repository/postgres/db.go:14`、`internal/repository/postgres/db.go:37`
4. 组装 Postgres 仓储：`internal/bootstrap/app.go:55`、`internal/bootstrap/app.go:56`
5. 组装鉴权能力（JWT + AuthService）：`internal/bootstrap/app.go:59`、`internal/bootstrap/app.go:60`
6. 组装缓存能力（Redis 或 no-op）：`internal/bootstrap/app.go:63`、`internal/bootstrap/app.go:64`  
   降级实现：`internal/repository/redis/noop_cache.go:10`
7. 组装 Todo 业务服务：`internal/bootstrap/app.go:86`
8. 组装 handler：`internal/bootstrap/app.go:89` 到 `internal/bootstrap/app.go:91`
9. 构建 Gin 路由：`internal/bootstrap/app.go:94`
10. 构建 `http.Server`：`internal/bootstrap/app.go:101`

### 4. 前端类比理解
- 类比前端 `createApp()`：在一个地方统一注入配置、接口客户端、状态模块、路由。
- `repository` 接口注入（`internal/repository/interfaces.go:10`）类似“面向接口编程”，便于替换实现和测试 mock。
- Redis 开关是“能力可选但主链路不崩”：不开 Redis 时走 no-op 缓存，业务仍能跑。

### 5. 前端可感知影响
1. 启动即校验依赖，DB/Redis 异常会尽早失败，不把错误拖到首个请求。
2. 缓存是可插拔的：即使 Redis 不可用，主流程仍可回源数据库。
3. 中间件、路由、业务服务在启动时一次装配，运行时路径稳定、排障更清晰。

### 6. 可验证动作（建议做一次）
1. 默认模式启动（Redis 开启）：
```bash
make run
```
观察日志是否出现 `redis cache enabled`。

2. 关闭 Redis 再启动（验证 no-op 降级）：
```bash
REDIS_ENABLED=false make run
```
服务应仍能启动并处理接口，只是 Todo 列表不走 Redis 缓存。

### 7. 本讲涉及文件
- `internal/bootstrap/app.go`
- `internal/config/config.go`
- `internal/repository/postgres/db.go`
- `internal/repository/interfaces.go`
- `internal/repository/redis/client.go`
- `internal/repository/redis/noop_cache.go`

## 补充：repository 层是什么？

### 1. 一句话结论
repository 层是“数据访问适配层”：负责把业务需要的读写操作，转换成具体的数据库/缓存操作，对上层屏蔽存储细节。

### 2. 在本项目里的职责边界
- `service` 只依赖接口，不直接写 SQL：`internal/service/auth_service.go:25`、`internal/service/todo_service.go:30`
- repository 接口定义在：`internal/repository/interfaces.go:10`
- Postgres 实现在：`internal/repository/postgres/user_repository.go:14`、`internal/repository/postgres/todo_repository.go:13`
- Redis 缓存作为独立缓存仓储能力：`internal/repository/interfaces.go:28`、`internal/repository/redis/todo_cache.go:12`

### 3. 它解决的问题
1. 把 SQL、行扫描、存储错误细节集中管理（例如 `sql.ErrNoRows` -> `repository.ErrNotFound`）。
2. 让业务层聚焦“规则”，不关心“怎么查库”。
3. 提高可测试性：service 测试时可替换 mock repository。

### 4. 前端类比
可以类比前端里“API Client + 数据适配器”：
- service 像业务 usecase，不直接 `fetch/axios`；
- repository 像统一的数据访问层，负责请求细节和错误归一化。

### 5. 一个最小调用链例子（登录）
`handler.Login -> service.Login -> users.GetByEmail(repository) -> service 生成 JWT -> response`

对应代码：
- service 调用仓储：`internal/service/auth_service.go:91`
- 仓储执行 SQL：`internal/repository/postgres/user_repository.go:50`
