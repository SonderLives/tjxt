> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/pay/rpc/internal/logic/*.go`

---

# Pay Business Rules

## 1. 支付渠道管理

**核心规则**：渠道编码 `channel_code` 是支付实现的路由键，一旦创建**永不可改**，以避免与历史订单对不上。

| 规则 | 说明 |
|------|------|
| 新增不带 id | `AddPayChannel` 中 `id != 0` 直接 `BadRequest("新增渠道不应携带 id")` |
| 名称与编码必填 | `name` 或 `channel_code` 为空 → `BadRequest` |
| 编码唯一性 | 新增前 `FindByCode` 探测，查到即 `Conflict("渠道编码已存在")` |
| 编码不可修改 | 更新时 `channel_code` 非空且与库中值不同 → `BadRequest("渠道编码不允许修改")` |
| 新增即启用 | `buildPayChannel` 固定写入 `Status: PayChannelStatusEnabled`（1），无法在新增时指定为停用 |
| 更新为非零覆盖 | `name` / `channel_priority` / `channel_icon` 仅在非零时覆盖原值 |
| 状态白名单 | `UpdatePayChannelStatus` 仅接受 1（使用中）与 2（停用），其它 → `BadRequest("渠道状态非法")` |
| ID 由数据库生成 | 渠道表为 `AUTO_INCREMENT`，用 `res.LastInsertId()` 取回，**不走雪花 ID** |

```
流程（AddPayChannel）:
  1. id != 0 → BadRequest
  2. name / channel_code 为空 → BadRequest
  3. FindByCode(channel_code)
       err == nil        → Conflict（编码已存在）
       !isNotFound(err)  → Internal
  4. buildPayChannel(in)，status 固定为 1
  5. Insert → LastInsertId 作为返回 id
```

> **⚠️ 并发风险**：第 3、5 步之间无锁且 `channel_code` 无数据库唯一索引，并发新增同一编码会插入重复行。

**查询规则**：

| 方法 | 排序 / 过滤 |
|------|-----------|
| `ListPayChannels` | `FindAllEnabled` 硬编码 `where status = 1`，按 `channel_priority asc` |
| `QueryPayChannelByCode` | `channel_code` 为空 → `BadRequest`；未查到 → `NotFound("支付渠道不存在")` |
| `PageQueryPayChannels` | `name` 模糊、`channel_code` 精确、`status > 0` 才过滤；统一按 `channel_priority asc` |

> `PageList` 的状态过滤条件是 `status > 0` 而非 `>= 0`，因此**无法筛选 status=0 的渠道**（DDL 中 status 的合法值仅 1/2，故实际无影响）。

---

## 2. 支付单申请（ApplyPayOrder）

**核心规则**：以 `biz_order_no` 为幂等键，同一业务订单永远只对应一张支付单。

| 规则 | 说明 |
|------|------|
| 参数校验 | `biz_order_no` / `biz_user_id` / `amount` 任一 `<= 0` → `BadRequest` |
| 渠道必填 | `pay_channel_code` 为空 → `BadRequest` |
| 渠道存在性 | `FindByCode` 未查到 → `NotFound("支付渠道不存在")` |
| 渠道可用性 | `channel.Status != 1` → `Conflict("支付渠道已停用")` |
| 超时缺省 | `pay_over_seconds <= 0` → 30 分钟（1800 秒） |
| 支付类型缺省 | `pay_type <= 0` → `PayTypeNative`（4，扫码） |
| 支付单号 | `idgen.NextID()` 雪花生成，与自增主键 `id` 相互独立 |
| 初始状态 | `status = PayOrderStatusPaying`（1，待支付），**跳过 0-待提交** |
| 二维码 | 当前为 `mockQrCodeUrl` 生成的占位链接 `tjxt://mock-pay?order_no=..&amount=..` |

**幂等三分支**（命中已有单时按状态分流）：

| 原单状态 | 行为 |
|---------|------|
| 1 待支付 | 直接返回原 `qr_code_url`，不新建单 |
| 3 支付成功 | `Conflict("订单已支付，请勿重复支付")` |
| 其它（0 待提交 / 2 已关闭） | `Conflict("订单已关闭，请重新下单")` |

```
流程（ApplyPayOrder）:
  1. 校验 biz_order_no / biz_user_id / amount > 0，pay_channel_code 非空
  2. FindByCode → 渠道存在且 status == 1
  3. FindOneByBizOrderNo(biz_order_no)
       命中 → 按状态三分支返回（幂等出口）
       ErrNotFound → 继续
       其它错误 → Internal
  4. overSeconds = pay_over_seconds 或 1800
     payOverTime = now + overSeconds
     payType     = pay_type 或 4
  5. poNo = idgen.NextID()
  6. Insert pay_order { status=1, notify_times=0, notify_status=0, qr_code_url=mock }
  7. 返回 qr_code_url
```

> **⚠️ 并发风险**：第 3 步查询与第 6 步插入之间无锁，但 `biz_order_no` 有**数据库唯一索引**兜底——并发重复下单会在 Insert 处报唯一键冲突并返回 `Internal`，而非产生重复单。
>
> **⚠️ Mock 实现**：`mockQrCodeUrl` 源码注释明确「真实项目里应该是调用微信/支付宝下单接口拿到的 code_url/prepay_id」。当前**未对接任何真实支付网关**。
>
> **超时未落地**：`pay_over_time` 已写入，但仓库中**没有扫描超时单并自动关单的定时任务**，超时关单目前依赖外部调用 `ClosePayOrder` 或 `NotifyPayFailed`。

---

## 3. 支付单状态机

```
                 ApplyPayOrder
       ┌──────────────────────────┐
       │                          ↓
   (0 待提交)  ──MarkToPaying──→  (1 待支付)
                                   │
              NotifyPaySuccess     │     NotifyPayFailed / ClosePayOrder
                  ┌────────────────┴────────────────┐
                  ↓                                 ↓
            (3 支付成功)  ←── 终态             (2 已关闭)  ←── 终态
```

**状态流转守卫**（三个变更接口共用同一套判定模式）：

| 接口 | 遇 3 支付成功 | 遇 2 已关闭 | 遇 0/1 |
|------|-------------|------------|--------|
| `NotifyPaySuccess` | 返回空响应（**幂等**） | `Conflict("订单已关闭，无法标记支付成功")` | `MarkToSuccess` |
| `NotifyPayFailed` | `Conflict("订单已支付成功，不能再标记失败")` | 返回空响应（**幂等**） | `MarkToClosed` |
| `ClosePayOrder` | `Conflict("订单已支付成功，无法关单")` | 返回空响应（**幂等**） | `MarkToClosed` |

**幂等设计要点**：重复回调「已是目标状态」时**静默成功**，而「已是另一终态」时**报冲突**——这保证了第三方网关的重试不会失败，同时阻止了终态之间的非法翻转。

```
流程（NotifyPaySuccess）:
  1. pay_order_no <= 0 → BadRequest
  2. FindOneByPayOrderNo → 未找到则 NotFound("支付单不存在")
  3. switch status:
       3 成功 → return Empty{}         // 幂等，不重复落库
       2 关闭 → Conflict
       0/1    → 继续
  4. MarkToSuccess(id, result_code, result_msg)
       UPDATE status=3, result_code, result_msg,
              pay_success_time=now(), update_time=now()
```

| 关单来源 | 写入的 result_code / result_msg |
|---------|-------------------------------|
| `NotifyPayFailed` | 取自入参 `result_code` / `result_msg` |
| `ClosePayOrder` | 硬编码 `"MANUAL_CLOSE"` / `"业务端主动关单"` |

> **⚠️ 缓存失效缺口**：`MarkToSuccess` / `MarkToClosed` 使用 `ExecNoCacheCtx` 直改数据库，**未调用 `DelCacheCtx`**；而 `FindOneByBizOrderNo` / `FindOneByPayOrderNo` 是 goctl 生成的**带缓存**查询。状态变更后缓存中的旧记录不会被清理，`QueryPayResult` 存在读到过期状态的风险。
>
> **⚠️ 回调通知未实现**：`notify_url`、`notify_times`、`notify_status` 三个字段仅在建单时写入初值；模型层虽提供 `IncrNotifyTimes` / `SetNotifyStatus`，但**没有任何 logic 调用它们**。「支付成功后通知业务端」这一环尚未落地（源码注释：「→ MQ/HTTP 通知业务端」）。

---

## 4. 退款申请（ApplyRefund）

**核心规则**：以 `biz_refund_order_no` 为幂等键；退款前必须校验原支付单已成功，且**累计退款金额不得超过原支付金额**。

| 规则 | 说明 |
|------|------|
| 参数校验 | `biz_order_no` / `biz_refund_order_no` / `refund_amount` 任一 `<= 0` → `BadRequest` |
| 幂等出口 | `FindOneByBizRefundOrderNo` 命中则**直接返回原退款单**，不做任何状态判断 |
| 原单存在性 | `FindOneByBizOrderNo` 未查到 → `NotFound("原支付单不存在")` |
| 原单必须已支付 | `payOrder.Status != 3` → `Conflict("原支付单未支付成功，无法退款")` |
| 单次金额上限 | `refund_amount > payOrder.Amount` → `BadRequest("退款金额超过原支付金额")` |
| 累计金额上限 | 已退（成功 + 退款中）+ 本次 > 原金额 → `BadRequest("累计退款金额超过原支付金额")` |
| 快照字段 | `total_amount` / `pay_channel_code` / `pay_order_no` 全部快照自原支付单 |
| 拆单标记 | `is_split` 固定写 0，**当前不支持拆单退款标记** |
| 初始状态 | `status = RefundStatusProcessing`（1，退款中），跳过 0-未提交 |

**累计退款金额的计算口径**：

```
refunded = Σ r.refund_amount
           where r.biz_order_no = 入参 biz_order_no
             and r.deleted = 0
             and r.status ∈ { 3 退款成功, 1 退款中 }
```

> **关键设计**：`RefundStatusProcessing`（退款中）被计入已退金额，属于**悲观占额**——防止多笔并发退款在都还未终态时各自通过校验导致超额退款。已失败（status=2）的退款单则**不占额**，金额自动释放。

```
流程（ApplyRefund）:
  1. 校验三个入参 > 0
  2. FindOneByBizRefundOrderNo(biz_refund_order_no)
       命中 → 原样返回 { refund_order_no, biz_refund_order_no,
                        refund_amount, status, result_msg }   // 幂等出口
       非 NotFound 错误 → Internal
  3. FindOneByBizOrderNo(biz_order_no) → 不存在则 NotFound
  4. payOrder.Status != 3 → Conflict
  5. refund_amount > payOrder.Amount → BadRequest
  6. FindListByBizOrderNo → 累加 status ∈ {1,3} 的 refund_amount
     refunded + refund_amount > payOrder.Amount → BadRequest
  7. roNo = idgen.NextID()
  8. Insert refund_order { status=1 退款中, is_split=0, notify_status=0,
                           total_amount / pay_channel_code / pay_order_no 快照 }
  9. [MOCK] MarkToSuccess(ro.Id, "MOCK_OK", "mock 退款成功", "mock")
       失败仅 l.Errorf 记录日志，不中断
 10. 返回 { status: 3 退款成功, result_msg: "mock 退款成功" }
```

> **⚠️ 并发风险**：`biz_refund_order_no` **无数据库唯一索引**（仅普通索引 `index_biz_order_id`），第 2 步的幂等查询与第 8 步的插入之间无锁，并发重复申请会产生**两张退款单**。累计金额校验（第 6 步）同样无锁，属于典型的 check-then-act 竞态。
>
> **⚠️ Mock 实现**：第 9 步源码注释明确「真实生产：调用第三方退款 API，再根据结果 async 通过 NotifyRefundSuccess/Failed 更新；demo 中暂时直接 mock 成功」。因此当前 `ApplyRefund` **总是同步返回退款成功**，`RefundStatusProcessing` 状态在实际运行中几乎不会被观测到。
>
> **⚠️ 状态不一致隐患**：第 9 步的 `MarkToSuccess` 失败只记日志不回滚，此时数据库中退款单停留在「退款中」，但 RPC 已返回「退款成功」。
>
> 第 9 步前的 `_ = sql.ErrNoRows` 是一行无实际作用的占位语句（用于保留 `database/sql` 的 import）。

---

## 5. 退款单状态机

```
                  ApplyRefund
        ┌───────────────────────────┐
        │                           ↓
   (0 未提交) ──MarkToProcessing──→ (1 退款中)
                                     │
           NotifyRefundSuccess       │      NotifyRefundFailed
                 ┌───────────────────┴───────────────────┐
                 ↓                                       ↓
           (3 退款成功) ←── 终态                    (2 退款失败) ←── 终态
```

**状态流转守卫**：

| 接口 | 遇 3 退款成功 | 遇 2 退款失败 | 遇 0/1 |
|------|-------------|-------------|--------|
| `NotifyRefundSuccess` | 返回空响应（**幂等**） | `Conflict("退款单已标记失败，不允许改为成功")` | `MarkToSuccess` |
| `NotifyRefundFailed` | `Conflict("退款单已成功，不能改为失败")` | 返回空响应（**幂等**） | `MarkToFailed` |

与支付单完全对称：定位键为 `refund_order_no`（非业务退款单号），`<= 0` → `BadRequest`，未查到 → `NotFound("退款单不存在")`。

| 接口 | 写入字段 |
|------|---------|
| `MarkToSuccess` | `status=3`, `result_code`, `result_msg`, `refund_channel`, `update_time` |
| `MarkToFailed` | `status=2`, `result_code`, `result_msg`, `update_time` |

> `refund_channel` 经 `sql.NullString{Valid: refundChannel != ""}` 处理，空串写入 NULL 而非空字符串。
>
> **⚠️ 退款通知未实现**：`notify_status` / `notify_failed_times` 与模型层的 `SetNotifyStatus` / `IncrNotifyFailedTimes` 同样**无 logic 调用**，退款结果通知业务端的链路尚未落地。

---

## 6. 查询接口

| 方法 | 定位键 | 返回粒度 |
|------|-------|---------|
| `QueryPayResult` | `biz_order_no` | 轻量：`pay_order_no` / `biz_order_no` / `status` 三字段 |
| `QueryPayOrderByBizOrderNo` | `biz_order_no` | 完整：19 字段，含通知次数、结果码、各类时间 |
| `QueryRefundResult` | `biz_refund_order_no` | 轻量：单号 / 金额 / 状态 / `result_msg` |
| `QueryRefundByBizRefundNo` | `biz_refund_order_no` | 完整：17 字段，含拆单标记、退款渠道、通知状态 |

| 通用规则 | 说明 |
|---------|------|
| 入参校验 | 单号 `<= 0` → `BadRequest` |
| 不存在 | `isNotFound` → `NotFound("支付单不存在")` / `NotFound("退款单不存在")` |
| 时间格式化 | `formatTime` 零值返回空串，否则 `2006-01-02 15:04:05` |
| 可空时间 | `formatNullTime` 无效返回空串（用于 `pay_success_time`） |
| 可空字符串 | `formatNullString` 无效返回空串（用于 `qr_code_url` / `refund_channel`） |
| bit 转 bool | `IsSplit: m.IsSplit == 1` |

---

## 7. 分页与错误码

**分页归一化**（`normalizePage` → `page.Normalize`）：

| 入参情况 | 归一化结果 |
|---------|-----------|
| `page_no < 1` | 置为 1 |
| `page_size < 1` | 置为 10 |
| `page_size > 100` | 置为 100 |
| 总页数 | `page.CalcPages(total, limit)`，`total <= 0` 返回 0 |

**错误码约定**：

| 场景 | 错误构造 |
|------|---------|
| 参数非法 | `xerr.BadRequestf(...)` |
| 记录不存在 | `xerr.NotFound(...)` |
| 状态冲突 / 唯一性冲突 | `xerr.Conflict(...)` |
| 数据库异常 | `xerr.Wrapf(err, xerr.CodeInternal, ...)` |

> pay 服务统一使用 `xerr.Wrapf`（user 服务用的是 `xerr.Wrap`），两者在错误包装语义上等价，仅格式化能力不同。

---

## 状态说明

### 支付渠道状态 `status`

| 值 | 含义 | 影响 |
|----|------|------|
| 1 | 使用中 | 可被 `ListPayChannels` 列出，可用于下单 |
| 2 | 停用 | `ApplyPayOrder` 返回 `Conflict("支付渠道已停用")` |

### 支付单状态 `status`

| 值 | 含义 | 备注 |
|----|------|------|
| 0 | 待提交 | 常量已定义，但 `ApplyPayOrder` 直接建为 1，实际不产生该状态 |
| 1 | 待支付 | 建单初始态，可流转至 2 或 3 |
| 2 | 支付超时或取消 | 终态 |
| 3 | 支付成功 | 终态，退款的前置条件 |

### 退款单状态 `status`

| 值 | 含义 | 备注 |
|----|------|------|
| 0 | 未提交 | 常量已定义，但 `ApplyRefund` 直接建为 1，实际不产生该状态 |
| 1 | 退款中 | 建单初始态，**计入累计已退金额** |
| 2 | 退款失败 | 终态，不占用退款额度 |
| 3 | 退款成功 | 终态，计入累计已退金额 |

### 支付类型 `pay_type`

| 值 | 含义 |
|----|------|
| 1 | h5 |
| 2 | 小程序 |
| 3 | 公众号 |
| 4 | 扫码（缺省值） |

### 回调状态

| 值 | `pay_order.notify_status` | `refund_order.notify_status` |
|----|--------------------------|------------------------------|
| 0 | 待回调 | 待通知 |
| 1 | 回调成功 | 通知成功 |
| 2 | 回调失败 | 通知中 |
| 3 | （无） | 通知失败 |

> 两张表的 `notify_status` **语义不同**，读写时不可混用。
