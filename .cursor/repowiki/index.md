# tjxt 天机学堂 — 项目规格文档索引

> 最后更新：2026-08-06 | 维护者：@team
> 覆盖：13 个微服务 · 79 篇文档 · 约 11,000 行

**tjxt** 是基于 **go-zero** 的在线教育微服务系统，采用 **API + RPC 分层**，通过 **etcd** 服务发现，**MySQL + Redis** 存储，**Redis Stream** 做事件总线。无独立网关。

---

## 快速导航

| 我想… | 去这里 |
|-------|--------|
| 了解整体架构 | [架构全览](00-architecture/overview.md) · [服务拓扑](00-architecture/service-topology.md) · [数据域划分](00-architecture/data-domains.md) |
| 开始写代码 | [本地开发快速上手](05-development/quickstart.md) |
| 知道能改什么、不能改什么 | [go-zero 开发约束](01-conventions/go-zero-rules.md) ⭐ |
| 定义接口 | [API 契约规范](01-conventions/api-contracts.md) · [代码风格](01-conventions/code-style.md) |
| 查错误码 / 事件 / 公共库 | [错误码](03-shared/error-codes.md) · [MQ 事件](03-shared/mq-events.md) · [pkg 契约](03-shared/pkg-contracts.md) |
| 搭本地环境 | [Docker 基础设施](04-infra/docker-compose.md) |
| **看当前进度和坑** | [**实现进度与已知缺口**](06-status/implementation-status.md) ⭐ |

---

## 服务级文档对照表

每个服务 5 篇：**api-spec**（HTTP 接口）· **rpc-spec**（gRPC 方法）· **data-model**（表结构与 Model）· **business-rules**（业务规则）· **configs**（配置项）

| 服务 | API | RPC | 数据模型 | 业务规则 | 配置 | 实现进度 |
|------|-----|-----|----------|----------|------|----------|
| **auth** 认证鉴权 | [📄](02-services/auth/api-spec.md) | [📄](02-services/auth/rpc-spec.md) | [📄](02-services/auth/data-model.md) | [📄](02-services/auth/business-rules.md) | [📄](02-services/auth/configs.md) | ✅ 37/37 |
| **user** 用户中心 | [📄](02-services/user/api-spec.md) | [📄](02-services/user/rpc-spec.md) | [📄](02-services/user/data-model.md) | [📄](02-services/user/business-rules.md) | [📄](02-services/user/configs.md) | ✅ 28/28 |
| **pay** 支付 | [📄](02-services/pay/api-spec.md) | [📄](02-services/pay/rpc-spec.md) | [📄](02-services/pay/data-model.md) | [📄](02-services/pay/business-rules.md) | [📄](02-services/pay/configs.md) | ✅ 33/33 |
| **promotion** 优惠券 | [📄](02-services/promotion/api-spec.md) | [📄](02-services/promotion/rpc-spec.md) | [📄](02-services/promotion/data-model.md) | [📄](02-services/promotion/business-rules.md) | [📄](02-services/promotion/configs.md) | ✅ 32/32 |
| **remark** 评论点赞 | [📄](02-services/remark/api-spec.md) | [📄](02-services/remark/rpc-spec.md) | [📄](02-services/remark/data-model.md) | [📄](02-services/remark/business-rules.md) | [📄](02-services/remark/configs.md) | ✅ 4/4 |
| **course** 课程 | [📄](02-services/course/api-spec.md) | [📄](02-services/course/rpc-spec.md) | [📄](02-services/course/data-model.md) | [📄](02-services/course/business-rules.md) | [📄](02-services/course/configs.md) | ✅ 76/76 |
| **trade** 交易订单 | [📄](02-services/trade/api-spec.md) | [📄](02-services/trade/rpc-spec.md) | [📄](02-services/trade/data-model.md) | [📄](02-services/trade/business-rules.md) | [📄](02-services/trade/configs.md) | ✅ 73/73 |
| **message** 站内信 | [📄](02-services/message/api-spec.md) | [📄](02-services/message/rpc-spec.md) | [📄](02-services/message/data-model.md) | [📄](02-services/message/business-rules.md) | [📄](02-services/message/configs.md) | ✅ 37/37 |
| **learning** 学习进度 | [📄](02-services/learning/api-spec.md) | [📄](02-services/learning/rpc-spec.md) | [📄](02-services/learning/data-model.md) | [📄](02-services/learning/business-rules.md) | [📄](02-services/learning/configs.md) | ✅ 20/20 |
| **media** 媒资文件 | [📄](02-services/media/api-spec.md) | [📄](02-services/media/rpc-spec.md) | [📄](02-services/media/data-model.md) | [📄](02-services/media/business-rules.md) | [📄](02-services/media/configs.md) | ✅ 20/20 |
| **exam** 考试题库 | [📄](02-services/exam/api-spec.md) | [📄](02-services/exam/rpc-spec.md) | [📄](02-services/exam/data-model.md) | [📄](02-services/exam/business-rules.md) | [📄](02-services/exam/configs.md) | ✅ 14/14 |
| **data** 统计大屏 | [📄](02-services/data/api-spec.md) | [📄](02-services/data/rpc-spec.md) | [📄](02-services/data/data-model.md) | [📄](02-services/data/business-rules.md) | [📄](02-services/data/configs.md) | ✅ 12/12 |
| **search** 全文检索 | [📄](02-services/search/api-spec.md) | [📄](02-services/search/rpc-spec.md) | [📄](02-services/search/data-model.md) | [📄](02-services/search/business-rules.md) | [📄](02-services/search/configs.md) | ✅ 8/8 |

**总计**：394/394 logic 已实现（100%），服务维度 13/13 完成。

---

## ⚠️ 阅读须知

1. **文档覆盖 ≠ 功能实现（仍适用）**。全部 13 个服务的 logic 均已落地并编译通过（2026-08-06 复核）。各 `business-rules.md` 开头的「实现状态」已校正为「已落地实现」，原「📋 设计意图（待实现）」章节标注已改为「📋 设计意图（契约推导）」——这些规则依据 proto/DDL/.api 推导，仍建议对照源码最终确认。已知缺口（media 对象存储 mock、pay 网关占位、跨域 RPC 未接线、data 配置缺段等）见 [实现进度与已知缺口](06-status/implementation-status.md)。
2. **契约先行**。proto、DDL、`.api`、`docs/tjxt.openapi.json` 是真实且完整的，因此 `rpc-spec` / `data-model` / `configs` 三类文档对所有服务都可信。
3. **已知坑请先读** [实现进度与已知缺口](06-status/implementation-status.md)，其中记录了目录结构偏离、缺表、缓存失效、跨服务未接线等 20+ 项问题。
4. 各服务 `api-spec.md` 由 OpenAPI 机械提取，响应体列残留 URL 编码（`%C2%AB` 即 `«`），可读性待改进。

---

## 目录结构

```
.cursor/repowiki/
├── index.md                    # 本文件
├── 00-architecture/            # 架构全览、服务拓扑、数据域
├── 01-conventions/             # go-zero 约束、API 契约、代码风格、Git 流程
├── 02-services/<svc>/          # 13 个服务 × 5 篇
├── 03-shared/                  # 错误码、MQ 事件、pkg 公共契约
├── 04-infra/                   # Docker 基础设施
├── 05-development/             # 快速上手
└── 06-status/                  # 实现进度与已知缺口
```

---

## 变更日志

- **2026-08-06（模块拆分）** 将「每服务单 module」重构为 **go-zero 官方标准结构**（每 api/rpc 独立 `go.mod`）：删除聚合 `apps/<svc>/go.mod` 与孤儿 `apps/data/api/go.mod`，新增 26 个服务模块 `go.mod`，重写 `go.work` 聚合 28 个 use 模块；全模块 `go build ./...` 通过。同步 [架构全览](00-architecture/overview.md)（新增模块结构章节、修正端口表为 13 服务）、[go-zero 规范](01-conventions/go-zero-rules.md)（模块路径表）、[快速上手](05-development/quickstart.md)（多模块工作区说明）、[实现进度](06-status/implementation-status.md)（v1.3，§2.1 孤儿 go.mod 已解决）。

- **2026-08-06（复核）** 全仓复核：media / exam / message / data / search 五个服务 logic 全部落地并编译通过，整体由 77.7%（8/13）校正为 **100%（13/13）**；修正 §2.4（media/exam 自定义 Model 已补齐，空壳已解决）；同步 5 篇 business-rules 状态为「已落地实现」、索引进度表与实现进度文档至 v1.2。

- **2026-08-06** course / trade / learning 三个服务全部 logic 落地（76/76、73/73、20/20），并校正其 business-rules 为「已落地实现」；刷新本索引进度（8/13 完成、303/390、77.7%），修正「阅读须知」的未落地服务数为 5。
- **2026-08-05（三）** 补齐 11 个服务的 rpc-spec / data-model / business-rules / configs 共 46 篇；新建 message、remark 两个缺失的服务目录；新增 `05-development/quickstart.md` 与 `06-status/implementation-status.md`；重写本索引，将失真的「全 ✅」修正为真实实现进度。文档量 31 → 79 篇。
- **2026-08-05（二）** 从 OpenAPI 规范提取各服务 HTTP API 接口清单
- **2026-08-05（一）** 初始化骨架，搬运 MEMORY.md 与 .cursor/rules 内容
