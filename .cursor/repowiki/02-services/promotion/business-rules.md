> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/promotion/rpc/internal/logic/*.go`, `apps/promotion/api/internal/logic/*.go`

---

# Promotion Business Rules

## 1. 优惠券建券与表单校验

**核心规则**：新建券一律落为 `draft` 草稿状态（`buildCoupon` 中硬编码 `Status: CouponStatusDraft`），必须再调 `CouponIssue` 才会开始发放。

`validateCouponForm` 按优惠类型分支校验：

| 优惠类型 | 校验规则 |
|----------|---------|
| `reduce` 满减 | `thresholdAmount > 0`；`0 < discountValue < thresholdAmount`（减的钱必须小于门槛） |
| `discount` 折扣 | `1 <= discountValue <= 99`（折扣百分比区间） |
| `no_threshold` 无门槛 | `discountValue > 0` |
| 其他 | 直接拒绝「优惠券类型非法」 |

通用校验：

| 规则 | 说明 |
|------|------|
| 名称非空 | `strings.TrimSpace(name)` 为空则拒绝 |
| 获取方式合法 | 必须是 `receive` / `exchange` / `assign` 之一 |
| 数量非负 | `totalNum` 与 `userLimit` 不能为负数 |
| 限定范围必填 | `specific=true` 时 `scopes` 不能为空 |
| scopes 序列化 | `specific=true` 时 `scopes` JSON 序列化后存入 `scopes` 字段；`specific=false` 时该字段留空 |

```
流程（CouponCreate）:
  1. validateCouponForm 校验表单
  2. buildCoupon 转 model，Status 固定 draft
  3. Insert → LastInsertId
  4. 回查 FindOne，返回 CouponDetailVO
```

## 2. 优惠券状态机

**核心规则**：状态流转由各 logic 前置校验硬性约束，非法流转返回 `Conflict`。

```
draft ──CouponIssue──> issued ──CouponPause──> paused ──CouponIssue──> issued
                                                  │
draft/paused ──CouponDelete──> deleted=1          │
                          ended（终态，不可再发放）
```

| 规则 | 说明 |
|------|------|
| 已结束不可再发放 | `status == ended` 时 `CouponIssue` 返回「优惠券已结束，无法再次发放」 |
| 只有发放中可暂停 | `status != issued` 时 `CouponPause` 返回「只有发放中的优惠券才能暂停」 |
| 只有草稿/暂停可删 | `SoftDelete` SQL 带 `status in ('draft','paused')` 条件，受影响行数为 0 即返回「优惠券不存在或当前状态不允许删除」 |
| 已删除按不存在处理 | `CouponGet` / `CouponPause` / `CouponIssue` 均在查到后判断 `Deleted == 1` 并返回 `NotFound` |
| 暂停不影响已领券 | 暂停仅改券模板状态，用户已持有的 `user_coupon` 不受影响 |

## 3. 发放与有效期设置（CouponIssue）

**核心规则**：发放期与有效期分别校验，有效期采用「绝对区间」或「相对天数」二选一。

| 规则 | 说明 |
|------|------|
| 时间格式兼容 | `parseTime` 依次尝试 `2006-01-02 15:04:05` / `2006-01-02T15:04:05` / RFC3339 / `2006-01-02`，全失败返回「时间格式非法」 |
| 发放期顺序 | `issueEndTime` 早于 `issueBeginTime` 时拒绝 |
| 有效期顺序 | `termEndTime` 早于 `termBeginTime` 时拒绝 |
| 有效期二选一 | `termDays <= 0` 且 `termEndTime` 为空时拒绝「请设置有效期天数或使用结束时间」 |
| 立即发放 | 未指定 `issueBeginTime` 时自动填 `now()` |
| 兑换码生成时机 | 仅在**首次发放**（原状态为 `draft`）且 `obtainWay == exchange` 且 `totalNum > 0` 时批量生成兑换码；再次发放（paused → issued）不会重复生成 |

## 4. 兑换码生成算法

**核心规则**：`generateCodes` 使用密码学安全随机源生成，字符集剔除易混淆字符。

| 参数 | 值 | 说明 |
|------|----|------|
| `codeAlphabet` | `23456789ABCDEFGHJKLMNPQRSTUVWXYZ` | 32 个字符，剔除 `0/O/1/I` |
| `codeLength` | `12` | 兑换码固定 12 位 |
| 随机源 | `crypto/rand.Int` | 非 `math/rand`，防止可预测 |

```
流程（generateCodes(n)）:
  1. n <= 0 直接返回 nil
  2. 循环直到收集满 n 个码：
     2.1 逐位从 codeAlphabet 中 crypto/rand 取字符，拼成 12 位
     2.2 命中 seen 集合（本批次内重复）则丢弃重来
     2.3 否则加入 seen 与结果集
  3. BatchInsert 一条 insert 多 values 写入 coupon_code
```

> 内存 `seen` 只能保证**单批次**不重复，跨批次靠 `coupon_code.uk_code` 唯一索引兜底。

## 5. 领取防超发（UserCouponReceive）

**核心规则**：库存扣减完全依赖 SQL 条件更新，**先抢库存再写用户券**，写失败则补偿回滚库存。

前置校验链：

| 顺序 | 校验 | 失败返回 |
|------|------|---------|
| 1 | `userId > 0` | `Unauthorized` |
| 2 | `id > 0` | 「优惠券 id 非法」 |
| 3 | 券存在且 `Deleted != 1` | `NotFound`「优惠券不存在」 |
| 4 | `obtainWay == receive` | `Conflict`「该优惠券不支持手动领取」 |
| 5 | `couponReceivable(coupon, now)` | `Conflict`「优惠券不在领取时间内或已被领完」 |
| 6 | `userLimit > 0` 时 `CountByUserAndCoupon < userLimit` | `Conflict`「已达到该优惠券的领取上限」 |

`couponReceivable` 的四项判定：`status == issued`、`now >= issueBeginTime`（若设置）、`now <= issueEndTime`（若设置）、`totalNum == 0 || issueNum < totalNum`。

**防超发关键 SQL**（`IncrIssueNum`）：

```sql
update coupon set issue_num = issue_num + 1, update_time = now()
where id = ? and deleted = 0 and status = 'issued'
  and (total_num = 0 or issue_num < total_num)
```

| 机制 | 说明 |
|------|------|
| 原子性来源 | InnoDB 行锁 + `issue_num < total_num` 条件，并发下只有满足条件的更新才会成功 |
| 抢空判定 | `RowsAffected == 0` 即判定「优惠券已被领完」，不做重试 |
| 库存回滚 | `user_coupon` 落库失败时调 `DecrIssueNum` 补偿；`DecrIssueNum` 带 `issue_num > 0` 条件，保证不会减成负数 |
| 回滚失败处理 | 回滚本身失败只 `l.Errorf` 记录，不阻断主流程错误返回 |

> **限领校验非强一致**：代码注释明确「非强一致，最终由库存与业务容忍度兜底」——`CountByUserAndCoupon` 是独立的 select，与后续 insert 之间存在并发窗口，极端并发下同一用户可能超出 `userLimit` 一张。

**过期时间计算**（`userCouponExpireTime`）：优先取券模板的绝对 `termEndTime`；否则按 `now + termDays` 天，并把末尾对齐到当天 `23:59:59`（避免按秒卡点）；两者都没有则留空（永不过期）。

## 6. 兑换码兑换（UserCouponExchange）

**核心规则**：兑换码核销与库存扣减均通过条件更新保证并发安全。

```
流程（UserCouponExchange）:
  1. userId > 0，code 归一化：TrimSpace + ToUpper，空则拒绝
  2. FindOneByCode 查码 → 不存在 / Deleted==1 → NotFound「兑换码不存在」
  3. status != unused → Conflict「兑换码已被使用」
  4. expireTime 已过 → Conflict「兑换码已过期」
  5. 查关联券 → 不存在 / Deleted==1 / status==ended → 「优惠券活动已结束」
  6. userLimit > 0 时校验 CountByUserAndCoupon
  7. MarkUsed 条件更新核销码（rows==0 → Conflict「兑换码已被使用」）
  8. IncrIssueNum 扣库存
  9. Insert user_coupon，code 字段回填兑换码
```

**并发安全关键 SQL**（`MarkUsed`）：

```sql
update coupon_code set status = 'used', user_id = ?, update_time = now()
where id = ? and status = 'unused' and deleted = 0
```

| 机制 | 说明 |
|------|------|
| 抢兑判定 | `RowsAffected == 0` 表示已被他人抢先兑换，返回 `Conflict` |
| 双重校验 | 步骤 3 的读校验只是快速失败，真正的并发防线是步骤 7 的条件更新 |
| 缓存失效 | `MarkUsed` 先 `FindOne` 拿到 `code`，更新后同时失效 `cache:couponCode:id:` 与 `cache:couponCode:code:` 两个键 |
| 库存扣减容错 | 第 8 步 `IncrIssueNum` 失败只 `l.Errorf` 记录，**不中断兑换**（码已核销，优先保证用户拿到券） |
| 用户券落库失败 | 返回 `Internal`「兑换失败」，此处**未回滚**已核销的兑换码与已扣减的库存 |

> 与领取流程不同，兑换流程没有对 `couponReceivable` 做发放期/库存判定——兑换码本身即库存凭证，只校验券未删除、未结束。

## 7. 核销幂等（UserCouponUse）

**核心规则**：状态流转 `unused → used`，靠带源状态条件的批量更新保证同一张券不会被重复核销。

```sql
update user_coupon
set status = 'used', use_time = now(), order_id = ?, update_time = now()
where id in (...) and user_id = ? and status = 'unused' and deleted = 0
```

| 规则 | 说明 |
|------|------|
| 归属校验 | `FindByIdsAndUser` 查回的条数必须等于入参 `ids` 长度，否则「优惠券不存在或不属于当前用户」 |
| 幂等保证 | `where status = 'unused'` 使重复调用的第二次更新影响 0 行 |
| 全成功语义 | `RowsAffected != len(ids)` 即判定「部分优惠券已被使用或已过期」并返回 `Conflict`——要么全部核销成功，要么整体报错 |
| 统计字段容错 | 逐张券 `AddUsedNum(+1)` 失败只 `l.Errorf`，不阻断主流程（注释：统计类字段，可离线校准） |
| 订单绑定 | 核销时把 `orderId` 写入 `user_coupon.order_id`，`use_time` 置 `now()` |

> **注意**：`RowsAffected != len(ids)` 判定发生在更新**之后**，此时已成功的那部分行不会回滚（无事务包裹），实际是「部分成功 + 报错」的语义。

## 8. 退还（UserCouponRefund）

**核心规则**：状态流转 `used → unused`，走同一个 `UpdateStatusByIds`，但走 default 分支清空使用痕迹。

```sql
update user_coupon
set status = 'unused', use_time = null, order_id = null, update_time = now()
where id in (...) and user_id = ? and status = 'used' and deleted = 0
```

| 规则 | 说明 |
|------|------|
| 归属校验 | `FindByIdsAndUser` 结果为空则「优惠券不存在或不属于当前用户」（与核销不同，此处只要求非空，不要求条数相等） |
| 幂等保证 | `where status = 'used'`，`RowsAffected == 0` 返回「优惠券未被使用，无需退还」 |
| 使用痕迹清理 | `use_time` 与 `order_id` 一并置 `null` |
| used_num 回滚 | 对原状态为 `used` 的券逐张 `AddUsedNum(-1)`；SQL 带 `used_num + delta >= 0` 条件防减成负数，失败只记日志 |
| 过期券退还 | 已过期的券仍会被恢复为 `unused`，仅打一条 `l.Infof` 日志——恢复状态但实际不可再用 |

## 9. 券组合穷举算法（`rpc/internal/logic/solution.go`）

这是 promotion 服务的核心算法，服务于 `UserCouponAvailable`（列所有方案）与 `UserCouponDiscount`（算指定方案）。

### 9.1 可用券筛选（`usableUserCoupons`）

逐张过滤，任一条不满足即剔除：

| 条件 | 说明 |
|------|------|
| `uc.Status == unused` | 已使用/已过期券不参与 |
| 未过期 | `uc.ExpireTime` 有效且 `now` 已超过则剔除 |
| 券规则存在 | `coupons[uc.CouponId]` 能查到且 `Deleted != 1` |
| 有效期已开始 | `c.TermBeginTime` 有效且 `now` 早于它则剔除 |

### 9.2 组合数量截断（`trimCoupons`）

| 常量 | 值 | 作用 |
|------|----|------|
| `maxCombineCoupons` | `12` | 参与组合运算的最大券数，组合方案数为 2^n，超出会拖垮接口 |
| `maxSolutions` | `30` | 最多返回的方案数 |

超过 12 张时，先对每张券单独调 `buildSolution` 算出「单券优惠力度」作为打分，`sort.SliceStable` 按分值降序，只保留前 12 张。无法生效的券打分为 0，自然排到末尾被淘汰。

### 9.3 位掩码穷举（`calcSolutions`）

```
流程（calcSolutions）:
  1. available 为空或 courses 为空 → 返回 nil
  2. trimCoupons 截断到 <= 12 张
  3. for mask := 1; mask < (1<<n); mask++    // 遍历所有非空子集
     3.1 按 mask 的位组装 combo
     3.2 buildSolution(combo, courses)，非法组合直接跳过
     3.3 solutionKey(ids) 去重（券 id 排序后拼串）
     3.4 加入 solutions
  4. 排序：优惠金额降序；金额相同时用券数量少的排前
  5. 截断到前 maxSolutions(30) 条
```

复杂度为 `O(2^n × n × m)`（n=券数≤12，m=课程数），最坏 4095 个子集。

### 9.4 单方案计算（`buildSolution`）

**核心规则**：券按数组顺序**依次作用于「剩余金额」**，避免同一笔钱被重复打折。

```
流程（buildSolution）:
  1. remaining[courseId] = 课程原价（累加同 id）
  2. 逐张券：
     2.1 matchedCourses：specific=0 全场通用；否则按 cateId 匹配 scopes
     2.2 matched 为空 → 组合非法，返回 false
     2.3 subtotal = 匹配课程的 remaining 之和
     2.4 discount = calcDiscount(券, subtotal)
     2.5 discount <= 0 → 组合非法，返回 false
     2.6 按 remaining 比例把 discount 分摊到各课程
  3. total <= 0 → 返回 false
  4. 返回 { discountAmount, discountDetail, ids, rules }
```

| 规则 | 说明 |
|------|------|
| 组合合法性 | **组合中任一张券无法生效（未达门槛或无适用课程）时整个组合作废**，不做部分生效 |
| 顺序敏感 | 券作用于剩余金额，因此不同顺序理论上结果不同；实现中按 `available` 数组顺序作用 |
| 比例分摊 | `share = discount × remaining[courseId] / subtotal` |
| 取整误差处理 | **最后一门课**用 `discount - allocated` 兜底吸收整数除法的余数 |
| 分摊封顶 | `share` 不超过该课程的 `remaining`，防止分摊为负余额 |

### 9.5 单券优惠计算（`calcDiscount`，`common.go`）

| 券类型 | 计算方式 |
|--------|---------|
| 门槛判定 | 非 `no_threshold` 券，`totalAmount < thresholdAmount` 时返回 0 |
| `discount` 折扣 | `discountValue` 落在 `[1,99]` 之外返回 0；否则 `discount = totalAmount × (100 - discountValue) / 100`，再受 `maxDiscountAmount > 0` 封顶 |
| `reduce` / `no_threshold` | `discount = discountValue`（固定金额） |
| 全局封顶 | 优惠额不超过 `totalAmount` |
| 边界 | `totalAmount <= 0` 直接返回 0 |

### 9.6 规则文案（`couponRule`）

| 券类型 | 文案模板 |
|--------|---------|
| `no_threshold` | 无门槛立减{X}元 |
| `discount` | 满{门槛}元打{discountValue×10 转元}折（`maxDiscountAmount > 0` 时追加「，最多减{X}元」） |
| `reduce`（默认分支） | 满{门槛}元减{X}元 |

`yuan` 函数把分转元：`fmt.Sprintf("%.2f", cent/100)` 后去掉尾随 0 与小数点。

## 10. 折扣方案复算（UserCouponDiscount）

**核心规则**：仅接受属于当前用户且未使用的券，防止越权占用他人优惠券。

| 顺序 | 校验 | 失败返回 |
|------|------|---------|
| 1 | `userId > 0` | `Unauthorized` |
| 2 | `courseList` 非空 | 「订单课程不能为空」 |
| 3 | `userCouponIds` 为空 | 返回空 `CouponDiscountDTO`（不报错） |
| 4 | `FindByIdsAndUser` 条数 == 入参条数 | 「优惠券不存在或不属于当前用户」 |
| 5 | `usableUserCoupons` 过滤后条数 == 入参条数 | `Conflict`「存在已使用或已过期的优惠券，请重新选择」 |
| 6 | `buildSolution` 返回 ok | `Conflict`「所选优惠券不满足使用条件」 |

## 11. C 端券列表标记（CouponList）

| 规则 | 说明 |
|------|------|
| 只展示可领券 | 遍历 `FindList` 结果，`couponReceivable == false` 的直接跳过 |
| 批量查已领 | 一次 `FindByUserAndStatus(userId, "")` 拉全部用户券，在内存里按 `couponId` 计数，避免逐张券查库 |
| `received` 标记 | 该券已领数量 `> 0` |
| `available` 标记 | `userLimit == 0 || got < userLimit`——达到限领数量后不可再领 |
| 未登录兼容 | `userId <= 0` 时不查用户券，全部券 `received=false`、`available=true` |

## 12. 我的券分页（UserCouponPage）

| 规则 | 说明 |
|------|------|
| 分页归一化 | `page.Normalize(pageNo, pageSize)`：pageNo 最小 1，pageSize 缺省 10、最大 100 |
| 批量加载券规则 | 收集 `couponId` 后一次 `FindByIds`，避免 N+1 |
| 到期时间覆盖 | 展示 `user_coupon.expire_time`（用户券实际到期时间），而非券模板的 `term_end_time` |
| 过期自动置灰 | `now` 超过 `expireTime` 时把 `available` 置为 `false`（不改 DB 状态） |
| 券规则缺失跳过 | `coupons[uc.CouponId]` 查不到的用户券直接跳过，不计入 list（但 `total` 已按 DB count 返回） |

## 13. API 层职责

**核心规则**：`promotion-api` 是纯透传层，全部业务判断下沉到 RPC。

| 规则 | 说明 |
|------|------|
| 身份注入 | 需要用户身份的接口统一 `auth.UserIdFromCtx(l.ctx)` 取 JWT 中的 userId，再塞进 RPC 请求 |
| 类型转换 | `convert.go` 承担 pb ↔ types 转换；`DiscountDetail` 的 int64 键在 API 层转成 string 键（JSON 兼容） |
| 分页页数计算 | API 层自行算 `pages = (total + limit - 1) / limit`，`pageSize <= 0` 时按 10 兜底 |
| 管理端接口 | `CouponCreate` / `CouponIssue` / `CouponPage` / `CouponCodePage` 不取 userId，直接透传 |

> `CouponCreate` 的 RPC 实现中 `buildCoupon(in, 0)` 的 operator 硬编码为 `0`，即新建券的 `creater` / `updater` 恒为 0，未接入操作人身份。

## 状态说明

### 优惠券状态

| 值 | 含义 | 可执行操作 |
|----|------|-----------|
| `draft` | 草稿（待发放） | CouponIssue、CouponDelete |
| `issued` | 已发放 | CouponPause、可被领取/兑换 |
| `paused` | 暂停发放 | CouponIssue（重新发放）、CouponDelete |
| `ended` | 已结束 | 终态，不可再发放 |

### 用户券状态

| 值 | 含义 | 流转 |
|----|------|------|
| `unused` | 未使用 | → `used`（UserCouponUse） |
| `used` | 已使用 | → `unused`（UserCouponRefund） |
| `expired` | 已过期 | 常量已定义，当前逻辑未主动写入（过期通过 `expire_time` 与 `now` 比较判定） |

### 兑换码状态

| 值 | 含义 | 流转 |
|----|------|------|
| `unused` | 未兑换 | → `used`（MarkUsed，单向） |
| `used` | 已兑换 | 终态 |

## 未实现说明

promotion 服务 api 层与 rpc 层的全部 logic 均已实现，**未发现** `todo: add your logic` 形式的 goctl 占位。
