> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/data/rpc/data/etc/data.yaml`, `apps/data/rpc/data/internal/svc/servicecontext.go`, `apps/data/rpc/data/data.proto`

---

# Data Data Model

## ⚠️ 无关系型表

data 服务**没有任何 MySQL 库与关系型表**，据实核查如下：

| 项 | 实际情况 | 核查方式 |
|----|---------|---------|
| DDL 文件 | **不存在** | `sql/ddl/` 下 14 个文件中无 `tj_data.sql`（其余 12 个服务均有对应 DDL） |
| MySQL 连接 | **未配置** | `apps/data/rpc/data/etc/data.yaml` 全文仅 `Name` / `ListenOn` / `Etcd` 三段，无 `DataSource` |
| 配置结构体 | **无存储字段** | `apps/data/rpc/data/internal/config/config.go` 仅 `zrpc.RpcServerConf` |
| model 目录 | **不存在** | `apps/data/rpc/data/internal/` 下只有 `config` / `logic` / `server` / `svc`，无 `model` |
| ServiceContext | **只有 Config** | `type ServiceContext struct { Config config.Config }` |
| Redis 连接 | **未配置** | `data.yaml` 无 `Cache` / `Redis` 段；`00-architecture/overview.md` 的服务表中 data 的「数据库」列亦为 `-` |

**结论**：截至当前提交，data 服务**没有接入任何持久化组件**——既无 MySQL，也无 Redis。6 个 RPC 方法与 6 个 HTTP 接口全部是 goctl 占位实现，返回零值。下文「Redis Key 设计」一节为按 proto 语义反推的设计意图，**不代表现有代码**。

---

## 数据来源与形态（📋 设计意图，推导，当前未落地）

大屏数据具备「高频读、低频写、可容忍丢失、无需事务」的特征，配合 proto 中的 `version` 字段，推导其目标存储为 **Redis**（本仓其余服务已在用 Redis，`docker-compose` 依赖已就绪）：

| 数据集 | 写入方 | 读取方 | 推导存储形态 |
|--------|--------|--------|-------------|
| 今日数据 | `SetTodayData` | `GetTodayData` | Hash 或 JSON String，单键 |
| 看板数据 | `SetBoardData`（按 `type` 逐类写） | `GetBoardData`（按 `types` 批量读） | 按 `type` 分键，读时 MGET |
| Top10 数据 | `SetTop10Data`（整体覆盖） | `GetTop10Data` | JSON String 或 List，单键 |

---

## Redis Key 设计（📋 设计意图，推导，当前未落地）

> 以下 key 命名**代码中不存在**，是依据 proto 的 `version` / `type` / `types` 字段与 go-zero 常见约定推导的方案，供实现时参考。

| 推导 Key | 类型 | 写入方法 | 读取方法 | 说明 |
|---------|------|---------|---------|------|
| `data:today:{version}` | Hash | `SetTodayData` | `GetTodayData` | field: `visits` / `orderAmount` / `orderNum` / `stuNewNum` |
| `data:board:{version}:{type}` | String(JSON) | `SetBoardData` | `GetBoardData` | 每次写入一个 `type`，值为 `repeated double` 序列 |
| `data:top10:{version}` | String(JSON) | `SetTop10Data` | `GetTop10Data` | 值为 `Top10DataSetUnit` 数组 |
| `data:version:current` | String | 三个 `Set*` 完成后更新 | 所有 `Get*` 先读此键 | 版本指针，用于「全量写新版 → 原子切读」 |

**推导依据**：

1. 三个 `Set*` 请求（`BoardDataSetReq` / `TodayDataSetReq` / `Top10DataSetReq`）**都带 `int32 version`**，而三个 `Get*` **都不带 version**（`GetBoardData` 只带 `types`，另两个入参是 `Empty`）——说明读侧不指定版本，只能由服务端维护一个「当前版本」指针，即上表的 `data:version:current`。
2. `BoardDataSetReq` 是**单个 type + 一组 double**，而 `BoardDataReq` 是**一组 types**——说明看板数据按 `type` 分片存储，读时批量取回后再拼装成 `EchartsVO`。
3. `EchartsVO` 的 `xAxis` / `yAxis` / `series` 是前端 ECharts 的原生结构，而 `SetBoardData` 只收 `repeated double` 裸数值——说明**轴配置与系列元信息由服务端在读路径上组装**，不落存储。
4. `Top10DataSetReq.data` 是扁平列表且无 hot/hotSales 标记，而 `Top10DataVO` 分两个榜——说明服务端需要按某种规则（如 `newStuNum` 排热门、`orderAmount` 排热销）从同一份数据派生出两个榜单。

**未定义的关键点**（实现前需要确认）：

| 缺口 | 说明 |
|------|------|
| 过期策略 | proto 无 TTL 相关字段，「今日数据」跨天后如何失效未定义 |
| version 生成 | 由调用方给定还是服务端自增，未定义 |
| 历史版本清理 | 旧 version 的 key 何时删除，未定义 |
| hot / hotSales 分流规则 | `Top10DataSetReq` 无区分字段，派生规则未定义 |
| 数据来源 | 各项指标（访问量、订单额、新增学员）由谁统计后回填，代码中无任何上游调用或 MQ 消费者 |

---

## 关系图

```
（无关系型表，无表间关系）

data-api :8818  ──HTTP──▶  data.rpc :8088  ──▶  ✗ 无任何存储组件
                                                （设计意图：Redis）

上游数据来源（设计意图，代码中不存在）：
  trade 域订单统计   ─┐
  user 域新增学员    ─┼─▶  离线任务/MQ  ─▶  Set* 系列 RPC
  网关访问量统计     ─┘
```

---

## ⚠️ 已知结构性问题

data 服务的目录布局与其余 12 个服务的约定**不一致**，在阅读本篇涉及的路径时需特别注意：

| 约定（其余 12 服务） | data 服务实际 |
|---------------------|--------------|
| `apps/<svc>/rpc/internal/model/` | **不存在**（无 model 层） |
| `apps/<svc>/rpc/<svc>.proto` | `apps/data/rpc/data/data.proto`（多一层 `data/`） |
| `sql/ddl/tj_<svc>.sql` | **不存在** |

完整的结构性偏差清单（含双 go.mod 问题）见 [configs.md](./configs.md) 的「⚠️ 已知结构性问题」一节。

---

## 模型扩展模式

**不适用**。data 服务无 `internal/model/` 目录、无 `*_gen.go`、无 `vars.go`，不存在 goctl model 生成与手写扩展的分层。

📋 **待补齐（设计意图）**：若采用 Redis 方案，需在 `apps/data/rpc/data/etc/data.yaml` 增加 `Redis` 段、在 `internal/config/config.go` 增加 `redis.RedisConf` 字段、在 `ServiceContext` 注入 `*redis.Redis`；若改用 MySQL，则需补 `sql/ddl/tj_data.sql` 并按约定生成 model。
