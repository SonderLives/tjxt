# remark（评论点赞）Specification

## Purpose

remark（评论点赞）是评论互动域的点赞微服务（RPC 服务名 `Remark`），以 `bizType + bizId` 多态泛化承载回复、笔记、问答等各业务对象的点赞关系，`bizType` 不做枚举校验、不感知具体业务表，点赞能力可被任意新业务复用。API 层 `remark-api`（HTTP :8813）为纯透传层，2 个 handler 全部转发到自身 RPC `remark.rpc`（gRPC :8093，etcd 服务发现 key `remark.rpc`）。数据落在 MySQL `tj_remark` 库唯一的 `like_record` 表（无逻辑删除字段，软取消语义），Redis 缓存主键与唯一索引查询，etcd `127.0.0.1:2379` 服务注册与发现；`remark.rpc` 不依赖任何其他服务的 RPC。

## Requirements

### Requirement: HTTP API 点赞接口

`remark-api` SHALL 提供以下点赞 HTTP 接口（认证列均为「是」，权限标签 like）：

| 方法 | 路径 | 说明 | 请求体 | 响应体 |
|------|------|------|--------|--------|
| POST | /likes | 点赞或取消点赞 | LikeRecordFormReq | OkVO |
| GET | /likes/list | 批量查询当前用户已点赞的业务id | LikeListReq (query) | LikeListResp |

请求/响应结构：

| 结构 | 字段 |
|------|------|
| `LikeRecordFormReq` | `bizId` int64（被点赞的业务 id）、`bizType` string（业务类型，如 reply / note / question）、`liked` bool（true=点赞，false=取消点赞） |
| `OkVO` | `success` bool（固定返回 true） |
| `LikeListReq`（query） | `bizType` string（业务类型）、`bizIds` []string（业务 id 列表，服务端逐个 `strconv.ParseInt` 转换，非法或 `<= 0` 的直接丢弃） |
| `LikeListResp` | `likedBizIds` []int64（入参 `bizIds` 中当前用户已点赞 liked=1 的子集） |

- 响应统一为 `pkg/response.R{Code,Msg,RequestId,Data any}`；错误码遵循共享错误码表（见 `openspec/specs/shared/error-codes`）
- 全局分页约定为 `PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`（remark 当前两个接口未用到分页）

#### Scenario: 点赞或取消点赞

- **WHEN** 客户端 POST /likes 提交 LikeRecordFormReq（bizId / bizType / liked）且携带有效 JWT
- **THEN** 用户身份由 `auth.UserIdFromCtx` 从上下文取出，三个参数原样透传 RPC `Like`
- **AND** RPC 成功后返回 `OkVO{Success: true}`，不回传点赞后的状态

#### Scenario: 批量查询点赞态

- **WHEN** 客户端 GET /likes/list 传 `bizType` 与 `bizIds`（query 参数，[]string）
- **THEN** 返回 `LikeListResp{likedBizIds}`，仅含入参 `bizIds` 中当前用户已点赞（liked=1）的子集
- **AND** 解析失败或 `<= 0` 的 `bizIds` 元素被静默丢弃，不报错

### Requirement: 鉴权与身份注入

remark-api SHALL 对全部接口启用全局 JWT 鉴权，并以 JWT 中的 userId 为唯一身份来源：

- `remark.api` 声明 `@server(jwt: Auth, group: like)`，`routes.go` 中两个路由由 `rest.WithJwt(serverCtx.Config.Auth.AccessSecret)` 统一包裹
- 无匿名接口：remark-api 不存在免鉴权路由，未带有效 JWT 的请求在中间件层即被拒绝
- 越权防护：所有查询与写入都以 JWT 中的 userId 为维度，用户无法读写他人的点赞记录
- 身份注入：API 层 `auth.UserIdFromCtx(l.ctx)` 取 userId，失败直接返回错误；proto 中虽有 `userId` 字段，但 API 层永远用它覆盖，不信任请求体中的 userId

JWT 配置（`remark-api.yaml`）：

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `Auth.AccessSecret` | `change-me-in-production` | 生产环境必须修改，否则 JWT 可被伪造；必须与签发方（auth 服务）保持一致，否则所有接口鉴权失败 |
| `Auth.AccessExpire` | `7200`（秒） | 在 remark-api 中仅作配置项声明，实际令牌过期校验由 go-zero JWT 中间件依据 token 内的 `exp` 完成 |

#### Scenario: 未携带有效 JWT 被拒绝

- **WHEN** 请求未携带有效 JWT 访问 /likes 或 /likes/list
- **THEN** 在 JWT 中间件层即被拒绝，不进入 logic 层
- **AND** remark-api 不存在任何免鉴权路由

#### Scenario: 身份注入覆盖请求身份

- **WHEN** API 层处理 /likes 或 /likes/list 请求
- **THEN** 以 `auth.UserIdFromCtx(l.ctx)` 取 userId 作为调用 RPC 的用户身份
- **AND** 请求体/query 中的任何 userId 声明均被忽略，不构成越权通道

### Requirement: RPC 服务契约与服务发现

remark 服务 SHALL 以 gRPC 服务名 `Remark` 注册于 etcd，为评论互动域的点赞微服务，通过 etcd 服务发现供 API 层调用：

| 项 | 值 |
|----|----|
| RPC 服务名 | `Remark` |
| RPC 监听地址 | `0.0.0.0:8093`（gRPC） |
| etcd Hosts / Key | `127.0.0.1:2379` / `remark.rpc` |
| API 服务 | `remark-api`，监听 `0.0.0.0:8813`（HTTP） |

消费方：

| 消费方 | 调用方式 | 说明 |
|--------|---------|------|
| `remark-api`（自身 API 层） | `apps/remark/api/internal/svc/servicecontext.go` 中 `remarkclient "tjxt/apps/remark/rpc/client/remark"` → `RemarkRpc` | 两个 HTTP handler（/likes、/likes/list）均透传到本 RPC |

> 注：截至当前版本，仅 remark-api 在 `servicecontext.go` 中注入 `remarkclient`。其他服务如 learning 的笔记/问答列表若需回填「我是否点赞」，可复用 `QueryLikedBizIds`，但目前尚未接入。
>
> `remark.rpc` 不依赖任何其他服务的 RPC；`bizType + bizId` 是多态关联，本服务不感知具体业务表结构，也不做跨服务校验——这是点赞能力可被任意业务复用的前提。
>
> 可观测性注入配置（Telemetry / Prometheus / Log 三段）：`remark-api.yaml` 与 `remark.yaml` 中已注入勿删（约定见 `openspec/specs/infra/observability`）。

#### Scenario: remark-api handler 转发

- **WHEN** remark-api 的 /likes 或 /likes/list handler 收到请求
- **THEN** logic 层调用 `RemarkRpc` 对应方法，客户端按 `RemarkRpc.Etcd`（`127.0.0.1:2379`，key `remark.rpc`）经 etcd 服务发现完成调用
- **AND** API 层为纯透传层，不落库、不做业务判断

### Requirement: RPC 点赞方法

`Remark` 服务 SHALL 提供 2 个点赞 RPC 方法：

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `Like` | `LikeReq { userId, bizId, bizType, liked }` | `Empty {}` | 点赞或取消点赞，同一 (用户,业务) 维度幂等 upsert |
| `QueryLikedBizIds` | `LikedReq { userId, bizType, bizIds }` | `LikedResp { likedBizIds }` | 批量查询当前用户在一组业务 id 中已点赞的子集 |

关键字段：

| Message.字段 | 类型 | 说明 |
|--------------|------|------|
| `LikeReq.userId` | int64 | 点赞人 ID（来自 JWT） |
| `LikeReq.bizId` | int64 | 被点赞的业务对象 ID |
| `LikeReq.bizType` | string | 业务类型，如 `reply` / `note` / `question` |
| `LikeReq.liked` | bool | true=点赞，false=取消点赞；服务端转成 `liked` 列的 1/0 |
| `LikedReq.userId` | int64 | 当前用户 ID（来自 JWT） |
| `LikedReq.bizType` | string | 业务类型 |
| `LikedReq.bizIds` | repeated int64 | 待查询的业务 id 列表，为空时直接返回空结果 |
| `LikedResp.likedBizIds` | repeated int64 | 入参 `bizIds` 中 `liked = 1` 的子集，未点赞与已取消的均不返回 |

#### Scenario: 列表页回填点赞态

- **WHEN** 业务列表（如回复列表）拿到一批 `bizId` 后调 `QueryLikedBizIds`
- **THEN** 一次调用返回当前用户已点赞的子集，前端按返回集合高亮已赞项，避免逐条查询
- **AND** 未点赞与已取消的 `bizId` 均不出现在 `likedBizIds` 中

#### Scenario: 重复点击容错

- **WHEN** 用户连点两次点赞，第二次调 `Like` 时 `existing.Liked == liked`
- **THEN** 直接返回 `Empty{}`，不产生任何写库

### Requirement: 点赞 upsert 规则

`Like` SHALL 以同一 `(user_id, biz_id, biz_type)` 三元组最多一行的方式落库——点赞与取消点赞都落在同一行上，通过 `liked` 列的 1/0 表达状态。具体约束：

- bool → tinyint 转换：`in.Liked == true` 转 `liked = 1`，否则 `liked = 0`
- 记录不存在则新增：`FindOneByUserIdBizIdBizType` 返回 `ErrNotFound` 时执行 `Insert`
- 状态相同则空转：`existing.Liked == liked` 时直接返回 `Empty{}`，不发起任何写库
- 状态不同则更新：只改 `existing.Liked` 后 `Update`，`user_id` / `biz_id` / `biz_type` 保持不变
- 软取消：取消点赞不删除行，只把 `liked` 置 0；重新点赞再置回 1

幂等性：

| 场景 | 行为 |
|------|------|
| 首次点赞 | Insert 一行，`liked=1` |
| 重复点赞 | 幂等短路，无写库，返回成功 |
| 点赞 → 取消 | Update 把 `liked` 改 0 |
| 重复取消 | 幂等短路，无写库，返回成功 |
| 取消 → 重新点赞 | Update 把 `liked` 改回 1 |

> 并发写入的兜底在 DB：`FindOneByUserIdBizIdBizType`（读）与 `Insert`（写）之间存在并发窗口，两个并发的首次点赞请求可能都判定为「记录不存在」并同时 Insert。此时 `uk_user_biz` 唯一索引会让其中一条 Insert 报错，错误直接返回给调用方——代码层面未做重复键的捕获与降级。
>
> 注（错误判定）：`likelogic.go` 中比较的是 `sqlc.ErrNotFound`，而 `likerecordmodel_gen.go` 的 `FindOneByUserIdBizIdBizType` 未命中时实际返回 `model.ErrNotFound`；两者在 `vars.go` 中通过 `var ErrNotFound = sqlx.ErrNotFound` 指向同一个 error 值，因此判定成立。

#### Scenario: Like 主流程

- **WHEN** `Like` 收到请求
- **THEN** 依次执行：`liked := 0`（`in.Liked` 为 true 时置 1）→ `FindOneByUserIdBizIdBizType(userId, bizId, bizType)`
- **AND** 未命中（ErrNotFound）→ `Insert{UserId, BizId, BizType, Liked}` 后返回；其它错误直接返回错误；命中且状态相同 → 直接返回 `Empty`（幂等短路）；命中且状态不同 → `Update` 把 `existing.Liked` 改为 `liked` 后返回

#### Scenario: 并发首次点赞由唯一索引兜底

- **WHEN** 两个并发的首次点赞请求都判定为「记录不存在」并同时 Insert
- **THEN** `uk_user_biz` 唯一索引使其中一条 Insert 报错，错误直接返回给调用方
- **AND** 代码层面未做重复键的捕获与降级

### Requirement: 批量点赞状态查询规则

`QueryLikedBizIds` SHALL 一次查询回填整页列表的点赞态、避免逐条查询：

- 纯透传：RPC logic 不做任何校验，直接调 `FindLikedBizIds` 并把结果包进 `LikedResp`
- 空入参短路：Model 层 `len(bizIds) == 0` 时返回 `nil, nil`，不发 SQL
- 只返回已赞子集：SQL 带 `liked = 1`，取消点赞的 `bizId` 不出现在结果中
- 不走缓存：`QueryRowsNoCacheCtx` 直读 DB，保证点赞后立即可见
- 结果为子集：返回的 `likedBizIds` 长度通常远小于入参 `bizIds`

#### Scenario: QueryLikedBizIds 流程

- **WHEN** 调 `QueryLikedBizIds(userId, bizType, bizIds)`
- **THEN** `FindLikedBizIds` 按 `bizIds` 长度动态拼接 `in (?,?,...)` 占位符，参数化执行 `select biz_id where user_id=? and biz_type=? and liked=1 and biz_id in (...)`，结果包装成 `LikedResp{LikedBizIds: ids}` 返回
- **AND** `bizIds` 为空时直接短路返回空结果，不发 SQL

### Requirement: API 层参数处理规则

remark-api SHALL 作为纯透传层处理参数（身份规则见「鉴权与身份注入」）：

/likes（LikeLogic）：

- 参数透传：`bizId` / `bizType` / `liked` 原样传给 RPC
- 固定响应：RPC 成功后返回 `OkVO{Success: true}`，不回传点赞后的状态
- 无参数校验：API 层不校验 `bizId > 0`，也不校验 `bizType` 是否在允许枚举内

/likes/list（ListLikedLogic）：

- 字符串转数值：请求中的 `bizIds` 是 `[]string`（query 参数），逐个 `strconv.ParseInt(s, 10, 64)`
- 非法值静默丢弃：解析失败（`perr != nil`）或 `id <= 0` 的元素直接跳过，不报错
- 全部非法的后果：转换后 `bizIds` 为空切片 → Model 层短路 → 返回空 `likedBizIds`
- 结果透传：直接把 `r.LikedBizIds` 包进 `LikeListResp` 返回

#### Scenario: ListLiked 流程

- **WHEN** /likes/list 收到请求
- **THEN** 依次执行：`auth.UserIdFromCtx` 取 userId → 遍历 `req.BizIds` 逐个 `strconv.ParseInt(s, 10, 64)`，仅 `perr == nil && id > 0` 的元素追加 → `RemarkRpc.QueryLikedBizIds(userId, bizType, bizIds)` → 返回 `LikeListResp{LikedBizIds: r.LikedBizIds}`

#### Scenario: 非法 bizIds 静默丢弃

- **WHEN** `bizIds` 中包含解析失败或 `<= 0` 的元素
- **THEN** 这些元素被直接跳过，不报错
- **AND** 全部非法时转换后为空切片，经 Model 层短路后返回空 `likedBizIds`

### Requirement: 数据模型表清单与关系

remark 服务 SHALL 将数据持久化于 MySQL `tj_remark` 库的唯一一张表 `like_record`，Model 由 goctl 生成（`*_gen.go` 禁止修改），扩展方法统一放在同名自定义 `.go` 文件中：

| 自定义 Model 文件 | 对应表 | 扩展方法 |
|------------------|--------|---------|
| `likerecordmodel.go` | `like_record` | FindLikedBizIds |

`vars.go` 定义 `ErrNotFound = sqlx.ErrNotFound`。remark 服务仅一张表、无独立的 `consts.go`；业务类型 `bizType` 由调用方自行传入字符串（如 `reply` / `note` / `question`），服务端不做枚举校验。

`like_record` 关键字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | bigint | 点赞记录 id，自增主键 |
| `user_id` | bigint | 点赞人 id（`uk_user_biz` 首列） |
| `biz_id` | bigint | 被点赞业务 id（`uk_user_biz` 次列 / `idx_biz` 次列） |
| `biz_type` | varchar(32) | 点赞业务类型，例如 reply / note / question（`uk_user_biz` 末列 / `idx_biz` 首列） |
| `liked` | tinyint | 1:点赞 0:取消，默认 1（查询过滤条件） |
| `create_time` / `update_time` | datetime | 创建 / 更新时间，自动填充 |

索引：

| 索引名 | 类型 | 列 | 用途 |
|--------|------|----|------|
| `PRIMARY` | 主键 | `id` | 主键查询，对应 `cache:likeRecord:id:` 缓存 |
| `uk_user_biz` | 唯一 | `user_id`, `biz_id`, `biz_type` | 保证同一用户对同一业务对象只有一行记录；`FindOneByUserIdBizIdBizType` 与 `cache:likeRecord:userId:bizId:bizType:` 索引缓存基于此建立 |
| `idx_biz` | 普通 | `biz_type`, `biz_id` | 按业务对象维度统计点赞（如某条回复的点赞数） |

关系（无数据库外键，仅逻辑关联）：

- `like_record.user_id` → user 域 `user.id`
- `like_record.(biz_type, biz_id)` → 多态关联：`biz_type = 'reply'` → 回复、`'note'` → 笔记、`'question'` → 问答

#### Scenario: 多态关联复用

- **WHEN** 新业务对象需要点赞能力
- **THEN** 业务侧自行约定 `bizType` 取值即可复用 `like_record`，remark 无需改表
- **AND** remark 不感知具体业务表结构，也不做跨服务校验

### Requirement: 存储与缓存不变量

remark 服务的存储层 SHALL 满足以下不变量：

- 软取消语义：`like_record` 无 `deleted` 逻辑删除列，也没有 `creater` / `updater`（与 auth / promotion 等域不同）；取消点赞把 `liked` 置 0 而非删除行，因此唯一索引不会被反复占用与释放，重新点赞只需再改回 1
- `uk_user_biz` 唯一索引是「同一用户对同一业务对象最多一行点赞」的数据库级保障，也是 `Like` 幂等 upsert 的并发兜底（见「点赞 upsert 规则」）
- `liked` 列语义：1=已点赞（首次 Insert 默认值即 1，或由 0 Update 而来），0=已取消（由 1 Update 而来，行仍保留不删除）

goctl 生成的基础方法（`likeRecordModel` 接口）缓存策略：

| 方法 | 缓存策略 |
|------|---------|
| `Insert` / `Update` | 写后失效 `cache:likeRecord:id:` 与 `cache:likeRecord:userId:bizId:bizType:` |
| `FindOne(ctx, id)` | 走 `cache:likeRecord:id:` 主键缓存 |
| `FindOneByUserIdBizIdBizType` | 走 `cache:likeRecord:userId:bizId:bizType:` 索引缓存 → 回主键查询（`Like` 逻辑判定「是否已有记录」的入口） |
| `Delete(ctx, id)` | 写后失效上述两个键（业务层未使用，取消点赞走 Update） |

- 扩展方法 `FindLikedBizIds` 刻意不走缓存（`QueryRowsNoCacheCtx`）：命中行数通常远小于入参 bizIds，且业务对实时性要求高，走 NoCache 查询避免缓存击穿与脏读；空入参短路不发 SQL；按 `bizIds` 长度动态生成 `?` 占位符，参数化查询无 SQL 注入风险；`liked = 1` 过滤保证取消的行不出现
- `parseTime=true` 为连接串必需项：`LikeRecord` 结构体的 `CreateTime` / `UpdateTime` 声明为 `time.Time`，未开启时驱动会返回 `[]byte` 导致扫描失败
- `DataSource` 与 `Cache` 未标记 `optional`（与 promotion 不同），缺失时服务启动即失败——配置缺陷能在启动阶段暴露
- ServiceContext 只注入一个 Model：`model.NewLikeRecordModel(sqlx.NewMysql(c.DataSource), c.Cache)`

连接与依赖配置（`remark.yaml`）：

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `Name` / `ListenOn` | `remark.rpc` / `0.0.0.0:8093` | RPC 服务名与监听地址 |
| `Etcd.Hosts[0]` / `Etcd.Key` | `127.0.0.1:2379` / `remark.rpc` | 服务注册 |
| `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_remark?charset=utf8mb4&parseTime=true&loc=Local` | MySQL 连接串 |
| `Cache[0]` | `127.0.0.1:6379`，Pass 空，Type `node` | Redis 缓存（单机模式） |

#### Scenario: 唯一索引兜底并发点赞

- **WHEN** 并发首次点赞同一业务对象且两请求的查询均未命中
- **THEN** Insert 阶段被 `uk_user_biz` 唯一索引拦截，仅一条成功
- **AND** 失败方的错误直接返回给调用方，代码未做重复键捕获与降级

#### Scenario: 点赞后列表页立即可见

- **WHEN** 用户点赞后立即批量查询点赞态
- **THEN** `FindLikedBizIds` 以 `QueryRowsNoCacheCtx` 直读 DB，返回包含最新点赞的 `liked=1` 子集
- **AND** 带缓存的 `FindOneByUserIdBizIdBizType` / `FindOne` 在 Insert / Update 后由 goctl 自动失效 `cache:likeRecord:id:` 与 `cache:likeRecord:userId:bizId:bizType:` 两个缓存键

