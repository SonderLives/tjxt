# search（搜索推荐）Specification

## Purpose

search（搜索推荐）定位为全文检索与课程推荐微服务（RPC 服务名 `Search`），但截至当前提交实际落地的职责仅为用户兴趣标签的存取。API 层 `search-api`（HTTP :8810）的 2 个 /interests 接口全部转发至自身 RPC `search.rpc`（gRPC :8090，etcd 服务发现 key `search.rpc`），proto 仅定义 `SaveInterests` / `GetInterests` 2 个兴趣方法。数据落在 MySQL `tj_search` 库的单表 `interests`（Redis 做主键缓存），未接入 Elasticsearch——全仓无 ES 依赖与配置，检索/推荐主职责（含 openapi 规划的 `GET /interests/{id}/courses` 课程 TOP10）均未落地，course → search 索引同步事件在 mq-events 中亦未定义。

## Requirements

### Requirement: HTTP API 兴趣域接口

`search-api` SHALL 提供用户兴趣的查询与保存 HTTP 接口。`docs/tjxt.openapi.json`（聚合文档，最后同步 2026-08-05）记录兴趣域接口：

| 方法 | 路径 | 说明 | 响应体 |
|------|------|------|--------|
| GET | /interests | 查询我的兴趣爱好 | R«List«CategoryBasicDTO»» |
| POST | /interests | 新增兴趣爱好 | R |

- 统一约定（引用全局规范）：响应格式 `pkg/response.R{Code,Msg,RequestId,Data any}`；分页 `PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`；错误码见共享错误码表（见 `openspec/specs/shared/error-codes`）
- openapi 聚合文档中三个接口认证列均标注「否」，但 `search.api` 实际为 `@server (jwt: Auth, group: interests)`，`routes.go` 内 `rest.WithJwt`——两个 /interests 接口均需登录

> 注（与 search.api 的出入）：`search.api` 中保存兴趣用的是 **`PUT /interests`**（入参 `SaveInterestsReq` 仅含 `Interests` 字段），而非 openapi 记录的 `POST /interests`；openapi 侧 `GET /interests` 响应为 `List<CategoryBasicDTO>`（分类对象列表），`search.api` 侧为单个 `InterestsVO`（逗号串）。二者尚未对齐。

#### Scenario: 查询我的兴趣爱好

- **WHEN** 已登录用户 GET /interests（JWT 校验通过）
- **THEN** API 层从 JWT 取 userId 调 RPC `GetInterests`，按 openapi 响应为 R«List«CategoryBasicDTO»»（分类对象列表）
- **AND** `search.api` 侧实际定义为单个 `InterestsVO{Id, Interests, CreateTime, UpdateTime}`（逗号串），与 openapi 尚未对齐

#### Scenario: 保存兴趣爱好

- **WHEN** 已登录用户提交保存兴趣请求
- **THEN** openapi 记录为 POST /interests，`search.api` 实际为 PUT /interests（入参 `SaveInterestsReq` 仅含 `Interests` 一个字段）
- **AND** 保存成功返回 `OkVO{Success: true}`（RPC 无错误即置 `Success = true`）

### Requirement: HTTP API 推荐域接口

`search-api` SHALL 按 openapi 约定提供根据二级分类查询课程 TOP10 的推荐接口，该能力当前为设计已规划、代码未落地：

| 方法 | 路径 | 说明 | 响应体 |
|------|------|------|--------|
| GET | /interests/{id}/courses | 根据二级分类id查询课程TOP10 | R«List«CourseVO»» |

> 已知缺口：`docs/tjxt.openapi.json` 记录了该接口，但 `apps/search/api/search.api` 与 `apps/search/rpc/search.proto` **均未定义**，属于设计已规划、代码未落地的部分；实现还需引入 course RPC 客户端（当前缺失）。

#### Scenario: 按二级分类查课程 TOP10（未落地）

- **WHEN** 前端按 openapi 约定 GET /interests/{id}/courses
- **THEN** `search.api` 与 `search.proto` 均未定义该路由与方法，返回 R«List«CourseVO»» 的实现不存在
- **AND** 属设计已规划、代码未落地部分，需先引入 course RPC 客户端方可实现

### Requirement: RPC 服务契约与服务发现

search 服务 SHALL 以 gRPC 服务名 `Search` 注册于 etcd（key `search.rpc`，监听 `0.0.0.0:8090`），当前 proto 只定义了用户兴趣（interests）相关的 2 个方法，检索/推荐类方法尚未定义。消费方：

| 消费方 | 调用方式 | 说明 |
|--------|---------|------|
| `search-api`（自身 API 层） | HTTP Handler → `searchclient.Search` RPC | `apps/search/api/internal/svc/servicecontext.go` 注入 `SearchRpc`，2 个 HTTP 接口全部指向自身 RPC |

> 已知缺口：全仓 Grep `apps/search/rpc/search` / `searchclient` / `SearchRpc`，除 `apps/search/` 自身外**无其它服务引用**；`specs/architecture/service-topology` 中标注的 `search-rpc --> AuthRPC` 依赖在代码中不存在——`apps/search/rpc/internal/config/config.go` 未定义任何 RPC 客户端配置项。

#### Scenario: search-api handler 转发

- **WHEN** search-api 任一 /interests handler 收到请求
- **THEN** logic 层直接调用 `l.svcCtx.SearchRpc` 对应 RPC 方法（API 层为纯转发层）
- **AND** SearchRpc 客户端按 `SearchRpc.Etcd`（`127.0.0.1:2379`，key `search.rpc`）经 etcd 发现 `search.rpc`

### Requirement: RPC 用户兴趣方法

`Search` 服务 SHALL 提供 2 个用户兴趣 RPC 方法：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `SaveInterests` | `SaveInterestsReq { id, interests }` | `Empty {}` | 保存用户兴趣（新增或覆盖） |
| `GetInterests` | `IdReq { id }` | `InterestsVO { id, interests, createTime, updateTime }` | 按用户 ID 查询兴趣 |

请求字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 主键，对应用户 id（`interests` 表主键即用户 id） |
| `interests` | string | 感兴趣的二级分类 id，以逗号分隔，例如：`120,220,330` |

响应字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 用户 id |
| `interests` | string | 逗号分隔的二级分类 id 串 |
| `createTime` | string | 创建时间 |
| `updateTime` | string | 更新时间 |

- `SaveInterests` 返回 `Empty` 不回传 ID——主键即调用方传入的用户 ID，无需服务端生成
- 个性化推荐（设计意图，proto 未定义）：取出用户 `interests` 中的二级分类 id，交由 course 域按分类查课程 TOP10（对应 openapi 的 `GET /interests/{id}/courses`）

> 实现状态（2026-08-06 复核）：8 个 logic 文件（RPC 4 + API 4）均已落地并编译通过；依据 proto/DDL/.api 契约推导的规则建议对照源码最终确认。

#### Scenario: 新用户选兴趣

- **WHEN** 注册引导页勾选二级分类后，search-api 从 JWT 取 userId
- **THEN** 调 `SaveInterests(id=userId, interests="120,220,330")` 保存，返回 `Empty` 不回传 ID
- **AND** 主键即调用方传入的用户 ID，无需服务端生成

#### Scenario: 回显已选兴趣

- **WHEN** 用户进入偏好设置页调 `GetInterests(id=userId)`
- **THEN** 返回 `InterestsVO{id, interests, createTime, updateTime}`
- **AND** 前端按逗号切分 id 串并回勾

### Requirement: 用户兴趣存取规则

用户兴趣存取 SHALL 满足以下核心规则：一个用户至多一行兴趣记录，主键即用户 ID，兴趣以逗号分隔的二级分类 id 串保存。

| 规则 | 依据 | 说明 |
|------|------|------|
| 主键即用户 ID | DDL `id` 注释「主键，对应用户 id」 | 非自增，由调用方传入；天然保证一人一行 |
| 兴趣为逗号分隔串 | DDL 注释「以逗号分隔，例如：120,220,330」 | 存 varchar(255)，非关联表 |
| 长度上限 255 | DDL `varchar(255)` | 按每个 id 约 4 位 + 逗号估算，最多约 50 个分类，超长需截断或拒绝 |
| 仅二级分类 | DDL 注释「感兴趣的二级分类 id」 | 一级分类不入库 |
| 可为空 | DDL `NULL DEFAULT NULL` | 用户可以不选任何兴趣；Go 侧为 `sql.NullString` |
| Save 为覆盖语义 | proto `SaveInterests` 返回 `Empty`，主键由调用方给定 | 全量替换而非追加，需 upsert |
| 物理删除 | DDL 无 `deleted` 字段 | 无删除接口，清空兴趣通过 Save 空串实现（推导） |

兴趣串形态：

| 形态 | 含义 |
|------|------|
| `NULL` | 从未设置过兴趣（`sql.NullString.Valid = false`） |
| `""` | 主动清空了兴趣 |
| `"120,220,330"` | 选中 3 个二级分类 |

> 以上规则为 📋 设计意图（契约推导），依据 proto 注释、search.api 类型定义、DDL 表结构与 openapi 推导。

#### Scenario: SaveInterests 保存流程（设计意图）

- **WHEN** API 层从 JWT 提取当前 userId（`search.api` 的 `SaveInterestsReq` 不含 id）并提交 interests 串
- **THEN** 依次校验：仅数字与逗号、去重去空项、长度 <= 255
- **AND** `FindOne(userId)` 返回 `ErrNotFound` 则 `Insert(&Interests{Id: userId, ...})`，命中则全字段 `Update`，成功后返回 `OkVO{Success: true}`

#### Scenario: GetInterests 首次查询返回空而非报错

- **WHEN** 用户从未设置过兴趣，RPC `GetInterests` 返回 `model.ErrNotFound`
- **THEN** 返回空 `InterestsVO` 而非报错（首次进入偏好页属正常）
- **AND** create_time / update_time 格式化为字符串返回

### Requirement: 身份与越权约束

search 服务的两个 /interests 接口 SHALL 均要求登录且 userId 只能来自 JWT：

- 两个接口均需登录：`search.api` 中 `@server (jwt: Auth, group: interests)`，`routes.go` 内 `rest.WithJwt`，无匿名接口
- userId 只能来自 JWT：`search.api` 的 `SaveInterestsReq` **只有 `Interests` 一个字段**、无 id，请求体无法指定他人 id，从设计上杜绝越权写
- RPC 层无鉴权：`apps/search/rpc/internal/config/config.go` 无 JWT 配置，RPC 信任调用方传入的 `id`，仅限集群内访问

接口鉴权清单：

| 接口 | JWT | 用户身份来源 |
|------|-----|-------------|
| `PUT /interests` | 必需 | JWT 载荷 |
| `GET /interests` | 必需 | JWT 载荷 |

#### Scenario: 越权写防护

- **WHEN** 客户端提交保存兴趣请求
- **THEN** userId 一律取自 JWT 载荷，请求体无法指定他人 id
- **AND** RPC 层不鉴权、信任调用方传入的 id，仅限集群内访问

### Requirement: HTTP 与 RPC 映射差异

search 服务 HTTP 与 RPC SHALL 按以下方式映射（API 层补 `id = JWT.userId`）：

| 项 | HTTP（`search.api`） | RPC（`search.proto`） | 处理方式 |
|----|---------------------|----------------------|---------|
| 保存方法 | `PUT /interests`，入参 `SaveInterestsReq{Interests}` | `SaveInterests(SaveInterestsReq{id, interests})` | API 层补 `id = JWT.userId` |
| 保存返回 | `OkVO{Success bool}` | `Empty{}` | RPC 无错误即置 `Success = true` |
| 查询方法 | `GET /interests`，无入参 | `GetInterests(IdReq{id})` | API 层补 `id = JWT.userId` |
| 查询返回 | `InterestsVO{Id, Interests, CreateTime, UpdateTime}` | 同名同构 | 直接字段映射 |

> 注（与 openapi 的出入）：`docs/tjxt.openapi.json` 记录 search 域有 3 个接口——`GET /interests`、`POST /interests`、`GET /interests/{id}/courses`；而 `search.api` 中保存兴趣用 `PUT /interests`，且没有 `/interests/{id}/courses`。openapi 侧 `GET /interests` 响应为 `List<CategoryBasicDTO>`（分类对象列表），`search.api` 侧为单个 `InterestsVO`（逗号串）。二者尚未对齐。

#### Scenario: 保存与查询的 id 注入

- **WHEN** API 层处理 `PUT /interests` 或 `GET /interests`
- **THEN** 一律以 `id = JWT.userId` 组装 RPC 请求（`SaveInterestsReq{id, interests}` / `IdReq{id}`）
- **AND** 保存返回 `OkVO{Success bool}`、查询返回 `InterestsVO` 同名同构直接字段映射

### Requirement: 检索与推荐规划能力（未落地）

search 服务命名所指向的主职责 SHALL 包含以下规划能力，但当前 proto、api、model、config 中均无任何对应定义（检索域无 HTTP 端点与 RPC 方法，仅为设计意图）：

| 规划能力 | 依据 | 落地缺口 |
|---------|------|---------|
| 按二级分类查课程 TOP10 | `docs/tjxt.openapi.json` 的 `GET /interests/{id}/courses` | api / proto 均未定义该路由与方法 |
| 全文检索 | 服务命名与 `00-architecture/overview.md` 的「ES 索引」描述 | 全仓无 ES 依赖（Grep `elastic` 零命中） |
| 索引同步 | 需 course 域发布事件 | `03-shared/mq-events.md` 未定义 search 域事件 |

索引结构（📋 设计意图，推导，当前未落地），由 proto、DDL 与 openapi 中 `GET /interests/{id}/courses` 反推：

| 推导项 | 内容 |
|--------|------|
| 检索对象 | 课程（course 域），按二级分类 id、名称、简介建全文索引 |
| 推荐入口 | 以 `interests.interests` 中的二级分类 id 集合作为召回条件 |
| 排序维度 | 按报名人数 / 销量取 TOP10（对应 openapi 摘要「课程 TOP10」） |
| 索引更新 | 需 course 域在课程发布/下架时同步索引；当前 mq-events 中未定义任何 search 域事件 |
| 与现状差距 | 需引入 ES 客户端依赖、在 `config.go` 增加 ES 配置段、在 `ServiceContext` 注入 ES client、在 proto 增加检索方法 |

> 已知缺口：检索/推荐主职责未落地——全仓无 ES 依赖、无 `/interests/{id}/courses` 路由，仅实现用户兴趣存取；course ↔ search 未接线。在上述能力落地前，search 服务的实际职责仅为「用户兴趣标签的存取」。

#### Scenario: 课程发布/下架索引同步（未落地）

- **WHEN** course 域课程发布或下架
- **THEN** 按设计意图应同步更新课程全文索引（按二级分类 id、名称、简介）
- **AND** 当前 mq-events 中未定义任何 search 域事件，course → search 事件接线未完成，该场景无实现

#### Scenario: 个性化推荐召回（设计意图）

- **WHEN** 需要为用户生成课程推荐
- **THEN** 以 `interests.interests` 中的二级分类 id 集合作为召回条件，按报名人数 / 销量取课程 TOP10
- **AND** 该链路当前无任何对应实现（proto 未定义检索方法、无 ES 依赖）

### Requirement: 数据模型表清单与关系

search 服务 SHALL 将数据持久化于 MySQL `tj_search` 库（当前仅 1 张表），Model 由 goctl 生成（`*_gen.go` 禁止修改），扩展方法统一放在同名自定义 `.go` 文件中：

| 自定义 Model 文件 | 对应表 | 扩展方法 |
|------------------|--------|---------|
| `interestsmodel.go` | `interests` | 无（仅 goctl 空壳） |

| 表 | 职责 | 关键字段 |
|----|------|---------|
| `interests` | 用户兴趣表，保存感兴趣的二级分类 id | `id`（bigint，主键即用户 id）、`interests`（varchar(255)，逗号分隔二级分类 id，可空）、`create_time`（自动填充）/ `update_time`（自动更新） |

- Model 在 `apps/search/rpc/internal/svc/servicecontext.go` 中通过 `sqlx.NewMysql(c.DataSource)` + `c.Cache` 注入，走带缓存的 `sqlc.CachedConn`
- 缓存 key 前缀（`interestsmodel_gen.go` 中定义）：`Interests` → `cache:interests:id:`
- 关系（跨库无外键）：`interests (1) ─── (1) user`（id 即 user 域的用户 id）；`interests` 字段内的二级分类 id 串逻辑指向 course 域的分类表（跨库、跨服务，无外键、无 RPC 调用）
- `tj_search` 库中只有一张表，库内无任何表间关系

#### Scenario: 兴趣与用户、分类的逻辑关联

- **WHEN** 查询或保存用户兴趣
- **THEN** `interests.id` 即 user 域用户 id，一人至多一行，跨库无外键
- **AND** `interests` 字段内的二级分类 id 逻辑指向 course 域分类表，无外键、无 RPC 调用

### Requirement: 数据模型关键不变量

search 服务的存储层 SHALL 满足以下不变量：

- 主键即用户 ID（bigint，非自增），由调用方传入，天然保证一个用户至多一行兴趣记录
- `interests` 为逗号分隔字符串（varchar(255)）而非关联表，无法按分类 id 反查用户；若需「某分类下有哪些用户」，当前表结构不支持
- Go 结构体 `Interests` 中 `interests` 映射为 `sql.NullString`（列 `NULL DEFAULT NULL`）
- 无 `deleted` 字段，删除为物理删除；无 `creater` / `updater` 审计字段
- 字符集 `utf8mb4` / `utf8mb4_0900_ai_ci`，与连接串 `charset=utf8mb4` 一致（本仓少数字符集完全对齐的服务）
- Model 当前为 goctl 空壳，可用方法仅 `Insert` / `FindOne`（按主键带缓存）/ `Update` / `Delete` 4 个；`SaveInterests` 的「存在则更新，不存在则插入」语义无 upsert 方法兜底，需在 logic 中先 `FindOne`、按 `model.ErrNotFound` 分支决定走 `Insert` 还是 `Update`，或在自定义 model 中补写 `Upsert`（📋 待补齐的设计意图）
- `vars.go` 仅定义 `ErrNotFound = sqlx.ErrNotFound`，无其它域那样的通用 SQL 工具函数
- 存储现状：无任何 Elasticsearch 索引结构落地——全仓 Grep `elastic` / `elasticsearch` / `olivere`（`*.go` / `*.yaml` / `*.mod`）零命中，`config.go` 无 ES 字段；截至当前提交数据来源仅为 MySQL `tj_search` 单表 `interests`（配合 Redis 做主键缓存）

#### Scenario: 一人一行不变量

- **WHEN** 任意用户保存兴趣
- **THEN** 以用户 id 为主键写入 `interests` 表（非自增主键由调用方传入）
- **AND** 一个用户至多一行兴趣记录，重复保存走 `Update` 覆盖而非新增

#### Scenario: 无 upsert 方法的实现约束

- **WHEN** `SaveInterests` 需要实现「存在则更新，不存在则插入」
- **THEN** model 层无 upsert 方法，logic 中先 `FindOne`、按 `model.ErrNotFound` 分支决定走 `Insert` 还是 `Update`
- **AND** 或在自定义 model（`interestsmodel.go` 手写扩展位，当前为空壳、无任何扩展方法）中补写 `Upsert`

### Requirement: 服务配置与依赖

search 服务 SHALL 按以下配置运行（`apps/search/api/etc/search-api.yaml`、`apps/search/rpc/etc/search.yaml`）：

| 服务 | 配置项 | 默认值 |
|------|--------|--------|
| `search-api` | `Name` / `Host` / `Port` | `search-api` / `0.0.0.0` / `8810`（HTTP，对外） |
| `search-api` | `Auth.AccessSecret` / `Auth.AccessExpire` | `change-me-in-production` / `7200`（秒） |
| `search-api` | `SearchRpc.Etcd.Hosts[0]` / `SearchRpc.Etcd.Key` | `127.0.0.1:2379` / `search.rpc` |
| `search.rpc` | `Name` / `ListenOn` | `search.rpc` / `0.0.0.0:8090`（gRPC，集群内） |
| `search.rpc` | `Etcd.Hosts[0]` / `Etcd.Key` | `127.0.0.1:2379` / `search.rpc` |
| `search.rpc` | `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_search?charset=utf8mb4&parseTime=true&loc=Local` |
| `search.rpc` | `Cache[0]` | `127.0.0.1:6379`，Pass 空，Type `node`（单机模式） |

- 依赖的外部服务：MySQL `tj_search` 库（单表 `interests`，仅 `InterestsModel` 一个模型使用）、Redis 缓存（节点模式，`InterestsModel` 走 `sqlc.CachedConn`，缓存前缀 `cache:interests:id:`）、etcd（RPC 端注册 `search.rpc`，API 端按同名 key 发现）
- API 层未配置 MySQL / Redis，不直连存储；auth 服务无直接连接，仅共享 `AccessSecret` 完成离线 JWT 校验，不调用 auth RPC
- JWT：search 只在 API 层校验不签发，`Auth.AccessSecret` 必须与 auth 服务的 `Jwt.AccessSecret` 保持一致，否则 `/interests` 两个接口全部 401；`Auth.AccessExpire` 仅作为配置项被反序列化，代码内无签发逻辑引用；生产环境必须修改默认值 `change-me-in-production`
- RPC config 结构体为 `zrpc.RpcServerConf` + `DataSource string` + `Cache cache.CacheConf`；未配置 JWT 段，也未配置任何 RPC 客户端——RPC 层不鉴权、不外调

> 已知缺口（缺失的外部依赖配置）：Elasticsearch 无任何配置项（`config.go` 无 ES 字段、`go.mod` 无 ES 客户端依赖），服务当前不具备任何全文检索能力，实际只是 `interests` 单表的 CRUD 外壳；course RPC 客户端缺失，`GET /interests/{id}/courses` 无法实现；MQ / Kafka 缺失，无法消费 course 域事件做索引同步。引入 ES 时需同步改动：`search.yaml`（新增 ES 段）、`apps/search/rpc/internal/config/config.go`（新增字段）、`servicecontext.go`（注入 client）、`apps/search/go.mod`（新增依赖）。
>
> 可观测性注入配置（Telemetry / Prometheus / Log 三段）：两份 yaml 中已注入勿删（约定见 `openspec/specs/infra/observability`）。

#### Scenario: API 层经 etcd 发现自身 RPC

- **WHEN** search-api 启动并处理 /interests 请求
- **THEN** `SearchRpc` 客户端按 `SearchRpc.Etcd`（`127.0.0.1:2379`，key `search.rpc`）发现 `search.rpc` 并完成转发调用
- **AND** API 层不落库、不做业务判断（未配置 MySQL / Redis）

#### Scenario: JWT 密钥不一致导致 401

- **WHEN** search-api 的 `Auth.AccessSecret` 与 auth 服务 `Jwt.AccessSecret` 不一致
- **THEN** `/interests` 两个接口全部 401（search 只在 API 层做离线校验，不签发令牌）
- **AND** 生产环境必须修改默认值 `change-me-in-production`

