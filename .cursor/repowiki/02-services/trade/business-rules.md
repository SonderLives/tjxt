> 版本：v1.1 | 更新：2026-08-05 | 来源：`apps/trade/rpc/internal/logic/*.go`, `apps/trade/api/internal/logic/*.go`（已落地实现）

---

# Trade Business Rules

## ✅ 实现状态

> **本服务业务逻辑已全部落地。** 73 个 logic 文件（RPC 37 + API 36）均为真实实现，`go build ./...` 两层均编译通过（rc=0），全仓已无 goctl 占位 stub。
>
> 本文档描述**当前代码的实际行为**。与原 proto/SQL 反推的「设计意图」存在若干偏差，统一收录于文末「已知缺口汇总」，实现时以代码为准。

### RPC 层（`apps/trade/rpc/internal/logic/`）— 37/37 已实现

| 业务分组 | Logic 方法 | 文件 | 实现要点 |
|---------|-----------|------|---------|
| 购物车 | `CartAdd` | `cartaddlogic.go` | 幂等加购：取 userId → 校验 courseId → `FindByUserIdAndCourseId` 去重 → `fetchCourseMap` 取课程快照(name/cover/price) → `Insert` |
| 购物车 | `CartList` | `cartlistlogic.go` | `ListByUserId` → 逐条 `toCartVO`（NowPrice=快照价，Expired=false） |
| 购物车 | `CartGet` | `cartgetlogic.go` | 按 id 取单条购物车条目 |
| 购物车 | `CartUpdate` | `cartupdatelogic.go` | 按 id 更新购物车条目（归属校验） |
| 购物车 | `CartDelete` | `cartdeletelogic.go` | 按 id 删除（带 user 归属校验） |
| 购物车 | `CartBatchDelete` | `cartbatchdeletelogic.go` | 按 `ids` 批量删除 |
| 订单 | `OrderPrePlace` | `orderpreplacelogic.go` | **仅试算**，不落库：累加课程价，返回 `OrderConfirmVO{OrderId:0, TotalAmount, Courses, Discounts:空}` |
| 订单 | `OrderPlace` | `orderplacelogic.go` | **真正下单**：雪花 order_id → `Insert` order(status=1,real=total,discount=0) → 逐课 `Insert` order_detail(status=1) → 返回 OrderId/PayAmount/Status=1/PayOutTime(now+15min) |
| 订单 | `OrderFreeCourse` | `orderfreecourselogic.go` | 免费课直通：Insert order(status=2 已支付)+detail(status=2)，金额 0 |
| 订单 | `OrderPageQuery` | `orderpagequerylogic.go` | `PageQueryByUser`（status/no_no 过滤+排序）→ OrderVO 列表+分页元数据 |
| 订单 | `OrderGet` | `ordergetlogic.go` | `FindOne` → `toOrderVO`（含 details + progressNodes） |
| 订单 | `OrderStatus` | `orderstatuslogic.go` | 查单条订单 status（含 desc） |
| 订单 | `OrderCancel` | `ordercancellogic.go` | 仅 status=1 可取消 → `UpdateStatus(status=3 关闭, message)` |
| 订单 | `OrderDelete` | `orderdeletelogic.go` | 逻辑删除（deleted=1） |
| 订单明细 | `OrderDetailGet` | `orderdetailgetlogic.go` | `FindOne` → `toOrderDetailItemVO`（含 can_refund） |
| 订单明细 | `OrderDetailCourseCheck` | `orderdetailcoursechecklogic.go` | 校验用户是否已购该课程（按 user_id+course_id 查 order_detail） |
| 订单明细 | `OrderDetailEnrollCourse` | `orderdetailenrollcourselogic.go` | `CountPaidByUserIds(studentIds)` → map[studentId]付费课程数（供他服务查报名数） |
| 订单明细 | `OrderDetailEnrollNum` | `orderdetailenrollnumlogic.go` | `StatByCourseId` → 报名数/实付金额/退款数（enrollNum/realPayAmount/refundNum） |
| 订单明细 | `OrderDetailPageQuery` | `orderdetailpagequerylogic.go` | `PageQuery(filter)` 管理后台分页 |
| 订单明细 | `OrderDetailPurchaseInfo` | `orderdetailpurchaseinfologic.go` | 查某课程购买信息（人数/状态等） |
| 支付渠道 | `PayChannelAdd` | `paychanneladdlogic.go` | 委托 `PayRpc.AddPayChannel` |
| 支付渠道 | `PayChannelList` | `paychannellistlogic.go` | 委托 `PayRpc.ListPayChannels` → 映射 PayChannelDTO |
| 支付渠道 | `PayChannelGet` | `paychannelgetlogic.go` | 委托 `PayRpc.ListPayChannels` 按 id 取单条 |
| 支付渠道 | `PayChannelDelete` | `paychanneldeletelogic.go` | 委托 `PayRpc.UpdatePayChannelStatus{Status:2}` 软删 |
| 支付 | `PayApply` | `payapplylogic.go` | 校验 order 存在 → `PayRpc.ApplyPayOrder{BizUserId:order.UserId, BizOrderNo:order.Id, Amount:order.TotalAmount, PayChannelCode, PayType:4}` → 返回 QrUrl |
| 支付 | `PayResultQuery` | `payresultquerylogic.go` | 委托 `PayRpc.QueryPayResult` → PayResultDTO |
| 支付 | `PayChannels` | `paychannelslogic.go` | 委托 `PayRpc.ListPayChannels` → 映射（学员侧渠道展示） |
| 退款 | `RefundApply` | `refundapplylogic.go` | 委托 `PayRpc.ApplyRefund{BizOrderNo, BizRefundOrderNo, RefundAmount}` → RefundResultDTO |
| 退款 | `RefundResultQuery` | `refundresultquerylogic.go` | 委托 `PayRpc.QueryRefundResult` → RefundResultDTO |
| 退款申请 | `RefundApplyCreate` | `refundapplycreatelogic.go` | 登录校验 → `FindOne` detail → `Insert` refund_apply(status=1) + `UpdateRefundStatus(detail,1)` 双写 |
| 退款申请 | `RefundApplyApprove` | `refundapplyapprovelogic.go` | `FindOne` → approve_type 1→status=3 / 2→status=4 → `UpdateApprove`(approver 硬编码 0) + `UpdateRefundStatus(detail,同status)` |
| 退款申请 | `RefundApplyCancel` | `refundapplycancellogic.go` | 取消 → status=2 + 同步 detail refund_status=2 |
| 退款申请 | `RefundApplyDetail` | `refundapplydetaillogic.go` | `FindOne` + 关联 order/detail → `toRefundApplyVO` |
| 退款申请 | `RefundApplyNext` | `refundapplynextlogic.go` | `FindNextPending` → 下一条 status=1 申请 |
| 退款申请 | `RefundApplyPageQuery` | `refundapplypagequerylogic.go` | `PageQuery(filter)` 后台分页 |
| 退款申请 | `RefundApplyGet` | `refundapplygetlogic.go` | 按 id 取申请 |
| 退款申请 | `RefundApplyDelete` | `refundapplydeletelogic.go` | 删除申请记录 |

**RPC 已实现 37 / 总计 37。**

### API 层（`apps/trade/api/internal/logic/`）— 36/36 已实现

| 业务分组 | Logic 方法 | 文件 | 实现要点 |
|---------|-----------|------|---------|
| 购物车 | `CartAdd` | `cartaddlogic.go` | 转发 `TradeRpc.CartAdd{CourseId}` → 返回 `NamePlaceVO{Existed:true, Message:"ok"}` |
| 购物车 | `CartList` | `cartlistlogic.go` | 转发 `TradeRpc.CartList` → 映射 CartVO 列表 |
| 购物车 | `CartGet` | `cartgetlogic.go` | 转发 `TradeRpc.CartGet` |
| 购物车 | `CartUpdate` | `cartupdatelogic.go` | 转发 `TradeRpc.CartUpdate` |
| 购物车 | `CartDelete` | `cartdeletelogic.go` | 转发 `TradeRpc.CartDelete` |
| 购物车 | `CartBatchDelete` | `cartbatchdeletelogic.go` | 转发 `TradeRpc.CartBatchDelete` |
| 订单 | `OrderPrePlace` | `orderpreplacelogic.go` | 转发 `TradeRpc.OrderPrePlace` → 映射 OrderConfirmVO |
| 订单 | `OrderPlace` | `orderplacelogic.go` | 转发 `TradeRpc.OrderPlace` → 映射 PlaceOrderResultVO |
| 订单 | `OrderFreeCourse` | `orderfreecourselogic.go` | 转发 `TradeRpc.OrderFreeCourse` |
| 订单 | `OrderPage` | `orderpagelogic.go` | 转发 `TradeRpc.OrderPageQuery` → 映射分页 |
| 订单 | `OrderGet` | `ordergetlogic.go` | 转发 `TradeRpc.OrderGet` → 逐字段映射 OrderVO（指针切片解引用为值切片） |
| 订单 | `OrderStatus` | `orderstatuslogic.go` | 转发 `TradeRpc.OrderStatus` |
| 订单 | `OrderCancel` | `ordercancellogic.go` | 转发 `TradeRpc.OrderCancel` |
| 订单 | `OrderDelete` | `orderdeletelogic.go` | 转发 `TradeRpc.OrderDelete` |
| 订单明细 | `OrderDetailGet` | `orderdetailgetlogic.go` | 转发 `TradeRpc.OrderDetailGet` |
| 订单明细 | `OrderDetailCourseCheck` | `orderdetailcoursechecklogic.go` | 转发 `TradeRpc.OrderDetailCourseCheck` |
| 订单明细 | `OrderDetailEnrollCourse` | `orderdetailenrollcourselogic.go` | 转发 `TradeRpc.OrderDetailEnrollCourse` |
| 订单明细 | `OrderDetailEnrollNum` | `orderdetailenrollnumlogic.go` | 转发 `TradeRpc.OrderDetailEnrollNum` |
| 订单明细 | `OrderDetailPage` | `orderdetailpagelogic.go` | 转发 `TradeRpc.OrderDetailPageQuery` |
| 订单明细 | `OrderDetailPurchaseInfo` | `orderdetailpurchaseinfologic.go` | 转发 `TradeRpc.OrderDetailPurchaseInfo` |
| 支付渠道 | `PayChannelAdd` | `paychanneladdlogic.go` | 转发 `TradeRpc.PayChannelAdd` |
| 支付渠道 | `PayChannelList` | `paychannellistlogic.go` | 转发 `TradeRpc.PayChannelList` |
| 支付渠道 | `PayChannelGet` | `paychannelgetlogic.go` | 转发 `TradeRpc.PayChannelGet` |
| 支付渠道 | `PayChannelDelete` | `paychanneldeletelogic.go` | 转发 `TradeRpc.PayChannelDelete` |
| 支付 | `PayChannels` | `paychannelslogic.go` | 转发 `TradeRpc.PayChannels` |
| 支付 | `PayOrderApply` | `payorderapplylogic.go` | 转发 `TradeRpc.PayApply` |
| 支付 | `PayResultQuery` | `payresultquerylogic.go` | 转发 `TradeRpc.PayResultQuery` |
| 退款 | `RefundResultQuery` | `refundresultquerylogic.go` | 转发 `TradeRpc.RefundResultQuery` |
| 退款申请 | `RefundApplyCreate` | `refundapplycreatelogic.go` | 转发 `TradeRpc.RefundApplyCreate` |
| 退款申请 | `RefundApplyApprove` | `refundapplyapprovelogic.go` | 转发 `TradeRpc.RefundApplyApprove` |
| 退款申请 | `RefundApplyCancel` | `refundapplycancellogic.go` | 转发 `TradeRpc.RefundApplyCancel` |
| 退款申请 | `RefundApplyDetail` | `refundapplydetaillogic.go` | 转发 `TradeRpc.RefundApplyDetail` |
| 退款申请 | `RefundApplyNext` | `refundapplynextlogic.go` | 转发 `TradeRpc.RefundApplyNext` |
| 退款申请 | `RefundApplyPage` | `refundapplypagelogic.go` | 转发 `TradeRpc.RefundApplyPageQuery` |
| 退款申请 | `RefundApplyGet` | `refundapplygetlogic.go` | 转发 `TradeRpc.RefundApplyGet` |
| 退款申请 | `RefundApplyDelete` | `refundapplydeletelogic.go` | 转发 `TradeRpc.RefundApplyDelete` |

**API 已实现 36 / 总计 36。**

### 汇总

| 层 | 已实现 | 总计 | 比例 |
|----|-------|------|------|
| RPC logic | 37 | 37 | 100% |
| API logic | 36 | 36 | 100% |
| **合计** | **73** | **73** | **100%** |

**前置阻塞项已解除**：`apps/trade/rpc/internal/model/` 下 4 个自定义 Model（cart/order/order_detail/refund_apply）已补齐分页、聚合、按 user_id 查询、状态更新、双写等方法，编译通过。

---

## 1. 购物车管理（已实现）

**核心规则**：购物车按 `user_id` 隔离，加购时写入课程快照（name/cover/price），避免每次列表回查 course 服务。

| 规则 | 依据 | 实际行为 |
|------|------|---------|
| 按用户隔离 | `cart.user_id` | 列表/删除均带 `user_id` 条件，禁止跨用户 |
| 课程快照冗余 | `cart.cover_url` / `course_name` / `price` | 加购时经 `CourseRpc.CourseSimpleInfoList` 写入，列表不再回查 |
| 幂等去重 | `CartAdd` | 同 userId+courseId 已存在仍返回成功，不重复插入 |
| 批量删除 | `CartBatchDeleteRequest.ids` | 支持一次多 id |

> ⚠️ **偏差**：设计意图要求列表展示时比对 course 服务当前售价、标记 `expired`（下架/涨价）。当前 `toCartVO` 中 `NowPrice` 恒等于加购快照 `price`、`Expired` 恒为 `false`，**未做现价比对**（见缺口 #8）。

---

## 2. 订单下单（已实现）

**核心规则**：下单经 `CourseRpc` 取价，雪花 ID 落库，一单多课拆明细。

| 规则 | 依据 | 实际行为 |
|------|------|---------|
| 两段式入口 | `OrderPrePlace` + `OrderPlace` | PrePlace 仅试算返回总价/课程列表；Place 真正生成订单与明细 |
| 金额 | `order.total_amount` / `real_amount` / `discount_amount` | 当前 `real_amount = total_amount`、`discount_amount=0`（优惠券未接入） |
| 一单多课 | `order`(1) → `order_detail`(N) | 每课一条明细，独立 `real_pay_amount` 与 `refund_status` |
| 免费课直通 | `OrderFreeCourse` | 金额 0，直接置 `status=2 已支付`（非设计意图的"已报名 5"） |
| 支付超时 | `PayOutTime = now+15min` | 下单下发支付截止时间（关单动作依赖支付回调，见缺口 #9） |

```
流程（OrderPlace）— 实际实现:
  1. 从 JWT 取 userId，校验 courseIds 非空
  2. fetchCourseMap 批量取课程快照
  3. 累加 total_amount（= real_amount，discount=0）
  4. 雪花 nextID() 生成 order.id，Insert order(status=1)
  5. 每 course_id 生成 detail.id，Insert order_detail(status=1, real_pay_amount=price)
  6. 返回 OrderId / PayAmount(total) / Status=1 / PayOutTime=now+15min
```

> ⚠️ **偏差**：设计意图的"预下单生成 order_id 并回传做幂等"未落地。`OrderPrePlace` 不生成/落库 order_id（返回 `OrderId:0`），`OrderPlace` 自行生成雪花 ID。幂等仅依赖购物车去重，**下单本身无 order_id 幂等保护**（见缺口 #2）。

---

## 3. 订单状态流转（已实现）

**核心规则**：`order.status` 与 `order_detail.status` 同步流转，状态枚举与 `common.go` 一致。

| 状态 | 值 | 触发动作 | 实现位置 |
|------|----|---------|----------|
| 待支付 | 1 | 下单成功 | `OrderPlace` / `OrderFreeCourse`(否) 置 1 |
| 已支付 | 2 | 支付回调（未落地） | `PayApply` 不置状态；回调在 pay 侧 |
| 已关闭 | 3 | 用户取消 / 超时 | `OrderCancel` 置 3 |
| 已完成 | 4 | 支付后 30 天 | 无写入路径（见缺口 #9） |
| 已报名 | 5 | 课程开通 | 无写入路径 |
| 申请退款 | 6 | 提交退款申请 | 无写入路径（退款审批未置 order=6） |

| 规则 | 实际行为 |
|------|---------|
| 取消仅限待支付 | `OrderCancel` 仅 status=1 时允许，否则拒绝 |
| 删除为逻辑删除 | `OrderDelete` 置 `deleted=1`，不物理删除 |
| 明细状态无 6 | 退款走 `order_detail.refund_status` 独立字段（1~6） |
| 进度节点 | `buildOrderProgressNodes` 按 order.status 组装 提交订单→支付成功/关闭→已完成/退款中 |

> ⚠️ **偏差**：`order.status` 的 2/4/5/6 在 trade 内**无写入路径**（支付回调、关单定时任务、退款联动均未落地），当前仅 1(下单)、3(取消) 由 trade 自身驱动（见缺口 #9）。

---

## 4. 支付（已实现 — 委托 PayRpc）

**核心规则**：trade 不落地支付流水，全部经 `PayRpc payclient.Pay` 代理 pay 服务（CourseRpc 之外已装配 `PayRpc`）。

| 规则 | 实际行为 |
|------|---------|
| 支付发起 | `PayApply` → `PayRpc.ApplyPayOrder{BizUserId:order.UserId, BizOrderNo:order.Id, Amount:order.TotalAmount, PayChannelCode, PayType:4(native 扫码)}` → 返回 `QrUrl` |
| 支付结果 | `PayResultQuery` → `PayRpc.QueryPayResult` → `PayResultDTO` |
| 渠道管理 | `PayChannel*` 4 方法全部委托 `PayRpc`（无本地 `pay_channel` 表） |
| 渠道展示 | `PayChannels`/`PayChannelList` → `PayRpc.ListPayChannels` 映射 VO，按 `channel_priority` 排序由客户端处理 |

> ⚠️ **偏差**：设计意图要求 `PayApply` 成功后回写 `order.pay_channel` / `order.pay_order_no`，并取 `real_amount` 作为支付金额。实际 `PayApply` **不回写 order**，且金额取自 `TotalAmount`（见缺口 #3）。

---

## 5. 支付成功事件（MQ）— 未发射

**设计意图**：支付/退款成功后经 RabbitMQ 通知 learning 服务开通/撤销课程。`servicecontext.go` 已装配 `MQProducer`（初始化失败仅 `logx.Errorf`，保持 nil 不阻塞启动）。

**实际行为**：经全仓 Grep 确认，`apps/trade/**/logic/*.go` 中**没有任何一处调用 `MQProducer.Publish`**。trade 不主动发射 `order.pay` / `order.refund` 事件，因此 learning 的 MQ 消费端**无生产者数据来源**（与 learning/business-rules.md 缺口呼应）。

---

## 6. 退款申请审批流（已实现 — 部分缺口）

**核心规则**：退款分「业务审批」（`refund_apply` 表）与「渠道退款」（pay 服务）两阶段，`refund_apply.status` 与 `order_detail.refund_status` 双写同步。

| 阶段 | 方法 | `refund_apply.status` | 实际行为 |
|------|------|----------------------|---------|
| 学员提交 | `RefundApplyCreate` | 1 待审批 | Insert(status=1) + `UpdateRefundStatus(detail,1)` 双写 |
| 学员撤回 | `RefundApplyCancel` | 2 取消退款 | status→2 + 同步 detail refund_status=2 |
| 管理员同意 | `RefundApplyApprove`(approve_type=1) | 3 同意退款 | `UpdateApprove`(approver 硬编码 0) + `UpdateRefundStatus(detail,3)` |
| 管理员拒绝 | `RefundApplyApprove`(approve_type=2) | 4 拒绝退款 | 同上，status→4 |
| 渠道退款 | `RefundApply`(独立 RPC) | — | 委托 `PayRpc.ApplyRefund`，返回 `RefundResultDTO` |
| 渠道结果 | `RefundResultQuery` | — | 委托 `PayRpc.QueryRefundResult` |

| 规则 | 实际行为 |
|------|---------|
| 申请粒度为明细 | `refund_apply.order_detail_id`，一单多课可单独退某门 |
| 可退判定 | `toOrderDetailItemVO.can_refund`：`status∈{2,4,5}` 且 `refund_status` 为空/0/1 |
| 状态双写 | `refund_apply.status` ↔ `order_detail.refund_status` 枚举一致（1~6）同步更新 |
| 待办取件 | `RefundApplyNext`：`FindNextPending` 取 status=1 下一条 |
| 审批留痕 | approver / approve_opinion / remark / approve_time 落库（approver 当前恒 0） |

```
流程（RefundApplyApprove，approve_type=1）— 实际实现:
  1. FindOne refund_apply（校验存在）
  2. status → 3（同意）/ 4（拒绝）
  3. UpdateApprove(id, status, approver=0, opinion, remark, approveTime=now, 0, 0)
  4. UpdateRefundStatus(detail, 同 status)   ← 仅到此为止
```

> ⚠️ **重大偏差**：设计意图中 `RefundApplyApprove` 是「业务审批」与「渠道退款」的衔接点，应在置 3 后调用 `PayRpc.ApplyRefund`、置 `order.status=6`、发布 `order.refund` MQ 事件。**实际 `RefundApplyApprove` 不触发任何 PayRpc 调用、不动 order 状态、不发射 MQ**。渠道退款只能由独立的 `RefundApply` RPC 另行触发（或由外部编排），审批→渠道退款的链路在 trade 内是断开的（见缺口 #4）。

---

## 状态说明

### 订单状态（`order.status` / `order_detail.status`）

| 值 | 含义 | 备注 |
|----|------|------|
| 1 | 待支付 | 默认值（下单写入） |
| 2 | 已支付 | 当前仅 `OrderFreeCourse` 写入；支付回调未落地 |
| 3 | 已关闭 | 取消或超时（`OrderCancel` 写入） |
| 4 | 已完成 | 支付后 30 天 — **无写入路径** |
| 5 | 已报名 | 课程已开通 — **无写入路径** |
| 6 | 已申请退款 | 仅 `order` 表定义 — **无写入路径** |

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

---

## 已知缺口汇总（相对设计意图的偏差）

1. **优惠券未接入**：`OrderPrePlace` / `OrderPlace` 折扣列表恒空，`real_amount` 恒等于 `total_amount`。`promotion` 服务已实现，但 trade 未装配 `PromotionRpc` 客户端。
2. **两段式下单未真正预留**：`OrderPrePlace` 不生成/落库 order_id（返回 `OrderId:0`），`OrderPlace` 自行生成雪花 ID；设计意图的"预下单 order_id 回传做幂等"未落地，下单本身无 order_id 幂等保护。
3. **PayApply 金额与回写缺失**：支付金额取自 `TotalAmount`（非 `real_amount`），且成功后未回写 `order.pay_channel` / `order.pay_order_no`。
4. **退款审批不触发渠道退款（重大）**：`RefundApplyApprove` 仅置 `refund_apply.status=3/4` 与 `order_detail.refund_status`，**未调用 `PayRpc.ApplyRefund`、未置 `order.status=6`、未发射 MQ**。渠道退款须由独立 `RefundApply` RPC 另行触发。
5. **审批人缺失**：`RefundApplyApprove.UpdateApprove` 的 `approver` 硬编码为 `0`（与 promotion 操作人恒 0 同类问题）。
6. **MQ 事件无生产者**：trade 全部 logic 均未调用 `MQProducer.Publish`，`order.pay` / `order.refund` 事件无发射点，learning 的 MQ 消费端无数据来源。
7. **OrderFreeCourse 状态偏差**：直接置 `status=2 已支付`，非设计意图的"已报名 5"。
8. **购物车现价比对未做**：`toCartVO` 中 `NowPrice` 恒等于加购快照价、`Expired` 恒 `false`，未与 course 服务当前价比对/标记下架。
9. **支付回调→后续动作缺失**：订单"已支付/已完成/已报名"状态流转、课程开通（learning 授权）、课程有效期（`valid_duration`/`course_expire_time`）计时等，在 trade 内均无实现入口，依赖 pay 服务回调链路——而该链路在 pay 服务侧同样未落地（见 pay 缺口）。

---

## 相关文档

- [实现进度与已知缺口](../06-status/implementation-status.md)
- [架构全览](../00-architecture/overview.md)
- [learning 业务规则](../learning/business-rules.md)（MQ 消费端缺口呼应）
