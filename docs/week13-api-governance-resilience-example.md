# Week13 API 治理与稳定性示例说明

## 目标

补齐生产化 API 治理能力：

- 统一响应结构（`request_id/data/error`）
- 幂等键重放（`Idempotency-Key`）
- 用户维度限流（429）
- 优雅停机入口

## 代码位置

- `examples/week13/api_governance_resilience.go`
- `examples/week13/api_governance_resilience_test.go`
- `examples/week13/cmd/main.go`

## 接口定义

- 方法：`POST`
- 路径：`/api/v1/orders`
- Header：`X-User-ID`（限流键）
- Header：`Idempotency-Key`（可选，开启幂等）
- Body：`{"item":"book"}`

## 状态码与错误码

1. 成功创建：`201 Created`
2. 入参非法：`400 Bad Request` + `INVALID_ITEM` 或 `INVALID_JSON`
3. 方法不允许：`405 Method Not Allowed` + `METHOD_NOT_ALLOWED`
4. 限流触发：`429 Too Many Requests` + `RATE_LIMITED`
5. 内部错误：`500 Internal Server Error` + `INTERNAL_ERROR`

## 运行与验证

运行示例：

```bash
go run ./examples/week13/cmd
```

预期输出：

```text
first => status=201 ...
replay => status=201 replay=1 ...
third-no-key => status=201 ...
```

运行测试：

```bash
go test -v ./examples/week13
```
