# 本地开发工作流规范 Specification

## Purpose

规范 tjxt 从拉起基础设施、初始化 go.work 多模块工作区、生成代码、启动 13 个服务到校验构建的完整本地开发流程，并固化已验证的 goctl 行为、数据库目录约定与已知开发陷阱，使任意开发者在 Windows + Docker Desktop 环境下可复现地拉起整套系统。

## Requirements

### Requirement: 环境依赖

本地开发环境 SHALL 具备以下组件：

| 组件 | 版本要求 | 说明 |
|------|---------|------|
| Go | 1.26.2+ | `go.work` 声明版本 |
| goctl | 1.10.x | go-zero 代码生成器，**api/保留 rpc 骨架代码必须由它生成** |
| buf | 1.75.0 | proto lint/breaking 门禁 + 折叠服务 pb 类型生成（`make pb`） |
| protoc-gen-go | 1.36.x | buf generate 的 pb 插件（折叠服务）；保留 rpc 另需 protoc / protoc-gen-go-grpc |
| Docker Desktop | — | 拉起 MySQL/Redis/RabbitMQ/etcd 及可观测性栈 |
| PowerShell | 5.1+ | Makefile 内部大量使用 `powershell -NoProfile`，**Windows 环境专用** |

```bash
go install github.com/zeromicro/go-zero/tools/goctl@latest
```

#### Scenario: 新机器初始化

- **WHEN** 开发者首次 clone 仓库
- **THEN** 依次安装 Go 1.26.2+、goctl 1.10.x、protoc 三件套、Docker Desktop
- **AND** Windows 环境需 PowerShell 5.1+ 以执行 Makefile 目标

### Requirement: 基础设施启动

本地基础设施 SHALL 通过 `make docker-up`（等价 `docker compose up -d`）拉起 8 个容器：

| 服务 | 镜像 | 端口 | 凭据 |
|------|------|------|------|
| MySQL | `mysql:8.0` | 3306 | root / `0000` |
| Redis | `redis:7-alpine` | 6379 | 无密码 |
| RabbitMQ | `rabbitmq:3.13-management` | 5672 / 15672 | rabbitmq / rabbitmq |
| etcd | `bitnami/etcd:3.5` | 2379 | 免认证 |
| Jaeger | `jaegertracing/all-in-one:1.57` | 16686 (UI) | 免认证 |
| otel-collector | `otel/opentelemetry-collector-contrib:latest` | 4317/4318/8889/13133 | — |
| Prometheus | `prom/prometheus:v2.53.1` | 9090 (UI) | — |
| Loki | `grafana/loki:latest` | 3100 | — |

- MySQL 容器首次启动挂载 `./sql/thirdparty`（第三方基建库脚本）与 `./sql/bootstrap.sql`（业务库建库引导）到 `/docker-entrypoint-initdb.d`；业务库表结构与种子数据由 goose 管理，容器就绪后执行 `make migrate-up`
- 需要重新初始化时须先 `docker compose down -v` 删除 `mysql-data` 卷
- 可观测性信号入口：Trace http://localhost:16686、Metrics http://localhost:9090（`tjxt-services` target 应 UP）、Logs http://localhost:3100

#### Scenario: 首次拉起环境

- **WHEN** 执行 `make docker-up`
- **THEN** 8 个容器启动，MySQL 自动执行第三方基建脚本与业务库建库引导（第三方库建库建表灌数据；业务库仅建库）
- **AND** 执行 `make migrate-up` 完成业务库建表与种子数据
- **AND** 可用 `make docker-logs` 跟踪容器日志（含 collector 健康检查 13133）

#### Scenario: 重新初始化数据库

- **WHEN** 需要清空并重灌初始数据
- **THEN** 先 `docker compose down -v` 删除 `mysql-data` 卷，再 `make docker-up`
- **AND** 仅重启容器不会触发初始化脚本重新执行

### Requirement: 工作区初始化与服务启动

工作区 SHALL 先 `make init`（`go work sync` + 逐模块 `go mod tidy`），再按「先 RPC 后 API（及 MQ 消费者）」的顺序启动服务。

- `go.work` 聚合 22 个 `use` 条目（13 api + 5 保留 rpc + 2 mq 消费者 + `pkg` + 根模块）
- 每个服务 `app/<svc>/api/`、保留服务 `app/<svc>/rpc/` 与 `app/mq/*/` 各有独立 `go.mod`；公共库经 `replace tjxt/pkg => ../../../pkg` 引用
- 启动命令：`make run-<svc>`（API）、`make run-<svc>-rpc`（仅 5 个保留 RPC）、`make run-trade-mq` / `make run-search-mq`（MQ 消费者）、`make run-all`（全部 API）/ `make run-all-rpc`（5 个保留 RPC + 2 个 MQ 消费者）
- 端口已统一分配且唯一：API 8801-8813、保留 RPC 8081/8083/8084/8088/8092、metrics：API 9201-9213、保留 RPC 9101/9103/9104/9108/9112、search-mq 9110

> 源冲突说明：wiki `05-development/quickstart.md` v1.2 仍称「media-api 与 pay-api 同占 8808」。实测 `etc/*.yaml` 端口已重排为 `media-api 8806` / `pay-api 8808`，API 端口 8801-8813 全局唯一，无冲突，可同时启动。

> 折叠说明：auth/exam/media/message/remark/search/trade/data 的 rpc 进程壳已删除，无 `run-<svc>-rpc` 目标；trade/search 原消费 MQ 的职责由 `app/mq/trade-mq`、`app/mq/search-mq` 承接，用 `run-trade-mq` / `run-search-mq` 启动。

#### Scenario: 冷启动全套服务

- **WHEN** 开发者启动整套系统
- **THEN** 先 `make init`，再启动 5 个保留 RPC（注册到 etcd）与 2 个 MQ 消费者，然后启动全部 API（从 etcd 发现 RPC）
- **AND** 顺序颠倒会导致 API 侧服务发现失败

### Requirement: 代码生成与校验命令

日常开发 SHALL 使用以下 Makefile 目标；单服务生成命令见 `conventions/goctl-codegen`。

| 目标 | 作用 |
|------|------|
| `make init` | `go work sync` + 逐模块 `go mod tidy` |
| `make api` | 按 `.api` 全量重生成 handler/types/routes（13 个服务） |
| `make rpc` | 按 `.proto` 全量重生成 pb/server/client（仅 5 个保留 RPC 服务，goctl） |
| `make pb` | buf generate 重生成 8 个折叠服务的 pb 类型到 `app/<svc>/api/backend/pb/` |
| `make model` | 按 `sql/migrations/<db>/00001_schema.sql` 生成带缓存的 model（exam/trade/message 自动剥除 goctl 不解析的表） |
| `make migrate-up` / `make migrate-down` / `make migrate-status` | goose 迁移 12 个业务库：全量应用 / 回退（MIGRATE_STEPS 步，默认 1）/ 查看状态 |
| `make generate` | `api + rpc + pb` 一键生成 |
| `make check-gen` | 重生成 api + pb 后 `git diff --exit-code -- app`（api types 与 pb 类型双守卫） |
| `make build` | 构建全部 20 个二进制（13 api + 5 rpc + 2 mq）到 `bin/` |
| `make fmt` | 全模块 `go fmt ./...` |
| `make lint` | golangci-lint 逐模块检查 |
| `make test` | 全模块 `go test ./...` |
| `make verify` | 校验所有服务的 .api/.proto 与 internal/ 目录齐全 |
| `make tidy` | 全模块 `go mod tidy` |
| `make docker-up` / `make docker-logs` | 启动 / 跟踪基础设施容器 |

单模块快速校验（api / rpc 各有独立 go.mod，需分别进入）：

```bash
cd app/<svc>/api && go mod tidy && go build ./... && go vet ./...
cd app/<svc>/rpc && go mod tidy && go build ./... && go vet ./...
```

#### Scenario: 改动契约后的最小闭环

- **WHEN** 修改 `.api` / `.proto` 后需要回归
- **THEN** 依次执行 `make api` 或 `make rpc`、`make build`、`make verify`
- **AND** 单模块问题可进入 `app/<svc>/api` 或 `app/<svc>/rpc` 单独 `go build ./... && go vet ./...` 定位

### Requirement: 数据库目录约定

SQL 脚本 SHALL 由 pressly/goose 管理，一业务库一迁移目录，改表时新增迁移：

| 目录 | 用途 |
|------|------|
| `sql/migrations/<db>/00001_schema.sql` | `tj_<db>` 库建表 DDL（含索引，与真实库表集合一致），Up 建表 / Down 逆序 DROP |
| `sql/migrations/<db>/00002_seed.sql` | 演示基线数据（仅 9 个有种子数据的库存在），Up INSERT / Down 按表 DELETE |
| `sql/migrations/<db>/00003_….sql` | 后续表结构变更用 `goose create` 追加的新迁移 |
| `sql/thirdparty/` | gogs/nacos/seata/xxl_job 第三方基建脚本，不归 goose 管理 |
| `sql/bootstrap.sql` | 业务库建库引导（goose 不建库，迁移前执行） |

一库一服务，库名 `tj_<domain>`。

#### Scenario: 表结构变更

- **WHEN** 某域新增字段或索引
- **THEN** 用 `goose -dir sql/migrations/<db> mysql "<DSN>/tj_<db>" create <name> sql` 新增迁移文件（编号递增），Up 写变更、Down 写逆操作
- **AND** `make model` 从迁移 schema 生成 model，只改迁移不重新生成 model 会导致两者漂移

### Requirement: 已知开发陷阱

开发者 SHALL 规避以下已记录陷阱：

| 陷阱 | 规避方式 |
|------|---------|
| Makefile 服务清单曾有缺口 | `SERVICES` 曾只列 11 个（缺 `promotion`/`remark`），现已是全量 13 个；折叠后新增服务须同步维护 `SERVICES`/`API_DIRS`/`RPC_DIRS`/`FOLDED_SVCS`/`MQ_DIRS` |
| `data` 服务目录曾有特例 | 旧布局多套一层（服务目录下再套服务名目录），现已扁平化为标准布局 `app/data/api`（目录迁移前为 `api/data`，rpc/data 已折叠删除），与其他服务一致 |
| Windows 下 `rm` 被沙箱拦截 | 清空目录准备重生成时 `rm -rf` 可能静默失败，需确认删除结果 |
| 统一响应易写错 | handler 必须用 `result.Write(w, r, data, err)`，不要用 goctl 生成的 Result 类型 |
| JWT 中间件 | 在 `.api` 中用 `@server jwt: Auth` 声明，userId 通过 `auth.UserIdFromCtx(l.ctx)` 获取 |
| 日志目录相对路径 | `Log.Path: logs/<svc>` 相对进程 CWD，须从仓库根启动服务，否则 collector 采集不到 |
| 批量改 yaml 转义 | 用脚本文件写 Python 再执行，勿用 git-bash heredoc（`\n` 会被字面化导致 YAML 解析失败） |

#### Scenario: 日志采集不到

- **WHEN** Jaeger 有 trace 但 Loki 查不到日志
- **THEN** 先确认服务是否从仓库根启动（相对路径 `logs/<svc>` 是否落在 `<repo>/logs/` 下）
- **AND** 确认 yaml 中 `Log.Encoding` 为 `json`（`plain` 会导致 filelog 无法解析结构化字段）
