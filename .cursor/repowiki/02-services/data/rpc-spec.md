> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/data/rpc/data/data.proto`

---

# Data RPC Spec

## 服务名

`Data` — 统计大屏数据微服务，通过 etcd 服务发现（key: `data.rpc`），监听 `0.0.0.0:8091`。

> ⚠️ proto 文件路径为 `apps/data/rpc/data/data.proto`，比其余 12 个服务多一层 `data/` 目录。详见 [configs.md](./configs.md) 的「⚠️ 已知结构性问题」一节。

## RPC 方法总览

### 看板数据（ECharts）

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `GetBoardData` | `BoardDataReq { types }` | `EchartsVO { xAxis, yAxis, series }` | 按数据类型列表拉取看板图表数据 |
| `SetBoardData` | `BoardDataSetReq { version, type, data }` | `OkReply { success }` | 写入某一类看板数据 |

**请求字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `types` | repeated int32 | 需要拉取的看板数据类型列表 |
| `version` | int32 | 数据版本号 |
| `type` | int32 | 看板数据类型 |
| `data` | repeated double | 该类型下的数值序列 |

**`EchartsVO` 结构**（直接对齐 ECharts 配置项）：

| 字段 | 类型 | 说明 |
|------|------|------|
| `xAxis` | repeated AxisVO | X 轴配置 |
| `yAxis` | repeated AxisVO | Y 轴配置 |
| `series` | repeated SerierVO | 数据系列 |

**`AxisVO` 字段**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `type` | string | 轴类型 |
| `max` / `min` | double | 轴最大值 / 最小值 |
| `average` | double | 平均值 |
| `data` | repeated string | 轴刻度数据 |
| `interval` | double | 刻度间隔 |

**`SerierVO` 字段**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `name` | string | 系列名称 |
| `type` | string | 系列图形类型 |
| `data` | repeated string | 系列数据 |
| `max` / `min` | string | 系列最大值 / 最小值（**字符串类型**，与 `AxisVO` 的 double 不同） |

> Message 名为 `SerierVO`（proto 中的拼写，非 `SeriesVO`），api 侧类型名同样为 `SerierVO`。

---

### 今日数据

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `GetTodayData` | `Empty {}` | `TodayDataVO { visits, orderAmount, orderNum, stuNewNum }` | 获取今日实时统计 |
| `SetTodayData` | `TodayDataSetReq { version, visits, orderAmount, orderNum, stuNewNum }` | `OkReply { success }` | 写入今日实时统计 |

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `version` | int32 | 数据版本号（仅 Set 请求有） |
| `visits` | double | 访问量 |
| `orderAmount` | double | 订单金额 |
| `orderNum` | int32 | 订单数 |
| `stuNewNum` | int32 | 新增学员数 |

> `GetTodayData` 入参为 `Empty`，无任何日期/维度参数——「今日」由服务端自行判定。

---

### Top10 数据

| 方法名 | 请求 Message | 响应 Message | 说明 |
|--------|-------------|-------------|------|
| `GetTop10Data` | `Empty {}` | `Top10DataVO { hot, hotSales }` | 获取热门榜与热销榜 |
| `SetTop10Data` | `Top10DataSetReq { version, data }` | `OkReply { success }` | 批量写入 Top10 数据 |

**`Top10DataVO` 结构**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `hot` | repeated CourseInfo | 热门课程榜 |
| `hotSales` | repeated CourseInfo | 热销课程榜 |

**`CourseInfo` / `Top10DataSetUnit` 字段**（两者字段完全一致）：

| 字段 | 类型 | 说明 |
|------|------|------|
| `category` | string | 课程分类 |
| `name` | string | 课程名称 |
| `newStuNum` | int32 | 新增学员数 |
| `orderAmount` | double | 订单金额 |

> proto 中为读、写各定义了一个结构相同的 message（`CourseInfo` 用于响应、`Top10DataSetUnit` 用于请求）。
> `Top10DataSetReq` 只有一个扁平的 `data` 列表，**未区分 hot / hotSales**，写入时如何分流到两个榜单在 proto 层面未表达。

---

## 被哪些 API 服务消费

| 消费方 | 调用方式 | 说明 |
|--------|---------|------|
| `data-api` (自身 API 层) | HTTP Handler → `dataclient.Data` RPC | `apps/data/api/data/internal/svc/servicecontext.go` 注入 `DataRpc`，import 路径为 `tjxt/apps/data/rpc/data/client/data` |

（注：全仓 Grep `apps/data/rpc/data/client` / `dataclient` / `DataRpc`，除 `apps/data/` 自身外**无其它服务引用**。`.cursor/repowiki/00-architecture/service-topology.md` 中标注的 `data-rpc --> AuthRPC` 依赖，在代码中不存在——`apps/data/rpc/data/internal/config/config.go` 只有 `zrpc.RpcServerConf`，无任何客户端配置。）

> data 是本仓唯一把 RPC 客户端生成到 `client/<svc>/` 子目录的服务（其余服务为 `apps/<svc>/rpc/<svc>/`），因此 import 路径为 `.../rpc/data/client/data` 而非 `.../rpc/data`。

---

## 调用典型场景

1. **大屏轮询刷新** → 前端定时调 `GET /data/today` → `data-api` 转发 `GetTodayData` → 返回访问量/订单额/订单数/新增学员
2. **图表渲染** → 前端带 `types` 调 `GET /data/board` → 转发 `GetBoardData` → 直接返回 ECharts 可用的 `xAxis`/`yAxis`/`series` 结构，前端无需二次组装
3. **榜单展示** → 调 `GET /data/top10` → 转发 `GetTop10Data` → 返回 `hot` 与 `hotSales` 两组课程
4. **数据回填** → 离线统计任务或运营后台调三个 `Set*` 方法，带 `version` 写入最新一版数据
5. **版本切换（推导）** → 三个 Set 请求均含 `version` 字段，推测用于「写新版本 → 读侧按最新版本取数」的原子切换，但 proto 的 Get 侧没有 `version` 入参，读写版本如何协商未定义

---

## 自定义 Model 方法

**无**。`apps/data/rpc/data/` 下**不存在 `internal/model/` 目录**，`ServiceContext` 中也只有 `Config` 一个字段：

```go
type ServiceContext struct {
	Config config.Config
}
```

data 服务既无 MySQL 连接（`data.yaml` 无 `DataSource`），也无 Redis 连接（`data.yaml` 无 `Cache`），因此没有任何 model 层。数据存取方式详见 [data-model.md](./data-model.md)。
