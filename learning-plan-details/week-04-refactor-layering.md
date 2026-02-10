# Week 04：重构与分层

## 本周目标

形成 `handler/service/store` 分层意识，减少 handler 业务膨胀。

## 详细步骤

1. 识别当前 handler 中的业务逻辑、校验逻辑、存储逻辑。
2. 提炼 service 接口，明确输入输出。
3. 把业务逻辑从 handler 移动到 service。
4. store 层只负责数据访问，不处理业务规则。
5. 通过依赖注入把 handler 与 service 解耦。
6. 为 service 关键路径增加单测。

## 建议实践清单

- handler 只做协议转换（HTTP <-> 业务对象）。
- service 只做业务规则。
- store 只做数据读写。

## 验收清单

- [ ] handler 代码显著简化。
- [ ] service 关键路径有单测覆盖。
- [ ] 依赖方向清晰（handler -> service -> store）。

## 产出物

- 分层重构后的代码。
- service 层测试用例。
- 一份架构变化说明。

## 常见风险与排查

- 过度抽象：接口数量控制在当前需求最小集合。
- 包循环依赖：按分层单向引用，必要时提取 `domain` 类型。
