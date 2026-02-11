# Week07 认证与鉴权示例说明

## 目标

实现最小认证鉴权闭环：

- 注册与登录
- Token 签发与校验
- 基于角色的中间件鉴权

## 代码位置

- `examples/week07/authn_authz.go`
- `examples/week07/authn_authz_test.go`
- `examples/week07/cmd/main.go`

## 输入输出与错误码映射

### 登录成功

- 输入：`username/password`
- 输出：token 字符串

### 中间件鉴权

- `Authorization` 缺失或 token 非法：`401 Unauthorized`
- 角色不匹配：`403 Forbidden`
- 鉴权通过：业务 handler 正常返回 `200`

## 运行与验证

运行示例：

```bash
go run ./examples/week07/cmd
```

预期输出（token 会变化）：

```text
issued token: <token>
authorized => status=200 body=access granted
missing token => status=401 body=unauthorized
```

运行测试：

```bash
go test -v ./examples/week07
```
