> 版本：v1.1 | 更新：2026-08-06 | 来源：`Makefile`, `docker-compose.yml`, `go.work`（本次模块拆分重构）

---

# 本地开发快速上手

## 1. 环境依赖

| 组件 | 版本要求 | 说明 |
|------|---------|------|
| Go | 1.26.2+ | `go.work` 声明版本 |
| goctl | 1.10.x | go-zero 代码生成器，**所有骨架代码必须由它生成** |
| protoc / protoc-gen-go / protoc-gen-go-grpc | 最新 | RPC 层 pb 生成 |
| Docker Desktop | — | 拉起 MySQL/Redis/RabbitMQ/etcd |
| PowerShell | 5.1+ | Makefile 内部大量使用 `powershell -NoProfile`，**Windows 环境专用** |

```bash
go install github.com/zeromicro/go-zero/tools/goctl@latest
```

---

## 2. 拉起基础设施

```bash
make docker-up      # 等价于 docker compose up -d
```

`docker-compose.yml` 会启动四个容器：

| 服务 | 镜像 | 端口 | 凭据 |
|------|------|------|------|
| MySQL | `mysql:8.0` | 3306 | root / `0000` |
| Redis | `redis:7-alpine` | 6379 | 无密码 |
| RabbitMQ | `rabbitmq:3.13-management` | 5672 / 15672 | rabbitmq / rabbitmq |
| etcd | `bitnami/etcd:3.5` | 2379 | 免认证 |

> MySQL 容器首次启动时会自动挂载 `./sql/migration` 到 `/docker-entrypoint-initdb.d`，**自动建库建表并灌入初始数据**。若需要重新初始化，须先 `docker compose down -v` 删除 `mysql-data` 卷。

⚠️ RabbitMQ 虽已在 compose 中启动，但当前代码的事件总线走的是 **Redis Stream**（见 `pkg/mq/`），RabbitMQ 暂未接线。

---

## 3. 初始化工作区

```bash
make init      # go work sync + 逐模块 go mod tidy
```

项目为 **多模块工作区（go-zero 官方标准结构）**：`go.work` 聚合了 28 个 `use` 条目（13 个服务的 `api`+`rpc` 共 26 个模块 + `pkg` + 根模块）。每个服务的 `api/` 与 `rpc/` 各有独立 `go.mod`（`module tjxt/apps/<svc>/api` 与 `module tjxt/apps/<svc>/rpc`），通过各自 `replace tjxt/pkg => ../../../pkg` 引用公共库；跨服务依赖（如 `trade` 调 `course`/`pay`）在调用方 `go.mod` 中 `require` + `replace` 到本地相对路径。

---

## 4. 启动服务

```bash
make run-user        # 启动 user-api
make run-user-rpc    # 启动 user.rpc
make run-all         # 启动全部服务
```

**启动顺序**：先起 RPC 层（注册到 etcd），再起 API 层（从 etcd 发现 RPC）。

⚠️ **端口冲突**：骨架阶段端口未统一分配，同时启动全部服务会失败。已知冲突见 [服务拓扑](../00-architecture/service-topology.md)：

- API：`media-api` 与 `pay-api` 同占 `8808`
- RPC 端口已唯一分配（auth 8082、course 8083、exam 8085、learning 8084、message 8087、trade 8089、data 8091、promotion 8092 等），无冲突，可同时启动。

同时启动前请先手工改对应 `etc/*.yaml`。

---

## 5. 代码生成（核心工作流）

> 🔒 **项目铁律**：只手写 `logic/`、自定义 model、业务扩展。**禁止手写** handler / types / routes / svc / config / pb 等骨架代码。详见 [go-zero 开发约束](../01-conventions/go-zero-rules.md)。

### 标准开发顺序

```
DDL 建表  →  Model(goctl)  →  Proto  →  RPC 生成  →  .api  →  API 生成  →  手写 Logic
```

### 常用命令

| 目标 | 命令 | 说明 |
|------|------|------|
| 全量 API | `make api` | 按 `.api` 重生成 handler/types/routes |
| 全量 RPC | `make rpc` | 按 `.proto` 重生成 pb/server/client |
| 全量 Model | `make model` | 按 `sql/ddl/*.sql` 生成带缓存的 model |

### 单服务手工生成

```bash
# API
cd apps/<svc>/api && goctl api go --api <svc>.api --dir . --style gozero

# RPC
cd apps/<svc>/rpc && goctl rpc protoc <svc>.proto \
  --go_out=. --go-grpc_out=. --zrpc_out=. --style gozero

# Model（带 Redis 缓存）
goctl model mysql ddl -src sql/ddl/tj_<svc>.sql \
  -dir apps/<svc>/rpc/internal/model -cache --style gozero
```

### 生成器行为要点（本机已验证）

| 行为 | 说明 |
|------|------|
| 不覆盖已有业务代码 | `goctl api go` / `goctl rpc protoc` **不会覆盖**已存在的 handler / logic / svc / config |
| 只刷新骨架 | 仅重写 `types.go` / `routes.go` / `pb` / `server` |
| 需重生成时 | 必须**先手工删除**目标文件，goctl 才会重新产出 |
| `.api` 不支持 `any` | 用 `interface{}` 代替 |
| 拒绝顶层数组请求体 | 报 `request body must be struct`，需包一层 struct（如 `courseList`） |
| Model 已存在时 | **不要重新生成覆盖**，改为扩展自定义 `<table>model.go` |

---

## 6. 校验与构建

```bash
make fmt        # go fmt 全模块
make lint       # go vet 全模块
make test       # 单元测试
make build      # 编译全部服务二进制
make verify     # fmt + lint + test 一站式
make tidy       # 逐模块 go mod tidy
```

单模块快速校验（api / rpc 各有独立 go.mod，需分别进入）：

```bash
cd apps/<svc>/api && go mod tidy && go build ./... && go vet ./...
cd apps/<svc>/rpc && go mod tidy && go build ./... && go vet ./...
```

---

## 7. 数据库约定

| 目录 | 用途 |
|------|------|
| `sql/ddl/tj_<domain>.sql` | **纯 DDL**，供 `goctl model` 解析，不含数据 |
| `sql/migration/tj_<domain>.sql` | 完整迁移脚本，含索引/外键/初始数据，供 Docker 初始化 |

一库一服务，库名 `tj_<domain>`。改表时**两个目录都要同步更新**，否则 model 生成与实际库结构会漂移。

---

## 8. 已知开发陷阱

| 陷阱 | 规避方式 |
|------|---------|
| Makefile 服务清单不全 | `SERVICES` 变量只列了 11 个，**漏了 `promotion` 和 `remark`**，这两个服务无法用 `make run-*` 启动，需手工 `cd apps/<svc>/api && go run .` |
| `data` 服务目录多套一层 | 路径为 `apps/data/api/data/`、`apps/data/rpc/data/`，与其余 12 个服务（`apps/<svc>/{api,rpc}`）相比多嵌套一级；模块数与其余服务一致（各含 api/rpc 两个 go.mod），仅路径更深，已在 `go.work` 中正确 `use` |
| Windows 下 `rm` 被沙箱拦截 | 清空目录准备重新生成时，`rm -rf` 可能静默失败导致删除不完整，需确认删除结果 |
| 统一响应易写错 | handler 必须用 `result.Write(w, r, data, err)`，**不要用 goctl 生成的 Result 类型** |
| JWT 中间件 | 在 `.api` 中用 `@server jwt: Auth` 声明，userId 通过 `auth.UserIdFromCtx(l.ctx)` 获取 |

---

## 相关文档

- [架构全览](../00-architecture/overview.md)
- [go-zero 开发约束](../01-conventions/go-zero-rules.md)
- [API 契约规范](../01-conventions/api-contracts.md)
- [Docker 基础设施](../04-infra/docker-compose.md)
- [实现进度与已知缺口](../06-status/implementation-status.md)
