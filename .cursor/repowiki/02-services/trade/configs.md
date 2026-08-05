> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/trade/api/etc/trade-api.yaml`, `apps/trade/rpc/etc/trade.yaml`

---

# Trade Configs

## API 服务配置 (`apps/trade/api/etc/trade-api.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `trade-api` | - | 服务名称 |
| `Host` | `0.0.0.0` | - | 监听地址 |
| `Port` | `8810` | - | 监听端口 |
| `Auth.AccessSecret` | `change-me-in-production` | - | JWT 签名密钥 |
| `Auth.AccessExpire` | `7200` | - | 访问令牌有效期（秒） |
| `TradeRpc.Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `TradeRpc.Etcd.Key` | `trade.rpc` | - | trade RPC 服务发现 key |

**依赖的外部服务**：
- 自身 RPC `trade.rpc`（HTTP handler → RPC client 调用）
- etcd `127.0.0.1:2379` — 服务发现

**配置结构体**：`apps/trade/api/internal/config/config.go`

```go
type Config struct {
    rest.RestConf
    Auth struct {
        AccessSecret string
        AccessExpire int64
    }
    // TradeRpc 交易域 RPC 客户端（通过 etcd 服务发现）
    TradeRpc zrpc.RpcClientConf
}
```

> API 层全部路由均声明 `@server (jwt: Auth)`，`Auth.AccessSecret` 必须与签发方 `auth.rpc` 的 `Jwt.AccessSecret` 保持一致，否则 token 校验失败。

---

## RPC 服务配置 (`apps/trade/rpc/etc/trade.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `trade.rpc` | - | RPC 服务名 |
| `ListenOn` | `0.0.0.0:8084` | - | RPC 监听地址 |
| `Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `Etcd.Key` | `trade.rpc` | - | 服务注册 key |
| `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_trade?charset=utf8mb4&parseTime=true&loc=Local` | - | MySQL 连接串 |
| `Cache[0].Host` | `127.0.0.1:6379` | - | Redis 地址 |
| `Cache[0].Type` | `node` | - | Redis 缓存类型（单机模式） |
| `Cache[0].Pass` | (空) | - | Redis 密码 |
| `RabbitMQ.Host` | `127.0.0.1` | - | RabbitMQ 主机 |
| `RabbitMQ.Port` | `5672` | - | RabbitMQ 端口 |
| `RabbitMQ.User` | `rabbitmq` | - | RabbitMQ 用户名 |
| `RabbitMQ.Pass` | `rabbitmq` | - | RabbitMQ 密码 |
| `RabbitMQ.Exchange` | `order.exchange` | - | 订单事件交换机 |
| `RabbitMQ.PayRoutingKey` | `order.pay` | - | 支付成功事件路由键 |
| `RabbitMQ.RefundRoutingKey` | `order.refund` | - | 退款事件路由键 |
| `PayRpc.Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `PayRpc.Etcd.Key` | `pay.rpc` | - | pay RPC 服务发现 key |

**依赖的外部服务**：
- MySQL `tj_trade` 库
- Redis 缓存（节点模式）
- RabbitMQ — 发布支付成功/退款事件到 learning
- `pay.rpc` — 支付/退款下单、关单、支付/退款结果查询

**配置结构体**：`apps/trade/rpc/internal/config/config.go`

```go
type Config struct {
    zrpc.RpcServerConf
    DataSource string
    Cache      cache.CacheConf

    // RabbitMQ 用于支付成功/退款事件发布到 learning
    RabbitMQ struct {
        Host             string
        Port             int
        User             string
        Pass             string
        Exchange         string
        PayRoutingKey    string
        RefundRoutingKey string
    }

    // PayRpc 调用 pay 服务（支付/退款下单、关单、支付/退款结果查询）
    PayRpc zrpc.RpcClientConf
}
```

---

## 配置项说明

### 数据库连接串

```
root:0000@tcp(127.0.0.1:3306)/tj_trade?charset=utf8mb4&parseTime=true&loc=Local
```

- 用户名: `root`
- 密码: `0000`
- 地址: `127.0.0.1:3306`
- 数据库: `tj_trade`
- 字符集: `utf8mb4`
- 时区: `Local`

对应 `sql/ddl/tj_trade.sql` 中的 5 张表：`cart`、`order`、`order_detail`、`refund_apply`、`undo_log`。

### 端口分配

| 服务 | 端口 | 协议 |
|------|------|------|
| `trade-api` | 8810 | HTTP |
| `trade.rpc` | 8084 | gRPC |

### RabbitMQ 事件拓扑

trade 是事件**生产者**，learning 是**消费者**，两端交换机与路由键必须一致：

| 项 | trade（生产端） | learning（消费端） |
|----|----------------|-------------------|
| 交换机 | `Exchange: order.exchange` | `PayExchange: order.exchange` / `RefundExchange: order.exchange` |
| 支付路由键 | `PayRoutingKey: order.pay` | `PayRoutingKey: order.pay` |
| 退款路由键 | `RefundRoutingKey: order.refund` | `RefundRoutingKey: order.refund` |
| 队列 | —（生产端不声明队列） | `PayQueue: learning.lesson.pay.queue` / `RefundQueue: learning.lesson.refund.queue` |

**Producer 容错**：`apps/trade/rpc/internal/svc/servicecontext.go:41-46` 中，DSN 拼接为 `amqp://{User}:{Pass}@{Host}:{Port}/`。若 `mq.NewProducer` 失败，仅打印 `init rabbitmq producer failed, will skip event publish` 日志，`MQProducer` 保持 nil，**服务照常启动**。业务代码发布事件前必须判空。

### 与外部服务的连接

| 服务 | 连接方式 | 说明 |
|------|---------|------|
| MySQL (tj_trade) | DataSource 配置 | 自建存储：购物车、订单、明细、退款申请 |
| Redis | Cache 配置 | goctl model 缓存层 |
| RabbitMQ | RabbitMQ 配置 | 发布 `order.pay` / `order.refund` 事件 |
| pay.rpc | RpcClient 配置 | 支付渠道管理、支付下单、退款下单、结果查询 |
| etcd | Etcd 配置 | 服务注册与发现 |

### JWT 密钥

生产环境**必须修改**默认值 `change-me-in-production`，否则 JWT 可被伪造。

`trade-api` 只做 token 校验（不签发），`Auth.AccessSecret` 需与 `apps/auth/rpc/etc/auth.yaml` 的 `Jwt.AccessSecret` 一致。

---

## 已知配置缺口

| 缺口 | 说明 |
|------|------|
| **API 层未配置 course/promotion RPC** | `trade-api.yaml` 只声明了 `TradeRpc`。但预下单需要课程信息与优惠券试算（`OrderConfirmVO.courses` / `discounts`），实现时需补 `CourseRpc` / `PromotionRpc` 客户端配置，或在 RPC 层补相应依赖。 |
| **RPC 层未配置 course RPC** | `trade.yaml` 只声明了 `PayRpc`。`CartVO.course_name` / `cover_url` / `now_price` 与 `OrderCourseVO` 均需课程侧数据，当前无 course 服务连接配置。 |
| **无 MQ 消费端配置** | trade 只有 Producer 侧的 `Exchange` + RoutingKey，没有 Queue 配置，符合「trade 只生产不消费」的定位。 |
