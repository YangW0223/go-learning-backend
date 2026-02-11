# Week10 测试体系示例说明

## 目标

演示可回归测试体系的最小闭环：

- service 层表驱动测试
- handler 场景测试（成功/400/404）
- 端到端风格流程测试（创建后查询）

## 代码位置

- `examples/week10/testing_system.go`
- `examples/week10/testing_system_test.go`
- `examples/week10/cmd/main.go`

## 输入输出与错误码

接口：

1. `POST /api/v1/todos`
- 成功：`201`
- title 非法：`400`

2. `GET /api/v1/todos/{id}`
- 成功：`200`
- id 非法：`400`
- 资源不存在：`404`

## 运行与验证

运行示例：

```bash
go run ./examples/week10/cmd
```

预期输出：

```text
POST /api/v1/todos => 201 ...
GET /api/v1/todos/1 => 200 ...
```

运行测试：

```bash
go test -v ./examples/week10
```
