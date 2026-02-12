# 声明伪目标，避免与同名真实文件冲突。
.PHONY: run test fmt tidy migrate-up migrate-down docker-up docker-down docker-logs

# 在本机直接运行 API（不使用 Docker）。
run:
	go run ./cmd/server

# 运行当前模块下所有 Go 测试。
test:
	go test ./...

# 格式化全部 Go 文件（排除 vendor 目录）。
fmt:
	gofmt -w $(shell find . -name '*.go' -not -path './vendor/*')

# 清理并同步 go.mod/go.sum 与当前依赖导入。
tidy:
	go mod tidy

# 执行待应用的数据库迁移。
migrate-up:
	go run ./scripts/migrate -direction up

# 按迁移工具规则回滚最近的迁移步骤。
migrate-down:
	go run ./scripts/migrate -direction down

# 启动完整本地栈（API + Postgres + Redis），并重新构建镜像。
# 使用 docker/.env.example 作为 compose 的 env-file 输入。
docker-up:
	docker compose --env-file docker/.env.example -f docker/docker-compose.yaml up --build

# 停止栈并删除容器及命名卷。
# 注意：`-v` 会删除 Postgres/Redis 的本地持久化数据。
docker-down:
	docker compose --env-file docker/.env.example -f docker/docker-compose.yaml down -v

# 持续跟踪核心服务日志输出。
docker-logs:
	docker compose --env-file docker/.env.example -f docker/docker-compose.yaml logs -f api postgres redis
