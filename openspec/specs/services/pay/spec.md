# pay（支付网关）Specification

## Purpose

pay（支付网关）是支付与退款微服务（RPC 服务名 `Pay`），管理支付渠道、支付单、退款单三类实体的全生命周期状态，对接第三方支付网关（当前为 mock 实现，未接入微信/支付宝）。API 层 `pay-api`（HTTP :8808）为纯转发层，13 个 handler 全部转发到自身 RPC `pay.rpc`（gRPC :8088，etcd 服务发现 key `pay.rpc`）；`trade-rpc` 已装配 PayRpc 客户端但 logic 层暂无实际调用点。数据落在 MySQL `tj_pay` 库的 `pay_channel` / `pay_order` / `refund_order` 三张表（无数据库外键，应用层按业务单号关联），Redis 缓存主键查询，etcd `127.0.0.1:2379` 服务注册与发现。

## Requirements

### Requirement: HTTP API 支付渠道管理接口

`pay-api` SHALL 提供支付渠道的管理与查询 HTTP 接口（接口清单中认证列均为「否」、权限标签为空）：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| POST | /pay-channels | 添加支付渠道 | PayChannelDTO | R«long» |
| GET | /pay-channels/list | 查询支付渠道列表 | - | R«List«PayChannelDTO»» |
| PUT | /pay-channels/{id} | 修改支付渠道 | PayChannelDTO | R |

- 响应统一为 `pkg/response.R{Code,Msg,RequestId,Data any}`；错误码遵循共享错误码表（见 `openspec/specs/shared/error-codes`）
- 三个接口为薄封装，实际行为由 `pay.rpc` 渠道方法承载（见「RPC 支付渠道方法」与「支付渠道管理规则」）

#### Scenario: 添加渠道后查询列表

- **WHEN** 客户端 POST /pay-channels 提交 PayChannelDTO 且校验通过
- **THEN** 渠道以启用状态（status=1）入库，响应 R«long» 返回新渠道 ID
- **AND** GET /pay-channels/list 可查询到该渠道（R«List«PayChannelDTO»»）

#### Scenario: 修改渠道

- **WHEN** 客户端 PUT /pay-channels/{id} 提交 PayChannelDTO
- **THEN** 按 RPC `UpdatePayChannel` 语义更新（`channel_code` 不可修改），返回统一响应 R

### Requirement: HTTP API 支付下单与结果查询接口

`pay-api` SHALL 提供支付下单、支付渠道展示与支付结果查询 HTTP 接口（认证列均为「否」）：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| POST | /pay-orders | 扫码支付申请支付单，返回支付url地址，用于生产二维码 | PayApplyDTO | R«string» |
| GET | /pay-orders/{bizOrderId}/status | 根据业务端订单id查询支付结果 | - | R«PayResultDTO» |
| GET | /pay/channels | 获取支付渠道列表接口 | - | R«List«PayChannelVO»» |
| POST | /pay/order | 支付申请,返回支付二维码url | PayApplyFormDTO | R«string» |

- HTTP 层未提供退款与对账接口，退款申请与完整单据查询仅经 RPC 暴露（见「RPC 退款订单方法」「查询、分页与错误码」）

#### Scenario: 扫码支付申请返回二维码 URL

- **WHEN** 客户端 POST /pay-orders 提交 PayApplyDTO
- **THEN** 转发 RPC `ApplyPayOrder`，响应 R«string» 返回支付 url 地址（用于生成二维码）
- **AND** 客户端可 GET /pay-orders/{bizOrderId}/status 按业务端订单 id 查询支付结果（R«PayResultDTO»）

#### Scenario: C 端支付页取渠道并发起支付

- **WHEN** 前端进入支付页调 GET /pay/channels
- **THEN** 返回启用渠道列表（R«List«PayChannelVO»»）
- **AND** POST /pay/order 提交 PayApplyFormDTO 后返回支付二维码 url（R«string»）

### Requirement: RPC 服务契约与服务发现

pay 服务 SHALL 以 gRPC 服务名 `Pay` 注册于 etcd（key `pay.rpc`，监听 `0.0.0.0:8088`），职责边界为管理支付渠道、支付单、退款单三类实体的全生命周期状态，对接第三方支付网关（当前为 mock 实现）。消费方：

| 消费方 | 调用方式 | 说明 |
|--------|---------|------|
| `pay-api`（自身 API 层） | `apps/pay/api/internal/svc/servicecontext.go` import `payclient "tjxt/apps/pay/rpc/pay"` | 渠道管理、下单、关单、回调、退款、结果查询共 13 个 handler 全部转发到自身 RPC |
| `trade-rpc` | `apps/trade/rpc/internal/svc/servicecontext.go:10` import `payclient "tjxt/apps/pay/rpc/pay"`，装配为 `PayRpc payclient.Pay` | 交易域调用支付/退款下单、关单、支付/退款结果查询（见 `apps/trade/rpc/internal/config/config.go:24` 注释） |

> 已知缺口：`trade.rpc` 已在 `servicecontext.go:38` 完成 `PayRpc` 装配、`trade.yaml` 也已配置 `PayRpc.Etcd.Key: pay.rpc`，但 `apps/trade/rpc/internal/logic/` 下暂无任何实际调用点，属于已接线未落地的依赖。

#### Scenario: pay-api handler 转发

- **WHEN** pay-api 任一 handler 收到请求
- **THEN** logic 层直接调用 `l.svcCtx.PayRpc` 对应 RPC 方法（API 层为纯转发层）
- **AND** PayRpc 客户端按 `PayRpc.Etcd`（`127.0.0.1:2379`，key `pay.rpc`）经 etcd 发现 `pay.rpc`

#### Scenario: 交易域依赖装配

- **WHEN** trade-rpc 启动
- **THEN** `PayRpc` 客户端完成装配并指向 etcd key `pay.rpc`
- **AND** trade-rpc logic 层当前无任何实际调用点（已接线未落地）

### Requirement: RPC 支付渠道方法

`Pay` 服务 SHALL 提供 6 个支付渠道 RPC 方法：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `ListPayChannels` | `ListPayChannelsRequest {}` | `ListPayChannelsResponse { list }` | 列出所有启用渠道，按 `channel_priority` 升序 |
| `AddPayChannel` | `PayChannelRequest { id, name, channel_code, channel_priority, channel_icon }` | `PayChannelIdResponse { id }` | 新增渠道，`channel_code` 唯一 |
| `UpdatePayChannel` | `PayChannelRequest` | `EmptyResponse {}` | 更新渠道，`channel_code` 不可改 |
| `UpdatePayChannelStatus` | `UpdatePayChannelStatusRequest { id, status }` | `EmptyResponse {}` | 启用/停用渠道 |
| `QueryPayChannelByCode` | `QueryPayChannelByCodeRequest { channel_code }` | `PayChannelResponse` | 按编码查渠道 |
| `PageQueryPayChannels` | `PageQueryPayChannelsRequest { page_no, page_size, name, channel_code, status }` | `PageQueryPayChannelsResponse { total, pages, list }` | 分页查询渠道 |

关键字段约束：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 渠道 ID，新增时必须为 0 |
| `name` | string | 渠道名称，新增时必填 |
| `channel_code` | string | 渠道编码，用于获取支付实现，唯一且不可修改 |
| `channel_priority` | int32 | 渠道优先级，数字越小优先级越高 |
| `channel_icon` | string | 渠道图标地址 |
| `status` | int32 | 1-使用中，2-停用 |

#### Scenario: 收银台渲染

- **WHEN** 前端进入支付页调 `ListPayChannels`
- **THEN** 返回所有 `status = 1` 的启用渠道，按 `channel_priority` 升序展示可用渠道

### Requirement: RPC 支付订单方法

`Pay` 服务 SHALL 提供 6 个支付订单 RPC 方法：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `ApplyPayOrder` | `ApplyPayOrderRequest { biz_user_id, biz_order_no, amount, pay_channel_code, pay_type, notify_url, expand_json, pay_over_seconds }` | `ApplyPayOrderResponse { qr_code_url }` | 申请支付单，按 `biz_order_no` 幂等 |
| `QueryPayResult` | `QueryPayResultRequest { biz_order_no }` | `PayResultResponse { pay_order_no, biz_order_no, status }` | 轻量查支付状态 |
| `NotifyPaySuccess` | `NotifyPaySuccessRequest { pay_order_no, result_code, result_msg, qr_code_url }` | `EmptyResponse {}` | 渠道回调：标记支付成功 |
| `NotifyPayFailed` | `NotifyPayFailedRequest { pay_order_no, result_code, result_msg }` | `EmptyResponse {}` | 渠道回调：标记支付失败并关单 |
| `ClosePayOrder` | `ClosePayOrderRequest { pay_order_no }` | `EmptyResponse {}` | 业务端主动关单 |
| `QueryPayOrderByBizOrderNo` | `QueryPayOrderRequest { biz_order_no }` | `PayOrderResponse` | 查支付单完整详情（19 个字段） |

`ApplyPayOrderRequest` 关键字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `biz_user_id` | int64 | 支付用户 ID，必须 > 0 |
| `biz_order_no` | int64 | 业务订单号，必须 > 0，幂等键 |
| `amount` | int64 | 支付金额，单位分，必须 > 0 |
| `pay_channel_code` | string | 支付渠道编码，必填且渠道须处于启用状态 |
| `pay_type` | int32 | 1-h5, 2-小程序, 3-公众号, 4-扫码；`<= 0` 缺省为 4 |
| `notify_url` | string | 业务端回调接口地址 |
| `expand_json` | string | 拓展字段，用于传递不同渠道单独处理的参数 |
| `pay_over_seconds` | int64 | 支付超时秒数，`<= 0` 缺省为 1800（30 分钟） |

回调类请求字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `pay_order_no` | int64 | 支付单号（非业务订单号），必须 > 0 |
| `result_code` | string | 第三方返回业务码 |
| `result_msg` | string | 第三方返回提示信息 |
| `qr_code_url` | string | `NotifyPaySuccessRequest` 中定义，**logic 未使用** |

#### Scenario: 发起支付与结果轮询

- **WHEN** 交易域生成 `biz_order_no` 后调 `ApplyPayOrder`
- **THEN** 返回 `qr_code_url` 供前端渲染二维码
- **AND** 前端定时调 `QueryPayResult`，读到 `status = 3` 后跳转成功页

#### Scenario: 渠道异步回调与主动关单

- **WHEN** 第三方网关异步回调平台 notify 接口（验签）后转发 `NotifyPaySuccess` / `NotifyPayFailed`
- **THEN** 支付单本地状态落地（支付成功 / 支付失败关单）
- **AND** 用户取消或超时场景由业务端调 `ClosePayOrder`，支付单置为已关闭（状态 2）

### Requirement: RPC 退款订单方法

`Pay` 服务 SHALL 提供 5 个退款订单 RPC 方法：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `ApplyRefund` | `ApplyRefundRequest { biz_order_no, biz_refund_order_no, refund_amount }` | `RefundResultResponse` | 申请退款，按 `biz_refund_order_no` 幂等 |
| `QueryRefundResult` | `QueryRefundResultRequest { biz_refund_order_no }` | `RefundResultResponse` | 轻量查退款状态 |
| `NotifyRefundSuccess` | `NotifyRefundSuccessRequest { refund_order_no, result_code, result_msg, refund_channel }` | `EmptyResponse {}` | 渠道回调：标记退款成功 |
| `NotifyRefundFailed` | `NotifyRefundFailedRequest { refund_order_no, result_code, result_msg }` | `EmptyResponse {}` | 渠道回调：标记退款失败 |
| `QueryRefundByBizRefundNo` | `QueryRefundRequest { biz_refund_order_no }` | `RefundOrderResponse` | 查退款单完整详情（17 个字段） |

请求字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `biz_order_no` | int64 | 业务端已支付的订单 ID，必须 > 0 |
| `biz_refund_order_no` | int64 | 业务端要退款的订单 ID（子订单 ID），必须 > 0，幂等键 |
| `refund_amount` | int64 | 本次退款金额，单位分，必须 > 0 |
| `refund_order_no` | int64 | 退款单号，每次退款的唯一标识，由服务端雪花生成 |
| `refund_channel` | string | 退款渠道，成功回调时写入 |

`RefundResultResponse` 字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `refund_order_no` | int64 | 退款单号 |
| `biz_refund_order_no` | int64 | 业务退款单号 |
| `refund_amount` | int64 | 退款金额（分） |
| `status` | int32 | 0-未提交, 1-退款中, 2-退款失败, 3-退款成功 |
| `result_msg` | string | 第三方交易信息 |

#### Scenario: 申请退款

- **WHEN** 交易域生成 `biz_refund_order_no` 后调 `ApplyRefund`
- **THEN** 校验原单已支付且累计退款不超额后创建退款单
- **AND** 同一 `biz_refund_order_no` 重复申请命中幂等，直接返回原退款单

#### Scenario: 退款结果确认

- **WHEN** 需要确认退款结果
- **THEN** 可轮询 `QueryRefundResult`，或由渠道回调 `NotifyRefundSuccess` / `NotifyRefundFailed` 推进状态

### Requirement: 支付渠道管理规则

渠道编码 `channel_code` 是支付实现的路由键，一旦创建 SHALL 永不可改，以避免与历史订单对不上。具体约束：

- 新增不带 id：`AddPayChannel` 中 `id != 0` 直接 `BadRequest("新增渠道不应携带 id")`
- 名称与编码必填：`name` 或 `channel_code` 为空 → `BadRequest`
- 编码唯一性：新增前 `FindByCode` 探测，查到即 `Conflict("渠道编码已存在")`
- 编码不可修改：更新时 `channel_code` 非空且与库中值不同 → `BadRequest("渠道编码不允许修改")`
- 新增即启用：`buildPayChannel` 固定写入 `Status = PayChannelStatusEnabled`（1），无法在新增时指定为停用
- 更新为非零覆盖：`name` / `channel_priority` / `channel_icon` 仅在非零时覆盖原值
- 状态白名单：`UpdatePayChannelStatus` 仅接受 1（使用中）与 2（停用），其它 → `BadRequest("渠道状态非法")`
- ID 由数据库生成：渠道表为 `AUTO_INCREMENT`，用 `res.LastInsertId()` 取回，不走雪花 ID

查询规则：

| 方法 | 排序 / 过滤 |
|------|-----------|
| `ListPayChannels` | `FindAllEnabled` 硬编码 `where status = 1`，按 `channel_priority asc` |
| `QueryPayChannelByCode` | `channel_code` 为空 → `BadRequest`；未查到 → `NotFound("支付渠道不存在")` |
| `PageQueryPayChannels` | `name` 模糊、`channel_code` 精确、`status > 0` 才过滤；统一按 `channel_priority asc` |

- 渠道状态影响：1-使用中可被 `ListPayChannels` 列出并可用于下单；2-停用则 `ApplyPayOrder` 返回 `Conflict("支付渠道已停用")`

> 已知缺口（并发风险）：`channel_code` 唯一性仅由应用层 `FindByCode` 前置查询保证，DDL 未建唯一索引且探测与插入之间无锁，并发新增同一编码会插入重复行。
>
> 注：`PageList` 的状态过滤条件是 `status > 0` 而非 `>= 0`，无法筛选 status=0 的渠道（DDL 中 status 合法值仅 1/2，实际无影响）。

#### Scenario: 新增渠道流程（AddPayChannel）

- **WHEN** `AddPayChannel` 收到请求
- **THEN** 依次执行：`id != 0` → `BadRequest("新增渠道不应携带 id")`；`name` / `channel_code` 为空 → `BadRequest`；`FindByCode(channel_code)` 查到 → `Conflict("渠道编码已存在")`、非 NotFound 错误 → `Internal`
- **AND** 通过校验后以 `buildPayChannel` 固定写入 status=1，Insert 后以 `LastInsertId` 作为返回 id（不走雪花 ID）

#### Scenario: 渠道编码不可修改

- **WHEN** `UpdatePayChannel` 提交的 `channel_code` 非空且与库中值不同
- **THEN** 返回 `BadRequest("渠道编码不允许修改")`
- **AND** `name` / `channel_priority` / `channel_icon` 仅在非零时覆盖原值

#### Scenario: 渠道状态白名单

- **WHEN** `UpdatePayChannelStatus` 提交的 status 不为 1 或 2
- **THEN** 返回 `BadRequest("渠道状态非法")`

### Requirement: 支付单申请与幂等

`ApplyPayOrder` SHALL 以 `biz_order_no` 为幂等键，同一业务订单永远只对应一张支付单。具体约束：

- 参数校验：`biz_order_no` / `biz_user_id` / `amount` 任一 `<= 0` → `BadRequest`
- 渠道必填：`pay_channel_code` 为空 → `BadRequest`
- 渠道存在性：`FindByCode` 未查到 → `NotFound("支付渠道不存在")`
- 渠道可用性：`channel.Status != 1` → `Conflict("支付渠道已停用")`
- 超时缺省：`pay_over_seconds <= 0` → 30 分钟（1800 秒）
- 支付类型缺省：`pay_type <= 0` → `PayTypeNative`（4，扫码）
- 支付单号：`idgen.NextID()` 雪花生成，与自增主键 `id` 相互独立
- 初始状态：`status = PayOrderStatusPaying`（1，待支付），跳过 0-待提交
- 二维码：当前为 `mockQrCodeUrl` 生成的占位链接 `tjxt://mock-pay?order_no=..&amount=..`

命中已有单时按原单状态三分支幂等分流：

| 原单状态 | 行为 |
|---------|------|
| 1 待支付 | 直接返回原 `qr_code_url`，不新建单 |
| 3 支付成功 | `Conflict("订单已支付，请勿重复支付")` |
| 其它（0 待提交 / 2 已关闭） | `Conflict("订单已关闭，请重新下单")` |

> 已知缺口（Mock 实现）：`mockQrCodeUrl` 源码注释明确「真实项目里应该是调用微信/支付宝下单接口拿到的 code_url/prepay_id」，当前未对接任何真实支付网关，也因此 configs 中无微信/支付宝渠道配置（appId / 商户号 / 私钥 / 回调验签均缺省）。
>
> 已知缺口（超时未落地）：`pay_over_time` 已写入，但仓库中没有扫描超时单并自动关单的定时任务，超时关单目前依赖外部调用 `ClosePayOrder` 或 `NotifyPayFailed`。
>
> 并发风险：幂等查询与插入之间无锁，但 `biz_order_no` 有数据库唯一索引兜底——并发重复下单会在 Insert 处报唯一键冲突并返回 `Internal`，而非产生重复单。

#### Scenario: 首次申请建单

- **WHEN** `ApplyPayOrder` 收到合法入参（三数值 > 0、渠道非空）且 `FindOneByBizOrderNo` 未命中
- **THEN** 取 `pay_over_seconds`（缺省 1800）与 `pay_type`（缺省 4），雪花生成支付单号，Insert `pay_order { status=1, notify_times=0, notify_status=0, qr_code_url=mock }`
- **AND** 返回 `qr_code_url`

#### Scenario: 重复下单幂等分流

- **WHEN** 同一 `biz_order_no` 再次调 `ApplyPayOrder` 且原单状态为 1（待支付）
- **THEN** 命中幂等，直接返回原二维码，不新建单
- **AND** 原单状态为 3 时返回 `Conflict("订单已支付，请勿重复支付")`；原单状态为 0/2 时返回 `Conflict("订单已关闭，请重新下单")`

#### Scenario: 渠道停用禁止下单

- **WHEN** `ApplyPayOrder` 指定的渠道 `Status != 1`
- **THEN** 返回 `Conflict("支付渠道已停用")`
- **AND** 渠道编码不存在时返回 `NotFound("支付渠道不存在")`

### Requirement: 支付单状态机与回调守卫

支付单状态 SHALL 按以下语义流转（`internal/logic/common.go` 常量）：

| 值 | 含义 | 备注 |
|----|------|------|
| 0 | 待提交 | 常量已定义，但 `ApplyPayOrder` 直接建为 1，实际不产生该状态 |
| 1 | 待支付 | 建单初始态，可流转至 2 或 3 |
| 2 | 支付超时或取消 | 终态 |
| 3 | 支付成功 | 终态，退款的前置条件 |

三个变更接口 SHALL 共用同一套状态流转守卫：

| 接口 | 遇 3 支付成功 | 遇 2 已关闭 | 遇 0/1 |
|------|-------------|------------|--------|
| `NotifyPaySuccess` | 返回空响应（幂等） | `Conflict("订单已关闭，无法标记支付成功")` | `MarkToSuccess` |
| `NotifyPayFailed` | `Conflict("订单已支付成功，不能再标记失败")` | 返回空响应（幂等） | `MarkToClosed` |
| `ClosePayOrder` | `Conflict("订单已支付成功，无法关单")` | 返回空响应（幂等） | `MarkToClosed` |

- 幂等设计要点：重复回调「已是目标状态」时静默成功，而「已是另一终态」时报冲突——保证第三方网关的重试不会失败，同时阻止终态之间的非法翻转
- 关单来源与写入值：`NotifyPayFailed` 写入入参 `result_code` / `result_msg`；`ClosePayOrder` 硬编码 `"MANUAL_CLOSE"` / `"业务端主动关单"`

> 已知缺口（缓存失效）：`MarkToSuccess` / `MarkToClosed` 使用 `ExecNoCacheCtx` 直改数据库且未调用 `DelCacheCtx`，而 `FindOneByBizOrderNo` / `FindOneByPayOrderNo` 是 goctl 生成的带缓存查询，状态变更后缓存旧记录不会被清理，`QueryPayResult` 存在读到过期状态的风险。
>
> 已知缺口（回调通知未实现）：`notify_url`、`notify_times`、`notify_status` 三个字段仅在建单时写入初值；模型层虽提供 `IncrNotifyTimes` / `SetNotifyStatus`，但没有任何 logic 调用它们，「支付成功后通知业务端」这一环尚未落地（源码注释：「→ MQ/HTTP 通知业务端」）。

#### Scenario: 渠道回调标记支付成功

- **WHEN** `NotifyPaySuccess` 收到 `pay_order_no <= 0` 的请求
- **THEN** 返回 `BadRequest`
- **AND** 单号合法但 `FindOneByPayOrderNo` 未找到 → `NotFound("支付单不存在")`；状态 0/1 → `MarkToSuccess`（UPDATE status=3, result_code, result_msg, pay_success_time=now(), update_time=now()）；状态 3 → 静默返回空响应（幂等，不重复落库）；状态 2 → `Conflict("订单已关闭，无法标记支付成功")`

#### Scenario: 已支付订单不能再标记失败

- **WHEN** `NotifyPayFailed` 定位到的支付单 status=3
- **THEN** 返回 `Conflict("订单已支付成功，不能再标记失败")`
- **AND** 状态 0/1 时 `MarkToClosed` 写入入参 result_code / result_msg；已关闭（status=2）时静默幂等

#### Scenario: 业务端主动关单

- **WHEN** `ClosePayOrder` 对状态 0/1 的支付单调用
- **THEN** `MarkToClosed` 置 status=2，result_code 硬编码 `"MANUAL_CLOSE"`、result_msg 硬编码 `"业务端主动关单"`
- **AND** 对 status=3 的单返回 `Conflict("订单已支付成功，无法关单")`；对已关闭单静默幂等

### Requirement: 退款申请与累计退款上限

`ApplyRefund` SHALL 以 `biz_refund_order_no` 为幂等键；退款前必须校验原支付单已成功，且累计退款金额不得超过原支付金额。具体约束：

- 参数校验：`biz_order_no` / `biz_refund_order_no` / `refund_amount` 任一 `<= 0` → `BadRequest`
- 幂等出口：`FindOneByBizRefundOrderNo` 命中则直接返回原退款单，不做任何状态判断
- 原单存在性：`FindOneByBizOrderNo` 未查到 → `NotFound("原支付单不存在")`
- 原单必须已支付：`payOrder.Status != 3` → `Conflict("原支付单未支付成功，无法退款")`
- 单次金额上限：`refund_amount > payOrder.Amount` → `BadRequest("退款金额超过原支付金额")`
- 累计金额上限：已退（成功 + 退款中）+ 本次 > 原金额 → `BadRequest("累计退款金额超过原支付金额")`
- 快照字段：`total_amount` / `pay_channel_code` / `pay_order_no` 全部快照自原支付单
- 拆单标记：`is_split` 固定写 0，当前不支持拆单退款标记
- 初始状态：`status = RefundStatusProcessing`（1，退款中），跳过 0-未提交

累计退款金额的计算口径：对 `biz_order_no` 等于入参、`deleted = 0` 且 `status ∈ {3 退款成功, 1 退款中}` 的退款单累加 `refund_amount`。`RefundStatusProcessing`（退款中）被计入已退金额属悲观占额——防止多笔并发退款在都还未终态时各自通过校验导致超额退款；已失败（status=2）的退款单不占额，金额自动释放。

> 已知缺口（Mock 实现）：`ApplyRefund` 建单后直接 `MarkToSuccess(ro.Id, "MOCK_OK", "mock 退款成功", "mock")` 并同步返回 `{ status: 3 退款成功, result_msg: "mock 退款成功" }`；源码注释明确「真实生产：调用第三方退款 API，再根据结果 async 通过 NotifyRefundSuccess/Failed 更新；demo 中暂时直接 mock 成功」。因此当前 `ApplyRefund` 总是同步返回退款成功，`RefundStatusProcessing` 状态在实际运行中几乎不会被观测到。
>
> 已知缺口（状态不一致隐患）：上述 `MarkToSuccess` 失败仅 `l.Errorf` 记录日志不回滚，此时数据库中退款单停留在「退款中」，但 RPC 已返回「退款成功」。
>
> 已知缺口（并发风险）：`biz_refund_order_no` 无数据库唯一索引（仅普通索引 `index_biz_order_id`），幂等查询与插入之间无锁，并发重复申请会产生两张退款单；累计金额校验同样无锁，属典型 check-then-act 竞态。
>
> 注：`ApplyRefund` 中第 9 步前的 `_ = sql.ErrNoRows` 是一行无实际作用的占位语句（用于保留 `database/sql` 的 import）。

#### Scenario: 幂等命中直接返回原退款单

- **WHEN** `ApplyRefund` 收到的 `biz_refund_order_no` 已存在退款单
- **THEN** 不做任何状态判断，原样返回 `{ refund_order_no, biz_refund_order_no, refund_amount, status, result_msg }`
- **AND** 查询出错且非 NotFound 时返回 `Internal`

#### Scenario: 原单未支付禁止退款

- **WHEN** `ApplyRefund` 定位的原支付单不存在或 `Status != 3`
- **THEN** 分别返回 `NotFound("原支付单不存在")` / `Conflict("原支付单未支付成功，无法退款")`

#### Scenario: 累计退款不超额

- **WHEN** 已退金额（`FindListByBizOrderNo` 中 status ∈ {1 退款中, 3 退款成功} 的 `refund_amount` 之和）+ 本次 `refund_amount` > 原支付金额
- **THEN** 返回 `BadRequest("累计退款金额超过原支付金额")`
- **AND** 单次 `refund_amount > payOrder.Amount` 时返回 `BadRequest("退款金额超过原支付金额")`；退款失败（status=2）的单不占用额度

### Requirement: 退款单状态机与回调守卫

退款单状态 SHALL 按以下语义流转：

| 值 | 含义 | 备注 |
|----|------|------|
| 0 | 未提交 | 常量已定义，但 `ApplyRefund` 直接建为 1，实际不产生该状态 |
| 1 | 退款中 | 建单初始态，计入累计已退金额 |
| 2 | 退款失败 | 终态，不占用退款额度 |
| 3 | 退款成功 | 终态，计入累计已退金额 |

回调接口 SHALL 使用与支付单对称的状态流转守卫：

| 接口 | 遇 3 退款成功 | 遇 2 退款失败 | 遇 0/1 |
|------|-------------|-------------|--------|
| `NotifyRefundSuccess` | 返回空响应（幂等） | `Conflict("退款单已标记失败，不允许改为成功")` | `MarkToSuccess` |
| `NotifyRefundFailed` | `Conflict("退款单已成功，不能改为失败")` | 返回空响应（幂等） | `MarkToFailed` |

- 定位键为 `refund_order_no`（非业务退款单号），`<= 0` → `BadRequest`，未查到 → `NotFound("退款单不存在")`
- 写入字段：`MarkToSuccess` 写 `status=3`, `result_code`, `result_msg`, `refund_channel`, `update_time`；`MarkToFailed` 写 `status=2`, `result_code`, `result_msg`, `update_time`
- `refund_channel` 经 `sql.NullString{Valid: refundChannel != ""}` 处理，空串写入 NULL 而非空字符串

> 已知缺口（退款通知未实现）：`notify_status` / `notify_failed_times` 与模型层的 `SetNotifyStatus` / `IncrNotifyFailedTimes` 同样无 logic 调用，退款结果通知业务端的链路尚未落地。

#### Scenario: 退款成功回调幂等

- **WHEN** `NotifyRefundSuccess` 对状态 0/1 的退款单调用
- **THEN** `MarkToSuccess` 置 status=3 并写入 result_code / result_msg / refund_channel
- **AND** 对已是退款成功（status=3）的单静默返回空响应（幂等）；对 status=2 的单返回 `Conflict("退款单已标记失败，不允许改为成功")`

#### Scenario: 已成功退款不能改为失败

- **WHEN** `NotifyRefundFailed` 定位到的退款单 status=3
- **THEN** 返回 `Conflict("退款单已成功，不能改为失败")`
- **AND** 状态 0/1 时 `MarkToFailed` 置 status=2；已失败（status=2）时静默幂等

### Requirement: 查询、分页与错误码

pay 服务 SHALL 提供以下四类查询接口：

| 方法 | 定位键 | 返回粒度 |
|------|-------|---------|
| `QueryPayResult` | `biz_order_no` | 轻量：`pay_order_no` / `biz_order_no` / `status` 三字段 |
| `QueryPayOrderByBizOrderNo` | `biz_order_no` | 完整：19 字段，含通知次数、结果码、各类时间 |
| `QueryRefundResult` | `biz_refund_order_no` | 轻量：单号 / 金额 / 状态 / `result_msg` |
| `QueryRefundByBizRefundNo` | `biz_refund_order_no` | 完整：17 字段，含拆单标记、退款渠道、通知状态 |

通用规则：

- 入参校验：单号 `<= 0` → `BadRequest`
- 不存在：`isNotFound` → `NotFound("支付单不存在")` / `NotFound("退款单不存在")`
- 时间格式化：`formatTime` 零值返回空串，否则 `2006-01-02 15:04:05`
- 可空时间：`formatNullTime` 无效返回空串（用于 `pay_success_time`）
- 可空字符串：`formatNullString` 无效返回空串（用于 `qr_code_url` / `refund_channel`）
- bit 转 bool：`IsSplit: m.IsSplit == 1`

分页归一化（`normalizePage` → `page.Normalize`）：

| 入参情况 | 归一化结果 |
|---------|-----------|
| `page_no < 1` | 置为 1 |
| `page_size < 1` | 置为 10 |
| `page_size > 100` | 置为 100 |
| 总页数 | `page.CalcPages(total, limit)`，`total <= 0` 返回 0 |

错误码约定：

| 场景 | 错误构造 |
|------|---------|
| 参数非法 | `xerr.BadRequestf(...)` |
| 记录不存在 | `xerr.NotFound(...)` |
| 状态冲突 / 唯一性冲突 | `xerr.Conflict(...)` |
| 数据库异常 | `xerr.Wrapf(err, xerr.CodeInternal, ...)` |

> 注：pay 服务统一使用 `xerr.Wrapf`（user 服务用的是 `xerr.Wrap`），两者在错误包装语义上等价，仅格式化能力不同。

#### Scenario: 对账取详情

- **WHEN** 运营/对账任务需要完整单据
- **THEN** 调 `QueryPayOrderByBizOrderNo` / `QueryRefundByBizRefundNo` 拉取含通知次数、结果码的完整单据（19 / 17 字段）

#### Scenario: 分页参数归一化

- **WHEN** `PageQueryPayChannels` 收到 `page_no < 1`、`page_size < 1` 或 `page_size > 100`
- **THEN** 分别归一化为 page_no=1、page_size=10、page_size=100
- **AND** 总页数按 `page.CalcPages(total, limit)` 计算，`total <= 0` 返回 0

#### Scenario: 查询不存在的单据

- **WHEN** 按不存在的单号调 `QueryPayResult` / `QueryPayOrderByBizOrderNo`
- **THEN** `isNotFound`（同时判定 `sql.ErrNoRows` 与 `model.ErrNotFound`）命中，返回 `NotFound("支付单不存在")`
- **AND** 退款查询对应返回 `NotFound("退款单不存在")`；单号 `<= 0` 一律 `BadRequest`

### Requirement: 数据模型表清单与关系

pay 服务 SHALL 将数据持久化于 MySQL `tj_pay` 库的三张表，Model 由 goctl 生成（`*_gen.go` 禁止修改），扩展方法统一放在同名自定义 `.go` 文件中：

| 自定义 Model 文件 | 对应表 | 扩展方法 |
|------------------|--------|---------|
| `paychannelmodel.go` | `pay_channel` | FindAllEnabled, FindByCode, PageList |
| `payordermodel.go` | `pay_order` | MarkToPaying, MarkToSuccess, MarkToClosed, IncrNotifyTimes, SetNotifyStatus |
| `refundordermodel.go` | `refund_order` | FindOneByBizRefundOrderNo, FindOneByRefundOrderNo, FindListByBizOrderNo, MarkToProcessing, MarkToSuccess, MarkToFailed, SetNotifyStatus, IncrNotifyFailedTimes |

表职责与关键字段：

| 表 | 职责 | 关键字段 |
|----|------|---------|
| `pay_channel` | 支付渠道 | `channel_code`（支付实现路由键）、`channel_priority`（越小越优先）、`status`（1-使用中/2-停用，默认 1） |
| `pay_order` | 支付订单 | `biz_order_no`、`pay_order_no`（默认 0）、`amount`（单位分）、`pay_type`（1-h5/2-小程序/3-公众号/4-扫码，默认 4）、`status`（0/1/2/3，默认 0）、`pay_over_time`（非空，超时扫描依据）、`qr_code_url` |
| `refund_order` | 退款订单 | `biz_refund_order_no`（子订单 ID）、`refund_order_no`、`refund_amount` / `total_amount`（单位分）、`status`（0/1/2/3，默认 0）、`refund_channel`（成功回调时写入）、`is_split`（默认 b'0'） |

表间关系（无数据库外键，全部通过业务单号在应用层关联）：

- `pay_channel (1) ──(channel_code)── (N) pay_order`
- `pay_order (1) ──(pay_order_no / biz_order_no)── (N) refund_order`
- `refund_order.total_amount` 快照自 `pay_order.amount`；`refund_order.pay_channel_code` 快照自 `pay_order.pay_channel_code`
- `pay_order.biz_order_no` → trade 域业务订单（跨库引用）；`pay_order.biz_user_id` → user 域用户（跨库引用）

状态常量映射（`internal/logic/common.go`）：

| 常量组 | 常量名 | 值 |
|--------|--------|-----|
| 渠道状态 | `PayChannelStatusEnabled` / `PayChannelStatusDisabled` | 1 / 2 |
| 支付单状态 | `PayOrderStatusPending` / `Paying` / `Closed` / `Success` | 0 / 1 / 2 / 3 |
| 退款单状态 | `RefundStatusInit` / `Processing` / `Failed` / `Success` | 0 / 1 / 2 / 3 |
| 支付回调状态 | `NotifyStatusPending` / `OK` / `Fail` | 0 / 1 / 2 |
| 退款通知状态 | `RefundNotifyStatusPending` / `Success` / `Processing` / `Failed` | 0 / 1 / 2 / 3 |
| 渠道类型 | `PayTypeH5` / `PayTypeMini` / `PayTypeMp` / `PayTypeNative` | 1 / 2 / 3 / 4 |

#### Scenario: 渠道与订单关联及退款快照

- **WHEN** 创建支付单或退款单
- **THEN** `pay_order` 通过 `pay_channel_code` 关联 `pay_channel`、通过 `biz_order_no` 关联 trade 域业务订单
- **AND** `refund_order` 的 `total_amount` / `pay_channel_code` / `pay_order_no` 快照自原支付单，三张表之间无数据库外键

### Requirement: 数据模型关键不变量

pay 服务的存储层 SHALL 满足以下不变量：

- `pay_order` 的 `biz_order_no` 与 `pay_order_no` 为双唯一索引，是 `ApplyPayOrder` 幂等与回调定位的数据库级保障
- `pay_channel` 的 `channel_code` 唯一性仅由应用层 `AddPayChannel` 中的 `FindByCode` 前置查询保证，DDL 未建唯一索引
- `refund_order` 的 `biz_refund_order_no` 与 `refund_order_no` 均为普通索引而非唯一索引（`index_biz_order_id` / `index_refund_order_id`），因此 `FindOneByBizRefundOrderNo` 需显式 `order by id desc limit 1` 兜底；`ApplyRefund` 的幂等仅靠应用层前置查询，无数据库级约束
- `notify_status` 语义不一致：`pay_order.notify_status` 为 0-待回调 / 1-回调成功 / 2-回调失败；`refund_order.notify_status` 为 0-待通知 / 1-通知成功 / 2-通知中 / 3-通知失败。两者在 `common.go` 中定义为两组独立常量，读写时不可混用
- `pay_order.amount` 与 `refund_order.refund_amount` 为 `int`（约 21 亿分 ≈ 2147 万元上限），而 proto 中声明为 `int64`，存在类型宽度差异
- `bit(1)` 字段（`is_split` / `deleted`）goctl 映射为 `int64`：写入时 `ApplyRefund` 显式赋 `IsSplit: 0`，读出时 `toRefundResp` 以 `m.IsSplit == 1` 转为 proto 的 `bool`，过滤时 `FindListByBizOrderNo` 拼接 `and deleted = 0`
- `pay_order.deleted` 虽有该列，但当前所有查询路径（含 goctl 生成的 `FindOneByBizOrderNo` / `FindOneByPayOrderNo`）均未加 `deleted = 0` 过滤，仅 `refund_order` 的 `FindListByBizOrderNo` 做了过滤
- 手写扩展方法一律使用 `ExecNoCacheCtx` / `QueryRowNoCacheCtx` / `QueryRowsNoCacheCtx` 绕过 goctl 缓存层；`pay_order` 的 `FindOneByBizOrderNo` / `FindOneByPayOrderNo` 由 goctl 依据唯一索引自动生成并带缓存，退款单的同类方法为手写且不走缓存（`QueryRowNoCacheCtx`）——状态流转后带缓存查询可能读到旧状态
- `IncrNotifyTimes` / `SetNotifyStatus` / `IncrNotifyFailedTimes` / `MarkToPaying` / `MarkToProcessing` 已定义但当前无 logic 调用，为回调重试机制预留
- `vars.go` 定义 `ErrNotFound = sqlx.ErrNotFound`；logic 层的 `isNotFound(err)` 同时判定 `sql.ErrNoRows` 与 `model.ErrNotFound` 两种情况

> 注：本规格的表清单以 data-model.md（来源 `sql/ddl/tj_pay.sql`）为准——实际仅 `pay_channel` / `pay_order` / `refund_order` 三张表；对账能力当前仅体现为 RPC 完整详情查询（`QueryPayOrderByBizOrderNo` / `QueryRefundByBizRefundNo`），无独立对账/流水表。

#### Scenario: 唯一索引兜底并发下单

- **WHEN** 并发对同一 `biz_order_no` 多次 `ApplyPayOrder` 且幂等查询均未命中
- **THEN** Insert 阶段被 `biz_order_no` 数据库唯一索引拦截，报唯一键冲突并返回 `Internal`，不产生重复支付单
- **AND** 退款单因 `biz_refund_order_no` 无唯一索引而不具备同等级兜底，并发重复申请会产生两张退款单

### Requirement: 服务配置与依赖

pay 服务 SHALL 按以下配置运行（`apps/pay/api/etc/pay-api.yaml`、`apps/pay/rpc/etc/pay.yaml`）：

| 服务 | 配置项 | 默认值 |
|------|--------|--------|
| `pay-api` | `Name` / `Host` / `Port` | `pay-api` / `0.0.0.0` / `8808`（HTTP） |
| `pay-api` | `Auth.AccessSecret` / `Auth.AccessExpire` | `change-me-in-production` / `7200`（秒） |
| `pay-api` | `PayRpc.Etcd.Hosts[0]` / `PayRpc.Etcd.Key` | `127.0.0.1:2379` / `pay.rpc` |
| `pay.rpc` | `Name` / `ListenOn` | `pay.rpc` / `0.0.0.0:8088`（gRPC） |
| `pay.rpc` | `Etcd.Hosts[0]` / `Etcd.Key` | `127.0.0.1:2379` / `pay.rpc` |
| `pay.rpc` | `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_pay?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai` |
| `pay.rpc` | `Cache[0]` | `127.0.0.1:6379`，Pass 空，Type `node`（单机模式） |

- 依赖的外部服务：MySQL `tj_pay` 库（含 `pay_channel` / `pay_order` / `refund_order` 三张表）、Redis 缓存（节点模式，缓存渠道与支付/退款单主键查询）、etcd（RPC 注册与发现）
- `parseTime=true` 为必需项：`PayOrder.PayOverTime` / `CreateTime` 声明为 `time.Time`，`PaySuccessTime` 为 `sql.NullTime`，均需驱动直接解析
- 作为被依赖方：`pay.rpc` 被 `apps/pay/api/etc/pay-api.yaml`（`PayRpc` 节）与 `apps/trade/rpc/etc/trade.yaml`（`PayRpc` 节）声明为客户端

> 已知缺口：第三方支付网关无任何配置项——当前为 mock 实现（`mockQrCodeUrl`），未接入微信/支付宝，故无 appId / 商户号 / 私钥 / 回调验签等配置。
>
> 已知缺口：`TablePrefix` 已在 `config.Config` 中声明并在 `pay.yaml` 置为空串，但全项目无任何读取点（模型层表名由 goctl 生成时硬编码），属预留但未接线的配置项。
>
> 注（时区差异）：`auth.rpc`、`user.rpc`、`trade.rpc` 均使用 `loc=Local`，唯独 pay 服务显式指定 `loc=Asia%2FShanghai`；支付涉及 `pay_over_time`（超时时间）、`pay_success_time`（成功时间）等强时间语义字段，跨服务对账时需注意这一时区配置差异。另 DDL 中三张表字符集为 `utf8` / `utf8_general_ci`，而连接串声明 `charset=utf8mb4`，存在字符集不一致（不影响纯数字与 ASCII 场景）。
>
> 可观测性注入配置（Telemetry / Prometheus / Log 三段）：两份 yaml 中已注入勿删（约定见 `openspec/specs/infra/observability`）。

#### Scenario: API 层经 etcd 发现自身 RPC

- **WHEN** pay-api 启动并处理首个请求
- **THEN** `PayRpc` 客户端按 `PayRpc.Etcd`（`127.0.0.1:2379`，key `pay.rpc`）发现 `pay.rpc` 并完成转发调用
- **AND** 全部 13 个 logic 直接调用 `l.svcCtx.PayRpc.*`，API 层不落库、不做业务判断

