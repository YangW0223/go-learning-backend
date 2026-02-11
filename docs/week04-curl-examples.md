# Week 04 curl 验证示例（PATCH）

## 1. 启动 Week04 示例服务

```bash
go run ./examples/week04/cmd -mode server -addr :18084
```

默认地址：`http://localhost:18084`

## 2. 标记完成成功示例（200）

```bash
curl -i -X PATCH http://localhost:18084/api/v1/todos/1/done
```

预期：

- 状态码：`200 OK`
- 响应体包含：`"done":true`

## 3. 非法 id 示例（400）

```bash
curl -i -X PATCH http://localhost:18084/api/v1/todos/abc/done
```

预期：

- 状态码：`400 Bad Request`
- 响应体：`{"data":null,"error":"invalid todo id"}`

## 4. 不存在 id 示例（404）

```bash
curl -i -X PATCH http://localhost:18084/api/v1/todos/99/done
```

预期：

- 状态码：`404 Not Found`
- 响应体：`{"data":null,"error":"todo not found"}`

## 5. 方法错误示例（405）

```bash
curl -i -X GET http://localhost:18084/api/v1/todos/1/done
```

预期：

- 状态码：`405 Method Not Allowed`
- 响应体：`{"data":null,"error":"method not allowed"}`
