# Week 02 删除接口行为说明

## 接口定义

- 方法：`DELETE`
- 路径：`/api/v1/todos/{id}`
- Content-Type：`application/json`

## 请求参数

- 路径参数 `id`：必须符合当前项目规则 `^\d{14}\.\d{9}$`
- 示例：`20260210112233.123456789`

## 成功响应

- 状态码：`200 OK`
- 响应体：

```json
{
  "data": {
    "id": "20260210112233.123456789",
    "deleted": true
  },
  "error": null
}
```

## 错误响应

1. 参数错误（400）

- 场景：`id` 为空或格式不合法
- 状态码：`400 Bad Request`
- 响应体：

```json
{
  "error": "invalid todo id"
}
```

2. 资源不存在（404）

- 场景：`id` 格式合法，但 todo 不存在
- 状态码：`404 Not Found`
- 响应体：

```json
{
  "error": "todo not found"
}
```

3. 服务内部错误（500）

- 场景：存储层出现未知错误
- 状态码：`500 Internal Server Error`
- 响应体：

```json
{
  "error": "failed to delete todo"
}
```

## 处理流程（实现对应）

1. 路由匹配：`internal/app/router.go`
2. 参数提取：从 `/api/v1/todos/{id}` 中提取 `id`
3. 参数校验：`internal/handler/todo.go` 的 `isValidTodoID`
4. 业务执行：`store.Delete(id)`
5. JSON 响应：按 200/400/404/500 返回
