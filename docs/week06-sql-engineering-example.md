# Week06 SQL 工程化示例说明

## 目标

演示 SQL 工程化中的选型思路与类型安全查询接口：

- `sqlc` 与 `gorm` 的最小选型决策
- 强类型参数（模拟 `sqlc` 生成 API）
- 查询行为可测试、可复用

## 代码位置

- `examples/week06/sql_engineering.go`
- `examples/week06/sql_engineering_test.go`
- `examples/week06/cmd/main.go`

## 关键行为

1. `ChooseDataAccessTool` 根据约束返回 `sqlc` 或 `gorm`。
2. `CreateUserParams` / `ListUsersByPrefixParams` 展示强类型查询参数。
3. 空用户名创建返回 `ErrInvalidName`。

## 运行与验证

运行示例：

```bash
go run ./examples/week06/cmd
```

预期输出：

```text
chosen tool: sqlc
prefix=al users: [{1 alice} {2 allen}]
```

运行测试：

```bash
go test -v ./examples/week06
```
