> 版本：v1.0 | 更新：2026-08-05 | 来源：`sql/ddl/tj_promotion.sql`

---

# Promotion Data Model

## 模型文件清单

| 自定义 Model 文件 | 对应表 | 扩展方法 |
|------------------|--------|---------|
| `couponmodel.go` | `coupon` | FindList, FindPage, FindByIds, IncrIssueNum, DecrIssueNum, AddUsedNum, UpdateStatus, SoftDelete |
| `usercouponmodel.go` | `user_coupon` | FindByUserAndStatus, FindPageByUser, FindByIdsAndUser, CountByUserAndCoupon, UpdateStatusByIds |
| `couponcodemodel.go` | `coupon_code` | FindPageByCoupon, MarkUsed, BatchInsert |

所有 `*_gen.go` 文件均由 goctl 生成，**禁止修改**。扩展方法统一放在同名自定义 `.go` 文件中（如 `couponmodel.go` 扩展 `CouponModel` 接口）。

`consts.go` 集中定义了状态与枚举常量，`vars.go` 定义 `ErrNotFound = sqlx.ErrNotFound`。

---

## 表清单与字段说明

### 1. `coupon` — 优惠券表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 优惠券 id，自增主键 | PK |
| `name` | varchar(100) | 优惠券名称 | - |
| `discount_type` | varchar(20) | 优惠类型：reduce-满减 discount-折扣 no_threshold-无门槛 | - |
| `discount_value` | int | 优惠值：满减金额/折扣百分比，单位分 | - |
| `max_discount_amount` | int | 最高优惠金额，单位分 | - |
| `threshold_amount` | int | 使用门槛金额，单位分，0 表示无门槛 | - |
| `obtain_way` | varchar(20) | 获取方式：receive-领取 exchange-兑换码 assign-发放 | - |
| `specific` | tinyint(1) | 是否限定适用范围：0-不限 1-限定 | - |
| `scopes` | json | 适用范围 id 列表（课程/分类），specific=1 时有效 | - |
| `total_num` | int | 发行总量 | 库存判定 |
| `issue_num` | int | 已领取数量 | 库存判定 |
| `used_num` | int | 已使用数量 | - |
| `user_limit` | int | 每人限领数量，0-不限 | - |
| `issue_begin_time` | datetime | 发放开始时间 | 可空 |
| `issue_end_time` | datetime | 发放结束时间 | 可空 |
| `term_begin_time` | datetime | 使用开始时间 | 可空 |
| `term_end_time` | datetime | 使用结束时间 | 可空 |
| `term_days` | int | 有效期天数（相对领取日），与 term_begin/end 二选一 | - |
| `status` | varchar(20) | 状态：draft-草稿 issued-已发放 paused-暂停 ended-结束 | `idx_status` |
| `remark` | varchar(255) | 备注 | - |
| `create_time` | datetime | 创建时间 | 自动填充 |
| `update_time` | datetime | 更新时间 | 自动更新 |
| `creater` | bigint | 创建人 | - |
| `updater` | bigint | 更新人 | - |
| `deleted` | tinyint | 0=正常, 1=逻辑删除 | `idx_deleted` |

> **金额单位统一为分**。`discount_value` 在 `discount` 类型下承载的是折扣百分比（如 80 表示 8 折），不再是金额。

---

### 2. `user_coupon` — 用户优惠券表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 用户优惠券 id，自增主键 | PK |
| `user_id` | bigint | 用户 id | `idx_user` |
| `coupon_id` | bigint | 关联优惠券 id | `idx_coupon` |
| `status` | varchar(20) | 状态：unused-未使用 used-已使用 expired-已过期 refunded-已退 | `idx_status` |
| `obtain_time` | datetime | 领取时间 | 可空 |
| `use_time` | datetime | 使用时间 | 可空，退还时置 null |
| `expire_time` | datetime | 过期时间 | 可空 |
| `order_id` | bigint | 关联订单 id（使用时记录） | 可空，退还时置 null |
| `code` | varchar(50) | 兑换码（兑换方式获取） | - |
| `create_time` | datetime | 创建时间 | 自动填充 |
| `update_time` | datetime | 更新时间 | 自动更新 |
| `creater` | bigint | 创建人 | - |
| `updater` | bigint | 更新人 | - |
| `deleted` | tinyint | 0=正常, 1=逻辑删除 | 过滤条件 |

> DDL 注释中列出了 `refunded-已退` 状态，但 `consts.go` 仅定义 `unused` / `used` / `expired` 三个常量，退还逻辑实际把状态改回 `unused`。

---

### 3. `coupon_code` — 优惠券兑换码表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 兑换码 id，自增主键 | PK |
| `coupon_id` | bigint | 关联优惠券 id | `idx_coupon` |
| `code` | varchar(50) | 兑换码 | `uk_code` 唯一约束 |
| `status` | varchar(20) | 状态：unused-未使用 used-已使用 | - |
| `user_id` | bigint | 兑换用户 id | 可空 |
| `expire_time` | datetime | 过期时间 | 可空 |
| `create_time` | datetime | 创建时间 | 自动填充 |
| `update_time` | datetime | 更新时间 | 自动更新 |
| `creater` | bigint | 创建人 | - |
| `updater` | bigint | 更新人 | - |
| `deleted` | tinyint | 0=正常, 1=逻辑删除 | 过滤条件 |

> `uk_code` 唯一索引是兑换码防重的最后一道兜底：应用层用 `crypto/rand` 生成 12 位码并在内存里去重，DB 唯一约束防止跨批次碰撞。

---

## 关系图

```
coupon (1) ──┬── (N) user_coupon ─── user_id → user 域
             │
             └── (N) coupon_code ─── user_id → user 域（兑换后回填）

user_coupon.order_id → trade 域 order（无外键约束，仅逻辑关联）
coupon.scopes(json) → course 域 三级分类 id（无外键约束）
```

兑换链路：`coupon_code.code` 兑换成功后写入 `user_coupon.code`，同一张券在两表中留痕。

---

## 状态枚举（`consts.go`）

| 常量组 | 常量名 | 值 | 含义 |
|--------|--------|----|------|
| 优惠券状态 | `CouponStatusDraft` | `draft` | 草稿（待发放） |
| | `CouponStatusIssued` | `issued` | 已发放 |
| | `CouponStatusPaused` | `paused` | 暂停发放 |
| | `CouponStatusEnded` | `ended` | 已结束 |
| 优惠券类型 | `DiscountTypeReduce` | `reduce` | 满减 |
| | `DiscountTypeDiscount` | `discount` | 折扣 |
| | `DiscountTypeNoThreshold` | `no_threshold` | 无门槛 |
| 获取方式 | `ObtainWayReceive` | `receive` | 手动领取 |
| | `ObtainWayExchange` | `exchange` | 兑换码兑换 |
| | `ObtainWayAssign` | `assign` | 后台发放 |
| 用户券状态 | `UserCouponStatusUnused` | `unused` | 未使用 |
| | `UserCouponStatusUsed` | `used` | 已使用 |
| | `UserCouponStatusExpired` | `expired` | 已过期 |
| 兑换码状态 | `CouponCodeStatusUnused` | `unused` | 未兑换 |
| | `CouponCodeStatusUsed` | `used` | 已兑换 |

---

## 模型扩展模式

所有 `*_gen.go` 由 goctl 自动生成，不可编辑。扩展方法统一放在同名的自定义 `.go` 文件中：

```
couponmodel_gen.go        ← goctl 生成，只读（Insert/FindOne/Update/Delete）
couponmodel.go            ← 手写扩展 FindList / FindPage / IncrIssueNum / ...
```

当前项目自定义 Model 模式：
- `couponmodel.go` — FindList, FindPage, FindByIds, IncrIssueNum, DecrIssueNum, AddUsedNum, UpdateStatus, SoftDelete
- `usercouponmodel.go` — FindByUserAndStatus, FindPageByUser, FindByIdsAndUser, CountByUserAndCoupon, UpdateStatusByIds
- `couponcodemodel.go` — FindPageByCoupon, MarkUsed, BatchInsert

goctl 生成的基础方法：`coupon` / `user_coupon` 为 `Insert` / `FindOne` / `Update` / `Delete`；`coupon_code` 额外生成了按唯一键查询的 `FindOneByCode`。

**缓存约定**：条件更新类扩展方法全部走 `ExecNoCacheCtx`，并在写后手工 `DelCacheCtx` 失效 `cache:coupon:id:` / `cache:userCoupon:id:` / `cache:couponCode:id:`、`cache:couponCode:code:` 等键；分页与批量查询走 `QueryRowsNoCacheCtx`，不进缓存。
