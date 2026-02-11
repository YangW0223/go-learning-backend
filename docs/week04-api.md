# Week 04 接口契约说明（分层示例）

## 接口定义

- 方法：`PATCH`
- 路径：`/api/v1/todos/{id}/done`
- Content-Type：`application/json`

## 请求参数

- 路径参数 `id`：必须匹配 `^[1-9][0-9]{0,8}$`
- 合法示例：`1`、`99`、`2026`
- 非法示例：`abc`、`01`、`1x`

## 成功响应

- 状态码：`200 OK`
- 响应体：

```json
{
  "data": {
    "id": "1",
    "title": "read layering notes",
    "done": true
  },
  "error": null
}
```

## 错误响应

1. 参数错误（400）
- 场景：`id` 非法，或路径不符合 `/api/v1/todos/{id}/done`
- 状态码：`400 Bad Request`
- 响应体：

```json
{
  "data": null,
  "error": "invalid todo id"
}
```

2. 资源不存在（404）
- 场景：`id` 格式合法，但 todo 不存在
- 状态码：`404 Not Found`
- 响应体：

```json
{
  "data": null,
  "error": "todo not found"
}
```

3. 服务内部错误（500）
- 场景：store 出现未知错误
- 状态码：`500 Internal Server Error`
- 响应体：

```json
{
  "data": null,
  "error": "internal server error"
}
```

## 分层职责映射

1. `handler`：`examples/week04/handler.go`
- 负责方法校验、路径解析、状态码映射、JSON 响应。

2. `service`：`examples/week04/service.go`
- 负责业务规则（id 校验、业务错误分类）。

3. `store`：`examples/week04/store.go`
- 负责数据读写（本示例为内存实现）。
