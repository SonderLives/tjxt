# 本地开发工作流规范 Specification

## Purpose

规范 tjxt 从拉起基础设施、初始化 go.work 多模块工作区、生成代码、启动 13 个服务到校验构建的完整本地开发流程，并固化已验证的 goctl 行为、数据库目录约定与已知开发陷阱，使任意开发者在 Windows + Docker Desktop 环境下可复现地拉起整套系统。

## Requirements

### Requirement: 环境依赖

本地开发环境 SHALL 具备以下组件：

| 组件 | 版本要求 | 说明 |
|------|---------|------|
| Go | 1.26.2+ | `go.work` 声明版本 |
| goctl | 1.10.x | go-zero 代码生成器，**所有骨架代码必须由它生成** |
| protoc / protoc-gen-go / protoc-gen-go-grpc | 最新 | RPC 层 pb 生成 |
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

- MySQL 容器首次启动挂载 `./sql/migration` 到 `/docker-entrypoint-initdb.d`，自动建库建表并灌入初始数据
- 需要重新初始化时须先 `docker compose down -v` 删除 `mysql-data` 卷
- 可观测性信号入口：Trace http://localhost:16686、Metrics http://localhost:9090（`tjxt-services` target 应 UP）、Logs http://localhost:3100

#### Scenario: 首次拉起环境

- **WHEN** 执行 `make docker-up`
- **THEN** 8 个容器启动，MySQL 自动执行 `sql/migration` 下的脚本建库建表并灌入初始数据
- **AND** 可用 `make docker-logs` 跟踪容器日志（含 collector 健康检查 13133）

#### Scenario: 重新初始化数据库

- **WHEN** 需要清空并重灌初始数据
- **THEN** 先 `docker compose down -v` 删除 `mysql-data` 卷，再 `make docker-up`
- **AND** 仅重启容器不会触发初始化脚本重新执行

### Requirement: 工作区初始化与服务启动

工作区 SHALL 先 `make init`（`go work sync` + 逐模块 `go mod tidy`），再按「先 RPC 后 API」的顺序启动服务。

- `go.work` 聚合 28 个 `use` 条目（13 服务 × api/rpc 共 26 个模块 + `pkg` + 根模块）
- 每个服务 `api/` 与 `rpc/` 各有独立 `go.mod`；公共库经 `replace tjxt/pkg => ../../../pkg` 引用
- 启动命令：`make run-<svc>`（API）、`make run-<svc>-rpc`（RPC）、`make run-all` / `make run-all-rpc`（全量）
- 端口已统一分配且唯一：API 8801-8813、RPC 8081-8093、metrics RPC 9101-9113 / API 9201-9213

> 源冲突说明：wiki `05-development/quickstart.md` v1.2 仍称「media-api 与 pay-api 同占 8808」。实测 26 份 `etc/*.yaml` 端口已重排为 `media-api 8806` / `pay-api 8808`，API 端口 8801-8813 全局唯一，无冲突，可同时启动。

> 已知缺口：`Makefile` 的 `SERVICES` 变量仅列出 11 个服务，**缺 `promotion` 与 `remark`**，这两个服务无法用 `make run-*` 启动，需手工 `cd apps/<svc>/api && go run . -f etc/<svc>-api.yaml`。

#### Scenario: 冷启动全套服务

- **WHEN** 开发者启动整套系统
- **THEN** 先 `make init`，再启动全部 RPC（注册到 etcd），然后启动全部 API（从 etcd 发现 RPC）
- **AND** 顺序颠倒会导致 API 侧服务发现失败

#### Scenario: 启动 promotion 或 remark

- **WHEN** 需要运行 promotion / remark 服务
- **THEN** 因 Makefile `SERVICES` 未包含，须手工进入服务目录 `go run`
- **AND** 从仓库根启动，保证 `Log.Path: logs/<svc>` 相对路径解析正确

### Requirement: 代码生成与校验命令

日常开发 SHALL 使用以下 Makefile 目标；单服务生成命令见 `conventions/goctl-codegen`。

| 目标 | 作用 |
|------|------|
| `make init` | `go work sync` + 逐模块 `go mod tidy` |
| `make api` | 按 `.api` 全量重生成 handler/types/routes |
| `make rpc` | 按 `.proto` 全量重生成 pb/server/client |
| `make model` | 按 `sql/ddl/*.sql` 生成带缓存的 model |
| `make generate` | `api + rpc` 一键生成 |
| `make build` | 构建全部 26 个二进制到 `bin/` |
| `make fmt` | 全模块 `go fmt ./...` |
| `make lint` | 全模块 `go vet ./...` |
| `make test` | 全模块 `go test ./...` |
| `make verify` | fmt + lint + test 一站式（另有 `make verify` 校验目录齐全的语义） |
| `make tidy` | 全模块 `go mod tidy` |
| `make docker-up` / `make docker-logs` | 启动 / 跟踪基础设施容器 |

单模块快速校验（api / rpc 各有独立 go.mod，需分别进入）：

```bash
cd apps/<svc>/api && go mod tidy && go build ./... && go vet ./...
cd apps/<svc>/rpc && go mod tidy && go build ./... && go vet ./...
```

#### Scenario: 改动契约后的最小闭环

- **WHEN** 修改 `.api` / `.proto` 后需要回归
- **THEN** 依次执行 `make api` 或 `make rpc`、`make build`、`make verify`
- **AND** 单模块问题可进入 `apps/<svc>/{api,rpc}` 单独 `go build ./... && go vet ./...` 定位

### Requirement: 数据库目录约定

SQL 脚本 SHALL 按用途分放两个目录，改表时同步更新：

| 目录 | 用途 |
|------|------|
| `sql/ddl/tj_<domain>.sql` | **纯 DDL**，供 `goctl model` 解析，不含数据 |
| `sql/migration/tj_<domain>.sql` | 完整迁移脚本，含索引/外键/初始数据，供 Docker 初始化 |

一库一服务，库名 `tj_<domain>`。

#### Scenario: 表结构变更

- **WHEN** 某域新增字段或索引
- **THEN** 同时更新 `sql/ddl/` 与 `sql/migration/` 下对应脚本
- **AND** 只改一处会导致 goctl 生成的 model 与实际库结构漂移

### Requirement: 已知开发陷阱

开发者 SHALL 规避以下已记录陷阱：

| 陷阱 | 规避方式 |
|------|---------|
| Makefile 服务清单不全 | `SERVICES` 只列 11 个，缺 `promotion`/`remark`，需手工启动 |
| `data` 服务目录多套一层 | 路径为 `apps/data/api/data/`、`apps/data/rpc/data/`，模块数与其余服务一致，仅路径更深 |
| Windows 下 `rm` 被沙箱拦截 | 清空目录准备重生成时 `rm -rf` 可能静默失败，需确认删除结果 |
| 统一响应易写错 | handler 必须用 `result.Write(w, r, data, err)`，不要用 goctl 生成的 Result 类型 |
| JWT 中间件 | 在 `.api` 中用 `@server jwt: Auth` 声明，userId 通过 `auth.UserIdFromCtx(l.ctx)` 获取 |
| 日志目录相对路径 | `Log.Path: logs/<svc>` 相对进程 CWD，须从仓库根启动服务，否则 collector 采集不到 |
| 批量改 yaml 转义 | 用脚本文件写 Python 再执行，勿用 git-bash heredoc（`\n` 会被字面化导致 YAML 解析失败） |

#### Scenario: 日志采集不到

- **WHEN** Jaeger 有 trace 但 Loki 查不到日志
- **THEN** 先确认服务是否从仓库根启动（相对路径 `logs/<svc>` 是否落在 `<repo>/logs/` 下）
- **AND** 确认 yaml 中 `Log.Encoding` 为 `json`（`plain` 会导致 filelog 无法解析结构化字段）
