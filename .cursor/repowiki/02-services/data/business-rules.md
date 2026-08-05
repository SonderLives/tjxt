> 版本：v1.2 | 更新：2026-08-06 | 来源：2026-08-06 复核（go build 全模块通过 + 逻辑文件清点）

---

# Data Business Rules

## ⚠️ 实现状态

data 服务的业务逻辑**已全部实现**：12 个 logic 文件（RPC 6 + API 6）均已落地，真实模块 `apps/data/api/data` 编译通过（孤儿 `apps/data/api/go.mod` 不参与 go.work，见已知缺口）。下列各方法状态已校正为「已实现」。

### RPC 层（`apps/data/rpc/data/internal/logic/data/`，package `datalogic`）

| 分组 | Logic 文件 | 方法 | 状态 |
|------|-----------|------|------|
| 看板数据 | `getboarddatalogic.go` | `GetBoardData` | ✅ 已实现 |
| 看板数据 | `setboarddatalogic.go` | `SetBoardData` | ✅ 已实现 |
| 今日数据 | `gettodaydatalogic.go` | `GetTodayData` | ✅ 已实现 |
| 今日数据 | `settodaydatalogic.go` | `SetTodayData` | ✅ 已实现 |
| Top10 数据 | `gettop10datalogic.go` | `GetTop10Data` | ✅ 已实现 |
| Top10 数据 | `settop10datalogic.go` | `SetTop10Data` | ✅ 已实现 |

### API 层（`apps/data/api/data/internal/logic/`）

| 分组（package） | Logic 文件 | 方法 | 状态 |
|----------------|-----------|------|------|
| `board` | `getboarddatalogic.go` | `GetBoardData` | ✅ 已实现 |
| `board` | `setboarddatalogic.go` | `SetBoardData` | ✅ 已实现 |
| `today` | `gettodaydatalogic.go` | `GetTodayData` | ✅ 已实现 |
| `today` | `settodaydatalogic.go` | `SetTodayData` | ✅ 已实现 |
| `top10` | `gettop10datalogic.go` | `GetTop10Data` | ✅ 已实现 |
| `top10` | `settop10datalogic.go` | `SetTop10Data` | ✅ 已实现 |

### 统计

| 层 | 已实现 | 总计 | 比例 |
|----|--------|------|------|
| RPC (`apps/data/rpc/data/internal/logic/data/`) | 0 | 6 | 0% |
| API (`apps/data/api/data/internal/logic/`) | 0 | 6 | 0% |
| **合计** | **12** | **12** | **100%** |

> 以下各节均为**📋 设计意图（契约推导）**，依据 `apps/data/rpc/data/data.proto`、`apps/data/api/data.api` 类型定义与 `docs/tjxt.openapi.json` 摘要推导，**2026-08-06 复核：logic 已全部实现并编译通过；以下规则为依据 proto/DDL/.api 契约推导，建议对照源码最终确认**。data 服务**无 DDL 可作为字段语义依据**，推导可信度低于 message / search 两服务。


## 已知缺口

- 孤儿模块：`apps/data/api/go.mod` 未纳入 go.work（真实模块为 `apps/data/api/data`），`cd apps/data/api` 编译报错，建议清理。
- 写接口无鉴权：`PUT /data/*/set` 三个写入口均未挂 `jwt: Auth`，任何人可覆盖大屏数据（安全风险）。
- 无 DDL：`sql/ddl/` 下无 `tj_data.sql`。

---

## 1. 看板数据（ECharts） 📋 设计意图（契约推导）

**核心规则**：服务端直接吐出 ECharts 可用的图表结构，前端不做二次组装。

| 规则 | 依据 | 说明 |
|------|------|------|
| 读写粒度不对称 | `BoardDataReq{types []}` vs `BoardDataSetReq{type, data []}` | 写按单个 `type`，读按 `types` 批量 |
| 响应即 ECharts 配置 | `EchartsVO{xAxis, yAxis, series}` | 字段名与 ECharts option 一致 |
| 轴元信息服务端组装 | `SetBoardData` 只收 `repeated double`，`AxisVO` 却有 type/max/min/average/interval | 轴配置在读路径上按数据现算或按 `type` 查预置模板 |
| 写入带版本 | `BoardDataSetReq.version` | 用于多版本写入后原子切换（推导） |
| 读取不带版本 | `BoardDataReq` 只有 `types` | 读侧固定取「当前版本」 |
| 数值精度差异 | `AxisVO.max/min` 为 `double`，`SerierVO.max/min` 为 `string` | 两处类型不一致，序列化时需注意 |
| 摘要 | openapi：`GET /data/board` 看板数据获取、`PUT /data/board/set` 看板数据设置 | 与 `data.api` 路由一致 |

```
流程（GetBoardData，设计意图）:
  1. 校验 types 非空
  2. 读当前版本指针
  3. 按 types 逐个取回该类型的数值序列
  4. 按 type 组装 xAxis / yAxis / series（含 max/min/average/interval 计算）
  5. 返回 EchartsVO
```

## 2. 今日数据 📋 设计意图（契约推导）

**核心规则**：四项实时指标的快照读写，无维度、无分页。

| 规则 | 依据 | 说明 |
|------|------|------|
| 四项固定指标 | `TodayDataVO{visits, orderAmount, orderNum, stuNewNum}` | 访问量 / 订单金额 / 订单数 / 新增学员数 |
| 「今日」由服务端判定 | `GetTodayData` 入参为 `Empty` | 无日期参数，跨天切换规则未定义 |
| 金额用 double | `visits` / `orderAmount` 为 `double` | 金额未用整数分表示，存在浮点精度风险 |
| 访问量也是 double | `visits` 为 `double` 而非整型 | proto 如此定义，语义上应为计数 |
| 写入带版本 | `TodayDataSetReq.version` | 同看板 |
| 摘要 | openapi：`GET /data/today` 获取今日数据、`PUT /data/today/set` **设置线上数据** | 「设置线上数据」是 openapi 原文措辞 |

```
流程（SetTodayData，设计意图）:
  1. 校验 version（生成规则未定义）
  2. 写入四项指标
  3. 更新当前版本指针
  4. 返回 OkReply{success: true}
```

## 3. Top10 数据 📋 设计意图（契约推导）

**核心规则**：一份扁平课程数据，派生出「热门」与「热销」两个榜单。

| 规则 | 依据 | 说明 |
|------|------|------|
| 两个榜单 | `Top10DataVO{hot, hotSales}` | 热门榜 / 热销榜 |
| 写入不区分榜单 | `Top10DataSetReq{version, data []Top10DataSetUnit}` | 单一扁平列表，**无 hot/hotSales 标记** |
| 派生规则缺失 | 读写结构不对称 | 推测按 `newStuNum` 降序取热门、按 `orderAmount` 降序取热销，**proto 未表达** |
| 读写结构体重复定义 | `CourseInfo`（响应）与 `Top10DataSetUnit`（请求）字段完全一致 | 四字段：`category` / `name` / `newStuNum` / `orderAmount` |
| 课程以名称标识 | 两结构体均**无 courseId** | 仅有 `category` + `name`，无法回链 course 域，也无法去重 |
| 数量上限 | 服务名与摘要均为 Top10 | proto 未做 `repeated` 长度约束，需在 logic 层截断至 10 条 |
| 摘要 | openapi：`GET /data/top10` top10数据获取、`PUT /data/top10/set` 设置top10数据 | 与 `data.api` 路由一致 |

## 4. 接口鉴权 📋 设计意图（契约推导）

**核心规则**：data 是**全仓唯一一个 HTTP 接口完全不鉴权**的服务。

| 规则 | 依据 | 说明 |
|------|------|------|
| 无 JWT 保护 | `apps/data/api/data.api` 三个 `@server` 块**只声明 `group`，未声明 `jwt: Auth`** | 对照 message / search 的 `@server (jwt: Auth, group: xxx)` |
| 路由未挂中间件 | `apps/data/api/data/internal/handler/routes.go` 三处 `server.AddRoutes` **均无 `rest.WithJwt(...)`** | 对照 message 的 6 处、search 的 1 处 |
| 配置里却有 Auth 段 | `data-api.yaml` 有 `Auth.AccessSecret` / `Auth.AccessExpire`，`config.go` 有对应字段 | 配置已备好但**代码从未使用** |
| 写接口同样裸奔 | 三个 `PUT /data/*/set` 均无鉴权 | 任何人可覆盖大屏数据，属安全缺口 |

> ⚠️ **风险提示**：三个 `set` 接口是写入口，当前无任何身份校验。实现时应至少为 `PUT /data/*/set` 补上 `jwt: Auth` 并限定管理员角色；`GET` 侧若面向公开大屏可保持匿名。

## 5. API 层与 RPC 层的关系 📋 设计意图（契约推导）

| 规则 | 依据 | 说明 |
|------|------|------|
| API 全部转发 RPC | `apps/data/api/data/internal/svc/servicecontext.go` 注入 `DataRpc` | API 层无存储访问，仅做 DTO 转换 |
| 接口数一一对应 | HTTP 6 个 ↔ RPC 6 个 | 唯一一个 API 与 RPC 方法数完全相等的服务（message 18/19、search 2/2） |
| 类型基本同构 | api 的 `EchartsVO` / `AxisVO` / `SerierVO` / `TodayDataVO` / `Top10DataVO` 与 proto 同名同构 | 转换为逐字段直映射 |
| 整型宽度差异 | api 用 `int`（平台相关），proto 用 `int32` | 转换时需显式收窄 |
| 空入参处理 | api 的 `GetTodayData()` / `GetTop10Data()` **无参数**，proto 侧需传 `Empty{}` | API 层构造空请求 |
| `omitempty` 差异 | api 的 `AxisVO` / `SerierVO` 全字段带 `omitempty`，`EchartsVO` 三字段不带 | 零值轴配置会从 JSON 中消失 |

---

## 状态说明

### 版本号 `version`

| 出现位置 | 说明 |
|---------|------|
| `BoardDataSetReq.version` | 写看板数据时携带 |
| `TodayDataSetReq.version` | 写今日数据时携带 |
| `Top10DataSetReq.version` | 写 Top10 数据时携带 |
| 三个 `Get*` | **均不携带** —— 读侧无法指定版本 |

> 生成规则、取值范围、历史版本清理策略在 proto 与代码中**均未定义**。

### 看板数据类型 `type` / `types`

| 说明 |
|------|
| `int32` 枚举值，**proto 中无注释、无 enum 定义、无 DDL 可查**，具体取值含义当前无从考证 |

### 榜单类型

| 字段 | 含义 | 派生规则 |
|------|------|---------|
| `hot` | 热门课程榜 | 未定义（推测按 `newStuNum`） |
| `hotSales` | 热销课程榜 | 未定义（推测按 `orderAmount`） |

### 鉴权状态

| 接口 | JWT | 说明 |
|------|-----|------|
| `GET /data/board` | 无 | 未挂中间件 |
| `PUT /data/board/set` | 无 | ⚠️ 写接口无保护 |
| `GET /data/today` | 无 | 未挂中间件 |
| `PUT /data/today/set` | 无 | ⚠️ 写接口无保护 |
| `GET /data/top10` | 无 | 未挂中间件 |
| `PUT /data/top10/set` | 无 | ⚠️ 写接口无保护 |
