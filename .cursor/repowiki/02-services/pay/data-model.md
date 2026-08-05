> 版本：v1.0 | 更新：2026-08-05 | 来源：`sql/ddl/tj_pay.sql`

---

# Pay Data Model

## 模型文件清单

| 自定义 Model 文件 | 对应表 | 扩展方法 |
|------------------|--------|---------|
| `paychannelmodel.go` | `pay_channel` | FindAllEnabled, FindByCode, PageList |
| `payordermodel.go` | `pay_order` | MarkToPaying, MarkToSuccess, MarkToClosed, IncrNotifyTimes, SetNotifyStatus |
| `refundordermodel.go` | `refund_order` | FindOneByBizRefundOrderNo, FindOneByRefundOrderNo, FindListByBizOrderNo, MarkToProcessing, MarkToSuccess, MarkToFailed, SetNotifyStatus, IncrNotifyFailedTimes |

所有 `*_gen.go` 文件均由 goctl 生成，**禁止修改**。扩展方法统一放在同名自定义 `.go` 文件中（如 `payordermodel.go` 扩展 `PayOrderModel` 接口）。

---

## 表清单与字段说明

### 1. `pay_channel` — 支付渠道表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 支付渠道 ID，自增 | PK |
| `name` | varchar(50) | 支付渠道名称 | - |
| `channel_code` | varchar(30) | 支付渠道编码，用于获取支付实现 | 应用层唯一（无 DDL 约束） |
| `channel_priority` | int | 渠道优先级，数字越小优先级越高 | 排序字段 |
| `channel_icon` | varchar(255) | 渠道图标 | - |
| `status` | int | 支付渠道状态，1-使用中，2-停用，默认 1 | 过滤条件 |
| `creater` | bigint | 创建人 | - |
| `updater` | bigint | 更新人 | - |
| `create_time` | datetime | 创建时间 | 自动填充 |
| `update_time` | datetime | 更新时间 | 自动更新 |

> **注意**：`channel_code` 的唯一性**仅由应用层** `AddPayChannel` 中的 `FindByCode` 前置查询保证，DDL 未建唯一索引，高并发下存在重复插入风险。

---

### 2. `pay_order` — 支付订单表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 主键，自增（AUTO_INCREMENT=1585168008500764674） | PK |
| `biz_order_no` | bigint | 业务订单号 | 唯一索引 `biz_order_no` |
| `pay_order_no` | bigint | 支付单号，默认 0 | 唯一索引 `pay_order_no` |
| `biz_user_id` | bigint | 支付用户 ID | - |
| `pay_channel_code` | varchar(30) | 支付渠道编码，默认 `'0'` | - |
| `amount` | int | 支付金额，**单位分** | - |
| `pay_type` | tinyint | 支付类型，1-h5, 2-小程序, 3-公众号, 4-扫码，默认 4 | - |
| `status` | tinyint | 支付状态，0-待提交, 1-待支付, 2-支付超时或取消, 3-支付成功，默认 0 | 状态机核心 |
| `expand_json` | varchar(1024) | 拓展字段，用于传递不同渠道单独处理的字段，默认空串 | - |
| `notify_url` | varchar(255) | 业务端回调接口，可空，默认空串 | - |
| `notify_times` | int | 业务端回调次数，默认 0 | - |
| `notify_status` | int | 回调状态，0-待回调, 1-回调成功, 2-回调失败，默认 0 | - |
| `result_code` | varchar(20) | 第三方返回业务码，可空 | - |
| `result_msg` | varchar(50) | 第三方返回提示信息，可空 | - |
| `pay_success_time` | datetime | 支付成功时间，可空 | 成功时回填 |
| `pay_over_time` | datetime | 支付超时时间，非空 | 超时扫描依据 |
| `qr_code_url` | varchar(255) | 支付二维码链接，可空 | - |
| `create_time` | datetime | 创建时间 | 自动填充 |
| `update_time` | datetime | 更新时间 | 自动更新 |
| `creater` | bigint | 创建人，默认 0 | - |
| `updater` | bigint | 更新人，默认 0 | - |
| `deleted` | bit(1) | 逻辑删除，默认 b'0' | - |

> **幂等基石**：`biz_order_no` 与 `pay_order_no` 双唯一索引，是 `ApplyPayOrder` 幂等与回调定位的数据库级保障。
>
> **注意**：`amount` 为 `int`（约 21 亿分 ≈ 2147 万元上限），而 proto 中声明为 `int64`，存在类型宽度差异。

---

### 3. `refund_order` — 退款订单表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 主键，自增（AUTO_INCREMENT=1588164197026390019） | PK |
| `biz_order_no` | bigint | 业务端已支付的订单 ID | - |
| `biz_refund_order_no` | bigint | 业务端要退款的订单 ID，也就是子订单 ID | 普通索引 `index_biz_order_id` |
| `pay_order_no` | bigint | 付款时传入的支付单号 | - |
| `refund_order_no` | bigint | 退款单号，每次退款的唯一标识 | 普通索引 `index_refund_order_id` |
| `refund_amount` | int | 本次退款金额，**单位分** | - |
| `total_amount` | int | 总金额，单位分（快照自原支付单 `amount`） | - |
| `is_split` | bit(1) | 是否是拆单退款，默认 b'0' | - |
| `pay_channel_code` | varchar(30) | 支付渠道编码，默认 `'0'`（快照自原支付单） | - |
| `result_code` | varchar(64) | 第三方交易编码，可空 | - |
| `result_msg` | varchar(255) | 第三方交易信息，可空 | - |
| `status` | int | 退款状态，0-未提交, 1-退款中, 2-退款失败, 3-退款成功，默认 0 | 状态机核心 |
| `refund_channel` | varchar(255) | 退款渠道，可空 | 成功回调时写入 |
| `notify_failed_times` | int | 业务端退款通知失败次数，默认 0 | - |
| `notify_status` | int | 退款接口通知状态，0-待通知, 1-通知成功, 2-通知中, 3-通知失败，默认 0 | - |
| `create_time` | datetime | 退款单据创建时间 | 普通索引 `index_create_time` |
| `update_time` | datetime | 退款单据修改时间 | 自动更新 |
| `creater` | bigint | 单据创建人，一般手动对账产生的单据才有值，默认 0 | - |
| `updater` | bigint | 单据修改人，一般手动对账产生的单据才有值，默认 0 | - |
| `deleted` | bit(1) | 逻辑删除，默认 b'0' | 过滤条件 |

> **注意**：`biz_refund_order_no` 与 `refund_order_no` 均为**普通索引而非唯一索引**，因此 `FindOneByBizRefundOrderNo` 需显式 `order by id desc limit 1` 兜底；`ApplyRefund` 的幂等仅靠应用层前置查询，无数据库级约束。
>
> **注意**：`notify_status` 的语义与 `pay_order.notify_status` **不一致** —— 退款单的 2 表示「通知中」、3 表示「通知失败」，支付单的 2 表示「回调失败」。两者在 `common.go` 中定义为两组独立常量。

---

## 状态常量映射（`internal/logic/common.go`）

| 常量组 | 常量名 | 值 |
|--------|--------|-----|
| 渠道状态 | `PayChannelStatusEnabled` / `PayChannelStatusDisabled` | 1 / 2 |
| 支付单状态 | `PayOrderStatusPending` / `Paying` / `Closed` / `Success` | 0 / 1 / 2 / 3 |
| 退款单状态 | `RefundStatusInit` / `Processing` / `Failed` / `Success` | 0 / 1 / 2 / 3 |
| 支付回调状态 | `NotifyStatusPending` / `OK` / `Fail` | 0 / 1 / 2 |
| 退款通知状态 | `RefundNotifyStatusPending` / `Success` / `Processing` / `Failed` | 0 / 1 / 2 / 3 |
| 渠道类型 | `PayTypeH5` / `PayTypeMini` / `PayTypeMp` / `PayTypeNative` | 1 / 2 / 3 / 4 |

---

## 关系图

```
pay_channel (1) ──(channel_code)── (N) pay_order
                                       │
                                       │ (pay_order_no / biz_order_no)
                                       ↓
pay_order (1) ────────────────────── (N) refund_order

refund_order.total_amount      ← 快照自 pay_order.amount
refund_order.pay_channel_code  ← 快照自 pay_order.pay_channel_code

pay_order.biz_order_no         → trade 域业务订单（跨库引用）
pay_order.biz_user_id          → user 域用户（跨库引用）
```

三张表之间**无数据库外键**，全部通过业务单号在应用层关联。

---

## bit(1) 字段的 Go 映射

`is_split` 与 `deleted` 在 DDL 中为 `bit(1)`，goctl 映射为 `int64`：

| 位置 | 处理方式 |
|------|---------|
| 写入 | `ApplyRefund` 中显式赋 `IsSplit: 0` |
| 读出 | `toRefundResp` 中 `IsSplit: m.IsSplit == 1` 转为 proto 的 `bool` |
| 过滤 | `FindListByBizOrderNo` 拼接 `and deleted = 0` |

> `pay_order.deleted` 虽有该列，但当前所有查询路径（含 goctl 生成的 `FindOneByBizOrderNo` / `FindOneByPayOrderNo`）**均未加 `deleted = 0` 过滤**，仅 `refund_order` 的 `FindListByBizOrderNo` 做了过滤。

---

## 模型扩展模式

所有 `*_gen.go` 由 goctl 自动生成，不可编辑。扩展方法统一放在同名的自定义 `.go` 文件中：

```
payordermodel_gen.go      ← goctl 生成，只读（含 FindOneByBizOrderNo / FindOneByPayOrderNo，带缓存）
payordermodel.go          ← 手写扩展 MarkToSuccess / MarkToClosed / IncrNotifyTimes 等状态流转方法
```

当前项目自定义 Model 模式：
- `paychannelmodel.go` — FindAllEnabled, FindByCode, PageList
- `payordermodel.go` — MarkToPaying, MarkToSuccess, MarkToClosed, IncrNotifyTimes, SetNotifyStatus
- `refundordermodel.go` — FindOneByBizRefundOrderNo, FindOneByRefundOrderNo, FindListByBizOrderNo, MarkToProcessing, MarkToSuccess, MarkToFailed, SetNotifyStatus, IncrNotifyFailedTimes

所有手写扩展方法一律使用 `ExecNoCacheCtx` / `QueryRowNoCacheCtx` / `QueryRowsNoCacheCtx`，**绕过 goctl 缓存层**——状态流转后需注意 `FindOneByBizOrderNo` 等带缓存查询可能读到旧状态。

`vars.go` 定义 `ErrNotFound = sqlx.ErrNotFound`；logic 层的 `isNotFound(err)` 同时判定 `sql.ErrNoRows` 与 `model.ErrNotFound` 两种情况。
