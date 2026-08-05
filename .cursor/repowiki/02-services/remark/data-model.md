> 版本：v1.0 | 更新：2026-08-05 | 来源：`sql/ddl/tj_remark.sql`

---

# Remark Data Model

## 模型文件清单

| 自定义 Model 文件 | 对应表 | 扩展方法 |
|------------------|--------|---------|
| `likerecordmodel.go` | `like_record` | FindLikedBizIds |

所有 `*_gen.go` 文件均由 goctl 生成，**禁止修改**。扩展方法统一放在同名自定义 `.go` 文件中（如 `likerecordmodel.go` 扩展 `LikeRecordModel` 接口）。

`vars.go` 定义 `ErrNotFound = sqlx.ErrNotFound`。remark 服务仅有一张表，无独立的 `consts.go`，业务类型 `bizType` 由调用方自行传入字符串（如 `reply` / `note` / `question`），服务端不做枚举校验。

---

## 表清单与字段说明

### 1. `like_record` — 点赞记录表

| 字段 | 类型 | 说明 | 索引 |
|------|------|------|------|
| `id` | bigint | 点赞记录 id，自增主键 | PK |
| `user_id` | bigint | 点赞人 id | `uk_user_biz` 首列 |
| `biz_id` | bigint | 被点赞业务 id | `uk_user_biz` 次列 / `idx_biz` 次列 |
| `biz_type` | varchar(32) | 点赞业务类型，例如 reply / note / question | `uk_user_biz` 末列 / `idx_biz` 首列 |
| `liked` | tinyint | 1:点赞 0:取消，默认 1 | 查询过滤条件 |
| `create_time` | datetime | 创建时间 | 自动填充 |
| `update_time` | datetime | 更新时间 | 自动更新 |

**索引说明**：

| 索引名 | 类型 | 列 | 用途 |
|--------|------|----|------|
| `PRIMARY` | 主键 | `id` | 主键查询，对应 `cache:likeRecord:id:` 缓存 |
| `uk_user_biz` | 唯一 | `user_id`, `biz_id`, `biz_type` | 保证同一用户对同一业务对象只有一行记录；`FindOneByUserIdBizIdBizType` 与 `cache:likeRecord:userId:bizId:bizType:` 索引缓存基于此建立 |
| `idx_biz` | 普通 | `biz_type`, `biz_id` | 按业务对象维度统计点赞（如某条回复的点赞数） |

> **无逻辑删除字段**：与 auth / promotion 等域不同，`like_record` 没有 `deleted` 列，也没有 `creater` / `updater`。取消点赞采用**软取消**语义——把 `liked` 置为 0 而非删除行，因此唯一索引不会被反复占用与释放，重新点赞只需再改回 1。

---

## 关系图

```
like_record.user_id → user 域 user.id（无外键约束，仅逻辑关联）

like_record.(biz_type, biz_id) → 多态关联，指向不同业务域的对象：
    biz_type = 'reply'    → 回复
    biz_type = 'note'     → 笔记
    biz_type = 'question' → 问答
```

`biz_type` + `biz_id` 构成多态外键，remark 服务不感知具体业务表结构，业务侧自行约定 `bizType` 取值——这使得点赞能力可以在不改表的前提下被任意新业务复用。

---

## 模型扩展模式

所有 `*_gen.go` 由 goctl 自动生成，不可编辑。扩展方法统一放在同名的自定义 `.go` 文件中：

```
likerecordmodel_gen.go    ← goctl 生成，只读
likerecordmodel.go        ← 手写扩展 FindLikedBizIds
```

当前项目自定义 Model 模式：
- `likerecordmodel.go` — FindLikedBizIds

**goctl 生成的基础方法**（`likeRecordModel` 接口）：

| 方法 | 缓存策略 |
|------|---------|
| `Insert(ctx, data)` | 写后失效 `cache:likeRecord:id:` 与 `cache:likeRecord:userId:bizId:bizType:` |
| `FindOne(ctx, id)` | 走 `cache:likeRecord:id:` 主键缓存 |
| `FindOneByUserIdBizIdBizType(ctx, userId, bizId, bizType)` | 走 `cache:likeRecord:userId:bizId:bizType:` 索引缓存 → 回主键查询 |
| `Update(ctx, data)` | 写后失效上述两个键 |
| `Delete(ctx, id)` | 写后失效上述两个键（业务层未使用，取消点赞走 Update） |

**扩展方法 `FindLikedBizIds` 的实现要点**：

```sql
select biz_id from like_record
where user_id = ? and biz_type = ? and liked = 1 and biz_id in (?, ?, ...)
```

| 要点 | 说明 |
|------|------|
| 空入参短路 | `len(bizIds) == 0` 直接返回 `nil, nil`，不发 SQL |
| 占位符拼接 | 按 `bizIds` 长度动态生成 `?` 占位符，参数化查询，无 SQL 注入风险 |
| 不走缓存 | 使用 `CachedConn.QueryRowsNoCacheCtx`，注释说明是为「避免缓存击穿与脏读」 |
| 只返回已赞 | `liked = 1` 过滤，取消点赞（`liked = 0`）的行不会出现在结果里 |

**缓存键前缀**（`likerecordmodel_gen.go`）：

| 前缀 | 承载对象 |
|------|---------|
| `cache:likeRecord:id:` | 点赞记录主键缓存 |
| `cache:likeRecord:userId:bizId:bizType:` | 唯一索引缓存，存主键 id |
