# API 统一契约规范 Specification

## Purpose

统一 tjxt 全部 HTTP 接口（`apps/*/api/*.api` 及其生成的 handler）的响应封装、分页、鉴权、错误码、路由与参数绑定规则，使所有服务的对外 API 表现一致：同一套 `pkg/response.R` 信封、同一套分页结构、同一套 `pkg/xerr` 错误码分段、同一套 RESTful 路由与 `/api/v1` 前缀。本规格是 API 层代码生成与 Code Review 的判定依据。

## Requirements

### Requirement: 统一响应封装

所有 HTTP 接口 SHALL 返回 `pkg/response.R` 结构，且 handler 必须通过 `result.Write(w, r, data, err)` 输出。

```go
// pkg/response/json.go
type R struct {
    Code      int    `json:"code"`       // 0=成功，非0=失败（见 specs/shared/error-codes）
    Msg       string `json:"msg"`
    RequestId string `json:"requestId"`  // 链路追踪 ID
    Data      any    `json:"data"`
}
```

- ✅ 正确：`func (h *XxxHandler) Get(w http.ResponseWriter, r *http.Request) { resp, err := h.logic.Get(r.Context()); result.Write(w, r, resp, err) }`
- ❌ 禁止：`httpx.OkJsonCtx`（无统一码、无 requestId）、`httpx.ErrorCtx`（无统一错误格式）、使用 goctl 生成的 `Result` 类型

#### Scenario: 成功响应

- **WHEN** logic 返回数据与 nil 错误
- **THEN** handler 调用 `result.Write(w, r, resp, err)` 输出 `R{Code:0, Msg:"", RequestId:<traceId>, Data:resp}`
- **AND** `Data` 字段在成功时非 nil

#### Scenario: 失败响应

- **WHEN** logic 返回 `xerr.NewErrCode(...)` 包装的错误
- **THEN** 响应体 `Code` 为该错误码，`Msg` 为人类可读提示，`RequestId` 可用于关联链路
- **AND** 不得绕过 `result.Write` 直接写 JSON

### Requirement: 统一分页契约

分页请求 SHALL 使用 `pkg/utils/page.PageRequest`，响应 SHALL 使用 `pkg/utils/page.PageResponse[T]`。

```go
type PageRequest struct {
    PageNo   int `form:"pageNo,default=1"    validate:"gte=1"`
    PageSize int `form:"pageSize,default=20" validate:"gte=1,lte=100"`
}

type PageResponse[T any] struct {
    Total    int64 `json:"total"`
    List     []T   `json:"list"`
    PageNo   int   `json:"pageNo"`
    PageSize int   `json:"pageSize"`
}
```

`.api` 中的等价声明：

```go
type PageReq {
    PageNo   int `json:"pageNo,default=1"`
    PageSize int `json:"pageSize,default=20"`
}
```

#### Scenario: 默认分页参数

- **WHEN** 客户端未传 `pageNo` / `pageSize`
- **THEN** 取默认值 `pageNo=1`、`pageSize=20`
- **AND** `pageSize` 上界为 100，超出按校验失败处理

#### Scenario: 分页响应结构一致

- **WHEN** 任意分页接口返回
- **THEN** 响应 `Data` 中含 `total` / `list` / `pageNo` / `pageSize` 四个字段
- **AND** 列表为空时 `list` 为 `[]` 而非 `null`

### Requirement: 认证与身份获取

JWT 鉴权 SHALL 在 `.api` 的 `@server` 指令中声明，身份 SHALL 通过 `pkg/auth` 从 context 取得。

- 声明：`@server(jwt: Auth)`，生成代码在 `routes.go` 中以 `rest.WithJwt(serverCtx.Config.Auth.AccessSecret)` 包裹路由
- 取身份：`userId := auth.GetUserId(ctx)`（int64）、`role := auth.GetRole(ctx)`（string）；API 层亦用 `auth.UserIdFromCtx(l.ctx)`
- 禁止在 handler 中手写 JWT 解析（读 `Authorization` 头并手动 parse）
- RBAC：接口级权限以 `<resource>:<action>` 标注（如 `course:create`、`user:delete`），实际校验在 logic 或中间件

#### Scenario: 未携带有效 Token

- **WHEN** 请求访问声明了 `jwt: Auth` 的路由且未带有效 JWT
- **THEN** 在 JWT 中间件层被拒绝，不进入 logic
- **AND** 返回统一错误信封中的鉴权类错误码

#### Scenario: 服务端身份覆盖请求身份

- **WHEN** 请求体中携带 `userId` 字段
- **THEN** logic 一律以 context 中的 userId 覆盖，不信任请求体声明
- **AND** 该行为是越权防护的硬性要求，任何服务不得例外

### Requirement: 错误码分段

错误码 SHALL 统一定义在 `pkg/xerr` 并按域分段，段间不重叠；业务错误一律以 `xerr.NewErrCode(...)` 返回。

| 分类 | 范围 | 说明 |
|------|------|------|
| 系统通用 | 100000-100999 | 参数校验、鉴权、限流、熔断 |
| 用户域 | 101000-101999 | 登录、注册、Token、角色 |
| 课程域 | 102000-102999 | 课程 CRUD、章节、资源 |
| 交易域 | 103000-103999 | 订单、支付、退款、分账 |
| 学习域 | 104000-104999 | 进度、笔记、证书 |
| 支付域 | 105000-105999 | 渠道、流水、对账 |
| 媒体域 | 106000-106999 | 上传、签名、转码 |
| 营销域 | 107000-107999 | 优惠券、活动 |
| 消息域 | 108000-108999 | 站内信、短信、模板 |
| 考试域 | 109000-109999 | 题库、试卷、考试 |
| 搜索域 | 110000-110999 | 兴趣、推荐 |
| 用户中心 | 111000-111999 | 档案、后台管理 |

#### Scenario: 新增领域错误码

- **WHEN** 某域需要新的业务错误码
- **THEN** 在该域所属区间内分配，不得越界占用他域号段
- **AND** 同步更新 `specs/shared/error-codes` 规格

### Requirement: .api 文件编写规范

`.api` 源文件 SHALL 遵守类型命名、路由分组与路由风格约定。

- 类型命名：`<动作><资源>Req` / `<动作><资源>Resp`（如 `CreateCourseReq`）
- 枚举：`type CourseStatus string` + 注释列出取值（`NORMAL=正常, DRAFT=草稿, OFFLINE=下架`）
- 任意类型：goctl 不支持 `any`，用 `interface{}`
- 路由分组：`@server(prefix: /api/v1, jwt: Auth, middleware: Logging,Recovery)`
- 路由风格：RESTful；非 CRUD 动作用动词后缀，如 `post /courses/:id/publish`、`post /courses/:id/offline`
- 文档注释：用 `@doc(summary:..., description:..., tag:...)` 供 Swagger 生成

#### Scenario: 新增动作型接口

- **WHEN** 需要「发布课程」这类非 CRUD 操作
- **THEN** 定义为 `post /courses/:id/publish (PublishCourseReq) returns (CommonResp)`
- **AND** 不改写为标准 REST 动词语义，保持动作可读

### Requirement: 参数绑定规则

请求参数 SHALL 按下表从 HTTP 请求中绑定：

| 位置 | 绑定方式 | 示例 |
|------|----------|------|
| Path | `:id` → `Id int64` | `/courses/:id` |
| Query | `form:"pageNo"` | `?pageNo=2` |
| Header | `header:"X-Request-Id"` | `RequestId string` |
| Body | JSON 自动绑定 | `CreateCourseReq` |

#### Scenario: 路径与查询参数混用

- **WHEN** 接口同时含路径参数与分页查询参数
- **THEN** 路径参数由 `:id` 绑定为对应字段，分页参数由 `form` tag 绑定
- **AND** 绑定失败返回系统通用的参数校验错误码

### Requirement: 版本与兼容

对外 API 版本 SHALL 以 URL 前缀表达，并在演进中保持向后兼容。

- URL 前缀固定 `/api/v1`；大版本升级新增 `/api/v2`，不原地修改 v1
- 响应字段**只增不减**；废弃字段标记 `deprecated: true` 后保留
- 新增必填字段须提供默认值或兼容旧客户端

#### Scenario: 破坏性变更

- **WHEN** 需要调整响应结构或更换认证方式
- **THEN** 以新大版本前缀（`/api/v2`）提供，v1 在约定周期内继续可用
- **AND** 不得在 v1 中删除已有字段
