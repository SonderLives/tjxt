# AGENTS.md

## Go-Zero 项目开发规范

> 本文件是项目级开发约定，供开发人员及 AI Coding Agent（如 Cursor、Claude Code、Codex）共同遵守。
>
> 目标是保持代码结构清晰、依赖关系可控、接口契约稳定、业务可测试、服务可运维。
>
> **规范优先级：**
>
> 安全性、数据正确性和明确的业务要求优先；其次遵循本文件；最后遵循个人编码偏好。
>
> Go-Zero/goctl
>
> 生成代码应优先保持生成器默认结构，避免为了风格统一而破坏再生成能力。

## AGENTS.md 目录

[一、项目结构规范](#一-项目结构规范)
[二、命名规范](#二-命名规范)
[三、API / RPC 契约规范](#三-API--RPC-契约规范)
[四、分层架构规范](#四-分层架构规范)
[五、依赖注入规范](#五-依赖注入规范)
[六、数据库设计规范](#六-数据库设计规范)
[七、Redis 缓存规范](#七-Redis-缓存规范)
[八、错误处理规范](#八-错误处理规范)
[九、Context 与并发规范](#九-Context-与并发规范)
[十、日志和可观测性规范](#十-日志和可观测性规范)
[十一、配置管理规范](#十一-配置管理规范)
[十二、认证与安全规范](#十二-认证与安全规范)
[十三、测试规范](#十三-测试规范)
[十四、Go-Zero 代码生成规范](#十四-Go-Zero-代码生成规范)
[十五、Git 与 CI/CD 规范](#十五-Git-与-CICD-规范)
[十六、服务生命周期规范](#十六-服务生命周期规范)
[十七、最容易犯的架构错误](#十七-最容易犯的架构错误)
## 一、项目结构规范

### 1. 通用原则

1.  优先采用简单、明确、可维护的设计，不为未来不确定的需求过度设计。

2.  优先使用 Go 标准库和 Go-Zero

    已提供的能力；引入第三方依赖前说明必要性。

3.  业务代码应易于测试；外部依赖应有清晰边界。

4.  任何改动都应尽量局部化，避免无关重构。

5.  不得为了通过测试而删除、跳过或弱化有效测试。

6.  不得臆造配置项、接口、表结构、业务规则或工具命令；不确定时先检查项目现状并说明。

7.  修改前先阅读相关目录、现有实现、配置和测试，遵循项目已有且合理的约定。

### 2. 项目结构与模块划分

#### 2.1 推荐结构
``` bash
tjxt/
├── service/
│   │
│   ├── user/
│   │   ├── api/                 # 用户 HTTP API
│   │   └── rpc/                 # 用户 RPC 服务
│   │
│   ├── course/
│   │   ├── api/                 # 课程 HTTP API
│   │   └── rpc/                 # 课程 RPC 服务
│   │
│   ├── trade/
│   │   ├── api/                 # 交易 HTTP API
│   │   └── rpc/                 # 交易 RPC 服务
│   │
│   └── learning/
│       └── rpc/                 # 学习 RPC 服务
│
├── app/
│   ├── mq/                      # 消息队列消费者
│   ├── cron/                    # 定时任务
│   └── job/                     # 后台任务
│
├── common/                      # 公共组件
│   ├── auth/
│   ├── cache/
│   ├── constants/
│   ├── errors/
│   ├── jwt/
│   ├── logger/
│   ├── middleware/
│   ├── response/
│   └── utils/
│
├── deploy/                      # 部署配置
│   ├── docker/
│   └── k8s/
│
├── docs/
├── scripts/
├── buf.yaml
├── Makefile
├── go.mod
├── go.sum
└── AGENTS.md
```
这是推荐结构，不要求所有项目机械照搬。应以当前仓库实际结构和 goctl生成结构为准。

官方目录最佳实践：https://go-zero.dev/zh-cn/concepts/project-structure/?utm_source=chatgpt.com

## 二、命名规范

Go-Zero 项目建议统一使用以下规则：
| 对象       | 命名规范       | 示例                |
| -------- | ---------- | ----------------- |
| 目录       | 小写         | `middleware`      |
| Go 文件    | snake_case | `user_service.go` |
| 包名       | 小写单词       | `middleware`      |
| 结构体      | PascalCase | `UserService`     |
| 接口       | PascalCase | `UserRepository`  |
| 函数       | PascalCase | `CreateUser`      |
| 局部变量     | camelCase  | `userID`          |
| 常量       | Go 风格      | `DefaultTimeout`  |
| 测试文件     | `\_test.go` | `user_test.go`    |
| API 文件   | 小写业务名      | `user.api`        |
| Proto 文件 | 小写业务名      | `user.proto`      |

额外注意：

- 避免 `util`、`common`、`base` 这种含义模糊的包。

- 包名不要使用复数，除非有充分理由。

- 不要使用 `GetUserInfoData` 这类冗长名称。

- 不要为了统一而给所有结构体添加 `I` 前缀。

- 不要随意使用 `UserModelImpl`、`UserServiceImpl` 这种没有实际区分意义的命名。

例如：
```go
// 推荐

type UserRepository interface {
    FindByID(ctx context.Context, id int64) (*User, error)
}

// 不推荐

type IUserRepository interface {
    GetUserInfoData(ctx context.Context, id int64) (*User, error)
}
```
## 三、API / RPC 契约规范

这是 Go-Zero 项目非常重要的一部分。

### 1. API 定义

必须坚持 `.api` 作为 HTTP 接口契约的来源。
```go
type CreateUserReq {
    Username string `json:"username"`
    Password string `json:"password"`
}

type CreateUserResp {
    UserID int64 `json:"userId"`
}

@server(
    prefix: /api/v1/users
)

service user-api {
    @handler CreateUser
    post / (CreateUserReq) returns (CreateUserResp)

}
```
建议：

- URL 使用 RESTful 风格。

- 使用 `/api/v1` 管理接口版本。

- GET 查询资源，POST 创建资源，PUT/PATCH 更新资源，DELETE 删除资源。

- 请求和响应结构体明确区分。

- 不直接将数据库 Model 暴露为 API 响应。

- 分页接口统一参数和返回格式。

- 对输入参数进行长度、范围、格式校验。

- 明确哪些接口需要 JWT 认证。

### 2. RPC 定义

使用 Protobuf 定义服务契约。
```Protobuf
syntax = "proto3";

package user;

option go_package = "./user";

service User {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
}

message GetUserRequest {
  int64 user_id = 1;
}

message GetUserResponse {
  int64 user_id = 1;
  string username = 2;
}
```
注意：

- 字段编号不能随意修改或复用。

- 删除字段时使用 `reserved` 保留编号和名称。

- 优先向后兼容地增加字段。

- 避免随意修改 RPC 方法签名。

- 明确 RPC 超时和错误语义。

- 不要在 Protobuf 中直接暴露数据库内部结构。

### 3. 接口版本管理

API 版本与 RPC 版本不一定需要同步。

例如：
```
HTTP API: /api/v1/users

RPC: user.User.GetUser
```
HTTP API 可以独立演进，RPC 则通过 Protobuf 向后兼容策略管理。

## 四、分层架构规范

建议先使用 Go-Zero 原生三层结构，不要一开始就引入完整 DDD。
```bash
Handler
   │
   ▼
Logic
   │
   ├── Model
   ├── RPC Client
   ├── Domain Service（复杂业务才引入）
   └── External Service
```
各层职责必须明确：
| 层          | 允许做什么         | 不应该做什么    |
| ---------- | ------------- | --------- |
| Handler    | 参数解析、响应处理     | 业务编排、SQL  |
| Logic      | 业务规则、流程编排     | HTTP 请求解析 |
| Model      | SQL、CRUD、数据映射 | 复杂业务判断    |
| Domain     | 核心业务规则        | 直接依赖 HTTP |
| Middleware | 认证、限流、日志      | 具体业务处理    |

#### 特别注意 5 个问题

****1. Logic 不要无限膨胀****

一个 Logic 对应一个业务用例，复杂业务可以抽取 Domain Service。

****2. 不要跨层调用****

例如 Handler 不应该直接调用 Model。

****3. 不要跨服务访问 Model****

User RPC 的 Model 不应该被 Trade RPC 直接引用。

****4. 不要让公共包依赖业务模块****

依赖方向应该从业务模块指向基础设施，而不是基础设施反过来依赖业务。

****5. 不要过度抽象接口****

只有在需要隔离依赖、支持替换实现或方便测试时才引入接口。不要每个结构体都配一个接口。

## 五、依赖注入规范

Go-Zero 使用 `ServiceContext` 管理服务依赖。

建议：
```go
type ServiceContext struct {
    Config config.Config
    DB    *sql.DB
    Redis *redis.Redis
    UserModel model.UserModel
}
```
需要注意：

- 所有共享资源在启动时初始化。

- 不允许在每次请求中创建 Redis 客户端。

- 不允许在 Logic 中读取环境变量。

- 不允许在 Logic 中直接初始化数据库连接。

- 不允许随意使用全局变量保存业务状态。

- 依赖通过构造函数显式传递。

例如：
```go
func NewCreateUserLogic(
    ctx context.Context,
    svcCtx *svc.ServiceContext,
) *CreateUserLogic {
    return &CreateUserLogic{
        ctx:   ctx
        svcCtx: svcCtx,
    }
}
```
****补充：**** 不要把所有依赖都塞进一个无限膨胀的 ServiceContext。随着项目增长，可以按职责拆分依赖结构，但要兼容 goctl 的生成结构。

## 六、数据库设计规范

数据库是项目后期最容易出现维护问题的地方之一。

建议重点关注：
| 项目   | 规范                        |
| ---- | ------------------------- |
| 主键   | 统一 ID 策略                  |
| 字段名  | snake_case                |
| 表名   | 小写、统一单复数约定                |
| 时间字段 | `created_at`、`updated_at` |
| 删除   | 明确软删除或物理删除                |
| 索引   | 根据查询条件设计                  |
| 外键   | 根据架构和数据一致性要求决定            |
| SQL  | 参数化查询                     |
| 分页   | 优先考虑游标分页或合理限制 OFFSET      |

#### 1. 不要在业务代码中拼接 SQL

不推荐：
```sql
query := "SELECT * FROM users WHERE username = '" + username + "'"
```
推荐：
```sql
query := "SELECT * FROM users WHERE username = ?"
```
#### 2. 避免 SELECT \\*

明确查询字段，降低不必要的数据读取和字段变更影响。

#### 3. 注意 N+1 查询

不要在循环中不断执行数据库查询。

#### 4. 数据库操作必须考虑 Context
```
db.QueryContext(ctx, query, args...)
```
#### 5. 事务边界必须清晰

事务应该覆盖一个明确的业务操作。

不要在事务中：

- 执行耗时的 HTTP 请求。

- 等待 MQ 消息。

- 执行长时间计算。

- 进行不必要的网络调用。

#### 6. Model 层不要滥用缓存

缓存策略最好由业务场景决定，不要为了方便在所有 Model 方法中强制加入 Redis。

## 七、Redis 缓存规范

制定统一缓存 Key 命名规则：
```bash
业务:模块:资源:标识
user:profile:10001
user:token:10001
course:detail:20001
trade:order:30001
```
需要注意：

- Key 必须有业务前缀。

- 设置合理 TTL。

- 缓存删除要考虑失败重试。

- 避免缓存与数据库双写不一致。

- 防止缓存穿透、击穿和雪崩。

- 分布式锁必须考虑过期时间和持有者身份。

- 不要把 Redis 当成分布式事务管理器。

尤其是库存、订单、支付等业务，不能简单依赖缓存保证最终数据正确性。

## 八、错误处理规范

建议建立统一错误体系：
```bash
pkg/
└── errorx/
    ├── code.go
    ├── error.go
    └── mapping.go
```
统一业务错误码：
```go
const (
    CodeSuccess          = 0
    CodeInvalidParams    = 40001
    CodeUnauthorized     = 40101
    CodeForbidden        = 40301
    CodeUserNotFound     = 40401
    CodeUserAlreadyExist = 40901
    CodeInternal         = 50000

)
```
注意以下原则：

1. 不要直接使用字符串判断错误类型。

2. 不要到处重复创建相同的业务错误。

3. 使用 `%w` 或 `errors.Join` 保留错误链。

4. 使用 `errors.Is` / `errors.As` 判断错误。

5. 内部错误与对外错误信息分离。

6. 不要将数据库原始错误返回给客户端。

7. 日志只记录一次主要错误，避免层层重复打印。

例如：
```go
if errors.Is(err, sql.ErrNoRows) {
    return nil, xerr.ErrUserNotFound

}

if err != nil {
    return nil, fmt.Errorf("find user: %w", err)

}
```
## 九、Context 与并发规范

这是 Go 项目特别需要重视的地方。

#### 1. Context 必须贯穿调用链
```bash
HTTP Request
    ↓
Handler
    ↓
Logic
    ↓
RPC Client
    ↓
RPC Server
    ↓
Model
```
需要注意：

- 不要随意使用 `context.Background()` 替代请求 Context。

- 不要将 Context 存放到长期存活的全局对象中。

- 不要把 Context 作为结构体字段长期保存，尤其是跨请求复用的对象。

- 数据库查询、RPC 调用、HTTP 请求都应该支持超时和取消。

#### 2. Goroutine 必须可控

不推荐：
```go
go func() {
    doSomething()
}()
```
如果任务可能长期运行，应该考虑：

- 生命周期管理。

- 超时和取消。

- 错误处理。

- 并发数量限制。

- Panic 恢复。

- 退出时等待任务结束。

#### 3. 不要随意使用全局协程池

先判断是否真的需要并发池，不要为了性能增加复杂度。

## 十、日志和可观测性规范

建议统一使用 Go-Zero 的 `logx`，并根据项目需要集成 OpenTelemetry。

日志级别建议：
| 级别    | 使用场景          |
| ----- | ------------- |
| Debug | 调试信息          |
| Info  | 正常业务事件        |
| Warn  | 可恢复异常         |
| Error | 业务或系统执行失败     |
| Fatal | 启动阶段无法继续的严重错误 |

日志内容建议结构化：
```go
logx.WithContext(ctx).Infow(
    "msg", "create user success",
    "userId", userID,

)
```
重点注意：

- 不打印密码、Token、密钥。

- 不在高频循环中打印大量日志。

- 记录 trace_id。

- 记录请求耗时。

- 区分业务失败和系统故障。

- 生产日志必须有轮转和保留策略。

建议监控以下指标：

- QPS。

- P50、P95、P99 延迟。

- HTTP 4xx / 5xx。

- RPC 错误率。

- 数据库连接池使用情况。

- Redis 命中率。

- Goroutine 数量。

- GC 和内存使用情况。

## 十一、配置管理规范

建议统一 YAML 配置结构：
```ymal
Name: user-api
Host: 0.0.0.0
Port: 8888
Timeout: 3000
Log:
  ServiceName: user-api
  Mode: console
  Level: info
Mysql:
  DataSource: ${MYSQL_DSN}
Redis:
  Host: ${REDIS_HOST}
```
这里需要特别注意：

- 配置结构使用强类型。

- 启动时进行配置校验。

- 敏感配置不提交 Git。

- 开发、测试、生产环境隔离。

- 不要在代码中硬编码服务地址。

- 不要让业务代码直接读取配置文件。

- 配置变更要考虑是否需要重启服务。

注意：`${MYSQL_DSN}` 是否能自动展开，取决于具体配置加载方式，不能默认认为 go-zero YAML 加载器会自动处理所有环境变量占位符。需要明确实现环境变量注入或配置替换。

## 十二、认证与安全规范

如果你的项目涉及用户、订单、管理后台，建议从一开始建立安全规范。
| 安全项   | 注意事项                 |
| ----- | -------------------- |
| JWT   | 校验签名、过期时间、签发者        |
| 密码    | 使用 bcrypt 或 Argon2id |
| RBAC  | 认证和授权分离              |
| SQL   | 参数化查询                |
| CORS  | 限制允许的来源              |
| 限流    | 登录、验证码、敏感接口重点保护      |
| 敏感数据  | 日志脱敏                 |
| 文件上传  | 校验类型、大小和内容           |
| HTTPS | 生产环境强制 HTTPS         |
| 防重放   | 对关键业务设计幂等机制          |

特别注意：JWT 不是权限系统，登录成功不代表用户拥有所有接口的访问权限。

## 十三、测试规范

建议建立完整的测试体系：
```go
logic/
├── create_user_logic.go
└── create_user_logic_test.go

model/
├── user_model.go
└── user_model_test.go
```
测试覆盖：

1. 正常流程。

2. 参数异常。

3. 数据不存在。

4. 数据库错误。

5. Redis 故障。

6. RPC 超时。

7. 并发冲突。

8. 事务回滚。

9. 权限不足。

10. 重复请求。

建议使用：

- `testing`：Go 标准测试框架。

- `testify`：断言和测试套件。

- `gomock`：依赖 Mock。

- `httptest`：HTTP 测试。

- `sqlmock`：SQL 交互测试。

- Testcontainers：需要真实数据库或 Redis 时使用。

另外，不要只追求代码覆盖率。业务关键路径和异常分支比单纯覆盖率数字更重要。

## 十四、Go-Zero 代码生成规范

这是 Go-Zero 项目中特别需要单独制定的规范。

建议坚持：

强制规范

&#x20;API、RPC 基础代码统一由 goctl 生成。

例如：
```bash
# 创建 API 项目

goctl api new user

# 根据 API 定义生成代码

goctl api go -api user.api -dir .

# 创建 RPC 项目

goctl rpc new user

# 根据 Proto 生成 RPC 代码

goctl rpc protoc user.proto \\
  --go_out=. \\
  --go-grpc_out=. \\
  --zrpc_out=.
```
注意：

1. 不要手工修改生成的 Handler、Config、Types 等基础文件。

2. 不要手工修改生成的 RPC Client、Server 代码。

3. 业务实现放在 Logic 中。

4. Model 自定义查询需要注意与 goctl 生成代码的边界。

5. 固定 goctl 版本。

6. 重新生成代码后必须执行测试。

7. 不要随意修改官方生成目录结构。

****关于文件命名，我补充一个重要原则：****

Go-Zero 生成的文件名优先遵循 goctl 默认规则，不要为了统一 snake_case 而批量重命名生成文件。自定义文件再遵循团队的 snake_case 规范。

## 十五、Git 与 CI/CD 规范

建议采用以下分支策略：
```bash
main
  │
  ├── develop
  │     ├── feature/user-register
  │     ├── feature/course-manage
  │     └── feature/trade-order
  │
  ├── release/v1.0.0
  │
  └── hotfix/login-error
```
Commit 规范：
```
feat: 新增用户注册

fix: 修复用户登录异常

refactor: 重构用户查询逻辑

perf: 优化课程列表查询

test: 增加用户注册测试

docs: 更新接口文档

chore: 更新 go-zero 依赖
```
CI 建议：
```bash
go mod tidy

go fmt ./...

go vet ./...

go test ./...

go build ./...
```
生产项目建议额外加入：

- golangci-lint。

- 依赖漏洞扫描。

- 单元测试。

- API / Proto 兼容性检查。

- Docker 镜像构建。

- 镜像安全扫描。

- 部署健康检查。

## 十六、服务生命周期规范

很多项目只关注服务启动，却忽略了服务关闭。

推荐统一管理：
```
启动
  ↓
加载配置
  ↓
初始化 Logger
  ↓
初始化 DB / Redis
  ↓
初始化 RPC Client
  ↓
启动 HTTP / RPC Server
  ↓
等待退出信号
  ↓
停止接收新请求
  ↓
等待正在执行的请求完成
  ↓
停止后台任务
  ↓
关闭 DB / Redis 等资源
  ↓

退出
```
需要注意：

- 支持 SIGTERM 和 SIGINT。

- HTTP Server 使用优雅关闭机制。

- 后台 Goroutine 能够正常退出。

- MQ Consumer 停止拉取新消息后处理剩余消息。

- 关闭资源时注意依赖顺序。

- 关闭过程必须有超时限制。

## 十七、最容易犯的架构错误

这些是我认为需要重点避免的反模式。
| 错误做法              | 正确做法                 |
| ----------------- | -------------------- |
| 一个 Logic 包含所有业务   | 一个用例一个 Logic         |
| 所有代码放在 common     | 按业务领域划分              |
| 所有 Model 都加 Redis | 按业务场景设计缓存            |
| 所有服务共享数据库 Model   | 服务拥有自己的数据访问边界        |
| 每个业务都创建接口         | 只在必要时抽象              |
| 所有业务都使用 DDD       | 复杂业务才引入领域层           |
| 到处使用全局变量          | 显式依赖注入               |
| 无限创建 Goroutine    | 控制并发和生命周期            |
| 忽略 Context        | 全链路传递 Context        |
| 捕获错误后直接忽略         | 明确处理或向上传递            |
| 只测试正常流程           | 覆盖异常和边界条件            |
| 直接修改生成代码          | 修改 API / Proto 后重新生成 ## AGENTS.md 目录
