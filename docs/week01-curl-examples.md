# Week 01 接口手动验证（curl）

## 1. 启动服务

在项目根目录执行：

```bash
go run ./cmd/api
```

默认监听端口是 `8080`。下面命令都假设服务地址为 `http://localhost:8080`。

## 2. 验证 `GET /ping`

```bash
curl -i http://localhost:8080/ping
```

预期：

- 状态码：`200 OK`
- 响应体：`{"message":"pong"}`

## 3. 创建 todo：`POST /api/v1/todos`

```bash
curl -i -X POST http://localhost:8080/api/v1/todos \
  -H 'Content-Type: application/json' \
  -d '{"title":"learn go basics"}'
```

预期：

- 状态码：`201 Created`
- 响应体里包含 `id`、`title`、`done`、`created_at`

你可以从响应体里复制 `id`，用于后续“标记完成”。

## 4. 查询 todo 列表：`GET /api/v1/todos`

```bash
curl -i http://localhost:8080/api/v1/todos
```

预期：

- 状态码：`200 OK`
- 响应体：JSON 数组（至少包含刚刚创建的 todo）

## 5. 标记完成：`PATCH /api/v1/todos/{id}/done`

把下面 `{id}` 替换成你在第 3 步拿到的真实 id：

```bash
curl -i -X PATCH http://localhost:8080/api/v1/todos/{id}/done
```

预期：

- 状态码：`200 OK`
- 响应体中的 `done` 字段为 `true`

## 6. 常见错误示例

空标题：

```bash
curl -i -X POST http://localhost:8080/api/v1/todos \
  -H 'Content-Type: application/json' \
  -d '{"title":"   "}'
```

预期：

- 状态码：`400 Bad Request`
- 响应体：`{"error":"title is required"}`

不存在的 id：

```bash
curl -i -X PATCH http://localhost:8080/api/v1/todos/not-exist-id/done
```

预期：

- 状态码：`404 Not Found`
- 响应体：`{"error":"todo not found"}`
