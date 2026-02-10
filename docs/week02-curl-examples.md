# Week 02 curl 验证示例（DELETE）

## 1. 启动服务

```bash
go run ./cmd/api
```

默认地址：`http://localhost:8080`

## 2. 先创建一条 todo（用于后续删除）

```bash
curl -i -X POST http://localhost:8080/api/v1/todos \
  -H 'Content-Type: application/json' \
  -d '{"title":"week02 delete demo"}'
```

从响应里记录 `id`（例如 `20260210112233.123456789`）。

## 3. 删除成功示例（200）

```bash
curl -i -X DELETE http://localhost:8080/api/v1/todos/{id}
```

预期：

- 状态码：`200 OK`
- 响应体包含：`"deleted": true`

## 4. 非法 id 示例（400）

```bash
curl -i -X DELETE http://localhost:8080/api/v1/todos/abc
```

预期：

- 状态码：`400 Bad Request`
- 响应体：`{"error":"invalid todo id"}`

## 5. 不存在 id 示例（404）

```bash
curl -i -X DELETE http://localhost:8080/api/v1/todos/20000101000000.000000000
```

预期：

- 状态码：`404 Not Found`
- 响应体：`{"error":"todo not found"}`
