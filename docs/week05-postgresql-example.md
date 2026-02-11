# Week05 PostgreSQL 基础示例说明

## 目标

通过内存实现模拟 PostgreSQL 的关键能力：

- 事务边界（提交/回滚）
- `Create` 与 `MarkDone` 在同一事务中的一致性
- 分页查询（`page/size`）

## 代码位置

- `examples/week05/postgresql_basics.go`
- `examples/week05/postgresql_basics_test.go`
- `examples/week05/cmd/main.go`

## 关键行为

1. `WithTx` 成功时提交，失败时回滚。
2. `ListTodos(page,size)` 对非法分页参数返回 `ErrInvalidPage`。
3. `MarkDone` 对不存在资源返回 `ErrTodoNotFound`。

## 运行与验证

运行示例：

```bash
go run ./examples/week05/cmd
```

预期输出：

```text
id=2 title="implement mark done" done=false
id=1 title="design todos table" done=true
```

运行测试：

```bash
go test -v ./examples/week05
```
