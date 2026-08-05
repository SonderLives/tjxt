> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/promotion/rpc/promotion.proto`

---

# Promotion RPC Spec

## 服务名

`Promotion` — 优惠券 / 促销中心微服务，通过 etcd 服务发现（key: `promotion.rpc`）。覆盖 coupons / user-coupons / codes 三组业务。

## RPC 方法总览

### 优惠券管理

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `CouponCreate` | `CouponFormDTO { id, name, discountType, ... }` | `CouponDetailVO` | 新增优惠券，落库后固定为 `draft` 草稿状态 |
| `CouponList` | `CouponListRequest { userId }` | `CouponListReply { list }` | C 端查询发放中的券，标记 available / received |
| `CouponPage` | `CouponPageRequest { page, name, status, type }` | `CouponPageReply { total, list }` | 管理端分页查询，支持名称模糊 + 状态 + 类型过滤 |
| `CouponGet` | `IdRequest { id, userId }` | `CouponDetailVO` | 按 ID 查询券详情，已删除券按不存在处理 |
| `CouponDelete` | `IdRequest { id, userId }` | `Empty {}` | 逻辑删除，仅 `draft` / `paused` 可删 |
| `CouponIssue` | `CouponIssueFormDTO { id, issueBeginTime, issueEndTime, termBeginTime, termDays, termEndTime }` | `Empty {}` | 发放券，写入发放期与有效期并置为 `issued` |
| `CouponPause` | `IdRequest { id, userId }` | `Empty {}` | 暂停发放，仅 `issued` 可暂停 |

**请求字段说明**（`CouponFormDTO`）：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 优惠券 ID，新增时省略 |
| `name` | string | 优惠券名称，不可空 |
| `discountType` | string | `reduce`-满减 / `discount`-折扣 / `no_threshold`-无门槛 |
| `discountValue` | int64 | 满减金额（分）或折扣百分比（如 80 表示 8 折） |
| `maxDiscountAmount` | int64 | 最高优惠金额（分），仅折扣券生效，0=不封顶 |
| `thresholdAmount` | int64 | 使用门槛金额（分），0=无门槛 |
| `obtainWay` | string | `receive`-领取 / `exchange`-兑换码 / `assign`-发放 |
| `specific` | bool | 是否限定适用范围 |
| `scopes` | repeated int64 | 适用范围 id 列表（课程三级分类 id），`specific=true` 时必填 |
| `totalNum` | int64 | 发行总量，0=无上限 |
| `userLimit` | int64 | 每人限领数量，0=不限 |

**请求字段说明**（`CouponIssueFormDTO`）：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 优惠券 ID |
| `issueBeginTime` | string | 发放开始时间，为空视为立即发放 |
| `issueEndTime` | string | 发放结束时间，不得早于开始时间 |
| `termBeginTime` | string | 使用开始时间 |
| `termDays` | int64 | 相对领取日的有效天数 |
| `termEndTime` | string | 使用结束时间，与 `termDays` 二选一必填 |

---

### 用户优惠券

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `UserCouponAvailable` | `OrderCourseListRequest { courseList, userId }` | `CouponDiscountListReply { list }` | 穷举用户所有可用券组合，按优惠金额降序返回方案 |
| `UserCouponDiscount` | `OrderCouponDTO { courseList, userCouponIds, userId }` | `CouponDiscountDTO` | 按用户选定的券方案计算优惠明细 |
| `UserCouponPage` | `UserCouponPageRequest { page, status, userId }` | `UserCouponPageReply { total, list }` | 分页查询我的优惠券 |
| `UserCouponRefund` | `IdsRequest { ids, userId, orderId }` | `Empty {}` | 退还券，状态 `used` → `unused` |
| `UserCouponRules` | `IdsRequest { ids, userId, orderId }` | `RulesReply { rules }` | 查询券规则文案，用于订单页「已优惠」明细 |
| `UserCouponUse` | `IdsRequest { ids, userId, orderId }` | `Empty {}` | 核销券，状态 `unused` → `used` |
| `UserCouponExchange` | `ExchangeRequest { code, userId }` | `Empty {}` | 兑换码兑换优惠券 |
| `UserCouponReceive` | `IdRequest { id, userId }` | `Empty {}` | 领取优惠券 |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `userId` | int64 | 当前用户 ID（来自 JWT），`<= 0` 直接返回未授权 |
| `ids` | repeated int64 | 用户券 ID 列表，仅接受属于当前用户的券 |
| `orderId` | int64 | 关联订单 ID，核销时写入 `user_coupon.order_id` |
| `code` | string | 兑换码，服务端统一 `TrimSpace` + `ToUpper` 后匹配 |
| `courseList` | repeated `OrderCourseDTO` | 订单课程列表，元素含 `cateId` / `id` / `price` |
| `userCouponIds` | repeated int64 | 用户选定的券方案对应的用户券 ID |
| `status` | string | 用户券状态过滤：`unused` / `used` / `expired`，空表示不限 |

**响应字段说明**（`CouponDiscountDTO`）：

| 字段 | 类型 | 说明 |
|------|------|------|
| `discountAmount` | int64 | 本方案最大优惠金额，单位分 |
| `discountDetail` | map<int64,int64> | key=课程 id，value=该课分摊到的优惠金额 |
| `ids` | repeated int64 | 本方案使用的用户券 id |
| `rules` | repeated string | 券规则文案，如「满100元减20元」 |

---

### 兑换码

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `CouponCodePage` | `CouponCodePageRequest { page, couponId, status }` | `CouponCodePageReply { total, list }` | 管理端分页查询兑换码 |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `couponId` | int64 | 所属优惠券 ID，`<= 0` 表示不过滤 |
| `status` | string | `unused`-未兑换 / `used`-已兑换，空表示不限 |

---

## 被哪些 API 服务消费

| 消费方 | 调用方式 | 说明 |
|--------|---------|------|
| `promotion-api` (自身 API 层) | `apps/promotion/api/internal/svc/servicecontext.go` 中 `promotionclient "tjxt/apps/promotion/rpc/promotion"` → `PromotionRpc` | 全部 16 个 HTTP handler 均透传到本 RPC |

（注：截至当前版本，仅 promotion-api 在 `servicecontext.go` 中注入 `promotionclient`。proto 头注释声明「供 promotion-api 与 trade 等内部服务通过 etcd 服务发现调用」，但 `apps/trade` 尚未 import 本服务 client。）

---

## 调用典型场景

1. **运营配置券** → 管理端调 `CouponCreate` 建券（draft）→ 调 `CouponIssue` 设置发放期与有效期并上线（issued）→ 兑换码类型券在首次发放时自动批量生成兑换码
2. **C 端领券** → 前端调 `CouponList` 拉取发放中券列表 → 用户点击领取调 `UserCouponReceive` → 服务端条件更新扣减 `issue_num` 后写入 `user_coupon`
3. **兑换码兑换** → 用户输入兑换码调 `UserCouponExchange` → 校验码状态/过期/限领 → `MarkUsed` 核销码 → 扣库存 → 发放用户券
4. **下单选券** → 结算页调 `UserCouponAvailable` 传入订单课程列表 → 返回按优惠金额降序的组合方案 → 用户选定后调 `UserCouponDiscount` 复算明细
5. **订单核销** → 订单支付成功后调 `UserCouponUse`（携带 `orderId`）核销券并累加 `used_num`
6. **订单取消/退款** → 调 `UserCouponRefund` 把券状态回滚为 `unused` 并回滚 `used_num`
7. **兑换码管理** → 管理端调 `CouponCodePage` 按券与状态分页查看兑换码发放与核销情况

---

## 自定义 Model 方法

`couponmodel.go` 扩展了：
- `FindList(ctx)` — 查询全部未删除优惠券
- `FindPage(ctx, name, status, discountType, offset, limit)` — 条件分页查询
- `FindByIds(ctx, ids)` — 批量查询券规则，避免折扣计算时 N+1
- `IncrIssueNum(ctx, id)` — 原子扣减库存（`issue_num + 1`），返回受影响行数
- `DecrIssueNum(ctx, id)` — 库存回滚补偿
- `AddUsedNum(ctx, id, delta)` — 已使用数量增减，delta 可为负
- `UpdateStatus(ctx, id, status, updater)` — 更新券状态
- `SoftDelete(ctx, id, updater)` — 逻辑删除，仅 `draft` / `paused` 生效

`usercouponmodel.go` 扩展了：
- `FindByUserAndStatus(ctx, userId, status)` — 按用户+状态查询
- `FindPageByUser(ctx, userId, status, offset, limit)` — 按用户+状态分页
- `FindByIdsAndUser(ctx, ids, userId)` — 批量按 id 查当前用户券，防越权
- `CountByUserAndCoupon(ctx, userId, couponId)` — 限领校验计数
- `UpdateStatusByIds(ctx, ids, userId, fromStatus, toStatus, orderId)` — 带源状态条件的批量状态流转

`couponcodemodel.go` 扩展了：
- `FindPageByCoupon(ctx, couponId, status, offset, limit)` — 兑换码分页
- `MarkUsed(ctx, id, userId)` — 条件更新核销兑换码，返回受影响行数
- `BatchInsert(ctx, couponId, codes, creater)` — 批量写入兑换码

`common.go` 提供逻辑层公共函数：`generateCodes`、`parseTime`、`parseScopes`、`matchScope`、`couponRule`、`calcDiscount`、`couponReceivable`、`userCouponExpireTime`、`validateCouponForm`。
