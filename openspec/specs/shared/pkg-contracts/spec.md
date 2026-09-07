# pkg-contracts（公共库契约）Specification

## Purpose

定义 `pkg/` 公共库对全部服务暴露的契约：JWT 身份上下文（pkg/auth）、统一响应（pkg/response）、分页（pkg/utils/page）、ID 生成（pkg/utils/idgen）、事件总线封装（pkg/mq）、错误码（pkg/xerr，见 [error-codes](../error-codes/spec.md)），以及共享库的版本兼容策略。

## Requirements

### Requirement: pkg/auth — JWT 身份上下文

Handler/Logic SHALL 只通过 `pkg/auth` 的函数获取身份，禁止直接解析 JWT：

- `GetUserId(ctx) int64` — 从 context 获取用户 ID
- `GetRole(ctx) string` — 从 context 获取角色编码
- `GetPermissions(ctx) []string` — 从 context 获取权限列表
- `SetUserContext(ctx, userId, role, perms) context.Context` — 中间件设置用户信息

#### Scenario: logic 获取当前用户

- **WHEN** logic 层需要当前登录用户 ID
- **THEN** 调用 `auth.GetUserId(ctx)` 获取
- **AND** 禁止在 handler/logic 中手动读取 `Authorization` header 并解析 token

### Requirement: pkg/response — 统一响应结构

所有 HTTP 接口 SHALL 通过 `response.Write(w, r, data, err)` 输出统一响应结构 `R{Code, Msg, RequestId, Data}`：

- `err == nil`：Code=0，Msg="success"，Data=业务数据
- `err` 为 `xerr.ErrCode`：Code=err.Code，Msg=err.Msg，Data=nil
- 其他 error：Code=100001（SERVER_ERROR），Msg=err.Error()，Data=nil
- `RequestId` 自动注入（从 context 或 header 获取）

#### Scenario: 业务错误自动映射

- **WHEN** logic 返回 `xerr.NewErrCode(xerr.PARAM_ERROR)`
- **THEN** 响应体为 `R{Code: 100002, Msg: <错误文案>, RequestId: <链路ID>, Data: null}`

#### Scenario: 禁止绕过统一响应

- **WHEN** handler 直接调用 `httpx.OkJsonCtx` / `httpx.ErrorCtx`
- **THEN** 视为违反契约（丢失统一码与 RequestId），代码评审必须拒绝

### Requirement: pkg/utils/page — 分页契约

分页请求/响应 SHALL 使用统一结构：

- 请求：`PageRequest{PageNo int(default=1, gte=1), PageSize int(default=20, gte=1, lte=100)}`
- 响应：`PageResponse[T]{Total int64, List []T, PageNo int, PageSize int}`

#### Scenario: 分页参数越界

- **WHEN** 请求携带 `pageSize=200`
- **THEN** 校验失败（PageSize 上限 100），返回参数错误

### Requirement: pkg/utils/idgen — 全局唯一 ID

全局唯一标识（订单号、流水号、券码、文件 ID 等）SHALL 使用 `idgen.NextId()`（int64 雪花 ID）或 `idgen.NextIdStr()`，禁止自行实现 ID 生成逻辑。

#### Scenario: 生成订单号

- **WHEN** trade 服务创建订单需要生成订单 ID
- **THEN** 调用 `idgen.NextId()` 获取雪花 ID，不在业务代码中自行实现 ID 生成

### Requirement: pkg/mq — 事件总线封装

事件收发 SHALL 通过 `pkg/mq` 封装：

- 生产者：`Producer.Publish(ctx, exchange, routingKey, msg)`（自动重试、确认机制）、`Producer.PublishDelay(..., delay)`（基于死信队列的延迟发布）
- 消费者：`Consumer.Consume(queue, handler func(ctx, msg []byte) error)` 注册处理器
- 事件结构体统一定义在 `pkg/mq/event/`

#### Scenario: 延迟发布

- **WHEN** 业务需要延时投递事件（如券过期提醒）
- **THEN** 使用 `PublishDelay`，由死信队列机制实现延迟

### Requirement: pkg 版本兼容策略

`pkg/` 共享库的导出符号变更 SHALL 遵循以下兼容策略：

| 变更类型 | 策略 |
|----------|------|
| 新增导出函数/类型 | 向下兼容，直接发布 |
| 修改函数签名 | 新增函数，旧函数标记 `// Deprecated` 保留 1 个大版本 |
| 删除导出符号 | 标记 Deprecated 1 个大版本后再删除 |
| 结构体新增字段 | 向下兼容（JSON 忽略未知字段） |
| 结构体删除字段 | 标记 Deprecated，保留字段但不再赋值 |

#### Scenario: 共享库变更发布

- **WHEN** 修改了 `pkg/` 中导出符号
- **THEN** 在 `pkg/CHANGELOG.md` 记录变更，发布时打 Tag `pkg/v<version>`

