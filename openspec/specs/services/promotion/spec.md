# promotion（营销优惠）Specification

## Purpose

promotion（营销优惠）微服务负责优惠券模板管理（建券 / 发放 / 暂停 / 删除）、用户优惠券全生命周期（领取 / 兑换码兑换 / 核销 / 退还 / 过期）、券组合穷举与订单优惠试算，以及兑换码批量生成与核销。HTTP 服务 `promotion-api` 监听 8812（纯透传层，全部业务判断下沉 RPC）；RPC 服务 `Promotion` 监听 8092（gRPC），经 etcd（`127.0.0.1:2379`，key：`promotion.rpc`）注册与发现。数据库为 MySQL `tj_promotion` 库（coupon / user_coupon / coupon_code 共 3 张表），Redis 缓存主键查询。proto 头注释声明供 promotion-api 与 trade 等内部服务调用，但截至当前版本仅 promotion-api 在 `servicecontext.go` 注入了 client，`apps/trade` 尚未 import（未接线）；本服务刻意不依赖任何其他服务的 RPC——券适用范围 `scopes` 仅回传课程三级分类 id，分类名称由 course 服务维护。

## Requirements

### Requirement: HTTP 优惠券管理接口

promotion-api SHALL 提供优惠券模板的管理接口（全部经 `PromotionRpc` 透传）：

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /coupons | 新增优惠券接口，请求体 `CouponFormDTO`，响应 `R{data: R}` |
| GET | /coupons/list | 查询发放中的优惠券列表，响应 `R{data: R«List«CouponVO»»}` |
| GET | /coupons/page | 分页查询优惠券接口，响应 `R{data: R«PageDTO«CouponPageVO»»}` |
| DELETE | /coupons/{id} | 删除优惠券 |
| GET | /coupons/{id} | 根据 id 查询优惠券接口，响应 `R{data: R«CouponDetailVO»}` |
| PUT | /coupons/{id}/issue | 发放优惠券接口，请求体 `CouponIssueFormDTO` |
| PUT | /coupons/{id}/pause | 暂停发放优惠券接口 |

补充约束：

- 响应统一为 `pkg/response.R{Code,Msg,RequestId,Data any}`；分页约定 `PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`；错误码遵循共享错误码规范
- 聚合文档 `docs/tjxt.openapi.json` 将上述端点认证列标注为「否」，但按 configs.md 记载，全部路由在 `routes.go` 中由 `rest.WithJwt(serverCtx.Config.Auth.AccessSecret)` 统一包裹，需携带有效 JWT，以配置记载为准
- rpc-spec 记载 promotion-api 共 16 个 HTTP handler，聚合文档仅收录 14 个端点

#### Scenario: 新增优惠券固定落为草稿

- **WHEN** 管理端携带 `CouponFormDTO` 请求 POST /coupons
- **THEN** 落库后固定为 `draft` 草稿状态，回查 `FindOne` 返回 `CouponDetailVO`
- **AND** 必须再调 PUT /coupons/{id}/issue 设置发放期与有效期后才开始发放

#### Scenario: 暂停发放不影响已领券

- **WHEN** 管理端对 `issued` 状态的券调 PUT /coupons/{id}/pause
- **THEN** 券模板状态置为 `paused`，C 端不可再领取
- **AND** 用户已持有的 `user_coupon` 不受影响

### Requirement: HTTP 用户优惠券接口

promotion-api SHALL 提供用户侧优惠券的领取、兑换、核销、退还与试算接口：

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /user-coupons/{couponId}/receive | 领取优惠券接口 |
| POST | /user-coupons/{code}/exchange | 兑换码兑换优惠券接口 |
| POST | /user-coupons/available | 查询我的优惠券可用方案，请求体 JSON（订单课程列表），响应 `R{data: R«List«CouponDiscountDTO»»}` |
| POST | /user-coupons/discount | 根据券方案计算订单优惠明细，请求体 `OrderCouponDTO`，响应 `R{data: R«CouponDiscountDTO»}` |
| GET | /user-coupons/page | 分页查询我的优惠券接口，响应 `R{data: R«PageDTO«CouponVO»»}` |
| PUT | /user-coupons/use | 核销指定优惠券 |
| PUT | /user-coupons/refund | 退还指定优惠券 |
| GET | /user-coupons/rules | 查询券规则文案，响应 `R{data: R«List«string»»}`（聚合文档摘要误标为「分页查询我的优惠券接口」，以 rpc-spec 的 `UserCouponRules` 语义为准） |

#### Scenario: 结算页查询可用方案

- **WHEN** 用户提交订单课程列表请求 POST /user-coupons/available
- **THEN** 穷举用户所有可用券组合，按优惠金额降序返回 `CouponDiscountDTO` 方案列表
- **AND** 用户选定后调 POST /user-coupons/discount 按选定券方案复算明细

#### Scenario: 订单支付成功核销优惠券

- **WHEN** 订单支付成功后携带用户券 id 调 PUT /user-coupons/use
- **THEN** 用户券状态 `unused` → `used`，`order_id` 与 `use_time` 写入记录

### Requirement: HTTP 兑换码分页查询接口

promotion-api SHALL 提供管理端兑换码分页查询端点 `GET /codes/page`（handler `CouponCodePageHandler`，logic `CouponCodePageLogic`）：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /codes/page | 管理端分页查询兑换码，透传 RPC `CouponCodePage` |

请求/响应结构：

| 结构 | 字段 |
|------|------|
| `CouponCodePageReq`（query，内嵌 `PageRequest`） | `pageNo` / `pageSize`、`couponId` int64（`form:"couponId,omitempty"`）、`status` string（`form:"status,omitempty"`） |
| `CouponCodePageResp` | `list` []`ExchangeCodeVO`、`total` int64、`pages` int64 |
| `ExchangeCodeVO` | `id` int64、`code` string |

- 过滤语义：`couponId <= 0` 不过滤券；`status` 取 `unused` / `used`，空串表示不限
- 分页计算：`limit = pageSize`，`pageSize <= 0` 时取 10；`pages = (total + limit - 1) / limit`
- 认证：与其余端点一致，由 `routes.go` 的 `rest.WithJwt(...)` 全局包裹

> 注：wiki 的 `promotion/api-spec.md` 未收录该端点（源于 `docs/tjxt.openapi.json` 提取不全），但 `apps/promotion/api/internal/handler/routes.go` 与 `couponcodepagelogic.go` 中确实存在，此处以代码实现为准补齐。

#### Scenario: 分页查询兑换码

- **WHEN** 管理端携带 `couponId` 与 `status` 请求 GET /codes/page
- **THEN** API 层透传 RPC `CouponCodePage(page, couponId, status)`，返回 `list` / `total` / `pages`
- **AND** `couponId <= 0` 时不过滤券，`status` 为空时不过滤状态

#### Scenario: 缺省分页大小

- **WHEN** 请求未传 `pageSize` 或传入 `<= 0`
- **THEN** 服务端按 `limit = 10` 计算总页数 `pages = (total + 10 - 1) / 10`
- **AND** 列表为空时 `list` 返回 `[]` 而非 null

### Requirement: RPC 服务契约

RPC 服务 `Promotion` SHALL 经 etcd（key：`promotion.rpc`，监听 `0.0.0.0:8092`）提供服务，方法按三组业务划分：

优惠券管理：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `CouponCreate` | `CouponFormDTO` | `CouponDetailVO` | 新增优惠券，落库后固定为 `draft` 草稿状态 |
| `CouponList` | `CouponListRequest { userId }` | `CouponListReply { list }` | C 端查询发放中的券，标记 available / received |
| `CouponPage` | `CouponPageRequest { page, name, status, type }` | `CouponPageReply { total, list }` | 管理端分页查询，支持名称模糊 + 状态 + 类型过滤 |
| `CouponGet` | `IdRequest { id, userId }` | `CouponDetailVO` | 按 ID 查询券详情，已删除券按不存在处理 |
| `CouponDelete` | `IdRequest { id, userId }` | `Empty {}` | 逻辑删除，仅 `draft` / `paused` 可删 |
| `CouponIssue` | `CouponIssueFormDTO` | `Empty {}` | 发放券，写入发放期与有效期并置为 `issued` |
| `CouponPause` | `IdRequest { id, userId }` | `Empty {}` | 暂停发放，仅 `issued` 可暂停 |

用户优惠券：

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

兑换码：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `CouponCodePage` | `CouponCodePageRequest { page, couponId, status }` | `CouponCodePageReply { total, list }` | 管理端分页查询兑换码（`couponId <= 0` 不过滤；`status`：`unused`/`used`，空不限） |

关键字段约束：

- `userId`（来自 JWT）`<= 0` 直接返回未授权；`ids` 仅接受属于当前用户的券；核销时 `orderId` 写入 `user_coupon.order_id`；兑换码 `code` 服务端统一 `TrimSpace` + `ToUpper` 后匹配
- `CouponFormDTO` 关键字段：`discountType`（`reduce`-满减 / `discount`-折扣 / `no_threshold`-无门槛）、`discountValue`（满减金额（分）或折扣百分比，如 80 表示 8 折）、`maxDiscountAmount`（最高优惠金额（分），仅折扣券生效，0=不封顶）、`thresholdAmount`（门槛金额（分），0=无门槛）、`obtainWay`（`receive`/`exchange`/`assign`）、`scopes`（课程三级分类 id 列表，`specific=true` 时必填）、`totalNum`（0=无上限）、`userLimit`（0=不限）
- `CouponIssueFormDTO` 关键字段：`issueBeginTime`（为空视为立即发放）、`issueEndTime`（不得早于开始时间）、`termBeginTime` / `termDays` / `termEndTime`（`termDays` 与 `termEndTime` 二选一必填）
- 响应 `CouponDiscountDTO`：`discountAmount`（本方案最大优惠金额，分）、`discountDetail`（map<课程 id, 分摊优惠金额>）、`ids`（使用的用户券 id）、`rules`（规则文案，如「满100元减20元」）

> 已知缺口（未接线）：proto 头注释声明「供 promotion-api 与 trade 等内部服务通过 etcd 服务发现调用」，但截至当前版本仅 promotion-api 在 `servicecontext.go` 注入了 `promotionclient`，`apps/trade` 尚未 import 本服务 client。

#### Scenario: C 端领券调用链

- **WHEN** 前端调 `CouponList` 拉取发放中券列表后用户点击领取
- **THEN** 调 `UserCouponReceive`，服务端条件更新扣减 `issue_num` 后写入 `user_coupon`

#### Scenario: 订单取消退还优惠券

- **WHEN** 订单取消或退款后调 `UserCouponRefund`（携带 `ids` 与 `orderId`）
- **THEN** 券状态回滚为 `unused` 并回滚 `used_num`

### Requirement: 建券与表单校验

系统 SHALL 在新建优惠券时一律落为 `draft` 草稿状态（`buildCoupon` 中硬编码 `Status: CouponStatusDraft`），并按优惠类型分支校验表单（`validateCouponForm`）：

| 优惠类型 | 校验规则 |
|----------|---------|
| `reduce` 满减 | `thresholdAmount > 0`；`0 < discountValue < thresholdAmount`（减的钱必须小于门槛） |
| `discount` 折扣 | `1 <= discountValue <= 99`（折扣百分比区间） |
| `no_threshold` 无门槛 | `discountValue > 0` |
| 其他 | 直接拒绝「优惠券类型非法」 |

通用校验：名称非空（`strings.TrimSpace(name)` 为空则拒绝）；获取方式必须是 `receive` / `exchange` / `assign` 之一；`totalNum` 与 `userLimit` 不能为负数；`specific=true` 时 `scopes` 不能为空且 JSON 序列化后存入 `scopes` 字段（`specific=false` 时该字段留空）。

#### Scenario: 满减券校验失败

- **WHEN** 新建 `reduce` 满减券时 `discountValue >= thresholdAmount` 或 `thresholdAmount <= 0`
- **THEN** 表单校验拒绝，优惠券不落库

#### Scenario: 建券流程

- **WHEN** `CouponCreate` 通过校验
- **THEN** `buildCoupon` 转 model（Status 固定 `draft`）→ `Insert` 取 `LastInsertId` → 回查 `FindOne` 返回 `CouponDetailVO`

### Requirement: 优惠券生命周期状态机

系统 SHALL 以状态机硬性约束优惠券模板的状态流转，非法流转返回 `Conflict`：`draft --CouponIssue--> issued --CouponPause--> paused --CouponIssue--> issued`；`draft`/`paused --CouponDelete--> deleted=1`；`ended` 为终态不可再发放。

| 规则 | 说明 |
|------|------|
| 已结束不可再发放 | `status == ended` 时 `CouponIssue` 返回「优惠券已结束，无法再次发放」 |
| 只有发放中可暂停 | `status != issued` 时 `CouponPause` 返回「只有发放中的优惠券才能暂停」 |
| 只有草稿/暂停可删 | `SoftDelete` SQL 带 `status in ('draft','paused')` 条件，受影响行数为 0 即返回「优惠券不存在或当前状态不允许删除」 |
| 已删除按不存在处理 | `CouponGet` / `CouponPause` / `CouponIssue` 均在查到后判断 `Deleted == 1` 并返回 `NotFound` |
| 暂停不影响已领券 | 暂停仅改券模板状态，用户已持有的 `user_coupon` 不受影响 |

#### Scenario: 结束状态为终态

- **WHEN** 对 `ended` 状态的券再次调 `CouponIssue`
- **THEN** 返回「优惠券已结束，无法再次发放」

#### Scenario: 逻辑删除条件更新

- **WHEN** 对非 `draft` / `paused` 状态的券调 `CouponDelete`
- **THEN** 条件更新影响 0 行，返回「优惠券不存在或当前状态不允许删除」

### Requirement: 发放与有效期设置

`CouponIssue` SHALL 分别校验发放期与有效期，有效期采用「绝对区间」或「相对天数」二选一：

| 规则 | 说明 |
|------|------|
| 时间格式兼容 | `parseTime` 依次尝试 `2006-01-02 15:04:05` / `2006-01-02T15:04:05` / RFC3339 / `2006-01-02`，全失败返回「时间格式非法」 |
| 发放期顺序 | `issueEndTime` 早于 `issueBeginTime` 时拒绝 |
| 有效期顺序 | `termEndTime` 早于 `termBeginTime` 时拒绝 |
| 有效期二选一 | `termDays <= 0` 且 `termEndTime` 为空时拒绝「请设置有效期天数或使用结束时间」 |
| 立即发放 | 未指定 `issueBeginTime` 时自动填 `now()` |
| 兑换码生成时机 | 仅在首次发放（原状态为 `draft`）且 `obtainWay == exchange` 且 `totalNum > 0` 时批量生成兑换码；再次发放（paused → issued）不会重复生成 |

#### Scenario: 发放兑换码类型券

- **WHEN** 首次对 `obtainWay=exchange` 且 `totalNum > 0` 的草稿券调 `CouponIssue`
- **THEN** 批量生成兑换码并将券置为 `issued`
- **AND** 后续 paused → issued 再次发放不会重复生成兑换码

#### Scenario: 有效期缺失拒绝发放

- **WHEN** 发放时 `termDays <= 0` 且 `termEndTime` 为空
- **THEN** 拒绝并返回「请设置有效期天数或使用结束时间」

### Requirement: 兑换码生成算法

系统 SHALL 使用密码学安全随机源生成兑换码，字符集剔除易混淆字符：

| 参数 | 值 | 说明 |
|------|----|------|
| `codeAlphabet` | `23456789ABCDEFGHJKLMNPQRSTUVWXYZ` | 32 个字符，剔除 `0/O/1/I` |
| `codeLength` | `12` | 兑换码固定 12 位 |
| 随机源 | `crypto/rand.Int` | 非 `math/rand`，防止可预测 |

生成流程：`n <= 0` 直接返回 nil；循环收集满 n 个 12 位码，命中本批次 `seen` 集合的重复码丢弃重来；`BatchInsert` 一条 insert 多 values 写入 `coupon_code`。内存 `seen` 只保证单批次不重复，跨批次靠 `coupon_code.uk_code` 唯一索引兜底。

#### Scenario: 批量生成去重

- **WHEN** 首次发放兑换码类型券触发 `generateCodes(n)`
- **THEN** 生成 n 个 12 位、剔除易混淆字符的随机码，单批次内无重复
- **AND** 跨批次碰撞由 `uk_code` 唯一约束兜底

### Requirement: 领取防超发

`UserCouponReceive` SHALL 完全依赖 SQL 条件更新扣减库存（先抢库存再写用户券），写失败则补偿回滚库存。前置校验链按序执行：`userId > 0`（否则 `Unauthorized`）→ `id > 0`（否则「优惠券 id 非法」）→ 券存在且未删除（否则 `NotFound`「优惠券不存在」）→ `obtainWay == receive`（否则 `Conflict`「该优惠券不支持手动领取」）→ `couponReceivable(coupon, now)`（否则 `Conflict`「优惠券不在领取时间内或已被领完」）→ `userLimit > 0` 时 `CountByUserAndCoupon < userLimit`（否则 `Conflict`「已达到该优惠券的领取上限」）。

`couponReceivable` 四项判定：`status == issued`、`now >= issueBeginTime`（若设置）、`now <= issueEndTime`（若设置）、`totalNum == 0 || issueNum < totalNum`。

防超发关键 SQL（`IncrIssueNum`）：`update coupon set issue_num = issue_num + 1 where id = ? and deleted = 0 and status = 'issued' and (total_num = 0 or issue_num < total_num)`。原子性来自 InnoDB 行锁 + 条件；`RowsAffected == 0` 即判定「优惠券已被领完」，不做重试。`user_coupon` 落库失败时调 `DecrIssueNum` 补偿（带 `issue_num > 0` 条件防减成负数）；回滚本身失败只记日志，不阻断主流程错误返回。

> 已知缺口（非强一致）：限领校验 `CountByUserAndCoupon` 是独立 select，与后续 insert 之间存在并发窗口，极端并发下同一用户可能超出 `userLimit` 一张；代码注释明确「非强一致，最终由库存与业务容忍度兜底」。

过期时间计算（`userCouponExpireTime`）：优先取券模板的绝对 `termEndTime`；否则按 `now + termDays` 天，并把末尾对齐到当天 `23:59:59`（避免按秒卡点）；两者都没有则留空（永不过期）。

#### Scenario: 并发抢券只有一个成功

- **WHEN** 多个用户并发领取库存仅剩 1 张的优惠券
- **THEN** 仅一个 `IncrIssueNum` 条件更新影响 1 行并成功写入 `user_coupon`
- **AND** 其余请求 `RowsAffected == 0`，返回「优惠券已被领完」，不做重试

#### Scenario: 用户券落库失败补偿回滚

- **WHEN** 库存扣减成功后 `user_coupon` 落库失败
- **THEN** 调 `DecrIssueNum` 回滚库存（`issue_num > 0` 条件保证不会减成负数）
- **AND** 回滚本身失败只记录错误日志，不阻断主流程错误返回

#### Scenario: 相对有效期对齐到日末

- **WHEN** 领取的券模板仅设置了 `termDays`
- **THEN** 用户券 `expire_time = now + termDays` 天，末尾对齐到当天 `23:59:59`

### Requirement: 兑换码兑换

`UserCouponExchange` SHALL 通过条件更新保证兑换码核销与库存扣减的并发安全，流程为：`userId > 0` 且 code 归一化（`TrimSpace` + `ToUpper`，空则拒绝）→ `FindOneByCode` 查码（不存在 / 已删除返回 `NotFound`「兑换码不存在」）→ 读校验 `status != unused` 返回 `Conflict`「兑换码已被使用」→ `expireTime` 已过返回 `Conflict`「兑换码已过期」→ 查关联券（不存在 / 已删除 / `ended` 返回「优惠券活动已结束」）→ `userLimit > 0` 时校验限领 → `MarkUsed` 条件更新核销码（`rows == 0` 返回 `Conflict`「兑换码已被使用」）→ `IncrIssueNum` 扣库存 → `Insert user_coupon` 并回填 `code` 字段。

并发安全关键 SQL（`MarkUsed`）：`update coupon_code set status = 'used', user_id = ? where id = ? and status = 'unused' and deleted = 0`；`RowsAffected == 0` 表示已被他人抢先兑换。读校验只是快速失败，真正的并发防线是条件更新。`MarkUsed` 更新后同时失效 `cache:couponCode:id:` 与 `cache:couponCode:code:` 两个键。

容错语义：`IncrIssueNum` 扣库存失败只记日志、不中断兑换（码已核销，优先保证用户拿到券）；用户券落库失败返回 `Internal`「兑换失败」，此处未回滚已核销的兑换码与已扣减的库存。与领取流程不同，兑换不做 `couponReceivable` 发放期/库存判定——兑换码本身即库存凭证，只校验券未删除、未结束。

> 已知缺口：用户券落库失败时未回滚已核销的兑换码与已扣减的库存。

#### Scenario: 并发抢兑同一码

- **WHEN** 两个用户并发兑换同一兑换码
- **THEN** 仅一个 `MarkUsed` 条件更新影响 1 行成功拿到券
- **AND** 另一个返回 `Conflict`「兑换码已被使用」

#### Scenario: 兑换链路两表留痕

- **WHEN** 用户成功兑换兑换码
- **THEN** `coupon_code` 置 `used` 并回填 `user_id`，同时写入 `user_coupon`（`code` 字段回填兑换码），同一张券在两表中留痕

### Requirement: 核销幂等

`UserCouponUse` SHALL 以带源状态条件的批量更新实现 `unused → used` 流转，保证同一张券不会被重复核销：`update user_coupon set status = 'used', use_time = now(), order_id = ? where id in (...) and user_id = ? and status = 'unused' and deleted = 0`。

| 规则 | 说明 |
|------|------|
| 归属校验 | `FindByIdsAndUser` 查回的条数必须等于入参 `ids` 长度，否则「优惠券不存在或不属于当前用户」 |
| 幂等保证 | `where status = 'unused'` 使重复调用的第二次更新影响 0 行 |
| 全成功语义 | `RowsAffected != len(ids)` 即判定「部分优惠券已被使用或已过期」并返回 `Conflict`——要么全部核销成功，要么整体报错 |
| 统计字段容错 | 逐张券 `AddUsedNum(+1)` 失败只记日志，不阻断主流程（统计类字段，可离线校准） |
| 订单绑定 | 核销时把 `orderId` 写入 `user_coupon.order_id`，`use_time` 置 `now()` |

> 已知缺口：`RowsAffected != len(ids)` 判定发生在更新之后，此时已成功的那部分行不会回滚（无事务包裹），实际是「部分成功 + 报错」的语义。

#### Scenario: 重复核销幂等失败

- **WHEN** 对已 `used` 的券再次调 `UserCouponUse`
- **THEN** 条件更新影响 0 行，返回 `Conflict`「部分优惠券已被使用或已过期」

#### Scenario: 部分核销不回滚

- **WHEN** 批量核销多张券且其中部分已非 `unused`
- **THEN** 更新执行后按影响行数判定返回 `Conflict`，已成功的那部分行不回滚

### Requirement: 退还

`UserCouponRefund` SHALL 以 `used → unused` 的带源状态条件批量更新实现退还，并清空使用痕迹：`update user_coupon set status = 'unused', use_time = null, order_id = null where id in (...) and user_id = ? and status = 'used' and deleted = 0`。

| 规则 | 说明 |
|------|------|
| 归属校验 | `FindByIdsAndUser` 结果为空则「优惠券不存在或不属于当前用户」（与核销不同，只要求非空，不要求条数相等） |
| 幂等保证 | `where status = 'used'`，`RowsAffected == 0` 返回「优惠券未被使用，无需退还」 |
| 使用痕迹清理 | `use_time` 与 `order_id` 一并置 `null` |
| used_num 回滚 | 对原状态为 `used` 的券逐张 `AddUsedNum(-1)`；SQL 带 `used_num + delta >= 0` 条件防减成负数，失败只记日志 |
| 过期券退还 | 已过期的券仍会被恢复为 `unused`，仅打一条日志——恢复状态但实际不可再用 |

#### Scenario: 退还清理使用痕迹

- **WHEN** 对 `used` 状态的券调 `UserCouponRefund`
- **THEN** 状态回滚为 `unused`，`use_time` 与 `order_id` 置 `null`，并 `AddUsedNum(-1)` 回滚统计

#### Scenario: 未使用券退还幂等拒绝

- **WHEN** 对非 `used` 状态的券调 `UserCouponRefund`
- **THEN** 条件更新影响 0 行，返回「优惠券未被使用，无需退还」

### Requirement: 券组合穷举与单方案计算

系统 SHALL 以位掩码穷举为核心算法（`solution.go`，服务 `UserCouponAvailable` 与 `UserCouponDiscount`）：

可用券筛选（`usableUserCoupons`，逐张过滤）：`uc.Status == unused`；未过期（`ExpireTime` 有效且 `now` 未超过）；券规则存在（`coupons[uc.CouponId]` 可查到且未删除）；有效期已开始（`TermBeginTime` 有效时 `now` 不得早于它）。

组合数量截断（`trimCoupons`）：`maxCombineCoupons = 12`（参与组合运算的最大券数，组合方案数为 2^n）；`maxSolutions = 30`（最多返回方案数）。超过 12 张时先对每张券单独调 `buildSolution` 算「单券优惠力度」作为打分，`sort.SliceStable` 按分值降序只保留前 12 张；无法生效的券打分为 0，自然排到末尾被淘汰。

位掩码穷举（`calcSolutions`）：available 或 courses 为空返回 nil；`for mask := 1; mask < (1<<n); mask++` 遍历所有非空子集，按位组装 combo 后 `buildSolution`（非法组合跳过），按 `solutionKey(ids)`（券 id 排序后拼串）去重；排序规则为优惠金额降序、金额相同时用券数量少的排前；截断到前 30 条。复杂度 `O(2^n × n × m)`（n=券数≤12，m=课程数），最坏 4095 个子集。

单方案计算（`buildSolution`）：券按数组顺序依次作用于「剩余金额」，避免同一笔钱被重复打折。`remaining[courseId]` 初始化为课程原价（累加同 id）；逐张券：`specific=0` 全场通用，否则按 `cateId` 匹配 `scopes`；matched 为空则组合非法；`subtotal` 为匹配课程 remaining 之和；`discount = calcDiscount(券, subtotal)`，`discount <= 0` 则组合非法；按 remaining 比例把 discount 分摊到各课程（`share = discount × remaining[courseId] / subtotal`）；`total <= 0` 返回 false；返回 `{ discountAmount, discountDetail, ids, rules }`。

| 规则 | 说明 |
|------|------|
| 组合合法性 | 组合中任一张券无法生效（未达门槛或无适用课程）时整个组合作废，不做部分生效 |
| 顺序敏感 | 券作用于剩余金额，不同顺序理论上结果不同；实现中按 `available` 数组顺序作用 |
| 取整误差处理 | 最后一门课用 `discount - allocated` 兜底吸收整数除法的余数 |
| 分摊封顶 | `share` 不超过该课程的 `remaining`，防止分摊为负余额 |

#### Scenario: 超过 12 张可用券时截断打分

- **WHEN** 用户可用券超过 12 张时查询可用方案
- **THEN** 按单券优惠力度打分降序保留前 12 张参与组合穷举
- **AND** 无法生效的券打分为 0 排到末尾被淘汰，最终方案列表截断到最多 30 条

#### Scenario: 组合中任一券无法生效则整体作废

- **WHEN** 某个券组合中存在一张未达门槛或无适用课程的券
- **THEN** 该组合整体非法被跳过，不做部分生效

#### Scenario: 优惠按剩余金额比例分摊

- **WHEN** 一张券匹配订单中的多门课程
- **THEN** 优惠金额按各课程剩余金额比例分摊，最后一门课兜底吸收取整余数
- **AND** 单课分摊不超过其剩余金额

### Requirement: 单券优惠计算与规则文案

系统 SHALL 按以下规则计算单券优惠（`calcDiscount`，`common.go`）：

| 券类型 | 计算方式 |
|--------|---------|
| 门槛判定 | 非 `no_threshold` 券，`totalAmount < thresholdAmount` 时返回 0 |
| `discount` 折扣 | `discountValue` 落在 `[1,99]` 之外返回 0；否则 `discount = totalAmount × (100 - discountValue) / 100`，再受 `maxDiscountAmount > 0` 封顶 |
| `reduce` / `no_threshold` | `discount = discountValue`（固定金额） |
| 全局封顶 | 优惠额不超过 `totalAmount` |
| 边界 | `totalAmount <= 0` 直接返回 0 |

规则文案模板（`couponRule`）：`no_threshold` 为「无门槛立减{X}元」；`discount` 为「满{门槛}元打{discountValue×10 转元}折」（`maxDiscountAmount > 0` 时追加「，最多减{X}元」）；`reduce`（默认分支）为「满{门槛}元减{X}元」。`yuan` 函数把分转元：`fmt.Sprintf("%.2f", cent/100)` 后去掉尾随 0 与小数点。

#### Scenario: 折扣券计算并封顶

- **WHEN** 计算一张 `discountValue=80`（8 折）、`maxDiscountAmount > 0` 的折扣券优惠
- **THEN** 优惠 = `totalAmount × 20 / 100`，再受最高优惠金额封顶，且不超过 `totalAmount`

#### Scenario: 未达门槛返回 0

- **WHEN** 非 `no_threshold` 券的 `totalAmount < thresholdAmount`
- **THEN** 优惠计算返回 0，该券在组合穷举中被判为无法生效

### Requirement: 折扣方案复算校验

`UserCouponDiscount` SHALL 仅接受属于当前用户且未使用的券，防止越权占用他人优惠券，校验链按序执行：`userId > 0`（否则 `Unauthorized`）→ `courseList` 非空（否则「订单课程不能为空」）→ `userCouponIds` 为空时返回空 `CouponDiscountDTO`（不报错）→ `FindByIdsAndUser` 条数 == 入参条数（否则「优惠券不存在或不属于当前用户」）→ `usableUserCoupons` 过滤后条数 == 入参条数（否则 `Conflict`「存在已使用或已过期的优惠券，请重新选择」）→ `buildSolution` 返回 ok（否则 `Conflict`「所选优惠券不满足使用条件」）。

#### Scenario: 越权或失效券复算被拒

- **WHEN** 复算方案中包含不属于当前用户、或已使用/已过期的券
- **THEN** 校验链返回相应错误，不计算优惠明细

#### Scenario: 空券方案返回空结果

- **WHEN** `userCouponIds` 为空但课程列表非空
- **THEN** 返回空 `CouponDiscountDTO`，不报错

### Requirement: C 端券列表与我的券分页

系统 SHALL 在 C 端列表与分页查询中应用以下规则：

`CouponList`（C 端券列表标记）：

| 规则 | 说明 |
|------|------|
| 只展示可领券 | 遍历 `FindList` 结果，`couponReceivable == false` 的直接跳过 |
| 批量查已领 | 一次 `FindByUserAndStatus(userId, "")` 拉全部用户券，在内存里按 `couponId` 计数，避免逐张券查库 |
| `received` 标记 | 该券已领数量 `> 0` |
| `available` 标记 | `userLimit == 0 || got < userLimit`——达到限领数量后不可再领 |
| 未登录兼容 | `userId <= 0` 时不查用户券，全部券 `received=false`、`available=true` |

`UserCouponPage`（我的券分页）：

| 规则 | 说明 |
|------|------|
| 分页归一化 | `page.Normalize(pageNo, pageSize)`：pageNo 最小 1，pageSize 缺省 10、最大 100 |
| 批量加载券规则 | 收集 `couponId` 后一次 `FindByIds`，避免 N+1 |
| 到期时间覆盖 | 展示 `user_coupon.expire_time`（用户券实际到期时间），而非券模板的 `term_end_time` |
| 过期自动置灰 | `now` 超过 `expireTime` 时把 `available` 置为 `false`（不改 DB 状态） |
| 券规则缺失跳过 | `coupons[uc.CouponId]` 查不到的用户券直接跳过不计入 list（但 `total` 已按 DB count 返回） |

#### Scenario: 达到限领后列表置为不可领

- **WHEN** 用户已领取数量达到 `userLimit` 且 `userLimit > 0`
- **THEN** 该券在 C 端列表中 `available=false`、`received=true`

#### Scenario: 我的券过期置灰

- **WHEN** 分页查询我的优惠券且某张券 `now > expire_time`
- **THEN** 该券 `available` 置为 `false`（仅展示层置灰，不修改 DB 状态）
- **AND** 到期时间展示 `user_coupon.expire_time` 而非券模板的 `term_end_time`

### Requirement: API 层透传职责

promotion-api SHALL 作为纯透传层，全部业务判断下沉到 RPC：

| 规则 | 说明 |
|------|------|
| 身份注入 | 需要用户身份的接口统一 `auth.UserIdFromCtx(l.ctx)` 取 JWT 中的 userId，再塞进 RPC 请求 |
| 类型转换 | `convert.go` 承担 pb ↔ types 转换；`DiscountDetail` 的 int64 键在 API 层转成 string 键（JSON 兼容） |
| 分页页数计算 | API 层自行算 `pages = (total + limit - 1) / limit`，`pageSize <= 0` 时按 10 兜底 |
| 管理端接口 | `CouponCreate` / `CouponIssue` / `CouponPage` / `CouponCodePage` 不取 userId，直接透传 |

> 已知缺口：`CouponCreate` 的 RPC 实现中 `buildCoupon(in, 0)` 的 operator 硬编码为 `0`，即新建券的 `creater` / `updater` 恒为 0，未接入操作人身份。

#### Scenario: 用户身份由 API 层注入

- **WHEN** C 端用户调用需要身份的用户券接口
- **THEN** API 层从 JWT 解出 userId 塞进 RPC 请求，RPC 层 `userId <= 0` 直接返回未授权

### Requirement: 数据模型与不变量

数据层 SHALL 由 MySQL `tj_promotion` 库的 3 张表承载，配合 goctl 生成的基础方法与手写扩展方法：

| 表 | 说明 | 关键索引 / 唯一键 | 手写扩展方法（`*_gen.go` 由 goctl 生成禁止修改） |
|----|------|------------------|------------------------------------------------|
| `coupon` | 优惠券模板（含优惠类型、门槛、总量/已领/已用数量、发放期与有效期、状态） | `idx_status`、`idx_deleted` | FindList / FindPage / FindByIds / IncrIssueNum / DecrIssueNum / AddUsedNum / UpdateStatus / SoftDelete |
| `user_coupon` | 用户优惠券（用户、券、状态、领取/使用/过期时间、订单、兑换码） | `idx_user`、`idx_coupon`、`idx_status` | FindByUserAndStatus / FindPageByUser / FindByIdsAndUser / CountByUserAndCoupon / UpdateStatusByIds |
| `coupon_code` | 兑换码（券、码、状态、兑换用户、过期时间） | `uk_code`（唯一）、`idx_coupon` | FindPageByCoupon / MarkUsed / BatchInsert（goctl 基础方法额外含按唯一键查询的 `FindOneByCode`） |

关键不变量与约定：

- 优惠券状态枚举（`consts.go`）：`draft` / `issued` / `paused` / `ended`；优惠类型：`reduce` / `discount` / `no_threshold`；获取方式：`receive` / `exchange` / `assign`；用户券状态：`unused` / `used` / `expired`；兑换码状态：`unused` / `used`
- 金额单位统一为分；`discount_value` 在 `discount` 类型下承载折扣百分比（如 80 表示 8 折），不再是金额
- 所有表带 `deleted`（0=正常，1=逻辑删除）与 `creater` / `updater` 审计字段；`uk_code` 唯一索引是兑换码防重的最后一道兜底
- 兑换链路两表留痕：`coupon_code.code` 兑换成功后写入 `user_coupon.code`
- 逻辑关联无外键约束：`user_coupon.order_id → trade 域 order`；`coupon.scopes(json) → course 域三级分类 id`
- 缓存约定：条件更新类扩展方法全部走 `ExecNoCacheCtx` 并在写后手工 `DelCacheCtx` 失效 `cache:coupon:id:` / `cache:userCoupon:id:` / `cache:couponCode:id:`、`cache:couponCode:code:` 键；分页与批量查询走 `QueryRowsNoCacheCtx` 不进缓存

> 冲突记载（以 data-model.md 表述为准）：DDL 注释中 `user_coupon.status` 列出了 `refunded-已退` 状态，但 `consts.go` 仅定义 `unused` / `used` / `expired` 三个常量，退还逻辑实际把状态改回 `unused`。另：`expired` 常量已定义但当前逻辑未主动写入，过期通过 `expire_time` 与 `now` 比较判定。

#### Scenario: 库存相关字段的原子更新

- **WHEN** 领取 / 兑换 / 核销 / 退还触发库存与统计变更
- **THEN** `issue_num` / `used_num` 均经带条件的原子 SQL 更新（`issue_num < total_num`、`used_num + delta >= 0`），不会出现负数或超发

#### Scenario: 兑换码唯一性兜底

- **WHEN** 不同批次生成的兑换码在应用层内存去重之外发生跨批次碰撞
- **THEN** `uk_code` 唯一索引拒绝重复写入，作为防重的最后一道兜底

### Requirement: 配置与部署形态

服务部署 SHALL 遵循以下配置（默认值来自 `apps/promotion/api/etc/promotion-api.yaml` 与 `apps/promotion/rpc/etc/promotion.yaml`）：

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| promotion-api `Port` | `8812` | HTTP 监听端口（Host `0.0.0.0`） |
| promotion.rpc `ListenOn` | `0.0.0.0:8092` | gRPC 监听地址 |
| `Etcd.Hosts[0]` / `Etcd.Key` | `127.0.0.1:2379` / `promotion.rpc` | 服务注册与发现（API 侧为 `PromotionRpc.Etcd`） |
| `Auth.AccessSecret` | `change-me-in-production` | JWT 签名密钥，生产环境必须修改且须与签发方（auth 服务）一致 |
| `Auth.AccessExpire` | `7200` | 访问令牌有效期（秒） |
| RPC `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_promotion?charset=utf8mb4&parseTime=true&loc=Local` | MySQL 连接串，`parseTime=true` 与 `loc=Local` 为必需项（`sql.NullTime` 与 `parseTime` / `userCouponExpireTime` 均按 `time.Local` 计算） |
| RPC `Cache[0]` | `127.0.0.1:6379`，Type `node`，Pass 空 | Redis 单节点缓存，三个 Model 共用 |

行为性约束：

- `DataSource` 与 `Cache` 均标记为 `optional`，缺省时服务仍可启动，但所有 Model 调用都会失败——生产环境必须显式配置
- promotion-api 全部路由由 `rest.WithJwt(serverCtx.Config.Auth.AccessSecret)` 统一包裹（需携带有效 JWT）
- promotion.rpc 不依赖任何其他服务的 RPC：券适用范围 `scopes` 中的课程分类名称由 course 服务维护，本服务仅回传 id（`toCouponDetailVO` 中 `CouponScopeVO` 只填 `Id`、不填 `Name`），刻意避免跨服务强依赖
- Telemetry / Prometheus / Log 等可观测性注入配置已注入勿删

#### Scenario: 缺省存储配置仍可启动

- **WHEN** RPC 配置缺失 `DataSource` / `Cache`
- **THEN** 服务正常启动（optional），但所有 Model 调用失败
- **AND** 生产环境必须显式配置 MySQL 与 Redis

#### Scenario: JWT 密钥不一致导致鉴权失败

- **WHEN** promotion-api 的 `Auth.AccessSecret` 与 auth 服务签发密钥不一致
- **THEN** 所有接口鉴权失败

