> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/pay/api/etc/pay-api.yaml`, `apps/pay/rpc/etc/pay.yaml`

---

# Pay Configs

## API 服务配置 (`apps/pay/api/etc/pay-api.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `pay-api` | - | 服务名称 |
| `Host` | `0.0.0.0` | - | 监听地址 |
| `Port` | `8808` | - | 监听端口 |
| `Auth.AccessSecret` | `change-me-in-production` | - | JWT 签名密钥 |
| `Auth.AccessExpire` | `7200` | - | 访问令牌有效期（秒） |
| `PayRpc.Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `PayRpc.Etcd.Key` | `pay.rpc` | - | pay RPC 服务发现 key |

**依赖的外部服务**：
- 自身 RPC `pay.rpc`（HTTP handler → RPC client 调用）
- etcd `127.0.0.1:2379` — 服务发现

> 结构体定义见 `apps/pay/api/internal/config/config.go`：内嵌 `rest.RestConf`，`Auth` 为匿名结构体，`PayRpc` 为 `zrpc.RpcClientConf`。API 层为纯转发层，13 个 logic 全部直接调用 `l.svcCtx.PayRpc.*`。

---

## RPC 服务配置 (`apps/pay/rpc/etc/pay.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `pay.rpc` | - | RPC 服务名 |
| `ListenOn` | `0.0.0.0:8088` | - | RPC 监听地址 |
| `Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `Etcd.Key` | `pay.rpc` | - | 服务注册 key |
| `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_pay?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai` | - | MySQL 连接串 |
| `Cache[0].Host` | `127.0.0.1:6379` | - | Redis 地址 |
| `Cache[0].Pass` | (空) | - | Redis 密码 |
| `Cache[0].Type` | `node` | - | Redis 缓存类型（单机模式） |
| `TablePrefix` | (空字符串) | - | 表名前缀，当前未被任何代码读取 |

**依赖的外部服务**：
- MySQL `tj_pay` 库
- Redis 缓存（节点模式）
- etcd `127.0.0.1:2379` — 服务注册

> 结构体定义见 `apps/pay/rpc/internal/config/config.go`：内嵌 `zrpc.RpcServerConf`，外加 `DataSource string`、`Cache cache.CacheConf`、`TablePrefix string`。

---

## 配置项说明

### 数据库连接串

```
root:0000@tcp(127.0.0.1:3306)/tj_pay?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai
```

- 用户名: `root`
- 密码: `0000`
- 地址: `127.0.0.1:3306`
- 数据库: `tj_pay`
- 字符集: `utf8mb4`
- 时区: `Asia/Shanghai`（URL 编码为 `Asia%2FShanghai`）

> **⚠️ 与其它服务不一致**：`auth.rpc`、`user.rpc`、`trade.rpc` 均使用 `loc=Local`，唯独 pay 服务显式指定 `loc=Asia%2FShanghai`。支付涉及 `pay_over_time`（超时时间）、`pay_success_time`（成功时间）等强时间语义字段，跨服务对账时需注意这一时区配置差异。
>
> `parseTime=true` 为必需项：`PayOrder.PayOverTime` / `CreateTime` 声明为 `time.Time`，`PaySuccessTime` 为 `sql.NullTime`，均需驱动直接解析。
>
> DDL 中三张表的字符集为 `utf8` / `utf8_general_ci`，而连接串声明 `charset=utf8mb4`，存在字符集不一致（不影响纯数字与 ASCII 场景）。

### Redis 缓存

`Cache` 由 `svc.NewServiceContext` 一次性传给三个模型：

| 模型 | 缓存使用情况 |
|------|------------|
| `PayChannelModel` | goctl 生成的 `FindOne(id)` 带缓存；手写的 `FindAllEnabled` / `FindByCode` / `PageList` 全部 `NoCache` |
| `PayOrderModel` | goctl 生成的 `FindOne` / `FindOneByBizOrderNo` / `FindOneByPayOrderNo` 带缓存；手写的 `MarkTo*` 系列全部 `NoCache` |
| `RefundOrderModel` | goctl 生成的 `FindOne(id)` 带缓存；手写的所有查询与状态流转方法全部 `NoCache` |

> **⚠️ 缓存一致性缺口**：`MarkToSuccess` / `MarkToClosed` 等状态流转方法用 `ExecNoCacheCtx` 直改数据库，未调用 `DelCacheCtx`，而 `FindOneByBizOrderNo` / `FindOneByPayOrderNo` 是带缓存查询。详见 business-rules.md 第 3 节。

### TablePrefix

`TablePrefix` 在 `config.Config` 中已声明并在 `pay.yaml` 中置为空串，但**全项目无任何读取点**——模型层的表名由 goctl 生成时硬编码（如 `` `pay_order` ``）。属于预留但未接线的配置项。

### 端口分配

| 服务 | 端口 | 协议 |
|------|------|------|
| `pay-api` | 8808 | HTTP |
| `pay.rpc` | 8088 | gRPC |

### 与外部服务的连接

| 服务 | 连接方式 | 说明 |
|------|---------|------|
| MySQL (tj_pay) | DataSource 配置 | 自建存储，含 `pay_channel` / `pay_order` / `refund_order` 三张表 |
| Redis | Cache 配置 | 缓存渠道与支付/退款单主键查询 |
| etcd | Etcd 配置 | RPC 注册与发现 |
| 第三方支付网关 | **无配置项** | 当前为 mock 实现（`mockQrCodeUrl`），未接入微信/支付宝，故无 appId / 商户号 / 私钥 / 回调验签等配置 |

### 作为被依赖方

`pay.rpc` 被以下服务在配置中声明为客户端：

| 消费方配置文件 | 配置节 | Key |
|---------------|--------|-----|
| `apps/pay/api/etc/pay-api.yaml` | `PayRpc` | `pay.rpc` |
| `apps/trade/rpc/etc/trade.yaml` | `PayRpc` | `pay.rpc` |

> `trade.rpc` 已完成配置与 `servicecontext.go` 装配，但其 logic 层尚无实际调用点。
