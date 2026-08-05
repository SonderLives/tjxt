> 版本：v1.0 | 更新：2026-08-05 | 来源：`sql/ddl/tj_search.sql`, `apps/search/rpc/internal/model/`, `apps/search/rpc/etc/search.yaml`

---

# Search Data Model

## ⚠️ 存储现状说明

search 服务在架构文档中被定位为「搜索推荐 / 兴趣标签 / ES 索引 / 推荐」（见 `.cursor/repowiki/00-architecture/overview.md`），但**代码与配置的实际情况与「用 Elasticsearch」的预期不符**，据实记录如下：

| 项 | 实际情况 | 核查方式 |
|----|---------|---------|
| Elasticsearch 客户端依赖 | **不存在** | 全仓 Grep `elastic` / `elasticsearch` / `olivere`（`*.go` / `*.yaml` / `*.mod`）零命中 |
| ES 连接配置 | **不存在** | `apps/search/rpc/internal/config/config.go` 仅有 `zrpc.RpcServerConf` + `DataSource` + `Cache` |
| MySQL 库 | **存在**，`tj_search` | `apps/search/rpc/etc/search.yaml` 的 `DataSource` |
| 关系型表 | **存在 1 张**，`interests` | `sql/ddl/tj_search.sql` |
| Redis | 存在，仅作 model 层缓存 | `search.yaml` 的 `Cache` 段，`sqlc.CachedConn` |

**结论**：截至当前提交，search 服务的数据来源是 **MySQL `tj_search` 库的单表 `interests`**（配合 Redis 做主键缓存），**没有任何 ES 索引结构落地**。下文「索引结构（推导）」一节为按 proto/DDL 语义反推的设计意图，**不代表现有代码**。

---

## 模型文件清单

| 自定义 Model 文件 | 对应表 | 扩展方法 |
|------------------|--------|---------|
| `interestsmodel.go` | `interests` | 无（仅 goctl 空壳） |

所有 `*_gen.go` 文件均由 goctl 生成，**禁止修改**。扩展方法统一放在同名自定义 `.go` 文件中（如 `interestsmodel.go` 扩展 `InterestsModel` 接口）——**当前该文件为空壳，尚未追加任何方法**。

Model 在 `apps/search/rpc/internal/svc/servicecontext.go` 中通过 `sqlx.NewMysql(c.DataSource)` + `c.Cache` 注入，走带缓存的 `sqlc.CachedConn`。

**缓存 key 前缀**（`interestsmodel_gen.go` 中定义）：

| Model | 缓存前缀 |
|-------|---------|
| `Interests` | `cache:interests:id:` |

---

## 表清单与字段说明

### 1. `interests` — 用户兴趣表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 主键，对应用户 id | PK |
| `interests` | varchar(255) | 感兴趣的二级分类 id，以逗号分隔，例如：`120,220,330`（可空） | - |
| `create_time` | datetime | 创建时间 | 自动填充 |
| `update_time` | datetime | 更新时间 | 自动更新 |

> - 表注释：「用户兴趣表，保存感兴趣的二级分类 id」
> - 字符集 `utf8mb4` / `utf8mb4_0900_ai_ci`
> - **主键即用户 ID**，非自增，由调用方传入；因此一个用户至多一行兴趣记录
> - Go 结构体 `Interests` 中 `interests` 映射为 `sql.NullString`
> - `interests` 为**逗号分隔的字符串**而非关联表，无法按分类 id 反查用户；若需「某分类下有哪些用户」，当前表结构不支持
> - 无 `deleted` 字段，删除为物理删除；无 `creater`/`updater` 审计字段

---

## 关系图

```
interests (1) ─── (1) user      （id 即 user 域的用户 id，跨库无外键）
                    │
                    └── interests 字段内的二级分类 id 串
                          → 逻辑指向 course 域的分类表（跨库、跨服务，无外键、无 RPC 调用）
```

> `tj_search` 库中**只有一张表**，库内无任何表间关系。

---

## 索引结构（📋 设计意图，推导，当前未落地）

以下内容由 `apps/search/rpc/search.proto`、`sql/ddl/tj_search.sql` 与 `docs/tjxt.openapi.json` 中 `GET /interests/{id}/courses`（根据二级分类 id 查询课程 TOP10）反推，**代码中不存在对应实现**，标注为推导：

| 推导项 | 内容 |
|--------|------|
| 检索对象 | 课程（course 域），按二级分类 id、名称、简介建全文索引 |
| 推荐入口 | 以 `interests.interests` 中的二级分类 id 集合作为召回条件 |
| 排序维度 | 按报名人数 / 销量取 TOP10（对应 openapi 摘要「课程 TOP10」） |
| 索引更新 | 需 course 域在课程发布/下架时同步索引；当前 `.cursor/repowiki/03-shared/mq-events.md` 中**未定义**任何 search 域事件 |
| 与现状差距 | 需引入 ES 客户端依赖、在 `config.go` 增加 ES 配置段、在 `ServiceContext` 注入 ES client、在 proto 增加检索方法 |

> 在上述能力落地前，search 服务的实际职责仅为「用户兴趣标签的存取」。

---

## 模型扩展模式

所有 `*_gen.go` 由 goctl 自动生成，不可编辑。扩展方法统一放在同名的自定义 `.go` 文件中：

```
interestsmodel_gen.go     ← goctl 生成，只读（Insert / FindOne / Update / Delete）
interestsmodel.go         ← 手写扩展位（当前为空壳，无任何扩展方法）
```

当前项目自定义 Model 模式：

- `interestsmodel.go` — 无扩展

`vars.go` 仅定义 `ErrNotFound = sqlx.ErrNotFound`，**无** auth 域那样的 `sqlhelper.go` 通用 SQL 工具函数。

📋 **待补齐（设计意图）**：`SaveInterests` 的「存在即更新」语义需要一个 `Upsert(ctx, id, interests)` 方法；若后续要支持按分类反查用户，还需将逗号串拆为关联表 `user_interests(user_id, category_id)`。
