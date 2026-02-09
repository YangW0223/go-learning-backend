# Go Learning Backend

一个面向前端开发者的 Go 后端学习项目。目标是从 0 到 1 完成一个可运行、可测试、可扩展的 API 服务。

## 当前能力

- HTTP 服务（标准库）
- 健康检查: `GET /healthz`
- Todo API（内存存储）
  - `POST /api/v1/todos`
  - `GET /api/v1/todos`
  - `PATCH /api/v1/todos/{id}/done`
- 单元测试
- Makefile 常用命令

## 快速开始

```bash
make run
```

默认监听 `:8080`。

### 示例请求

```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H 'Content-Type: application/json' \
  -d '{"title":"learn go basics"}'

curl http://localhost:8080/api/v1/todos

curl -X PATCH http://localhost:8080/api/v1/todos/<todo_id>/done
```

## 学习路线

12 周学习计划见 `LEARNING_PLAN_12_WEEKS.md`。
# go-learning-backend
