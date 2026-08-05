# 代码风格与包布局规范

> 版本：v1.0 | 更新：2026-08-05

## Go 代码风格

### 基础规范
- **gofmt/goimports**：提交前必须 `make fmt`
- **行长度**：建议 ≤ 120 字符
- **命名**：
  - 包名：小写单词，无下划线（`usercoupon`、`svc`）
  - 接口：单词 + `er` 后缀（`UserStore`、`CouponIssuer`）
  - 实现：接口名去掉 `er`（`userStore`、`couponIssuer`）
  - 私有：小写开头（`userModel`、`calcDiscount`）
  - 常量：全大写下划线（`CouponTypeDiscount`、`MaxRetryTimes`）

### 错误处理
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

// ❌ 禁止 fmt.Errorf 包装标准库错误（丢失堆栈）
return nil, fmt.Errorf("db error: %w", err)
```

### Context 传递
- 入参第一位永远是 `ctx context.Context`
- 跨层透传，不丢弃，不新建空 context
- 超时控制在 RPC client 调用处用 `context.WithTimeout`

### 依赖注入
- `svc.ServiceContext` 持有所有外部依赖（DB、Redis、RPC Client）
- Logic 通过 `NewXxxLogic(ctx, svcCtx)` 构造，不直接 `sqlx.NewMysql` 等

---

## 包布局标准（go-zero 默认）

```
apps/<svc>/
├── api/
│   ├── <svc>.api              # API 定义（源头）
│   ├── etc/
│   │   └── <svc>-api.yaml     # API 配置
│   ├── internal/
│   │   ├── config/            # goctl 生成，不手改
│   │   ├── handler/           # goctl 生成路由，handler 只做参数绑定+result.Write
│   │   ├── logic/             # ★ 手写业务逻辑核心
│   │   ├── svc/               # goctl 生成，手改注入依赖
│   │   └── types/             # goctl 生成，请求响应结构体
│   └── <svc>.go               # main 入口
├── rpc/
│   ├── <svc>.proto            # Proto 定义（源头）
│   ├── etc/
│   │   └── <svc>.yaml         # RPC 配置
│   ├── internal/
│   │   ├── config/            # goctl 生成
│   │   ├── logic/             # ★ 手写 RPC 业务逻辑
│   │   ├── model/             # ★ Custom Model (xxxmodel.go) + 生成 (*_gen.go)
│   │   ├── server/            # goctl 生成 gRPC server 注册
│   │   └── svc/               # goctl 生成，手改注入依赖
│   ├── <svc>.go               # RPC client 对外暴露（goctl 1.10 命名）
│   └── pb/                    # 生成的 pb.go
└── go.mod                     # module tjxt/apps/<svc>
```

### 公共库 `pkg/`
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

---

## Import 规则

```go
// 1. 标准库
import (
    "context"
    "time"
)

// 2. 第三方库（按字母序）
import (
    "github.com/zeromicro/go-zero/core/logx"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
    "gorm.io/gorm"
)

// 3. 项目内部（按模块路径）
import (
    "tjxt/apps/auth/rpc/auth"           // RPC client
    "tjxt/pkg/response"
    "tjxt/pkg/xerr"
)
```

**禁止**：
- 相对路径 import（`import "../logic"`）
- 跨服务 import internal（`import "tjxt/apps/user/internal/logic"`）
- 循环依赖

---

## 文件头注释（可选，建议加）

```go
// Package logic implements business logic for coupon service.
package logic
```

---

## 单元测试命名

| 文件 | 约定 |
|------|------|
| `*_test.go` | 同包测试，白盒 |
| `*_test.go` (外部包) | `package logic_test`，黑盒测试 |
| `testdata/` | 测试固定数据 |

运行：`make test`（等同于 `go test ./...`）