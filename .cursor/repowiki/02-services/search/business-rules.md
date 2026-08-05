> 版本：v1.2 | 更新：2026-08-06 | 来源：2026-08-06 复核（go build 全模块通过 + 逻辑文件清点）

---

# Search Business Rules

## ⚠️ 实现状态

search 服务的业务逻辑**已全部实现**：8 个 logic 文件（RPC 4 + API 4）均已落地并编译通过。下列各方法状态已校正为「已实现」。

### RPC 层（`apps/search/rpc/internal/logic/`）

| 分组 | Logic 文件 | 方法 | 状态 |
|------|-----------|------|------|
| 用户兴趣 | `saveinterestslogic.go` | `SaveInterests` | ✅ 已实现 |
| 用户兴趣 | `getinterestslogic.go` | `GetInterests` | ✅ 已实现 |

### API 层（`apps/search/api/internal/logic/`）

| 分组（package） | Logic 文件 | 方法 | 状态 |
|----------------|-----------|------|------|
| `interests` | `saveinterestslogic.go` | `SaveInterests` | ✅ 已实现 |
| `interests` | `getinterestslogic.go` | `GetInterests` | ✅ 已实现 |

### 统计

| 层 | 已实现 | 总计 | 比例 |
|----|--------|------|------|
| RPC (`apps/search/rpc/internal/logic/`) | 4 | 4 | 100% |
| API (`apps/search/api/internal/logic/`) | 4 | 4 | 100% |
| **合计** | **8** | **8** | **100%** |

> 以下各节均为**📋 设计意图（契约推导）**，依据 `apps/search/rpc/search.proto` 注释、`apps/search/api/search.api` 类型定义、`sql/ddl/tj_search.sql` 表结构与 `docs/tjxt.openapi.json` 推导，**2026-08-06 复核：logic 已全部实现并编译通过；以下规则为依据 proto/DDL/.api 契约推导，建议对照源码最终确认**。


## 已知缺口

- 检索/推荐主职责未落地：全仓无 ES 依赖、无 `/interests/{id}/courses` 路由，仅实现用户兴趣存取。
- course ↔ search 未接线。

---

## 1. 用户兴趣存取 📋 设计意图（契约推导）

**核心规则**：一个用户至多一行兴趣记录，主键即用户 ID，兴趣以逗号分隔的二级分类 id 串保存。

| 规则 | 依据 | 说明 |
|------|------|------|
| 主键即用户 ID | DDL `id` 注释「主键，对应用户 id」 | 非自增，由调用方传入；天然保证一人一行 |
| 兴趣为逗号分隔串 | DDL 注释「以逗号分隔，例如：120,220,330」 | 存 varchar(255)，非关联表 |
| 长度上限 255 | DDL `varchar(255)` | 按每个 id 约 4 位 + 逗号估算，最多约 50 个分类，超长需截断或拒绝 |
| 仅二级分类 | DDL 注释「感兴趣的二级分类 id」 | 一级分类不入库 |
| 可为空 | DDL `NULL DEFAULT NULL` | 用户可以不选任何兴趣；Go 侧为 `sql.NullString` |
| Save 为覆盖语义 | proto `SaveInterests` 返回 `Empty`，主键由调用方给定 | 全量替换而非追加，需 upsert |
| 物理删除 | DDL 无 `deleted` 字段 | 无删除接口，清空兴趣通过 Save 空串实现（推导） |

```
流程（SaveInterests，设计意图）:
  1. API 层从 JWT 提取当前 userId（search.api 的 SaveInterestsReq 不含 id）
  2. 校验 interests 串格式：仅数字与逗号，去重、去空项
  3. 校验长度 <= 255
  4. FindOne(userId)：
       ErrNotFound → Insert(&Interests{Id: userId, Interests: 串})
       命中        → Update(全字段)
  5. 返回 OkVO{Success: true}
```

```
流程（GetInterests，设计意图）:
  1. API 层从 JWT 提取 userId
  2. RPC GetInterests(IdReq{id: userId})
  3. ErrNotFound → 返回空 InterestsVO（而非报错，首次进入偏好页属正常）
  4. create_time / update_time 格式化为字符串返回
```

## 2. 身份与越权 📋 设计意图（契约推导）

| 规则 | 依据 | 说明 |
|------|------|------|
| 两个接口均需登录 | `search.api` 中 `@server (jwt: Auth, group: interests)`，`routes.go` 内 `rest.WithJwt` | 无匿名接口 |
| userId 只能来自 JWT | `search.api` 的 `SaveInterestsReq` **只有 `Interests` 一个字段**，无 id | 请求体无法指定他人 id，从设计上杜绝越权写 |
| RPC 层无鉴权 | `apps/search/rpc/internal/config/config.go` 无 JWT 配置 | RPC 信任调用方传入的 `id`，仅限集群内访问 |

## 3. HTTP 与 RPC 的映射差异 📋 设计意图（契约推导）

| 项 | HTTP（`search.api`） | RPC（`search.proto`） | 处理方式 |
|----|---------------------|----------------------|---------|
| 保存方法 | `PUT /interests`，入参 `SaveInterestsReq{Interests}` | `SaveInterests(SaveInterestsReq{id, interests})` | API 层补 `id = JWT.userId` |
| 保存返回 | `OkVO{Success bool}` | `Empty{}` | RPC 无错误即置 `Success = true` |
| 查询方法 | `GET /interests`，无入参 | `GetInterests(IdReq{id})` | API 层补 `id = JWT.userId` |
| 查询返回 | `InterestsVO{Id, Interests, CreateTime, UpdateTime}` | 同名同构 | 直接字段映射 |

> **与 openapi 的出入**：`docs/tjxt.openapi.json` 记录 search 域有 3 个接口——`GET /interests`、**`POST /interests`**、**`GET /interests/{id}/courses`**；而 `search.api` 中保存兴趣用的是 **`PUT /interests`**，且**没有** `/interests/{id}/courses`。openapi 侧 `GET /interests` 的响应为 `List<CategoryBasicDTO>`（分类对象列表），`search.api` 侧为单个 `InterestsVO`（逗号串）。二者尚未对齐。

## 4. 检索与推荐 📋 设计意图（契约推导）

**核心规则**：这是 search 服务命名所指向的主职责，但**当前 proto、api、model、config 中均无任何对应定义**。

| 规划能力 | 依据 | 落地缺口 |
|---------|------|---------|
| 按二级分类查课程 TOP10 | `docs/tjxt.openapi.json` 的 `GET /interests/{id}/courses` | api / proto 均未定义该路由与方法 |
| 全文检索 | 服务命名与 `00-architecture/overview.md` 的「ES 索引」描述 | 全仓无 ES 依赖（Grep `elastic` 零命中） |
| 索引同步 | 需 course 域发布事件 | `03-shared/mq-events.md` 未定义 search 域事件 |

详见 [data-model.md](./data-model.md) 的「⚠️ 存储现状说明」与「索引结构（📋 设计意图，推导，当前未落地）」两节。

---

## 状态说明

### 兴趣串格式

| 形态 | 含义 |
|------|------|
| `NULL` | 从未设置过兴趣（`sql.NullString.Valid = false`） |
| `""` | 主动清空了兴趣 |
| `"120,220,330"` | 选中 3 个二级分类 |

### 接口鉴权

| 接口 | JWT | 用户身份来源 |
|------|-----|-------------|
| `PUT /interests` | 必需 | JWT 载荷 |
| `GET /interests` | 必需 | JWT 载荷 |
