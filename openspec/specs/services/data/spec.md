# data（统计大屏）Specification

## Purpose

data（统计大屏）微服务面向运营统计大屏提供三组展示数据的读写能力：看板图表数据（ECharts 结构）、今日实时指标（访问量 / 订单金额 / 订单数 / 新增学员数）与 Top10 课程榜单（热门榜 / 热销榜）。HTTP 服务 `data-api` 监听 8811，RPC 服务 `data.rpc` 监听 `0.0.0.0:8091`（gRPC），经 etcd（`127.0.0.1:2379`，key `data.rpc`）注册与发现。本服务无独立数据库：无 MySQL、无 Redis、无 model 层，目标存储为 Redis（设计意图，推导，当前未落地）；与其他服务的关系上，代码中仅有 data-api 自身消费 `data.rpc`（6 个 HTTP 接口与 6 个 RPC 方法一一对应转发），无任何上游 RPC 调用或 MQ 消费者——上游统计来源（trade 域订单统计、user 域新增学员、网关访问量统计）按设计意图经离线任务 / MQ 回填 `Set*` 系列 RPC，wiki 拓扑标注的 `data-rpc --> AuthRPC` 依赖在代码中不存在。12 个 logic（RPC 6 + API 6）按 business-rules v1.2 复核均已实现并编译通过。

## Requirements

### Requirement: HTTP 看板数据接口（ECharts）

data-api SHALL 提供看板图表数据的获取与设置接口：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /data/board | 看板数据获取，响应 `R{data: R«EchartsVO»}` |
| PUT | /data/board/set | 看板数据设置，请求体 `BoardDataSetDTO`，响应 `R{data: R}` |

补充约束：

- 响应统一为 `pkg/response.R{Code,Msg,RequestId,Data any}`；错误码遵循共享错误码规范
- 两接口在聚合文档 `docs/tjxt.openapi.json` 中认证列均标注「否」（全服务无鉴权，见「接口鉴权状态」）
- 核心规则（📋 设计意图，契约推导）：服务端直接吐出 ECharts 可用的图表结构，前端不做二次组装

#### Scenario: 大屏图表渲染

- **WHEN** 前端带 `types` 调用 GET /data/board
- **THEN** 返回 ECharts 可直接使用的 `xAxis` / `yAxis` / `series` 结构，前端无需二次组装

#### Scenario: 看板数据回填

- **WHEN** 离线统计任务或运营后台按单个 `type` 携带 `version` 调用 PUT /data/board/set
- **THEN** 写入该类型下的数值序列，供读侧按当前版本取数

### Requirement: HTTP 今日数据接口

data-api SHALL 提供今日实时统计指标的获取与设置接口：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /data/today | 获取今日数据，响应 `R{data: R«TodayDataVO»}` |
| PUT | /data/today/set | 设置线上数据（openapi 原文措辞），请求体 `TodayDataDTO`，响应 `R{data: R}` |

补充约束：

- `TodayDataVO` 含四项固定指标：`visits`（访问量）/ `orderAmount`（订单金额）/ `orderNum`（订单数）/ `stuNewNum`（新增学员数）
- 大屏轮询刷新为典型场景：前端定时调 GET /data/today
- 两接口认证列在聚合文档中均标注「否」

#### Scenario: 大屏轮询刷新

- **WHEN** 前端定时轮询 GET /data/today
- **THEN** 返回今日访问量 / 订单金额 / 订单数 / 新增学员数四项实时指标

#### Scenario: 设置线上数据

- **WHEN** 回填方携带 `version` 与四项指标调用 PUT /data/today/set
- **THEN** 写入最新一版今日数据
- **AND** `version` 的生成规则在 proto 与代码中均未定义

### Requirement: HTTP Top10 榜单接口

data-api SHALL 提供热门榜与热销榜数据的获取与设置接口：

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /data/top10 | top10数据获取，响应 `R{data: R«Top10DataVO»}` |
| PUT | /data/top10/set | 设置top10数据，请求体 `Top10DataSetDTO`，响应 `R{data: R}` |

#### Scenario: 榜单展示

- **WHEN** 前端调用 GET /data/top10
- **THEN** 返回 `hot`（热门课程榜）与 `hotSales`（热销课程榜）两组课程数据

#### Scenario: 批量写入榜单

- **WHEN** 回填方携带 `version` 与扁平课程列表调用 PUT /data/top10/set
- **THEN** 批量写入 Top10 数据
- **AND** 请求未区分 hot / hotSales，写入时如何分流到两个榜单在 proto 层面未表达

### Requirement: RPC 服务契约

`Data` 服务 SHALL 经 etcd 注册（key `data.rpc`，监听 `0.0.0.0:8091`）对外提供 6 个 RPC 方法，按 3 个业务域分组：

| 业务域 | 方法 | 请求 Message | 响应 Message | 说明 |
|--------|------|-------------|-------------|------|
| 看板数据 | GetBoardData | `BoardDataReq{types}` | `EchartsVO{xAxis, yAxis, series}` | 按数据类型列表拉取看板图表数据 |
| 看板数据 | SetBoardData | `BoardDataSetReq{version, type, data}` | `OkReply{success}` | 写入某一类看板数据 |
| 今日数据 | GetTodayData | `Empty{}` | `TodayDataVO{visits, orderAmount, orderNum, stuNewNum}` | 获取今日实时统计 |
| 今日数据 | SetTodayData | `TodayDataSetReq{version, visits, orderAmount, orderNum, stuNewNum}` | `OkReply{success}` | 写入今日实时统计 |
| Top10 数据 | GetTop10Data | `Empty{}` | `Top10DataVO{hot, hotSales}` | 获取热门榜与热销榜 |
| Top10 数据 | SetTop10Data | `Top10DataSetReq{version, data}` | `OkReply{success}` | 批量写入 Top10 数据 |

关键字段约定：

- `BoardDataReq.types`：repeated int32，需要拉取的看板数据类型列表；`BoardDataSetReq`：version（int32 数据版本号）、type（看板数据类型）、data（repeated double 该类型下的数值序列）
- `AxisVO`：type（轴类型）、max/min（double）、average、data（repeated string 轴刻度）、interval（刻度间隔）
- `SerierVO`：name、type（系列图形类型）、data（repeated string）、max/min（**字符串类型**，与 `AxisVO` 的 double 不同）；Message 名为 `SerierVO`（proto 原始拼写，非 `SeriesVO`），api 侧类型名相同
- `CourseInfo`（响应）与 `Top10DataSetUnit`（请求）字段完全一致：category / name / newStuNum（int32）/ orderAmount（double）
- `GetTodayData` 入参为 `Empty`，无任何日期 / 维度参数——「今日」由服务端自行判定
- 三个 `Set*` 请求均带 `int32 version`，三个 `Get*` 均不带——读侧固定取「当前版本」

消费关系：当前仅 `data-api` 自身消费该 RPC（HTTP Handler → `dataclient.Data`，`servicecontext.go` 注入 `DataRpc`，import 路径 `tjxt/apps/data/rpc/data/client/data`）；全仓 Grep 除 `apps/data/` 自身外无其它服务引用。`specs/architecture/service-topology` 中标注的 `data-rpc --> AuthRPC` 依赖在代码中不存在（`apps/data/rpc/data/internal/config/config.go` 只有 `zrpc.RpcServerConf`，无任何客户端配置）。

> 注：data 是本仓唯一把 RPC 客户端生成到 `client/<svc>/` 子目录的服务（其余服务为 `apps/<svc>/rpc/<svc>/`），因此 import 路径为 `.../rpc/data/client/data` 而非 `.../rpc/data`；proto 文件路径 `apps/data/rpc/data/data.proto` 也比其余 12 个服务多一层 `data/` 目录。

#### Scenario: 大屏轮询刷新链路

- **WHEN** 前端定时调 GET /data/today
- **THEN** data-api 转发 `GetTodayData`（proto 侧传 `Empty{}`），返回访问量 / 订单额 / 订单数 / 新增学员

#### Scenario: 数据回填与版本切换（推导）

- **WHEN** 离线统计任务或运营后台调三个 `Set*` 方法并携带 `version` 写入最新一版数据
- **THEN** 推测用于「写新版本 → 读侧按最新版本取数」的原子切换
- **AND** proto 的 Get 侧没有 `version` 入参，读写版本如何协商未定义

### Requirement: 看板数据业务规则

看板数据（ECharts）SHALL 遵循以下规则（📋 设计意图，契约推导；data 服务无 DDL 可作为字段语义依据，推导可信度低于 message / search 两服务，建议对照源码最终确认）：

- 读写粒度不对称：写按单个 `type`（`BoardDataSetReq{type, data []}`），读按 `types` 批量（`BoardDataReq{types []}`）
- 响应即 ECharts 配置：`EchartsVO{xAxis, yAxis, series}` 字段名与 ECharts option 一致
- 轴元信息服务端组装：`SetBoardData` 只收 `repeated double` 裸数值，`AxisVO` 的 type / max / min / average / interval 在读路径上按数据现算或按 `type` 查预置模板，不落存储
- 写入带版本、读取不带版本：`BoardDataSetReq.version` 用于多版本写入后原子切换（推导），`BoardDataReq` 只有 `types`，读侧固定取「当前版本」
- 数值精度差异：`AxisVO.max/min` 为 double，`SerierVO.max/min` 为 string，序列化时需注意
- 看板数据类型 `type` / `types` 为 int32 枚举值，proto 中无注释、无 enum 定义、无 DDL 可查，具体取值含义当前无从考证

#### Scenario: GetBoardData 读取流程

- **WHEN** 收到 `GetBoardData` 请求（校验 types 非空）
- **THEN** 读当前版本指针，按 types 逐个取回该类型的数值序列
- **AND** 按 type 组装 xAxis / yAxis / series（含 max/min/average/interval 计算），返回 EchartsVO

#### Scenario: 按类型分片写入

- **WHEN** 收到 `SetBoardData` 请求（单个 type + 一组 double 数值）
- **THEN** 按单个 `type` 分片写入该类数据（携带 version）
- **AND** 轴配置与系列元信息不落存储，由服务端在读路径上组装

### Requirement: 今日数据业务规则

今日数据 SHALL 遵循以下规则（📋 设计意图，契约推导）：

- 四项固定指标的快照读写，无维度、无分页：`TodayDataVO{visits, orderAmount, orderNum, stuNewNum}`（访问量 / 订单金额 / 订单数 / 新增学员数）
- 「今日」由服务端判定：`GetTodayData` 入参为 `Empty`，无日期参数，跨天切换规则未定义
- 金额用 double：`visits` / `orderAmount` 为 double，金额未用整数分表示，存在浮点精度风险
- 访问量也是 double：`visits` 为 double 而非整型（proto 如此定义，语义上应为计数）
- 写入带版本：`TodayDataSetReq.version`，同看板数据

#### Scenario: SetTodayData 写入流程

- **WHEN** 收到 `SetTodayData` 请求（校验 version，生成规则未定义）
- **THEN** 写入四项指标并更新当前版本指针
- **AND** 返回 `OkReply{success: true}`

#### Scenario: 跨天语义

- **WHEN** 大屏跨天后继续调用 GET /data/today
- **THEN** 「今日」由服务端自行判定（请求无日期入参）
- **AND** 跨天切换规则未定义

### Requirement: Top10 数据业务规则

Top10 数据 SHALL 遵循以下规则（📋 设计意图，契约推导）：

- 两个榜单：`Top10DataVO{hot, hotSales}`——热门榜 / 热销榜，由一份扁平课程数据派生
- 写入不区分榜单：`Top10DataSetReq{version, data []Top10DataSetUnit}` 为单一扁平列表，无 hot / hotSales 标记；推测按 `newStuNum` 降序取热门、按 `orderAmount` 降序取热销，proto 未表达
- 读写结构体重复定义：`CourseInfo`（响应）与 `Top10DataSetUnit`（请求）字段完全一致（category / name / newStuNum / orderAmount 四字段）
- 课程以名称标识：两结构体均无 courseId，仅有 category + name，无法回链 course 域，也无法去重
- 数量上限：服务名与摘要均为 Top10，proto 未做 repeated 长度约束，需在 logic 层截断至 10 条

#### Scenario: 扁平数据派生双榜单

- **WHEN** 收到 `GetTop10Data` 请求（入参 `Empty{}`）
- **THEN** 从同一份数据按派生规则返回 hot 与 hotSales 两组榜单
- **AND** 派生规则未定义（推测按 `newStuNum` 排热门、按 `orderAmount` 排热销）

#### Scenario: 榜单条目截断

- **WHEN** 写入的课程条目超过 10 条
- **THEN** logic 层需截断至 10 条（proto 未做 repeated 长度约束）

### Requirement: 接口鉴权状态

data SHALL 保持「全仓唯一一个 HTTP 接口完全不鉴权」的现状，并按以下规则约束后续演进（📋 设计意图，契约推导）：

- 无 JWT 保护：`apps/data/api/data.api` 三个 `@server` 块只声明 `group`，未声明 `jwt: Auth`（对照 message / search 的 `@server (jwt: Auth, group: xxx)`）
- 路由未挂中间件：`apps/data/api/data/internal/handler/routes.go` 三处 `server.AddRoutes` 均无 `rest.WithJwt(...)`
- 配置里却有 Auth 段：`data-api.yaml` 有 `Auth.AccessSecret` / `Auth.AccessExpire`，config.go 有对应字段——配置已备好但代码从未使用（悬空项）
- 写接口同样无保护：三个 `PUT /data/*/set` 均无鉴权，任何人可覆盖大屏数据，属安全缺口

> ⚠️ 风险提示（源文档原文）：三个 set 接口是写入口，当前无任何身份校验。实现时应至少为 `PUT /data/*/set` 补上 `jwt: Auth` 并限定管理员角色；GET 侧若面向公开大屏可保持匿名。

#### Scenario: 写入口未鉴权风险

- **WHEN** 任意调用方在无 token 情况下请求 PUT /data/today/set
- **THEN** 请求不被拦截，可直接覆盖大屏数据（当前实现的安全缺口）
- **AND** 改进方向：为写接口补 `jwt: Auth` 并限定管理员角色，GET 侧若面向公开大屏可保持匿名

#### Scenario: 六接口鉴权现状

- **WHEN** 检查 GET/PUT /data/board、GET/PUT /data/today、GET/PUT /data/top10 六个接口
- **THEN** 均未声明 `jwt: Auth`、未挂 `rest.WithJwt(...)` 中间件（聚合文档认证列为「否」）

### Requirement: API 层与 RPC 层关系

data-api SHALL 仅做 DTO 转换并全部转发 RPC，遵循以下规则（📋 设计意图，契约推导）：

- API 全部转发 RPC：`apps/data/api/data/internal/svc/servicecontext.go` 注入 `DataRpc`，API 层无存储访问
- 接口数一一对应：HTTP 6 个 ↔ RPC 6 个，data 是唯一一个 API 与 RPC 方法数完全相等的服务（message 18/19、search 2/2）
- 类型基本同构：api 的 `EchartsVO` / `AxisVO` / `SerierVO` / `TodayDataVO` / `Top10DataVO` 与 proto 同名同构，转换为逐字段直映射
- 整型宽度差异：api 用 `int`（平台相关），proto 用 `int32`，转换时需显式收窄
- 空入参处理：api 的 `GetTodayData()` / `GetTop10Data()` 无参数，proto 侧需传 `Empty{}`
- `omitempty` 差异：api 的 `AxisVO` / `SerierVO` 全字段带 `omitempty`，`EchartsVO` 三字段不带——零值轴配置会从 JSON 中消失

#### Scenario: 同名转发链路

- **WHEN** API 层处理任一 /data 请求
- **THEN** 经 `DataRpc` 转发同名 RPC 方法，api 类型与 proto 类型逐字段直映射（int 显式收窄至 int32）

#### Scenario: 空入参构造

- **WHEN** api 层调用 `GetTodayData()` / `GetTop10Data()`（无参数）
- **THEN** RPC 层构造 `Empty{}` 作为请求消息

### Requirement: 数据模型与聚合模式

data 服务 SHALL 不依赖任何持久化组件（无 MySQL / 无 Redis / 无 model 层），以「上游统计回填 → 大屏快照读写」的聚合模式提供数据存取：

- 无关系型表：`sql/ddl/` 下 14 个文件中无 `tj_data.sql`（其余 12 个服务均有对应 DDL）
- 无 MySQL 连接：`apps/data/rpc/data/etc/data.yaml` 全文仅 `Name` / `ListenOn` / `Etcd` 三段，无 `DataSource`
- 无 Redis 连接：`data.yaml` 无 `Cache` / `Redis` 段（`00-architecture/overview.md` 服务表中 data 的「数据库」列亦为 `-`）
- 无 model 层：`apps/data/rpc/data/internal/` 下只有 `config` / `logic` / `server` / `svc`，无 model 目录；`ServiceContext` 仅有 `Config` 一个字段
- 结论（data-model.md 原文）：截至当前提交，data 服务没有接入任何持久化组件，6 个 RPC 方法与 6 个 HTTP 接口全部是 goctl 占位实现，返回零值
  - 注（源冲突）：business-rules.md v1.2（2026-08-06 复核，go build 全模块通过）将 12 个 logic 状态校正为「已实现」；「无持久化组件 / 无 DDL」为两篇一致的事实，实现状态以较新复核为准

数据来源与形态（📋 设计意图，推导，当前未落地）——大屏数据具备「高频读、低频写、可容忍丢失、无需事务」特征，配合 proto 的 `version` 字段，推导目标存储为 Redis（本仓其余服务已在用 Redis，docker-compose 依赖已就绪）：

| 数据集 | 写入方 | 读取方 | 推导存储形态 |
|--------|--------|--------|-------------|
| 今日数据 | SetTodayData | GetTodayData | Hash 或 JSON String，单键 |
| 看板数据 | SetBoardData（按 type 逐类写） | GetBoardData（按 types 批量读） | 按 type 分键，读时 MGET |
| Top10 数据 | SetTop10Data（整体覆盖） | GetTop10Data | JSON String 或 List，单键 |

Redis Key 设计（📋 设计意图，推导，代码中不存在，供实现参考）：`data:today:{version}`（Hash，field: visits / orderAmount / orderNum / stuNewNum）、`data:board:{version}:{type}`（String(JSON)，值为 repeated double 序列）、`data:top10:{version}`（String(JSON)，值为 Top10DataSetUnit 数组）、`data:version:current`（String 版本指针，三个 `Set*` 完成后更新，所有 `Get*` 先读此键，用于「全量写新版 → 原子切读」）。

未定义的关键点（实现前需确认）：过期策略（proto 无 TTL 相关字段，「今日数据」跨天后如何失效未定义）、version 生成（调用方给定还是服务端自增）、历史版本清理（旧 version 的 key 何时删除）、hot / hotSales 分流规则、数据来源（各项指标由谁统计后回填，代码中无任何上游调用或 MQ 消费者）。

上游数据来源（设计意图，代码中不存在）：trade 域订单统计、user 域新增学员、网关访问量统计 → 离线任务 / MQ → `Set*` 系列 RPC 回填。

#### Scenario: 数据存取现状

- **WHEN** 追查任一 Get / Set 方法的存储访问
- **THEN** 服务既无 MySQL 也无 Redis 连接、无 model 层，`ServiceContext` 只有 Config 一个字段
- **AND** 按 data-model.md 记载接口为 goctl 占位实现返回零值（business-rules v1.2 复核校正为已实现；「无存储组件」为两篇一致）

#### Scenario: 设计目标的版本切换

- **WHEN** 按 Redis 设计意图落地数据存取
- **THEN** 三个 `Set*` 完成后更新 `data:version:current` 版本指针，所有 `Get*` 先读此键，实现「全量写新版 → 原子切读」

#### Scenario: 存储扩展待补齐

- **WHEN** 需要为 data 服务接入真实存储
- **THEN** 若采用 Redis：在 `data.yaml` 增加 `Redis` 段、`internal/config/config.go` 增加 `redis.RedisConf` 字段、`ServiceContext` 注入 `*redis.Redis`
- **AND** 若改用 MySQL：补 `sql/ddl/tj_data.sql` 并按约定生成 model

### Requirement: 服务配置与已知结构性缺口

data SHALL 按以下配置部署运行：

| 配置项 | 取值 |
|--------|------|
| data-api Name / Host / Port | data-api / 0.0.0.0 / 8811（HTTP） |
| data.rpc Name / ListenOn | data.rpc / 0.0.0.0:8091（gRPC） |
| etcd | 127.0.0.1:2379（API 侧 `DataRpc.Etcd`，RPC 侧 `Etcd`），key `data.rpc` |
| Auth.AccessSecret | `change-me-in-production`（占位值，生产须替换） |
| Auth.AccessExpire | 7200（秒） |
| DataSource / Cache | 无——`data.yaml` 缺失这两段（与「不使用关系型表、不依赖 Model」一致，但意味着目前 RPC 没有可落地的数据源） |
| Telemetry / Prometheus / Log | 可观测性注入配置已注入勿删（统一约定见 `specs/infra/observability/spec.md`） |

已知结构性缺口（configs.md ⚠️ 逐条如实记录）：

1. 目录多套一层：`data.proto` 在 `apps/data/rpc/data/data.proto`，API 代码在 `apps/data/api/data/...`，比其余 12 个服务的约定（`apps/data/api/`、`apps/data/rpc/`）多嵌套一层 `data/`
2. 三份 go.mod：`apps/data/api/go.mod`（module 名 `api`，仅 3 行、无 require，残缺占位 / 遗留废弃文件）；`apps/data/api/data/go.mod`（module `tjxt/apps/data/api/data`）；`apps/data/rpc/data/go.mod`（module `tjxt/apps/data/rpc/data`）；go.work 仅引用后两者
3. 无 DDL：`sql/ddl/` 下存在 `tj_message.sql` / `tj_message_model.sql` / `tj_search.sql`，但没有 `tj_data.sql`
4. 配置缺 DataSource / Cache：`data.yaml` 仅有 Name / ListenOn / Etcd 三段，RPC 暂无可用数据源配置
5. JWT 配置悬空：`data-api.yaml` 声明了 `Auth.AccessSecret` / `AccessExpire`，但 `data.api` 路由无 `jwt: Auth`，鉴权未生效

#### Scenario: 端口与服务发现

- **WHEN** 部署 data 服务
- **THEN** data-api 监听 8811（HTTP 0.0.0.0）、data.rpc 监听 8091（gRPC 0.0.0.0）
- **AND** API 与 RPC 经 etcd 服务发现对接（`DataRpc.Etcd.Key: data.rpc`，etcd 127.0.0.1:2379）

#### Scenario: 构建 data 模块

- **WHEN** 在 go.work 下构建 data 服务
- **THEN** 引用的真实模块为 `./apps/data/api/data` 与 `./apps/data/rpc/data`
- **AND** `apps/data/api/go.mod` 不参与编译（孤儿模块，`cd apps/data/api` 编译报错，建议清理）

#### Scenario: JWT 密钥占位

- **WHEN** 生产环境部署 data-api
- **THEN** 必须替换 `Auth.AccessSecret` 默认占位值 `change-me-in-production`
- **AND** 当前路由未启用 JWT，该配置为悬空项（启用前需先为路由补 `jwt: Auth`）

