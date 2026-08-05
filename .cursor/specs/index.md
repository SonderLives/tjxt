# tjxt 项目规格文档索引

> 最后更新：2026-08-05 | 维护者：@team

## 快速导航
- [架构全览](00-architecture/overview.md)
- [核心约束](01-conventions/go-zero-rules.md)
- [服务清单](02-services/)
- [共享契约](03-shared/)
- [本地开发指南](05-development/quickstart.md)

## 服务级文档对照表

| 服务 | API Spec | RPC Spec | 数据模型 | 业务规则 | 配置 |
|------|----------|----------|----------|----------|------|
| auth | [✅](02-services/auth/api-spec.md) | [✅](02-services/auth/rpc-spec.md) | [✅](02-services/auth/data-model.md) | [✅](02-services/auth/business-rules.md) | [✅](02-services/auth/configs.md) |
| course | [✅](02-services/course/api-spec.md) | [✅](02-services/course/rpc-spec.md) | [✅](02-services/course/data-model.md) | [✅](02-services/course/business-rules.md) | [✅](02-services/course/configs.md) |
| trade | [✅](02-services/trade/api-spec.md) | [✅](02-services/trade/rpc-spec.md) | [✅](02-services/trade/data-model.md) | [✅](02-services/trade/business-rules.md) | [✅](02-services/trade/configs.md) |
| learning | [✅](02-services/learning/api-spec.md) | [✅](02-services/learning/rpc-spec.md) | [✅](02-services/learning/data-model.md) | [✅](02-services/learning/business-rules.md) | [✅](02-services/learning/configs.md) |
| pay | [✅](02-services/pay/api-spec.md) | [✅](02-services/pay/rpc-spec.md) | [✅](02-services/pay/data-model.md) | [✅](02-services/pay/business-rules.md) | [✅](02-services/pay/configs.md) |
| media | [✅](02-services/media/api-spec.md) | [✅](02-services/media/rpc-spec.md) | [✅](02-services/media/data-model.md) | [✅](02-services/media/business-rules.md) | [✅](02-services/media/configs.md) |
| promotion | [✅](02-services/promotion/api-spec.md) | [✅](02-services/promotion/rpc-spec.md) | [✅](02-services/promotion/data-model.md) | [✅](02-services/promotion/business-rules.md) | [✅](02-services/promotion/configs.md) |
| message | [✅](02-services/message/api-spec.md) | [✅](02-services/message/rpc-spec.md) | [✅](02-services/message/data-model.md) | [✅](02-services/message/business-rules.md) | [✅](02-services/message/configs.md) |
| exam | [✅](02-services/exam/api-spec.md) | [✅](02-services/exam/rpc-spec.md) | [✅](02-services/exam/data-model.md) | [✅](02-services/exam/business-rules.md) | [✅](02-services/exam/configs.md) |
| search | [✅](02-services/search/api-spec.md) | [✅](02-services/search/rpc-spec.md) | [✅](02-services/search/data-model.md) | [✅](02-services/search/business-rules.md) | [✅](02-services/search/configs.md) |
| user | [✅](02-services/user/api-spec.md) | [✅](02-services/user/rpc-spec.md) | [✅](02-services/user/data-model.md) | [✅](02-services/user/business-rules.md) | [✅](02-services/user/configs.md) |
| data | [✅](02-services/data/api-spec.md) | [✅](02-services/data/rpc-spec.md) | [✅](02-services/data/data-model.md) | [✅](02-services/data/business-rules.md) | [✅](02-services/data/configs.md) |

## 变更日志
- 2026-08-05 初始化骨架，搬运 MEMORY.md 与 .cursor/rules 内容
- 2026-08-05 从 OpenAPI 规范提取各服务 HTTP API 接口清单