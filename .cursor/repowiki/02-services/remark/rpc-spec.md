> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/remark/rpc/remark.proto`

---

# Remark RPC Spec

## 服务名

`Remark` — 评论互动域的点赞微服务，通过 etcd 服务发现（key: `remark.rpc`）。以 `bizType + bizId` 泛化承载各业务对象（回复、笔记、问答等）的点赞关系。

## RPC 方法总览

### 点赞

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `Like` | `LikeReq { userId, bizId, bizType, liked }` | `Empty {}` | 点赞或取消点赞，同一 (用户,业务) 维度幂等 upsert |
| `QueryLikedBizIds` | `LikedReq { userId, bizType, bizIds }` | `LikedResp { likedBizIds }` | 批量查询当前用户在一组业务 id 中已点赞的子集 |

**请求字段说明**（`LikeReq`）：

| 字段 | 类型 | 说明 |
|------|------|------|
| `userId` | int64 | 点赞人 ID（来自 JWT） |
| `bizId` | int64 | 被点赞的业务对象 ID |
| `bizType` | string | 业务类型，如 `reply` / `note` / `question` |
| `liked` | bool | true=点赞，false=取消点赞；服务端转成 `liked` 列的 1/0 |

**请求字段说明**（`LikedReq`）：

| 字段 | 类型 | 说明 |
|------|------|------|
| `userId` | int64 | 当前用户 ID（来自 JWT） |
| `bizType` | string | 业务类型 |
| `bizIds` | repeated int64 | 待查询的业务 id 列表，为空时直接返回空结果 |

**响应字段说明**（`LikedResp`）：

| 字段 | 类型 | 说明 |
|------|------|------|
| `likedBizIds` | repeated int64 | 入参 `bizIds` 中 `liked = 1` 的子集，未点赞与已取消的均不返回 |

---

## 被哪些 API 服务消费

| 消费方 | 调用方式 | 说明 |
|--------|---------|------|
| `remark-api` (自身 API 层) | `apps/remark/api/internal/svc/servicecontext.go` 中 `remarkclient "tjxt/apps/remark/rpc/client/remark"` → `RemarkRpc` | 两个 HTTP handler（`/likes`、`/likes/list`）均透传到本 RPC |

（注：截至当前版本，仅 remark-api 在 `servicecontext.go` 中注入 `remarkclient`。其他服务如 learning 的笔记/问答列表若需回填「我是否点赞」，可复用 `QueryLikedBizIds`，但目前尚未接入。）

---

## 调用典型场景

1. **用户点赞** → 前端 POST `/likes` 带 `bizId` / `bizType` / `liked=true` → API 层 `auth.UserIdFromCtx` 取 userId → 调 `Like` → 无记录则 insert，有记录且状态不同则 update
2. **取消点赞** → 同一接口传 `liked=false` → 走 update 把 `liked` 列改成 0（**软取消**，不删行）
3. **列表页回填点赞态** → 业务列表（如回复列表）拿到一批 `bizId` → 一次调 `QueryLikedBizIds` 批量查询 → 前端按返回集合高亮已赞项，避免逐条查询
4. **重复点击容错** → 用户连点两次点赞，第二次 `existing.Liked == liked` 直接返回，不产生任何写库

---

## 自定义 Model 方法

`likerecordmodel.go` 扩展了：
- `FindLikedBizIds(ctx, userId, bizType, bizIds)` — 批量查询某用户在一组 `bizId` 中 `liked=1` 的集合

该方法刻意走 `CachedConn.QueryRowsNoCacheCtx` 不带缓存：注释说明「命中行数通常远小于入参 bizIds，且业务对实时性要求高，走 NoCache 查询避免缓存击穿与脏读」。

goctl 生成的基础方法（`likerecordmodel_gen.go`，只读）：
- `Insert(ctx, data)` / `FindOne(ctx, id)` / `Update(ctx, data)` / `Delete(ctx, id)`
- `FindOneByUserIdBizIdBizType(ctx, userId, bizId, bizType)` — 由 `uk_user_biz` 唯一索引生成，带 `cache:likeRecord:userId:bizId:bizType:` 索引缓存，是 `Like` 逻辑判定「是否已有记录」的入口
