# Week11 交付与部署示例说明

## 目标

演示“可部署、可回滚”的最小链路：

- 环境变量加载
- 健康检查与就绪检查
- 版本信息与回滚说明

## 代码位置

- `examples/week11/delivery_deployment.go`
- `examples/week11/delivery_deployment_test.go`
- `examples/week11/cmd/main.go`

## 关键行为

1. `LoadConfigFromEnv` 缺失关键配置时返回 `ErrMissingConfig`。
2. `/healthz` 反映进程活性。
3. `/readyz` 反映依赖配置是否就绪。
4. `BuildRollbackPlan` 给出最小回滚描述。

## 运行与验证

运行示例：

```bash
go run ./examples/week11/cmd
```

预期输出：

```text
loaded config: port=8080 env=dev version=v1.1.0
GET /healthz => 200 ok
GET /readyz => 200 ready
GET /version => 200 v1.1.0
rollback from demo:v2 to demo:v1
```

运行测试：

```bash
go test -v ./examples/week11
```
