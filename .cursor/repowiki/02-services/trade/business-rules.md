> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/trade/rpc/internal/logic/*.go`, `apps/trade/api/internal/logic/*.go`

---

# Trade Business Rules

## ⚠️ 实现状态

> **本服务业务逻辑尚未实现。** 全部 logic 文件仍为 goctl 生成的占位实现（函数体只有 `// todo: add your logic here and delete this line` 与零值返回）。
>
> 因此本文档中除本节外的所有规则，均**不是**对现有代码行为的描述，而是依据 `apps/trade/rpc/trade.proto` 注释、`sql/ddl/tj_trade.sql` 表结构与 `docs/tjxt.openapi.json` 契约反推出的**设计意图**，统一标注为「📋 设计意图（待实现）」。实现时须以实际代码为准并回写本文档。

### RPC 层（`apps/trade/rpc/internal/logic/`）

| 业务分组 | Logic 方法 | 文件 | 实现状态 |
|---------|-----------|------|---------|
| 购物车 | `CartAdd` | `cartaddlogic.go` | 未实现-goctl占位 |
| 购物车 | `CartList` | `cartlistlogic.go` | 未实现-goctl占位 |
| 购物车 | `CartGet` | `cartgetlogic.go` | 未实现-goctl占位 |
| 购物车 | `CartUpdate` | `cartupdatelogic.go` | 未实现-goctl占位 |
| 购物车 | `CartDelete` | `cartdeletelogic.go` | 未实现-goctl占位 |
| 购物车 | `CartBatchDelete` | `cartbatchdeletelogic.go` | 未实现-goctl占位 |
| 订单 | `OrderPrePlace` | `orderpreplacelogic.go` | 未实现-goctl占位 |
| 订单 | `OrderPlace` | `orderplacelogic.go` | 未实现-goctl占位 |
| 订单 | `OrderFreeCourse` | `orderfreecourselogic.go` | 未实现-goctl占位 |
| 订单 | `OrderPageQuery` | `orderpagequerylogic.go` | 未实现-goctl占位 |
| 订单 | `OrderGet` | `ordergetlogic.go` | 未实现-goctl占位 |
| 订单 | `OrderStatus` | `orderstatuslogic.go` | 未实现-goctl占位 |
| 订单 | `OrderCancel` | `ordercancellogic.go` | 未实现-goctl占位 |
| 订单 | `OrderDelete` | `orderdeletelogic.go` | 未实现-goctl占位 |
| 订单明细 | `OrderDetailGet` | `orderdetailgetlogic.go` | 未实现-goctl占位 |
| 订单明细 | `OrderDetailCourseCheck` | `orderdetailcoursechecklogic.go` | 未实现-goctl占位 |
| 订单明细 | `OrderDetailEnrollCourse` | `orderdetailenrollcourselogic.go` | 未实现-goctl占位 |
| 订单明细 | `OrderDetailEnrollNum` | `orderdetailenrollnumlogic.go` | 未实现-goctl占位 |
| 订单明细 | `OrderDetailPageQuery` | `orderdetailpagequerylogic.go` | 未实现-goctl占位 |
| 订单明细 | `OrderDetailPurchaseInfo` | `orderdetailpurchaseinfologic.go` | 未实现-goctl占位 |
| 支付渠道 | `PayChannelAdd` | `paychanneladdlogic.go` | 未实现-goctl占位 |
| 支付渠道 | `PayChannelList` | `paychannellistlogic.go` | 未实现-goctl占位 |
| 支付渠道 | `PayChannelGet` | `paychannelgetlogic.go` | 未实现-goctl占位 |
| 支付渠道 | `PayChannelDelete` | `paychanneldeletelogic.go` | 未实现-goctl占位 |
| 支付 | `PayApply` | `payapplylogic.go` | 未实现-goctl占位 |
| 支付 | `PayResultQuery` | `payresultquerylogic.go` | 未实现-goctl占位 |
| 支付 | `PayChannels` | `paychannelslogic.go` | 未实现-goctl占位 |
| 退款 | `RefundApply` | `refundapplylogic.go` | 未实现-goctl占位 |
| 退款 | `RefundResultQuery` | `refundresultquerylogic.go` | 未实现-goctl占位 |
| 退款申请 | `RefundApplyCreate` | `refundapplycreatelogic.go` | 未实现-goctl占位 |
| 退款申请 | `RefundApplyApprove` | `refundapplyapprovelogic.go` | 未实现-goctl占位 |
| 退款申请 | `RefundApplyCancel` | `refundapplycancellogic.go` | 未实现-goctl占位 |
| 退款申请 | `RefundApplyDetail` | `refundapplydetaillogic.go` | 未实现-goctl占位 |
| 退款申请 | `RefundApplyNext` | `refundapplynextlogic.go` | 未实现-goctl占位 |
| 退款申请 | `RefundApplyPageQuery` | `refundapplypagequerylogic.go` | 未实现-goctl占位 |
| 退款申请 | `RefundApplyGet` | `refundapplygetlogic.go` | 未实现-goctl占位 |
| 退款申请 | `RefundApplyDelete` | `refundapplydeletelogic.go` | 未实现-goctl占位 |

**RPC 已实现 0 / 总计 37。**

### API 层（`apps/trade/api/internal/logic/`）

| 业务分组 | Logic 方法 | 文件 | 实现状态 |
|---------|-----------|------|---------|
| 购物车 | `CartAdd` | `cartaddlogic.go` | 未实现-goctl占位 |
| 购物车 | `CartList` | `cartlistlogic.go` | 未实现-goctl占位 |
| 购物车 | `CartGet` | `cartgetlogic.go` | 未实现-goctl占位 |
| 购物车 | `CartUpdate` | `cartupdatelogic.go` | 未实现-goctl占位 |
| 购物车 | `CartDelete` | `cartdeletelogic.go` | 未实现-goctl占位 |
| 购物车 | `CartBatchDelete` | `cartbatchdeletelogic.go` | 未实现-goctl占位 |
| 订单 | `OrderPrePlace` | `orderpreplacelogic.go` | 未实现-goctl占位 |
| 订单 | `OrderPlace` | `orderplacelogic.go` | 未实现-goctl占位 |
| 订单 | `OrderFreeCourse` | `orderfreecourselogic.go` | 未实现-goctl占位 |
| 订单 | `OrderPage` | `orderpagelogic.go` | 未实现-goctl占位 |
| 订单 | `OrderGet` | `ordergetlogic.go` | 未实现-goctl占位 |
| 订单 | `OrderStatus` | `orderstatuslogic.go` | 未实现-goctl占位 |
| 订单 | `OrderCancel` | `ordercancellogic.go` | 未实现-goctl占位 |
| 订单 | `OrderDelete` | `orderdeletelogic.go` | 未实现-goctl占位 |
| 订单明细 | `OrderDetailGet` | `orderdetailgetlogic.go` | 未实现-goctl占位 |
| 订单明细 | `OrderDetailCourseCheck` | `orderdetailcoursechecklogic.go` | 未实现-goctl占位 |
| 订单明细 | `OrderDetailEnrollCourse` | `orderdetailenrollcourselogic.go` | 未实现-goctl占位 |
| 订单明细 | `OrderDetailEnrollNum` | `orderdetailenrollnumlogic.go` | 未实现-goctl占位 |
| 订单明细 | `OrderDetailPage` | `orderdetailpagelogic.go` | 未实现-goctl占位 |
| 订单明细 | `OrderDetailPurchaseInfo` | `orderdetailpurchaseinfologic.go` | 未实现-goctl占位 |
| 支付渠道 | `PayChannelAdd` | `paychanneladdlogic.go` | 未实现-goctl占位 |
| 支付渠道 | `PayChannelList` | `paychannellistlogic.go` | 未实现-goctl占位 |
| 支付渠道 | `PayChannelGet` | `paychannelgetlogic.go` | 未实现-goctl占位 |
| 支付渠道 | `PayChannelDelete` | `paychanneldeletelogic.go` | 未实现-goctl占位 |
| 支付 | `PayChannels` | `paychannelslogic.go` | 未实现-goctl占位 |
| 支付 | `PayOrderApply` | `payorderapplylogic.go` | 未实现-goctl占位 |
| 支付 | `PayResultQuery` | `payresultquerylogic.go` | 未实现-goctl占位 |
| 退款 | `RefundResultQuery` | `refundresultquerylogic.go` | 未实现-goctl占位 |
| 退款申请 | `RefundApplyCreate` | `refundapplycreatelogic.go` | 未实现-goctl占位 |
| 退款申请 | `RefundApplyApprove` | `refundapplyapprovelogic.go` | 未实现-goctl占位 |
| 退款申请 | `RefundApplyCancel` | `refundapplycancellogic.go` | 未实现-goctl占位 |
| 退款申请 | `RefundApplyDetail` | `refundapplydetaillogic.go` | 未实现-goctl占位 |
| 退款申请 | `RefundApplyNext` | `refundapplynextlogic.go` | 未实现-goctl占位 |
| 退款申请 | `RefundApplyPage` | `refundapplypagelogic.go` | 未实现-goctl占位 |
| 退款申请 | `RefundApplyGet` | `refundapplygetlogic.go` | 未实现-goctl占位 |
| 退款申请 | `RefundApplyDelete` | `refundapplydeletelogic.go` | 未实现-goctl占位 |

**API 已实现 0 / 总计 36。**

### 汇总

| 层 | 已实现 | 总计 | 比例 |
|----|-------|------|------|
| RPC logic | 0 | 37 | 0% |
| API logic | 0 | 36 | 0% |
| **合计** | **0** | **73** | **0%** |

**前置阻塞项**：`apps/trade/rpc/internal/model/` 下 4 个自定义 Model 均为空壳，仅有 `Insert` / `FindOne` / `Update` / `Delete`。分页、聚合、按 user_id 查询等能力需先补 Model 扩展方法。

---

## 1. 购物车管理

📋 **设计意图（待实现）**

**核心规则**：购物车按 `user_id` 隔离，条目记录加购时点的课程快照。

| 规则 | 依据 | 说明 |
|------|------|------|
| 按用户隔离 | `cart.user_id` 字段 | 所有查询/删除须带 `user_id` 条件，禁止跨用户操作 |
| 课程快照冗余 | `cart.cover_url` / `course_name` / `price` | 加购时写入，避免每次列表都回查 course 服务 |
| 现价与失效标识 | `CartVO.now_price` / `CartVO.expired` | 列表展示时需与 course 服务比对当前售价，判断是否涨价/下架 |
| 批量删除 | `CartBatchDeleteRequest.ids` | 支持一次传入多个条目 ID |

```
流程（CartAdd）— 设计意图:
  1. 从 JWT 上下文取 userId
  2. 校验 courseId > 0
  3. 查 course 服务获取 name / coverUrl / price
  4. 幂等：同 userId + courseId 已存在则不重复插入
  5. 插入 cart，写入课程快照字段
```

---

## 2. 订单下单

📋 **设计意图（待实现）**

**核心规则**：采用「预下单 → 正式下单」两段式，预下单先生成 `order_id` 并计算可用优惠券。

| 规则 | 依据 | 说明 |
|------|------|------|
| 两段式下单 | `OrderPrePlace` → `OrderPlace` 的 `order_id` 传递 | 预下单返回的 `order_id` 需在正式下单时回传，实现下单幂等 |
| 优惠券试算 | `OrderConfirmVO.discounts` (`CouponDiscountDTO`) | 预下单阶段返回 `discount_amount` 与 `rule_desc` 供前端展示 |
| 金额三元组 | `order.total_amount` / `discount_amount` / `real_amount` | 满足 `real_amount = total_amount - discount_amount`，单位分 |
| 优惠券记录 | `order.coupon_ids` (json) | 下单时把使用的优惠券 ID 数组落库 |
| 一单多课 | `order` (1) → `order_detail` (N) | 每门课程一条明细，明细独立记录 `real_pay_amount` 与退款状态 |
| 免费课直通 | `OrderFreeCourse` | 单课程报名，跳过支付环节直接置为已支付/已报名 |
| 支付超时 | `PlaceOrderResultVO.pay_out_time` | 下单时下发支付截止时间，超时应关单（`order.close_time`） |

```
流程（OrderPrePlace → OrderPlace）— 设计意图:
  1. PrePlace: 传入 course_ids
  2.   → 校验课程存在、未下架、用户未重复购买（order_detail.idx_user_course）
  3.   → 生成雪花 order_id，累加 total_amount
  4.   → 查 promotion 服务获取可用优惠券，返回 discounts
  5. Place: 回传 order_id + course_ids + coupon_ids
  6.   → 按 coupon_ids 核销优惠，计算 discount_amount / real_amount
  7.   → 落 order + N 条 order_detail，status=1（待支付）
  8.   → 返回 pay_amount / pay_out_time / status
```

---

## 3. 订单状态流转

📋 **设计意图（待实现）**

**核心规则**：`order.status` 与 `order_detail.status` 同步流转，各状态对应独立的时间戳字段。

| 状态 | 值 | 触发动作 | 时间戳字段 |
|------|----|---------|----------|
| 待支付 | 1 | 下单成功 | `create_time` |
| 已支付 | 2 | 支付回调成功 | `pay_time` |
| 已关闭 | 3 | 用户取消 / 支付超时 | `close_time` |
| 已完成 | 4 | 支付后 30 天（见 `finish_time` 注释） | `finish_time` |
| 已报名 | 5 | 课程开通成功 | - |
| 已申请退款 | 6 | 提交退款申请 | `refund_time` |

| 规则 | 依据 | 说明 |
|------|------|------|
| 取消仅限待支付 | `OrderCancel` + status 枚举 | status≠1 时应拒绝取消 |
| 删除为逻辑删除 | `order.deleted` 字段 | `OrderDelete` 置 `deleted=1`，不物理删除 |
| 明细状态无 6 | `order_detail.status` 注释止于 5 | 退款走 `order_detail.refund_status` 独立字段表达 |
| 完成时间 = 支付后 30 天 | `order.finish_time` 注释 | 30 天后订单不可再退款 |

---

## 4. 支付

📋 **设计意图（待实现）**

**核心规则**：trade 不落地支付流水，通过 `PayRpc payclient.Pay` 代理 pay 服务。

| 规则 | 依据 | 说明 |
|------|------|------|
| 支付能力代理 | `servicecontext.go:38` `PayRpc: payclient.NewPay(...)` | `PayApply` → `pay.ApplyPayOrder`，`PayResultQuery` → `pay.QueryPayResult` |
| 渠道数据在 pay 库 | `sql/ddl/tj_pay.sql` 含 `pay_channel` 表 | `PayChannel*` 5 个方法均为 pay 服务的透传 |
| 渠道编码 | `PayApplyRequest.pay_channel_code` 注释 | 取值 `wxPay` / `aliPay` |
| 二维码支付 | `PayApplyReply.qr_url` | 返回扫码 URL，前端渲染二维码 |
| 支付结果三态 | `PayResultDTO.status` 注释 | 1支付中 2失败 3成功 |
| 渠道排序 | `PayChannelVO.channel_priority` | 学员侧渠道列表按优先级展示 |
| 流水号回写 | `order.pay_order_no` / `order.pay_channel` | 支付成功后回写到 order 表 |

```
流程（PayApply）— 设计意图:
  1. 校验 order 存在、status=1（待支付）、归属当前用户
  2. 按 pay_channel_code 查渠道有效性
  3. 调 PayRpc.ApplyPayOrder(bizOrderNo=order.id, amount=real_amount, channel)
  4. 回写 order.pay_channel
  5. 返回 qr_url
```

---

## 5. 支付成功事件（MQ）

📋 **设计意图（待实现）**

**核心规则**：支付/退款成功后通过 RabbitMQ 通知 learning 服务开通/撤销课程，trade 与 learning 之间不做 RPC 直连。

| 配置项 | 值 | 说明 |
|--------|----|------|
| `RabbitMQ.Exchange` | `order.exchange` | 订单事件交换机 |
| `RabbitMQ.PayRoutingKey` | `order.pay` | 支付成功事件路由键 |
| `RabbitMQ.RefundRoutingKey` | `order.refund` | 退款成功事件路由键 |

| 规则 | 依据 | 说明 |
|------|------|------|
| Producer 容错 | `servicecontext.go:42-46` | MQ 初始化失败只 `logx.Errorf`，`MQProducer` 保持 nil，**不阻塞服务启动** |
| 发布前判空 | `MQProducer` 注释「可为 nil：MQ 未就绪时不阻塞启动」 | 业务代码发布事件前必须判 nil |
| 事件载荷 | learning `GrantCoursesRequest { user_id, course_ids }` | 消息体应携带 userId 与课程 ID 列表 |
| 课程有效期 | `order_detail.valid_duration` / `course_expire_time` | 支付成功开始计时，空表示永久有效 |

---

## 6. 退款申请审批流

📋 **设计意图（待实现）**

**核心规则**：退款分「业务审批」（`refund_apply` 表）与「渠道退款」（pay 服务）两阶段，`RefundApplyApprove` 是两者的衔接点。

| 阶段 | 方法 | `refund_apply.status` |
|------|------|----------------------|
| 学员提交申请 | `RefundApplyCreate` | 1 待审批 |
| 学员撤回 | `RefundApplyCancel` | 2 取消退款 |
| 管理员同意（`approve_type=1`） | `RefundApplyApprove` | 3 同意退款 |
| 管理员拒绝（`approve_type=2`） | `RefundApplyApprove` | 4 拒绝退款 |
| 渠道退款成功 | `RefundApply` → 回调 | 5 退款成功 |
| 渠道退款失败 | `RefundApply` → 回调 | 6 退款失败 |

| 规则 | 依据 | 说明 |
|------|------|------|
| 申请粒度为明细 | `refund_apply.order_detail_id` | 一单多课时可单独退某一门课 |
| 可退判定 | `OrderDetailItemVO.can_refund` | 需综合 `order_detail.status` 与 `order.finish_time`（支付后 30 天） |
| 审批留痕 | `approver` / `approve_opinion` / `remark` / `approve_time` | 审批人、意见、备注、时间四字段落库 |
| 失败留痕 | `failed_reason` / `finish_time` | 渠道退款失败原因与完成时间 |
| 状态双写 | `refund_apply.status` ↔ `order_detail.refund_status` | 两处枚举完全一致（1~6），须同步更新 |
| 待办取件 | `RefundApplyNext` | 无参，取下一条 `status=1` 的申请给当前审批人 |
| 渠道退款代理 | `RefundApply` → `pay.ApplyRefund` | 流水号回写 `refund_apply.refund_order_no` |

```
流程（RefundApplyApprove，approve_type=1）— 设计意图:
  1. 查 refund_apply，校验 status=1（待审批）
  2. 写 approver / approve_opinion / remark / approve_time
  3. status → 3（同意退款）
  4. 调 RefundApply → PayRpc.ApplyRefund
  5. 渠道受理后 status → 5 或 6，写 finish_time
  6. 同步 order_detail.refund_status、order.status=6
  7. 发布 order.refund 事件 → learning 撤销课程
```

---

## 状态说明

### 订单状态（`order.status` / `order_detail.status`）

| 值 | 含义 | 备注 |
|----|------|------|
| 1 | 待支付 | 默认值 |
| 2 | 已支付 | - |
| 3 | 已关闭 | 取消或超时 |
| 4 | 已完成 | 支付后 30 天 |
| 5 | 已报名 | 课程已开通 |
| 6 | 已申请退款 | 仅 `order` 表定义 |

### 退款状态（`refund_apply.status` / `order_detail.refund_status`）

| 值 | 含义 |
|----|------|
| 1 | 待审批 |
| 2 | 取消退款 |
| 3 | 同意退款 |
| 4 | 拒绝退款 |
| 5 | 退款成功 |
| 6 | 退款失败 |

### 支付/退款结果状态（`PayResultDTO.status` / `RefundResultDTO.status`）

| 值 | 含义 |
|----|------|
| 1 | 支付中 / 退款中 |
| 2 | 失败 |
| 3 | 成功 |

### 审批类型（`ApproveRequest.approve_type`）

| 值 | 含义 |
|----|------|
| 1 | 同意 |
| 2 | 拒绝 |
