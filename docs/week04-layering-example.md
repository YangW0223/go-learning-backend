# Week04 分层重构示例说明

## 目标

通过一个最小 Todo 完成接口示例，演示 `handler -> service -> store` 分层职责：

- `handler`：只做 HTTP 协议解析与状态码映射。
- `service`：只做业务规则（参数校验、业务错误分类）。
- `store`：只做数据读写（本示例使用内存实现）。

## 代码位置

- `examples/week04/handler.go`
- `examples/week04/service.go`
- `examples/week04/store.go`
- `examples/week04/cmd/main.go`

接口契约与手动验证文档：

- `docs/week04-api.md`
- `docs/week04-curl-examples.md`

## 运行示例

在项目根目录执行：

```bash
go run ./examples/week04/cmd
```

预期输出示例：

```text
PATCH /api/v1/todos/1/done => status=200 body={"data":{"id":"1","title":"read layering notes","done":true},"error":null}
PATCH /api/v1/todos/abc/done => status=400 body={"data":null,"error":"invalid todo id"}
PATCH /api/v1/todos/99/done => status=404 body={"data":null,"error":"todo not found"}
```

## 测试命令

```bash
go test -v ./examples/week04
```

测试覆盖：

- 成功路径（200）
- 参数非法（400）
- 资源不存在（404）
- 内部错误（500）
