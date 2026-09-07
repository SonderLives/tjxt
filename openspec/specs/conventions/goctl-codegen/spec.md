# goctl 代码生成规范 Specification

## Purpose

规范 tjxt 微服务在 go-zero v1.10.3 下的代码生成职责边界与 goctl 工作流，确保骨架代码一律由 goctl 产出、人工只编写业务逻辑（`internal/logic`）、Custom Model 扩展与业务扩展文件，从而在重生成骨架时不丢失业务实现、不破坏模块结构。本规格取代任何「手写 handler / types / routes / pb / server / config 骨架」的做法，是新增服务或改动契约时的唯一权威流程。

## Requirements

### Requirement: 生成优先与手写边界

代码产出 SHALL 遵循「生成优先、手写最小化」：仅允许手写三类文件，其余骨架一律由 goctl 生成。

| 类别 | 文件 | 说明 |
|------|------|------|
| 可手写 | `internal/logic/*.go` | 业务逻辑核心 |
| 可手写 | `internal/model/*model.go` | Custom Model 扩展方法，**不含 `*_model_gen.go`** |
| 可手写 | 业务扩展文件 | 如 `solution.go`、常量 `consts.go`、`vars.go` |
| 禁止手写 | `handler/`、`internal/types/`、`routes.go` | goctl api go 产出 |
| 禁止手写 | `pb/`、`internal/server/`、client 代码 | goctl rpc protoc 产出 |
| 禁止手写 | `internal/config/`、`internal/svc/` 骨架 | 生成后仅允许在其中注入依赖 |

#### Scenario: 新增接口只需重生成骨架

- **WHEN** 在 `.api` 或 `.proto` 中新增接口定义
- **THEN** 仅运行对应 goctl 命令重生成骨架，`internal/logic` 下的已有业务实现保持不动
- **AND** 新接口新增 logic 文件由人工补齐业务实现

#### Scenario: 禁止手改生成物

- **WHEN** 开发者需要调整请求/响应结构体或路由
- **THEN** 必须回到 `.api` / `.proto` 源文件修改并重生成
- **AND** 直接编辑 `types.go` / `routes.go` / `*.pb.go` 的行为视为违规，改动会在下次生成时被覆盖

### Requirement: 标准开发顺序

新增或改造一个服务 SHALL 严格按以下顺序执行：

```
1. 设计 DB 表结构
2. 写入 sql/ddl/tj_<domain>.sql（纯 DDL）
3. goctl model 生成 Model 层
4. 设计 apps/<svc>/rpc/<svc>.proto
5. goctl rpc 生成 RPC 骨架（pb/server/client/logic/svc/config）
6. 设计 apps/<svc>/api/<svc>.api
7. goctl api 生成 API 骨架（handler/logic/types/routes）
8. 手写 Logic 实现业务
9. 构建测试（go build / go vet / go test）
```

#### Scenario: DDL 先行

- **WHEN** 新增业务域能力
- **THEN** 先落 `sql/ddl/tj_<domain>.sql` 纯 DDL，再生成 Model
- **AND** 改表时同步更新 `sql/migration/tj_<domain>.sql`，避免 model 与实际库结构漂移

### Requirement: goctl API 生成

API 骨架 SHALL 由以下命令产出（工作目录为 `apps/<svc>/api`）：

```bash
cd apps/<svc>/api
goctl api go --api <svc>.api --dir . --style gozero
```

| 行为 | 说明 |
|------|------|
| 输出目录 | 与 `.api` 文件所在目录一致（`.`） |
| 不覆盖 | `handler/`、`logic/`、`internal/svc/`、`internal/config/` |
| 仅刷新 | `internal/types/`、`internal/handler/routes.go` |
| 需强制重生成 | 必须先手工删除目标文件，goctl 才会重新产出 |
| 类型限制 | `.api` 不支持 `any`，须用 `interface{}` |
| 请求体限制 | 拒绝顶层数组请求体（报 `request body must be struct`），须包一层 struct |

#### Scenario: 修改 API 契约后重生成

- **WHEN** 修改 `<svc>.api` 中的类型或路由
- **THEN** 执行 `goctl api go` 后 `internal/types/` 与 `routes.go` 被刷新，既有 handler/logic 保留
- **AND** 若需让 goctl 重新产出某个已存在的骨架文件，须先手工删除该文件

#### Scenario: `.api` 类型写法避坑

- **WHEN** 需要泛型/任意类型响应字段
- **THEN** 使用 `interface{}` 而非 `any`
- **AND** 需要数组请求体时包一层 struct（如 `courseList`），不得直接声明顶层数组

### Requirement: goctl RPC 生成

RPC 骨架 SHALL 由以下命令产出（工作目录为 `apps/<svc>/rpc`）：

```bash
cd apps/<svc>/rpc
goctl rpc protoc <svc>.proto \
  --go_out=. --go-grpc_out=. --zrpc_out=. \
  --style gozero --client=true
```

| 行为 | 说明 |
|------|------|
| `--client=true` | 必须固定带上，否则不生成供他方调用的 RPC client |
| `-m` 标志 | 可生成 logic 模板（首次使用） |
| 不覆盖 | `internal/logic/`、`internal/svc/`、`internal/config/` |
| 仅刷新 | `pb/`、`internal/server/`、client 代码 |

> 已知缺口：goctl 重跑 `rpc protoc` 会在扁平 `internal/logic/`（包 `logic`）重新生成 todo 桩，并把 `internal/server/<svc>/<svc>server.go` 改写为 import 扁平包。若真实实现位于子包（如 `internal/logic/search/`，包 `searchlogic`），须删除扁平桩并保留子包实现，同时确认 `<svc>server.go` 的 import 与真实 logic 所在包一致——否则服务会「编译通过但返回空结果」。

#### Scenario: 生成后校验接线正确性

- **WHEN** 重跑 `goctl rpc protoc` 之后
- **THEN** 必须检查 `grep -n "internal/logic" internal/server/*/<svc>server.go`，确认 server import 的包与真实 logic 文件所在包一致
- **AND** 若 `internal/logic/` 出现同名扁平桩文件，删除之，保留子包中的真实实现

#### Scenario: 缺少 RPC client

- **WHEN** 其他服务需要调用本服务
- **THEN** 本服务生成时必须带 `--client=true`，产出 `apps/<svc>/rpc/<svc>/<svc>.go`
- **AND** 调用方通过 `go.mod` 的 `require` + 本地相对路径 `replace` 引入，禁止 import 对方 `internal` 包

### Requirement: goctl Model 生成与扩展

Model 层 SHALL 由 DDL 生成，并只允许在 Custom Model 中扩展：

```bash
goctl model mysql ddl \
  -src sql/ddl/tj_<domain>.sql \
  -dir apps/<svc>/rpc/internal/model \
  -cache --style gozero
```

| 项 | 约定 |
|----|------|
| `-cache` | 启用 Redis 缓存 |
| 产出 | `<table>model.go`（可编辑）、`<table>_model_gen.go`（**禁止手改**）、`vars.go` |
| Model 已存在 | **不得重新生成覆盖**，只在 `<table>model.go` 中新增扩展方法 |
| 生成 model 的局限 | 只有 CRUD；业务所需的 List/Count/Page/FindByX/Upsert/级联 Delete 需手写扩展 |
| CachedConn 避坑 | 自定义 model 嵌入 `CachedConn`，只能用 `QueryRowsNoCacheCtx` / `QueryRowNoCacheCtx` / `ExecNoCacheCtx`；`QueryRowsCtx` 被遮蔽会编译失败 |

扩展模式示例：

```go
// apps/auth/rpc/internal/model/usermodel.go（手写，可编辑）
type UserModel interface {
    UserModelGen // 嵌入生成的接口
    FindByPhone(ctx context.Context, phone string) (*User, error)
}

func (m *defaultUserModel) FindByPhone(ctx context.Context, phone string) (*User, error) {
    var resp User
    query := "select * from `user` where `phone` = ? limit 1"
    err := m.conn.QueryRowCtx(ctx, &resp, query, phone)
    return &resp, err
}
```

#### Scenario: 已有 Model 增加查询方法

- **WHEN** 业务需要生成物未提供的查询（如按手机号查用户）
- **THEN** 在 `<table>model.go` 的接口中追加方法并在 `default<Table>Model` 上实现
- **AND** 不运行 `goctl model` 重新生成，避免覆盖已有扩展

#### Scenario: 自定义 model 的缓存调用

- **WHEN** 在 Custom Model 中编写带缓存连接的查询
- **THEN** 只使用 `QueryRowsNoCacheCtx` / `QueryRowNoCacheCtx` / `ExecNoCacheCtx`
- **AND** 使用被遮蔽的 `QueryRowsCtx` 导致的编译错误必须在提交前修复

### Requirement: 命名与模块结构约定

生成的目录、文件与模块路径 SHALL 遵守下表：

| 项目 | 规范 | 示例 |
|------|------|------|
| 模块路径 | `module tjxt/apps/<svc>/api`、`module tjxt/apps/<svc>/rpc` | `module tjxt/apps/auth/api` |
| 依赖组织 | api 与 rpc 各自独立 `go.mod` | `apps/<svc>/api/go.mod`、`apps/<svc>/rpc/go.mod` |
| 公共库引用 | 各模块 `replace tjxt/pkg => ../../../pkg`（data 服务四级 `../../../../pkg`） | — |
| 跨服务依赖 | 调用方 `go.mod` 中 `require` + 本地相对路径 `replace`，伪版本 `v0.0.0-00010101000000-000000000000` | `tjxt/apps/course/rpc => ../../course/rpc` |
| 目录命名 | gozero 风格（小驼峰，无下划线） | `internal/logic`、`userCoupon` |
| 文件命名 | 小驼峰 | `userCouponLogic.go`、`couponModel.go` |
| 包名 | 与目录名一致 | `package logic` |
| Proto 服务名 | `<Svc>Service` | `service AuthService` |
| API 文件 | `<svc>.api` | `auth.api` |

#### Scenario: 新增服务的模块落位

- **WHEN** 新建服务 `<svc>`
- **THEN** 建立 `apps/<svc>/api` 与 `apps/<svc>/rpc` 两个 go.mod，并在 `go.work` 中 `use`
- **AND** 公共库与跨服务依赖一律通过 `replace` 指向本地相对路径，不发布到远端

### Requirement: 常见错误与规避

生成与手写过程中 SHALL 规避下列已记录错误：

| 错误现象 | 原因 | 规避 |
|----------|------|------|
| 手改 `*_gen.go` 被覆盖 | 重新生成 model | 只在 `*_model.go` 扩展 |
| `.api` 里用 `any` 类型报错 | goctl 不支持 `any` | 用 `interface{}` |
| handler 里直接返回 JSON | 违反统一响应规范 | 用 `result.Write(w, r, data, err)` |
| 跨服务 import internal 包 | 破坏模块边界 | 只能调用对方 RPC client |
| RPC client 未生成 | 忘加 `--client=true` | 命令固定带上 |
| 配置文件不生效 | yaml 字段名与 struct tag 不匹配 | 对照 `config.go` 检查 |

#### Scenario: 配置不生效排查

- **WHEN** yaml 中新增配置项但服务未读取到
- **THEN** 对照 `internal/config/config.go` 的 struct tag 核对字段名与层级
- **AND** 确认该 yaml 位于服务实际加载的 `etc/` 路径（data 服务为 `apps/data/api/data/etc/`、`apps/data/rpc/data/etc/`）

### Requirement: 可观测性配置保护

每个服务 `apps/<svc>/{api,rpc}/etc/*.yaml` 中已注入的三段配置 SHALL 被视为不可删除、不可改坏：

```yaml
Telemetry:                 # 链路追踪：go-zero 自动 StartAgent，推 OTLP 到 127.0.0.1:4318
  Name: <svc>
  Endpoint: 127.0.0.1:4318
  Sampler: 1.0
  Batcher: otlphttp

Prometheus:               # 指标：各自暴露 /metrics，由 otel-collector 抓取后聚合
  Host: 0.0.0.0
  Port: <唯一>            # RPC 9101-9113 / API 9201-9213
  Path: /metrics

Log:                      # 日志：写 JSON 文件，由 collector filelog 采集进 Loki
  Mode: file
  Encoding: json
  Path: logs/<svc>        # 相对仓库根；服务须从仓库根启动
  Level: info
```

#### Scenario: 重生成骨架后配置仍在

- **WHEN** 执行 `goctl api go` 或 `goctl rpc protoc` 重生成骨架
- **THEN** 上述三段配置保持不变（goctl 只刷新 types/routes/pb/server）
- **AND** 若发现配置丢失，须按上表原样恢复后再启动服务

#### Scenario: 日志路径约束

- **WHEN** 服务启动时加载 `Log.Path: logs/<svc>`
- **THEN** 该路径相对进程工作目录解析，服务必须从仓库根启动
- **AND** 从服务子目录启动会导致日志散落、collector 的 filelog 采集不到
