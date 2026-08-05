> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/pay/rpc/pay.proto`

---

# Pay RPC Spec

## 服务名

`Pay` — 支付与退款微服务，通过 etcd 服务发现（key: `pay.rpc`）。

职责边界：管理支付渠道、支付单、退款单三类实体的全生命周期状态，对接第三方支付网关（当前为 mock 实现）。

## RPC 方法总览

### 支付渠道

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `ListPayChannels` | `ListPayChannelsRequest {}` | `ListPayChannelsResponse { list }` | 列出所有启用渠道，按 `channel_priority` 升序 |
| `AddPayChannel` | `PayChannelRequest { id, name, channel_code, channel_priority, channel_icon }` | `PayChannelIdResponse { id }` | 新增渠道，`channel_code` 唯一 |
| `UpdatePayChannel` | `PayChannelRequest` | `EmptyResponse {}` | 更新渠道，`channel_code` 不可改 |
| `UpdatePayChannelStatus` | `UpdatePayChannelStatusRequest { id, status }` | `EmptyResponse {}` | 启用/停用渠道 |
| `QueryPayChannelByCode` | `QueryPayChannelByCodeRequest { channel_code }` | `PayChannelResponse` | 按编码查渠道 |
| `PageQueryPayChannels` | `PageQueryPayChannelsRequest { page_no, page_size, name, channel_code, status }` | `PageQueryPayChannelsResponse { total, pages, list }` | 分页查询渠道 |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 渠道 ID，新增时必须为 0 |
| `name` | string | 渠道名称，新增时必填 |
| `channel_code` | string | 渠道编码，用于获取支付实现，唯一且不可修改 |
| `channel_priority` | int32 | 渠道优先级，数字越小优先级越高 |
| `channel_icon` | string | 渠道图标地址 |
| `status` | int32 | 1-使用中，2-停用 |

---

### 支付订单

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `ApplyPayOrder` | `ApplyPayOrderRequest { biz_user_id, biz_order_no, amount, pay_channel_code, pay_type, notify_url, expand_json, pay_over_seconds }` | `ApplyPayOrderResponse { qr_code_url }` | 申请支付单，按 `biz_order_no` 幂等 |
| `QueryPayResult` | `QueryPayResultRequest { biz_order_no }` | `PayResultResponse { pay_order_no, biz_order_no, status }` | 轻量查支付状态 |
| `NotifyPaySuccess` | `NotifyPaySuccessRequest { pay_order_no, result_code, result_msg, qr_code_url }` | `EmptyResponse {}` | 渠道回调：标记支付成功 |
| `NotifyPayFailed` | `NotifyPayFailedRequest { pay_order_no, result_code, result_msg }` | `EmptyResponse {}` | 渠道回调：标记支付失败并关单 |
| `ClosePayOrder` | `ClosePayOrderRequest { pay_order_no }` | `EmptyResponse {}` | 业务端主动关单 |
| `QueryPayOrderByBizOrderNo` | `QueryPayOrderRequest { biz_order_no }` | `PayOrderResponse` | 查支付单完整详情（19 个字段） |

**`ApplyPayOrderRequest` 字段说明**：

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

**回调类请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `pay_order_no` | int64 | 支付单号（非业务订单号），必须 > 0 |
| `result_code` | string | 第三方返回业务码 |
| `result_msg` | string | 第三方返回提示信息 |
| `qr_code_url` | string | `NotifyPaySuccessRequest` 中定义，**logic 未使用** |

---

### 退款订单

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `ApplyRefund` | `ApplyRefundRequest { biz_order_no, biz_refund_order_no, refund_amount }` | `RefundResultResponse` | 申请退款，按 `biz_refund_order_no` 幂等 |
| `QueryRefundResult` | `QueryRefundResultRequest { biz_refund_order_no }` | `RefundResultResponse` | 轻量查退款状态 |
| `NotifyRefundSuccess` | `NotifyRefundSuccessRequest { refund_order_no, result_code, result_msg, refund_channel }` | `EmptyResponse {}` | 渠道回调：标记退款成功 |
| `NotifyRefundFailed` | `NotifyRefundFailedRequest { refund_order_no, result_code, result_msg }` | `EmptyResponse {}` | 渠道回调：标记退款失败 |
| `QueryRefundByBizRefundNo` | `QueryRefundRequest { biz_refund_order_no }` | `RefundOrderResponse` | 查退款单完整详情（17 个字段） |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `biz_order_no` | int64 | 业务端已支付的订单 ID，必须 > 0 |
| `biz_refund_order_no` | int64 | 业务端要退款的订单 ID（子订单 ID），必须 > 0，幂等键 |
| `refund_amount` | int64 | 本次退款金额，单位分，必须 > 0 |
| `refund_order_no` | int64 | 退款单号，每次退款的唯一标识，由服务端雪花生成 |
| `refund_channel` | string | 退款渠道，成功回调时写入 |

**`RefundResultResponse` 字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `refund_order_no` | int64 | 退款单号 |
| `biz_refund_order_no` | int64 | 业务退款单号 |
| `refund_amount` | int64 | 退款金额（分） |
| `status` | int32 | 0-未提交, 1-退款中, 2-退款失败, 3-退款成功 |
| `result_msg` | string | 第三方交易信息 |

---

## 被哪些 API 服务消费

| 消费方 | 调用方式 | 说明 |
|--------|---------|------|
| `pay-api` (自身 API 层) | `apps/pay/api/internal/svc/servicecontext.go` import `payclient "tjxt/apps/pay/rpc/pay"` | 渠道管理、下单、关单、回调、退款、结果查询共 13 个 handler 全部转发到自身 RPC |
| `trade-rpc` | `apps/trade/rpc/internal/svc/servicecontext.go:10` import `payclient "tjxt/apps/pay/rpc/pay"`，装配为 `PayRpc payclient.Pay` | 交易域调用支付/退款下单、关单、支付/退款结果查询（见 `apps/trade/rpc/internal/config/config.go:24` 注释） |

> **⚠️ 现状说明**：`trade.rpc` 已在 `servicecontext.go:38` 完成 `PayRpc` 装配、`trade.yaml` 也已配置 `PayRpc.Etcd.Key: pay.rpc`，但 `apps/trade/rpc/internal/logic/` 下**暂无任何实际调用点**，属于已接线未落地的依赖。

---

## 调用典型场景

1. **收银台渲染** → 前端进入支付页调 `ListPayChannels` → 按 `channel_priority` 升序展示可用渠道
2. **发起支付** → 交易域生成 `biz_order_no` 后调 `ApplyPayOrder` → 返回 `qr_code_url` 供前端渲染二维码
3. **重复下单** → 同一 `biz_order_no` 再次调 `ApplyPayOrder` → 命中幂等，直接返回原二维码（状态为待支付时）
4. **支付结果轮询** → 前端定时调 `QueryPayResult` → 读到 `status = 3` 后跳转成功页
5. **渠道异步回调** → 第三方网关 → 平台 notify 接口（验签）→ `NotifyPaySuccess` / `NotifyPayFailed` 落地本地状态
6. **用户取消 / 超时** → 业务端调 `ClosePayOrder` → 支付单置为已关闭（状态 2）
7. **申请退款** → 交易域生成 `biz_refund_order_no` 后调 `ApplyRefund` → 校验原单已支付 + 累计退款不超额 → 创建退款单
8. **退款结果确认** → 轮询 `QueryRefundResult`，或由渠道回调 `NotifyRefundSuccess` / `NotifyRefundFailed` 推进状态
9. **对账取详情** → 运营/对账任务调 `QueryPayOrderByBizOrderNo` / `QueryRefundByBizRefundNo` 拉取含通知次数、结果码的完整单据

---

## 自定义 Model 方法

`paychannelmodel.go` 扩展了：
- `FindAllEnabled(ctx)` — 列出所有 `status = 1` 的渠道，按 `channel_priority asc` 排序
- `FindByCode(ctx, code)` — 按 `channel_code` 查询单条
- `PageList(ctx, name, channelCode, status, offset, limit)` — 分页查询，`name` 模糊 / `channel_code` 精确 / `status > 0` 才过滤

`payordermodel.go` 扩展了：
- `MarkToPaying(ctx, id, qrCodeUrl)` — 待提交 → 待支付（status=1），写入二维码
- `MarkToSuccess(ctx, id, resultCode, resultMsg)` — 待支付 → 支付成功（status=3），写 `pay_success_time`
- `MarkToClosed(ctx, id, resultCode, resultMsg)` — 待提交/待支付 → 关闭（status=2）
- `IncrNotifyTimes(ctx, id)` — `notify_times` 自增 1
- `SetNotifyStatus(ctx, id, status)` — 设置业务端回调状态

`refundordermodel.go` 扩展了：
- `FindOneByBizRefundOrderNo(ctx, bizRefundOrderNo)` — 按业务退款单号查询，`order by id desc limit 1`
- `FindOneByRefundOrderNo(ctx, refundOrderNo)` — 按退款单号查询
- `FindListByBizOrderNo(ctx, bizOrderNo)` — 列出某业务订单下所有未删除退款单
- `MarkToProcessing(ctx, id)` — 未提交 → 退款中（status=1）
- `MarkToSuccess(ctx, id, resultCode, resultMsg, refundChannel)` — 退款中 → 退款成功（status=3）
- `MarkToFailed(ctx, id, resultCode, resultMsg)` — 退款中 → 退款失败（status=2）
- `SetNotifyStatus(ctx, id, status)` — 设置退款回调状态
- `IncrNotifyFailedTimes(ctx, id)` — `notify_failed_times` 自增 1

> `pay_order` 的 `FindOneByBizOrderNo` / `FindOneByPayOrderNo` 由 goctl 依据唯一索引自动生成（见 `payordermodel_gen.go:113` / `:133`），带缓存；退款单的同类方法为手写且**不走缓存**（`QueryRowNoCacheCtx`）。
>
> `IncrNotifyTimes` / `SetNotifyStatus` / `IncrNotifyFailedTimes` / `MarkToPaying` / `MarkToProcessing` 已定义但当前无 logic 调用，为回调重试机制预留。
