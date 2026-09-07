# 代码风格与包布局规范 Specification

## Purpose

统一 tjxt 的 Go 代码风格、错误处理、依赖注入方式与包布局，使 26 个可独立构建模块（13 服务 × api/rpc）在命名、import 分组、上下文传递和目录结构上保持一致，避免跨服务 import `internal` 包、裸返回 error、循环依赖等破坏模块边界的写法。本规格与 `conventions/goctl-codegen` 互补：前者规定「骨架由谁生成」，本规格规定「手写代码长什么样」。

## Requirements

### Requirement: Go 基础风格

Go 源码 SHALL 遵守以下基础规范，并在提交前执行 `make fmt`。

| 项 | 规范 |
|----|------|
| 格式化 | gofmt / goimports，提交前必须 `make fmt` |
| 行长度 | 建议 ≤ 120 字符 |
| 包名 | 小写单词、无下划线（`usercoupon`、`svc`） |
| 接口名 | 单词 + `er` 后缀（`UserStore`、`CouponIssuer`） |
| 实现名 | 接口名去 `er`（`userStore`、`couponIssuer`） |
| 私有标识 | 小写开头（`userModel`、`calcDiscount`） |
| 常量 | 全大写下划线（`CouponTypeDiscount`、`MaxRetryTimes`） |

#### Scenario: 提交前格式化

- **WHEN** 开发者准备提交代码
- **THEN** 先执行 `make fmt`（全模块 `go fmt ./...`）
- **AND** PR 检查中格式不符的改动不予合入

### Requirement: 错误处理

错误 SHALL 统一用 `pkg/xerr` 包装后返回，禁止裸返回 error 或用 `fmt.Errorf` 包装标准库错误。

```go
// ✅ 统一用 xerr 包装
if err != nil {
    return nil, xerr.NewErrCode(xerr.DB_ERROR).WithError(err)
}

// ✅ 业务校验失败返回具体码
if user.Status != 1 {
    return nil, xerr.NewErrCode(xerr.USER_DISABLED)
}

// ❌ 禁止裸返回 error
return nil, err

// ❌ 禁止 fmt.Errorf 包装（丢失堆栈）
return nil, fmt.Errorf("db error: %w", err)
```

#### Scenario: 基础设施错误

- **WHEN** DB / Redis / RPC 调用返回 error
- **THEN** 用 `xerr.NewErrCode(xerr.<对应码>).WithError(err)` 包装后向上返回
- **AND** 原始 error 被保留在包装结果中以便排查，但不直接暴露给客户端

#### Scenario: 业务规则不满足

- **WHEN** 业务校验失败（如用户被禁用）
- **THEN** 返回具体业务错误码（如 `xerr.USER_DISABLED`），不附带底层 error

### Requirement: Context 传递与超时

`ctx context.Context` SHALL 作为函数第一个入参并跨层透传。

- 不丢弃上游 ctx，不新建空 context
- 超时控制在调用 RPC client 处用 `context.WithTimeout` 施加（跨服务调用统一 3s）

#### Scenario: RPC 调用超时

- **WHEN** 服务 A 调用服务 B 的 RPC
- **THEN** 在调用处以 `context.WithTimeout(ctx, 3*time.Second)` 派生 ctx
- **AND** 派生 ctx 传入 RPC client，不得传入 `context.Background()`

### Requirement: 依赖注入

外部依赖 SHALL 由 `svc.ServiceContext` 统一持有，logic 通过 `NewXxxLogic(ctx, svcCtx)` 构造。

- `ServiceContext` 持有 DB、Redis、RPC Client、MQ Producer 等
- logic 内禁止直接 `sqlx.NewMysql(...)` 等自行建立连接

#### Scenario: 新增外部依赖

- **WHEN** 某 logic 需要访问新的中间件或下游 RPC
- **THEN** 在 `svc.ServiceContext` 增加字段并在 `NewServiceContext` 中初始化
- **AND** logic 通过 `l.svcCtx.<依赖>` 使用，不自行构造

### Requirement: 标准包布局

服务目录 SHALL 保持 go-zero 生成的标准布局（data 服务多嵌套一级）。

```
apps/<svc>/
├── api/
│   ├── <svc>.api              # API 定义（源头）
│   ├── etc/<svc>-api.yaml     # API 配置
│   ├── internal/
│   │   ├── config/            # goctl 生成，不手改
│   │   ├── handler/           # goctl 生成；handler 只做参数绑定 + result.Write
│   │   ├── logic/             # ★ 手写业务逻辑核心
│   │   ├── svc/               # goctl 生成，手改注入依赖
│   │   └── types/             # goctl 生成，请求响应结构体
│   └── <svc>.go               # main 入口
└── rpc/
    ├── <svc>.proto            # Proto 定义（源头）
    ├── etc/<svc>.yaml         # RPC 配置
    ├── internal/
    │   ├── config/、logic/、model/、server/、svc/
    ├── <svc>.go               # RPC client 对外暴露（goctl 1.10 命名）
    └── pb/                    # 生成的 pb.go
```

公共库 `pkg/` 约定：

```
pkg/
├── auth/          # JWT 上下文工具
├── mq/            # RabbitMQ 生产者/消费者封装
├── response/      # 统一响应 R{Code,Msg,RequestId,Data}
├── utils/
│   ├── idgen/     # 雪花算法 ID 生成
│   └── page/      # 分页请求/响应结构
└── xerr/          # 统一错误码定义
```

#### Scenario: data 服务的路径差异

- **WHEN** 定位 data 服务的配置与入口
- **THEN** 使用 `apps/data/api/data/etc/data-api.yaml` 与 `apps/data/rpc/data/etc/data.yaml`
- **AND** 其余 12 个服务为 `apps/<svc>/{api,rpc}/etc/`，data 仅路径更深一级，模块结构同构

### Requirement: Import 规则

import 分组 SHALL 按「标准库 → 第三方 → 项目内部」顺序排列，组内按字母序。

```go
import (
    "context"
    "time"
)

import (
    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

import (
    "tjxt/apps/auth/rpc/auth"   // RPC client
    "tjxt/pkg/response"
    "tjxt/pkg/xerr"
)
```

禁止项：

- 相对路径 import（`import "../logic"`）
- 跨服务 import `internal`（`import "tjxt/apps/user/internal/logic"`）
- 循环依赖

#### Scenario: 跨服务复用能力

- **WHEN** 服务 A 需要服务 B 的能力
- **THEN** 通过 B 的 RPC client（`tjxt/apps/<b>/rpc/<b>`）调用，或经 MQ 事件解耦
- **AND** import B 的 `internal` 包一律视为违规，Code Review 应拒绝

### Requirement: 测试与文件头约定

测试文件与包注释 SHALL 遵守命名约定，测试通过 `make test` 运行。

| 文件 | 约定 |
|------|------|
| `*_test.go` | 同包测试，白盒 |
| `*_test.go`（外部包） | `package logic_test`，黑盒测试 |
| `testdata/` | 测试固定数据 |
| 包注释 | 可选但建议加，如 `// Package logic implements business logic for coupon service.` |

#### Scenario: 新增单元测试

- **WHEN** 为 logic 编写测试
- **THEN** 白盒测试放同包 `*_test.go`，需黑盒时使用 `package logic_test`
- **AND** 通过 `make test`（`go test ./...`）验证
