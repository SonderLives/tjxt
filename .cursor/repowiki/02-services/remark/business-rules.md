> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/remark/rpc/internal/logic/remark/*.go`, `apps/remark/api/internal/logic/like/*.go`

---

# Remark Business Rules

## 1. 点赞 upsert（Like）

**核心规则**：同一 `(user_id, biz_id, biz_type)` 三元组在 `like_record` 中最多一行，点赞与取消点赞都落在同一行上，通过 `liked` 列的 1/0 表达状态。

| 规则 | 说明 |
|------|------|
| bool → tinyint 转换 | `in.Liked == true` 转 `liked = 1`，否则 `liked = 0` |
| 记录不存在则新增 | `FindOneByUserIdBizIdBizType` 返回 `sqlc.ErrNotFound` 时执行 `Insert` |
| 状态相同则空转 | `existing.Liked == liked` 时直接返回 `Empty{}`，**不发起任何写库** |
| 状态不同则更新 | 只改 `existing.Liked` 后 `Update`，`user_id` / `biz_id` / `biz_type` 保持不变 |
| 软取消 | 取消点赞不删除行，只把 `liked` 置 0；重新点赞再置回 1 |

```
流程（Like）:
  1. liked := 0；in.Liked 为 true 时 liked = 1
  2. FindOneByUserIdBizIdBizType(userId, bizId, bizType)
     2.1 err == sqlc.ErrNotFound → Insert{UserId, BizId, BizType, Liked} → 返回
     2.2 err != nil             → 直接返回错误
  3. existing.Liked == liked → 直接返回 Empty（幂等短路）
  4. existing.Liked = liked → Update → 返回
```

**幂等性分析**：

| 场景 | 行为 |
|------|------|
| 首次点赞 | Insert 一行，`liked=1` |
| 重复点赞 | 第 3 步短路，无写库，返回成功 |
| 点赞 → 取消 | Update 把 `liked` 改 0 |
| 重复取消 | 第 3 步短路，无写库，返回成功 |
| 取消 → 重新点赞 | Update 把 `liked` 改回 1 |

> **并发写入的兜底在 DB**：`FindOneByUserIdBizIdBizType`（读）与 `Insert`（写）之间存在并发窗口，两个并发的首次点赞请求可能都判定为「记录不存在」并同时 Insert。此时 `uk_user_biz` 唯一索引会让其中一条 Insert 报错，错误直接返回给调用方——代码层面**未做重复键的捕获与降级**。

> **错误判定用 `sqlc.ErrNotFound`**：`likelogic.go` 中比较的是 `sqlc.ErrNotFound`，而 `likerecordmodel_gen.go` 的 `FindOneByUserIdBizIdBizType` 在未命中时实际返回的是 `model.ErrNotFound`。两者在 `vars.go` 中通过 `var ErrNotFound = sqlx.ErrNotFound` 指向同一个 error 值，因此判定成立。

## 2. 批量点赞状态查询（QueryLikedBizIds）

**核心规则**：一次查询回填整页列表的点赞态，避免逐条查询。

| 规则 | 说明 |
|------|------|
| 纯透传 | RPC logic 不做任何校验，直接调 `FindLikedBizIds` 并把结果包进 `LikedResp` |
| 空入参短路 | Model 层 `len(bizIds) == 0` 时返回 `nil, nil`，不发 SQL |
| 只返回已赞子集 | SQL 带 `liked = 1`，取消点赞的 `bizId` 不出现在结果中 |
| 不走缓存 | `QueryRowsNoCacheCtx` 直读 DB，保证点赞后立即可见 |
| 结果为子集 | 返回的 `likedBizIds` 长度通常远小于入参 `bizIds` |

```
流程（QueryLikedBizIds）:
  1. FindLikedBizIds(userId, bizType, bizIds)
     1.1 bizIds 为空 → nil
     1.2 动态拼接 in (?,?,...) 占位符，参数化执行
     1.3 select biz_id where user_id=? and biz_type=? and liked=1 and biz_id in (...)
  2. 包装成 LikedResp{LikedBizIds: ids} 返回
```

## 3. API 层身份与参数处理

**核心规则**：`remark-api` 是纯透传层，用户身份统一从 JWT 上下文取，不信任请求体中的 userId（proto 中虽有 `userId` 字段，但 API 层永远用 `auth.UserIdFromCtx` 覆盖）。

### 3.1 `/likes`（LikeLogic）

| 规则 | 说明 |
|------|------|
| 身份注入 | `auth.UserIdFromCtx(l.ctx)` 取 userId，失败直接返回错误 |
| 参数透传 | `bizId` / `bizType` / `liked` 原样传给 RPC |
| 固定响应 | RPC 成功后返回 `OkVO{Success: true}`，不回传点赞后的状态 |
| 无参数校验 | API 层不校验 `bizId > 0`，也不校验 `bizType` 是否在允许枚举内 |

### 3.2 `/likes/list`（ListLikedLogic）

| 规则 | 说明 |
|------|------|
| 身份注入 | 同上，`auth.UserIdFromCtx` |
| 字符串转数值 | 请求中的 `bizIds` 是 `[]string`（query 参数），逐个 `strconv.ParseInt(s, 10, 64)` |
| 非法值静默丢弃 | 解析失败（`perr != nil`）或 `id <= 0` 的元素**直接跳过，不报错** |
| 全部非法的后果 | 转换后 `bizIds` 为空切片 → Model 层短路 → 返回空 `likedBizIds` |
| 结果透传 | 直接把 `r.LikedBizIds` 包进 `LikeListResp` 返回 |

```
流程（ListLiked）:
  1. userId := auth.UserIdFromCtx(ctx)
  2. bizIds := []int64{}
     for s := range req.BizIds:
         id, perr := strconv.ParseInt(s, 10, 64)
         if perr == nil && id > 0: append(bizIds, id)
  3. RemarkRpc.QueryLikedBizIds(userId, bizType, bizIds)
  4. 返回 LikeListResp{LikedBizIds: r.LikedBizIds}
```

## 4. 鉴权

| 规则 | 说明 |
|------|------|
| 全局 JWT | `remark.api` 声明 `@server(jwt: Auth, group: like)`，`routes.go` 中两个路由由 `rest.WithJwt(serverCtx.Config.Auth.AccessSecret)` 统一包裹 |
| 无匿名接口 | remark-api 不存在免鉴权路由，未带有效 JWT 的请求在中间件层即被拒绝 |
| 越权防护 | 所有查询与写入都以 JWT 中的 userId 为维度，用户无法读写他人的点赞记录 |

## 状态说明

### `liked` 列

| 值 | 含义 | 产生方式 |
|----|------|---------|
| 1 | 已点赞 | 首次 Insert（默认值即 1）或由 0 Update 而来 |
| 0 | 已取消 | 由 1 Update 而来；行仍保留，不删除 |

### `bizType` 取值

DDL 注释与 proto 注释中给出的示例为 `reply` / `note` / `question`。服务端**不做枚举校验**，取值由调用方业务侧约定——新增业务类型无需改动 remark 服务。

## 未实现说明

remark 服务 api 层（`LikeLogic` / `ListLikedLogic`）与 rpc 层（`LikeLogic` / `QueryLikedBizIdsLogic`）的全部 logic 均已实现，**未发现** `todo: add your logic` 形式的 goctl 占位。

已知能力边界（非占位，而是当前设计未覆盖）：
- 无点赞数聚合统计接口——`idx_biz(biz_type, biz_id)` 索引已为此预留，但尚无对应 RPC 方法
- 无点赞记录列表 / 分页查询接口
- 无点赞事件的消息投递（点赞数回写业务方需业务方自行统计）
