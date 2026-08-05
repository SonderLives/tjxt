> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/search/rpc/search.proto`

---

# Search RPC Spec

## 服务名

`Search` — 全文检索与课程推荐微服务，通过 etcd 服务发现（key: `search.rpc`），监听 `0.0.0.0:8090`。

> 当前 proto 中**只定义了用户兴趣（interests）相关的 2 个方法**，检索/推荐类方法尚未定义。

## RPC 方法总览

### 用户兴趣

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SaveInterests` | `SaveInterestsReq { id, interests }` | `Empty {}` | 保存用户兴趣（新增或覆盖） |
| `GetInterests` | `IdReq { id }` | `InterestsVO { id, interests, createTime, updateTime }` | 按用户 ID 查询兴趣 |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 主键，对应用户 id（`interests` 表主键即用户 id） |
| `interests` | string | 感兴趣的二级分类 id，以逗号分隔，例如：`120,220,330` |

**响应字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 用户 id |
| `interests` | string | 逗号分隔的二级分类 id 串 |
| `createTime` | string | 创建时间 |
| `updateTime` | string | 更新时间 |

> `SaveInterests` 返回 `Empty`，不回传 ID——因为主键即调用方传入的用户 ID，无需服务端生成。

---

## 被哪些 API 服务消费

| 消费方 | 调用方式 | 说明 |
|--------|---------|------|
| `search-api` (自身 API 层) | HTTP Handler → `searchclient.Search` RPC | `apps/search/api/internal/svc/servicecontext.go` 注入 `SearchRpc`，2 个 HTTP 接口全部指向自身 RPC |

（注：全仓 Grep `apps/search/rpc/search` / `searchclient` / `SearchRpc`，除 `apps/search/` 自身外**无其它服务引用**。`.cursor/repowiki/00-architecture/service-topology.md` 中标注的 `search-rpc --> AuthRPC` 依赖，在代码中不存在——`apps/search/rpc/internal/config/config.go` 未定义任何 RPC 客户端配置项。）

---

## 调用典型场景

1. **新用户选兴趣** → 注册引导页勾选二级分类 → `search-api` 从 JWT 取 userId → 调 `SaveInterests(id=userId, interests="120,220,330")`
2. **回显已选兴趣** → 用户进入偏好设置页 → 调 `GetInterests(id=userId)` → 前端按逗号切分 id 串并回勾
3. **个性化推荐（设计意图，proto 未定义）** → 取出用户 `interests` 中的二级分类 id → 交由 course 域按分类查课程 TOP10（对应 `docs/tjxt.openapi.json` 中的 `GET /interests/{id}/courses`）

> 📋 `docs/tjxt.openapi.json` 记录了 `GET /interests/{id}/courses`（根据二级分类 id 查询课程 TOP10），但 `apps/search/api/search.api` 与 `apps/search/rpc/search.proto` **均未定义该接口**，属于设计已规划、代码未落地的部分。

---

## 自定义 Model 方法

`apps/search/rpc/internal/model/interestsmodel.go` **当前为 goctl 生成的空壳**，接口体内仅内嵌 `interestsModel`，未追加任何扩展方法：

```go
InterestsModel interface {
    interestsModel
}
```

因此可用方法仅为 goctl 生成的 4 个基础方法：

- `Insert(ctx, data *Interests)` — 插入
- `FindOne(ctx, id int64)` — 按主键（用户 ID）查询，带缓存
- `Update(ctx, data *Interests)` — 全字段更新
- `Delete(ctx, id int64)` — 按主键删除

> `SaveInterests` 语义上是「存在则更新，不存在则插入」，但 model 层**没有 upsert 方法**；实现时需在 logic 中先 `FindOne`，按 `model.ErrNotFound` 分支决定走 `Insert` 还是 `Update`，或在自定义 model 中补写 `Upsert`。
