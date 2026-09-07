# OpenSpec 规格文档

本目录存放按 [OpenSpec](https://openspec.dev) 规范整理的项目能力规格（`specs/`），由项目原 repo wiki（80 篇）拆分转换而来。**源 wiki 目录已迁移完成并删除，本目录是项目规格的唯一事实源。** 转换基准日：2026-09-08（源 wiki 最后更新 2026-08-15）。

## 目录结构

```
openspec/
├── config.yaml            # schema 配置 + 项目上下文（架构/规范类 wiki 浓缩）
├── specs/                 # 主规格（当前系统「是什么」的契约）
│   ├── services/<svc>/    # 13 个微服务能力规格（每服务一份）
│   ├── architecture/       # 跨服务架构：service-topology（依赖拓扑/调用链/契约原则/部署视图）
│   ├── conventions/       # 工程约定：goctl-codegen / api-contracts / code-style / git-workflow
│   ├── shared/            # 跨服务共享契约
│   │   ├── error-codes/   # pkg/xerr 统一错误码
│   │   ├── mq-events/     # RabbitMQ 事件总线契约
│   │   └── pkg-contracts/ # pkg/ 公共库契约（auth/response/page/idgen/mq）
│   └── infra/             # 基础设施规格
│       ├── observability/ # otel-collector 三信号收口（trace/metrics/logs）
│       ├── local-env/     # docker-compose 本地环境（8 容器 + 端口 + 健康检查）
│       └── local-dev-workflow/ # 本地开发工作流（环境/启动/生成/校验/陷阱）
└── changes/               # 变更提案（spec-driven 工作流：proposal/specs/design/tasks）
```

## 迁移溯源：原 wiki → specs 映射

| 原 wiki 来源 | 去向 |
|-----------|------|
| `00-architecture/overview.md`、`service-topology.md`、`data-domains.md` | `config.yaml` 的 `context`（技术栈、端口表、依赖原则、库分布）+ `specs/architecture/service-topology/spec.md`（依赖图、三条调用链、契约原则、跨服务数据访问规则、SQL 版本管理、生产部署拓扑）；数据表细节并入各服务 spec 的数据模型 Requirement |
| `01-conventions/go-zero-rules.md` | `specs/conventions/goctl-codegen/spec.md` |
| `01-conventions/api-contracts.md` | `specs/conventions/api-contracts/spec.md` |
| `01-conventions/code-style.md` | `specs/conventions/code-style/spec.md` |
| `01-conventions/git-workflow.md` | `specs/conventions/git-workflow/spec.md`（同时仍在 `config.yaml` 浓缩引用） |
| `02-services/<svc>/`（api-spec / rpc-spec / data-model / business-rules / configs 五篇） | `specs/services/<svc>/spec.md` 一份，五篇内容合并为按功能域分组的 Requirement + Scenario |
| `03-shared/error-codes.md` | `specs/shared/error-codes/spec.md` |
| `03-shared/mq-events.md` | `specs/shared/mq-events/spec.md` |
| `03-shared/pkg-contracts.md` | `specs/shared/pkg-contracts/spec.md` |
| `04-infra/observability.md` | `specs/infra/observability/spec.md` |
| `04-infra/docker-compose.md` | `specs/infra/local-env/spec.md` |
| `05-development/quickstart.md` | `specs/infra/local-dev-workflow/spec.md` |
| `06-status/implementation-status.md` | 已随源 wiki 一并删除（进度快照，非系统契约）；其中 20+ 已知缺口已分散保留在各服务 spec 的 `> 已知缺口：` 引用块中 |

## 服务规格清单

| 服务 | 规格 | 服务 | 规格 |
|------|------|------|------|
| user 用户中心 | `services/user` | pay 支付网关 | `services/pay` |
| auth 认证授权 | `services/auth` | trade 交易订单 | `services/trade` |
| course 课程管理 | `services/course` | promotion 营销优惠 | `services/promotion` |
| learning 学习进度 | `services/learning` | remark 评论点赞 | `services/remark` |
| exam 考试题库 | `services/exam` | message 站内信 | `services/message` |
| media 媒资文件 | `services/media` | search 全文检索 | `services/search` |
| data 统计大屏 | `services/data` | | |

### 工程约定与基础设施规格

| 规格 | 内容 |
|------|------|
| `architecture/service-topology` | 服务依赖图（设计意图 vs 实际接线对照）、下单支付/学习进度/文件上传三条调用链、5 条服务间契约原则、跨服务数据访问正反例、数据域边界与 SQL 版本管理、K8s 生产部署拓扑 |
| `conventions/goctl-codegen` | 生成优先与手写边界、标准开发顺序、goctl api/rpc/model 命令与覆盖行为、命名与模块结构、Model 扩展模式、常见错误规避、可观测性配置保护 |
| `conventions/api-contracts` | 统一响应 `response.R` + `result.Write`、分页契约、JWT 与身份获取、错误码分段、`.api` 编写规范、参数绑定、版本兼容 |
| `conventions/code-style` | 命名/格式化、xerr 错误处理、ctx 透传与 3s 超时、ServiceContext 依赖注入、标准包布局、import 分组与禁止项、测试约定 |
| `conventions/git-workflow` | 分支策略、Conventional Commits、PR 模板与门禁、SemVer、发布流程、CHANGELOG 与分支保护 |
| `infra/local-dev-workflow` | 环境依赖、8 容器启动、go.work 初始化、先 RPC 后 API 启动顺序、Makefile 目标、SQL 目录约定、已知开发陷阱 |
| `infra/local-env` | docker-compose 容器清单、端口、凭据、健康检查 |
| `infra/observability` | trace/metrics/logs 三信号经 otel-collector 收口 |

## 转换约定

- **格式**：每份 spec 含 `## Purpose` 与 `## Requirements`；Requirement 用 `### Requirement:`，每个至少一个 `#### Scenario:`（WHEN/THEN/AND 加粗列表）；需求陈述保留 SHALL。全部通过 `openspec validate --specs --strict`（24/24）。
- **忠实转写**：端口、默认值、错误码、状态枚举原样保留；wiki 中标注的「已知缺口 / mock / 未接线 / 占位」以引用块（`> 已知缺口：…`）保留在对应 Requirement 处。
- **源冲突处理**：同一事实在多篇 wiki 中描述不一致时，以 `data-model.md`（源自 DDL，最可信）为准并加注；与代码实现不一致时以实测为准并加注。本次已修正两处：
  - `quickstart.md` v1.2 称「media-api 与 pay-api 同占 8808」——实测 26 份 `etc/*.yaml` 端口为 API 8801-8813 全局唯一（media 8806 / pay 8808），无冲突。
  - `quickstart.md` 记录的「Makefile `SERVICES` 缺 `promotion`/`remark`」经核对**仍然成立**（当前仅 11 个），已在 `infra/local-dev-workflow` 中作为已知缺口保留。
  - `service-topology.md` 的依赖图中 `auth-api` 标 8801、`auth-rpc` 标 8081 → 实测为 8802 / 8082（8801/8081 属 user），已在 `architecture/service-topology` 更正。
- **与代码冲突时以代码为准**：`promotion` 的 `GET /codes/page` 在 wiki api-spec 与 OpenAPI 中缺失，但 `routes.go` 与 `couponcodepagelogic.go` 中真实存在，已按代码补齐到 `services/promotion`。
- **api-spec 响应体 URL 编码**已解码（`%C2%AB`→`«`）。
- **mq-events 与实际实现的差异**：`shared/mq-events` 记录的是 wiki 规划口径（`trade.events` / `order.created` 等）；learning spec 中已加注实际代码使用 `order.exchange` + `order.pay`/`order.refund` 的命名差异。

## 常用命令

```bash
openspec list --type spec        # 列出全部规格
openspec show services/auth      # 查看某规格
openspec validate --specs --strict  # 全量校验
```

## 相关文档

- 完成度分析报告（代码完成度 / 功能完备度统计）：`docs/completion-analysis.md`
- OpenAPI 契约原文：`docs/tjxt.openapi.json`

后续新需求走 `openspec new change "<name>"`（spec-driven 工作流），在 `changes/` 下产出 proposal / delta specs / design / tasks，实现完成归档后同步回 `specs/`。
