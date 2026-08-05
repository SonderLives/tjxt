> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/trade/rpc/trade.proto`

---

# Trade RPC Spec

## 服务名

`Trade` — 交易订单微服务，覆盖 cart / order / order_detail / refund_apply，以及迁移自 pay 服务的 pay-channels / pay-orders / refund-orders 业务能力，通过 etcd 服务发现（key: `trade.rpc`）。

共 **37 个 RPC 方法**，按业务分为 7 组。

## RPC 方法总览

### 购物车

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `CartAdd` | `CartAddRequest { course_id }` | `Empty {}` | 添加课程到购物车 |
| `CartList` | `Empty {}` | `CartListReply { items }` | 查询购物车中的课程列表 |
| `CartGet` | `IdRequest { id }` | `CartVO` | 按购物车条目 ID 查询 |
| `CartUpdate` | `CartUpdateRequest { id, course_id }` | `Empty {}` | 更新购物车条目 |
| `CartDelete` | `IdRequest { id }` | `Empty {}` | 删除指定购物车条目 |
| `CartBatchDelete` | `CartBatchDeleteRequest { ids }` | `Empty {}` | 批量删除购物车条目 |

**响应字段说明（`CartVO`）**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 购物车条目 ID |
| `course_id` | int64 | 课程 ID |
| `course_name` | string | 课程名称 |
| `cover_url` | string | 课程封面路径 |
| `price` | int64 | 加入购物车时的单价（分） |
| `now_price` | int64 | 当前售价（分） |
| `expired` | bool | 课程是否已下架/失效 |

---

### 订单

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `OrderPrePlace` | `PrePlaceRequest { course_ids }` | `OrderConfirmVO { order_id, total_amount, courses, discounts }` | 预下单，生成订单 ID 并返回可用优惠券信息 |
| `OrderPlace` | `PlaceOrderRequest { order_id, course_ids, coupon_ids }` | `PlaceOrderResultVO` | 正式下单 |
| `OrderFreeCourse` | `FreeCourseRequest { course_id }` | `PlaceOrderResultVO` | 免费课立刻报名 |
| `OrderPageQuery` | `OrderPageRequest { page_no, page_size, is_asc, sort_by, no_no, status }` | `OrderPageReply { total, pages, list }` | 分页查询我的订单 |
| `OrderGet` | `IdRequest { id }` | `OrderVO` | 按 ID 查询订单详细信息 |
| `OrderStatus` | `IdRequest { id }` | `PlaceOrderResultVO` | 查询订单支付状态 |
| `OrderCancel` | `IdRequest { id }` | `Empty {}` | 取消订单 |
| `OrderDelete` | `IdRequest { id }` | `Empty {}` | 删除订单 |

**响应字段说明（`PlaceOrderResultVO`）**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `order_id` | int64 | 订单 ID |
| `pay_amount` | int64 | 应付金额（分） |
| `pay_out_time` | int64 | 支付超时时间 |
| `status` | int32 | 1待支付 2已支付 3已关闭 4已完成 5已报名 6申请退款 |

**响应字段说明（`OrderVO`）**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 订单 ID |
| `create_time` | string | 下单时间 |
| `total_amount` | int64 | 订单总金额（分） |
| `real_amount` | int64 | 实付金额（分） |
| `discount_amount` | int64 | 优惠金额（分） |
| `status` | int32 | 订单状态 |
| `status_desc` | string | 状态文案 |
| `message` | string | 状态备注 |
| `coupon_desc` | string | 优惠券说明 |
| `details` | repeated `OrderDetailItemVO` | 订单明细列表 |
| `progress_nodes` | repeated `OrderProgressNodeVO` | 订单进度节点（name / desc / status / time） |

---

### 订单明细

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `OrderDetailGet` | `IdRequest { id }` | `OrderDetailAdminVO` | 管理端查询订单明细（含退款信息） |
| `OrderDetailCourseCheck` | `IdRequest { id }` | `BoolReply { value }` | 校验课程是否已购买 |
| `OrderDetailEnrollCourse` | `EnrollCourseRequest { student_ids }` | `EnrollCourseReply { items }` | 批量查学员报名课程数，map: student_id → 数量 |
| `OrderDetailEnrollNum` | `EnrollNumRequest { course_id_list }` | `EnrollNumReply { items }` | 批量查课程报名人数，map: course_id → 人数 |
| `OrderDetailPageQuery` | `OrderDetailPageRequest` | `OrderDetailPageReply { total, pages, list }` | 管理端分页查询订单明细 |
| `OrderDetailPurchaseInfo` | `PurchaseInfoRequest { course_id }` | `PurchaseInfoReply { enroll_num, real_pay_amount, refund_num }` | 课程销售汇总信息 |

**请求字段说明（`OrderDetailPageRequest`）**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `page_no` / `page_size` | int64 | 分页参数 |
| `is_asc` | bool | 是否升序 |
| `sort_by` | string | 排序字段 |
| `no_no` | int64 | 订单号筛选 |
| `id` | int64 | 明细 ID 筛选 |
| `mobile` | string | 学员手机号筛选 |
| `status` | int64 | 明细状态筛选 |
| `refund_status` | int64 | 退款状态筛选 |
| `pay_channel` | string | 支付渠道筛选 |
| `order_start_time` / `order_end_time` | string | 下单时间区间 |

---

### 支付渠道（管理端）

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `PayChannelAdd` | `PayChannelDTO { id, name, channel_code, channel_icon, channel_priority, status }` | `IdReply { id }` | 新增支付渠道 |
| `PayChannelList` | `Empty {}` | `PayChannelListReply { items }` | 查询支付渠道列表 |
| `PayChannelGet` | `IdRequest { id }` | `PayChannelDTO` | 按 ID 查询支付渠道 |
| `PayChannelDelete` | `IdRequest { id }` | `Empty {}` | 删除支付渠道 |

**请求字段说明（`PayChannelDTO`）**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 渠道 ID |
| `name` | string | 渠道名称 |
| `channel_code` | string | 渠道编码（wxPay / aliPay） |
| `channel_icon` | string | 渠道图标 |
| `channel_priority` | int32 | 展示优先级 |
| `status` | int32 | 渠道状态 |

---

### 支付（学员侧）

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `PayApply` | `PayApplyRequest { order_id, pay_channel_code }` | `PayApplyReply { qr_url }` | 发起支付，返回支付二维码 URL |
| `PayResultQuery` | `PayResultQueryRequest { biz_order_id }` | `PayResultDTO` | 查询支付结果 |
| `PayChannels` | `Empty {}` | `PayChannelVOList { items }` | 学员侧可选支付渠道列表 |

**响应字段说明（`PayResultDTO`）**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `biz_order_id` | int64 | 业务订单 ID |
| `status` | int32 | 1支付中 2失败 3成功 |
| `pay_channel` | string | 支付渠道 |
| `pay_order_no` | int64 | 支付流水号 |
| `success_time` | string | 支付成功时间 |

---

### 退款（业务侧）

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `RefundApply` | `RefundApplyRequest { biz_order_no, biz_refund_order_no, refund_amount }` | `RefundResultDTO` | 向支付渠道发起退款 |
| `RefundResultQuery` | `RefundResultQueryRequest { biz_refund_order_id }` | `RefundResultDTO` | 查询退款结果 |

**响应字段说明（`RefundResultDTO`）**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `biz_pay_order_id` | int64 | 业务支付单 ID |
| `biz_refund_order_id` | int64 | 业务退款单 ID |
| `pay_order_no` | int64 | 支付流水号 |
| `refund_order_no` | int64 | 退款流水号 |
| `status` | int32 | 1退款中 2失败 3成功 |
| `pay_channel` | string | 支付渠道 |
| `refund_channel` | string | 退款渠道 |

---

### 退款申请（审批流）

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `RefundApplyCreate` | `RefundApplyFormRequest { order_detail_id, refund_reason, question_desc }` | `Empty {}` | 学员提交退款申请 |
| `RefundApplyApprove` | `ApproveRequest { id, approve_type, approve_opinion, remark }` | `Empty {}` | 管理员审批退款申请 |
| `RefundApplyCancel` | `RefundCancelRequest { id, order_detail_id }` | `Empty {}` | 学员取消退款申请 |
| `RefundApplyDetail` | `IdRequest { id }` | `RefundApplyVO` | 学员侧查询退款申请详情 |
| `RefundApplyNext` | `Empty {}` | `RefundApplyVO` | 取下一条待审批的退款申请 |
| `RefundApplyPageQuery` | `RefundApplyPageRequest` | `RefundApplyPageReply { total, pages, list }` | 管理端分页查询退款申请 |
| `RefundApplyGet` | `IdRequest { id }` | `RefundApplyVO` | 管理端按 ID 查询退款申请 |
| `RefundApplyDelete` | `IdRequest { id }` | `Empty {}` | 删除退款申请 |

**请求字段说明（`ApproveRequest`）**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 退款申请 ID |
| `approve_type` | int64 | 1 同意，2 拒绝 |
| `approve_opinion` | string | 审批意见 |
| `remark` | string | 审批备注 |

**请求字段说明（`RefundApplyPageRequest`）**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `page_no` / `page_size` | int64 | 分页参数 |
| `is_asc` | bool | 是否升序 |
| `sort_by` | string | 排序字段 |
| `id` | int64 | 退款申请 ID 筛选 |
| `order_detail_id` | int64 | 订单明细 ID 筛选 |
| `order_id` | int64 | 订单 ID 筛选 |
| `refund_status` | int64 | 退款状态筛选 |
| `mobile` | string | 学员手机号筛选 |
| `apply_start_time` / `apply_end_time` | string | 申请时间区间 |

---

## 被哪些 API 服务消费

| 消费方 | 调用方式 | 说明 |
|--------|---------|------|
| `trade-api` (自身 API 层) | HTTP Handler → `tradeclient.Trade` RPC | `apps/trade/api/internal/svc/servicecontext.go:17` 注入 `TradeRpc tradeclient.Trade` |

（注：截至当前版本，Grep 全部 `apps/*/api/internal/svc/servicecontext.go` 与 `apps/*/rpc/internal/svc/servicecontext.go`，**只有 trade 自身 API 层**引用 `tjxt/apps/trade/rpc/trade`。trade 对 learning 的课程开通不走 RPC 直连，而是通过 RabbitMQ `order.exchange` 交换机发布 `order.pay` / `order.refund` 事件，由 learning 侧消费。）

**trade 自身依赖的下游**：

| 下游 | 注入位置 | 说明 |
|------|---------|------|
| `pay.rpc` | `apps/trade/rpc/internal/svc/servicecontext.go:38` `PayRpc payclient.Pay` | 支付/退款下单、关单、支付/退款结果查询 |
| RabbitMQ Producer | `apps/trade/rpc/internal/svc/servicecontext.go:45` `MQProducer *mq.Producer` | 支付成功/退款事件发布；初始化失败仅记录日志，不阻塞启动 |

---

## 调用典型场景

1. **加购到下单** → `CartAdd` 加入购物车 → `CartList` 查看 → `OrderPrePlace` 预下单生成 `order_id` 与可用优惠券 → `OrderPlace` 正式下单
2. **支付** → `PayChannels` 拉取可选渠道 → `PayApply` 获取 `qr_url` 扫码 → 轮询 `PayResultQuery` / `OrderStatus` 直到 `status=3`（成功）
3. **免费课报名** → `OrderFreeCourse` 直接生成已支付订单，跳过支付环节
4. **课程开通** → 支付成功后 trade 通过 `MQProducer` 向 `order.exchange` 发布 `order.pay` 事件 → learning 消费并开通课程
5. **退款审批流** → 学员 `RefundApplyCreate` → 管理员 `RefundApplyNext` 取件 → `RefundApplyApprove`（approve_type=1 同意）→ 调 `RefundApply` 向渠道发起退款 → 轮询 `RefundResultQuery`
6. **管理端运营看板** → `OrderDetailPageQuery` 分页查明细 → `OrderDetailPurchaseInfo` 查课程销售汇总 → `OrderDetailEnrollNum` 查各课程报名人数

---

## 自定义 Model 方法

trade 的 4 个自定义 Model 文件（`cartmodel.go`、`ordermodel.go`、`orderdetailmodel.go`、`refundapplymodel.go`）**当前均为 goctl 生成的空壳**，仅做 `customXxxModel` 包装，接口体内只嵌入 `xxxModel`，**尚未添加任何扩展方法**。

可用方法仅为 `*_gen.go` 提供的 CRUD：

| Model | 可用方法 |
|-------|---------|
| `CartModel` | `Insert`, `FindOne`, `Update`, `Delete` |
| `OrderModel` | `Insert`, `FindOne`, `Update`, `Delete` |
| `OrderDetailModel` | `Insert`, `FindOne`, `Update`, `Delete` |
| `RefundApplyModel` | `Insert`, `FindOne`, `Update`, `Delete` |

> **待补缺口**：proto 中的分页查询（`OrderPageQuery` / `OrderDetailPageQuery` / `RefundApplyPageQuery`）、聚合统计（`OrderDetailEnrollNum` / `OrderDetailPurchaseInfo`）、按用户查购物车（`CartList`）等能力，均需在自定义 Model 中补写扩展方法后才能实现。参考 auth 域 `rolemodel.go` 的 `FindPage` 模式。
