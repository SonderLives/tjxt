> 版本：v1.0 | 更新：2026-08-05 | 来源：`sql/ddl/tj_trade.sql`

---

# Trade Data Model

## 模型文件清单

| 自定义 Model 文件 | 对应表 | 扩展方法 |
|------------------|--------|---------|
| `cartmodel.go` | `cart` | （无，仅 goctl 空壳包装） |
| `ordermodel.go` | `order` | （无，仅 goctl 空壳包装） |
| `orderdetailmodel.go` | `order_detail` | （无，仅 goctl 空壳包装） |
| `refundapplymodel.go` | `refund_apply` | （无，仅 goctl 空壳包装） |

所有 `*_gen.go` 文件均由 goctl 生成，**禁止修改**。扩展方法统一放在同名自定义 `.go` 文件中（如 `ordermodel.go` 扩展 `OrderModel` 接口）。

> **当前状态**：4 个自定义 Model 文件的接口体内均只有 `cartModel` / `orderModel` / `orderDetailModel` / `refundApplyModel` 一行嵌入，**尚未声明任何扩展方法**，可用能力仅限 `Insert` / `FindOne` / `Update` / `Delete` 四个 goctl 基础方法。

`vars.go` 定义 `ErrNotFound = sqlx.ErrNotFound`。

---

## 表清单与字段说明

### 1. `cart` — 购物车条目表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 购物车条目 id | PK |
| `user_id` | bigint | 用户 id | - |
| `course_id` | bigint | 课程 id | - |
| `cover_url` | varchar(255) | 课程封面路径 | - |
| `course_name` | varchar(255) | 课程名称 | - |
| `price` | int | 单价（分） | - |
| `create_time` | datetime | 创建时间，默认 CURRENT_TIMESTAMP | 自动填充 |
| `update_time` | datetime | 更新时间，ON UPDATE CURRENT_TIMESTAMP | 自动更新 |

> 表注释：购物车条目信息，也就是购物车中的课程。`cover_url` / `course_name` / `price` 为下单时点的课程快照冗余字段。

---

### 2. `order` — 订单表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 订单 id | PK |
| `user_id` | bigint | 用户 id | - |
| `pay_order_no` | bigint | 交易流水支付单号，可空 | - |
| `status` | tinyint | 订单状态，默认 1 | - |
| `message` | varchar(255) | 状态备注，默认 `''` | - |
| `total_amount` | int | 订单总金额，单位分 | - |
| `real_amount` | int | 实付金额，单位分 | - |
| `discount_amount` | int | 优惠金额，单位分，默认 0 | - |
| `pay_channel` | varchar(20) | 支付渠道，默认 `''` | - |
| `coupon_ids` | json | 优惠券 id 列表 | - |
| `create_time` | datetime | 创建订单时间 | 自动填充 |
| `pay_time` | datetime | 支付时间，可空 | - |
| `close_time` | datetime | 订单关闭时间，可空 | - |
| `finish_time` | datetime | 订单完成时间，支付后 30 天 | - |
| `refund_time` | datetime | 申请退款时间，可空 | - |
| `update_time` | datetime | 更新时间 | 自动更新 |
| `creater` | bigint | 创建人 | - |
| `updater` | bigint | 更新人 | - |
| `deleted` | tinyint | 逻辑删除，默认 0 | 过滤条件 |

> **注意**：表名 `order` 是 MySQL 保留字，SQL 中必须用反引号包裹。

---

### 3. `order_detail` — 订单明细表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 订单明细 id | PK |
| `order_id` | bigint | 订单 id | `idx_order` |
| `user_id` | bigint | 用户 id | `idx_user_course` |
| `course_id` | bigint | 课程 id | `idx_user_course` |
| `price` | int | 课程价格 | - |
| `name` | varchar(128) | 课程名称 | - |
| `cover_url` | varchar(255) | 封面地址 | - |
| `valid_duration` | int | 课程学习有效期，单位月，空则永久有效 | - |
| `course_expire_time` | datetime | 课程学习过期时间，支付成功开始计时 | `idx_course_expire_time` |
| `discount_amount` | int | 折扣金额，默认 0 | - |
| `real_pay_amount` | int | 实付金额 | - |
| `status` | tinyint | 明细状态 | - |
| `refund_status` | tinyint | 退款状态，可空 | - |
| `pay_channel` | varchar(50) | 支付渠道名称，默认 `''` | `idx_pay_channel` |
| `create_time` | datetime | 创建时间 | 自动填充 |
| `update_time` | datetime | 更新时间 | 自动更新 |
| `creater` | bigint | 创建人 | - |
| `updater` | bigint | 更新人 | - |

**索引清单**：

| 索引名 | 字段 | 用途 |
|--------|------|------|
| `idx_order` | `order_id` | 按订单查明细 |
| `idx_user_course` | `user_id`, `course_id` | 校验用户是否已购课程 |
| `idx_course_expire_time` | `course_expire_time` | 扫描到期课程 |
| `idx_pay_channel` | `pay_channel` | 按支付渠道统计 |

---

### 4. `refund_apply` — 退款申请表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 退款 id | PK |
| `order_detail_id` | bigint | 订单明细 id | - |
| `order_id` | bigint | 订单 id | - |
| `pay_order_no` | bigint | 流水支付单号，可空 | - |
| `refund_order_no` | bigint | 流水退款单号，可空 | - |
| `user_id` | bigint | 订单所属用户 id，默认 0 | - |
| `refund_amount` | bigint | 退款金额 | - |
| `status` | int | 退款状态，默认 1 | - |
| `refund_reason` | varchar(255) | 申请退款原因 | - |
| `message` | varchar(255) | 退款状态描述 | - |
| `approver` | bigint | 审批人 id，可空 | - |
| `approve_opinion` | varchar(255) | 审批意见，可空 | - |
| `remark` | varchar(255) | 审批备注，可空 | - |
| `failed_reason` | varchar(255) | 退款失败原因，可空 | - |
| `question_desc` | varchar(255) | 退款问题说明，可空 | - |
| `refund_channel` | varchar(50) | 退款渠道，可空 | - |
| `create_time` | datetime | 创建退款申请时间 | 自动填充 |
| `approve_time` | datetime | 审批时间，可空 | - |
| `finish_time` | datetime | 退款完成时间（成功或失败），可空 | - |
| `update_time` | datetime | 更新时间 | 自动更新 |
| `creater` | bigint | 创建人 | - |
| `updater` | bigint | 更新人 | - |

---

### 5. `undo_log` — 分布式事务回滚日志表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `branch_id` | bigint | branch transaction id | `ux_undo_log` |
| `xid` | varchar(100) | global transaction id | `ux_undo_log`（唯一） |
| `context` | varchar(128) | undo_log context，如序列化方式 | - |
| `rollback_info` | longblob | rollback info | - |
| `log_status` | int | 0=normal status, 1=defense status | - |
| `log_created` | datetime(6) | create datetime | - |
| `log_modified` | datetime(6) | modify datetime | - |

> 表注释：AT transaction mode undo table。属于 Seata AT 模式的基础设施表，字符集为 `utf8`（区别于其余表的 `utf8mb4`），**当前 Go 侧未生成对应 model，也未接入任何分布式事务框架**，属遗留自 Java 版本的结构。

---

## 关系图

```
cart (N) ─── course_id ──→ course 域（跨库，无外键）

order (1) ─── (N) order_detail
                      │
                      └── (1) refund_apply ── order_detail_id / order_id

order.pay_order_no      ──→ tj_pay.pay_order（跨库，无外键）
refund_apply.refund_order_no ──→ tj_pay.refund_order（跨库，无外键）

undo_log（独立，Seata AT 基础设施表）
```

---

## 状态枚举

### `order.status` / `order_detail.status` — 订单状态

| 值 | 含义 |
|----|------|
| 1 | 待支付 |
| 2 | 已支付 |
| 3 | 已关闭 |
| 4 | 已完成 |
| 5 | 已报名 |
| 6 | 已申请退款（仅 `order` 表定义，`order_detail` 注释止于 5） |

### `order_detail.refund_status` / `refund_apply.status` — 退款状态

| 值 | 含义 |
|----|------|
| 1 | 待审批 |
| 2 | 取消退款 |
| 3 | 同意退款 |
| 4 | 拒绝退款 |
| 5 | 退款成功 |
| 6 | 退款失败 |

---

## 已知结构性缺口

| 缺口 | 说明 |
|------|------|
| **支付渠道/流水表不在本库** | proto 提供 `PayChannelAdd` / `PayChannelList` / `PayChannelGet` / `PayChannelDelete` / `PayChannels` 等 5 个支付渠道方法，以及 `PayApply` / `PayResultQuery` / `RefundApply` / `RefundResultQuery`，但 `sql/ddl/tj_trade.sql` **不含 `pay_channel` / `pay_order` / `refund_order` 表**。这三张表定义在 `sql/ddl/tj_pay.sql`（pay 服务库），trade 通过 `PayRpc payclient.Pay` 代理调用 pay 服务，本库不落地支付流水。 |
| **无自定义 Model 扩展方法** | 4 个 model 均为空壳，proto 中的分页、聚合统计、按 user_id 列表查询等能力缺少 SQL 支撑。 |
| **`undo_log` 无对应 model** | 表存在但 Go 侧未生成 model，也未接入分布式事务框架。 |

---

## 模型扩展模式

所有 `*_gen.go` 由 goctl 自动生成，不可编辑。扩展方法统一放在同名的自定义 `.go` 文件中：

```
ordermodel_gen.go          ← goctl 生成，只读
ordermodel.go              ← 手写扩展位（当前为空壳，待补）
```

当前 trade 域自定义 Model 现状：
- `cartmodel.go` — 空壳，`CartModel interface { cartModel }`
- `ordermodel.go` — 空壳，`OrderModel interface { orderModel }`
- `orderdetailmodel.go` — 空壳，`OrderDetailModel interface { orderDetailModel }`
- `refundapplymodel.go` — 空壳，`RefundApplyModel interface { refundApplyModel }`

对照 learning 域 `learninglessonmodel.go`（已声明 11 个扩展方法）与 auth 域 `rolemodel.go`，trade 域的 Model 扩展是后续实现 37 个 RPC 方法的前置依赖。
