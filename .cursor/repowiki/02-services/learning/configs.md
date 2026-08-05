> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/learning/api/etc/learning-api.yaml`, `apps/learning/rpc/etc/learning.yaml`

---

# Learning Configs

## API 服务配置 (`apps/learning/api/etc/learning-api.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `learning-api` | - | 服务名称 |
| `Host` | `0.0.0.0` | - | 监听地址 |
| `Port` | `8813` | - | 监听端口 |
| `Auth.AccessSecret` | `change-me-in-production` | - | JWT 签名密钥 |
| `Auth.AccessExpire` | `7200` | - | 访问令牌有效期（秒） |
| `LearningRpc.Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `LearningRpc.Etcd.Key` | `learning.rpc` | - | learning RPC 服务发现 key |
| `CourseRpc.Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `CourseRpc.Etcd.Key` | `course.rpc` | - | course RPC 服务发现 key |

**依赖的外部服务**：
- 自身 RPC `learning.rpc`（HTTP handler → RPC client 调用）
- `course.rpc` — 课程域 RPC。yaml 内注释：「课程目录展示时，需要查询课程名/章节数等基础信息」
- etcd `127.0.0.1:2379` — 服务发现

**配置结构体**：`apps/learning/api/internal/config/config.go`

```go
type Config struct {
    rest.RestConf
    Auth struct {
        AccessSecret string
        AccessExpire int64
    }
    // LearningRpc 学习域 RPC 客户端（通过 etcd 服务发现）
    LearningRpc zrpc.RpcClientConf
    // CourseRpc 用于补全 LearningLessonVO 中 course_name/cover 等课程侧字段
    CourseRpc zrpc.RpcClientConf
}
```

> API 层全部路由均声明 `@server (jwt: Auth)`（见 `apps/learning/api/learning.api:95-97`），`Auth.AccessSecret` 必须与签发方 `auth.rpc` 的 `Jwt.AccessSecret` 保持一致，否则 token 校验失败。

**CourseRpc 的必要性**：`learning_lesson` 表不含课程侧数据，以下 `LearningLessonVO` 字段必须由 `CourseRpc` 补全：

| 字段 | 说明 |
|------|------|
| `course_name` | 课程名称 |
| `course_cover_url` | 课程封面 |
| `course_amount` | 课程数量 |
| `sections` | 课程总章节数 |
| `latest_section_name` | 最近学习小节名称 |
| `latest_section_index` | 最近学习小节序号 |

---

## RPC 服务配置 (`apps/learning/rpc/etc/learning.yaml`)

| 配置项 | 默认值 | 环境变量 | 说明 |
|--------|--------|----------|------|
| `Name` | `learning.rpc` | - | RPC 服务名 |
| `ListenOn` | `0.0.0.0:8085` | - | RPC 监听地址 |
| `Etcd.Hosts[0]` | `127.0.0.1:2379` | - | etcd 地址 |
| `Etcd.Key` | `learning.rpc` | - | 服务注册 key |
| `DataSource` | `root:0000@tcp(127.0.0.1:3306)/tj_learning?charset=utf8mb4&parseTime=true&loc=Local` | - | MySQL 连接串 |
| `Cache[0].Host` | `127.0.0.1:6379` | - | Redis 地址 |
| `Cache[0].Type` | `node` | - | Redis 缓存类型（单机模式） |
| `Cache[0].Pass` | (空) | - | Redis 密码 |
| `RabbitMQ.Host` | `127.0.0.1` | - | RabbitMQ 主机 |
| `RabbitMQ.Port` | `5672` | - | RabbitMQ 端口 |
| `RabbitMQ.User` | `rabbitmq` | - | RabbitMQ 用户名 |
| `RabbitMQ.Pass` | `rabbitmq` | - | RabbitMQ 密码 |
| `RabbitMQ.PayQueue` | `learning.lesson.pay.queue` | - | 支付成功事件消费队列 |
| `RabbitMQ.RefundQueue` | `learning.lesson.refund.queue` | - | 退款事件消费队列 |
| `RabbitMQ.PayExchange` | `order.exchange` | - | 支付事件交换机 |
| `RabbitMQ.RefundExchange` | `order.exchange` | - | 退款事件交换机 |
| `RabbitMQ.PayRoutingKey` | `order.pay` | - | 支付事件路由键 |
| `RabbitMQ.RefundRoutingKey` | `order.refund` | - | 退款事件路由键 |

**依赖的外部服务**：
- MySQL `tj_learning` 库
- Redis 缓存（节点模式）
- RabbitMQ — 消费 trade 发布的支付成功/退款事件

**配置结构体**：`apps/learning/rpc/internal/config/config.go`

```go
// Config learning.rpc 配置
type Config struct {
    zrpc.RpcServerConf
    DataSource string
    Cache      cache.CacheConf

    // RabbitMQ 消费：订单支付/退款事件触达 learning 开通/撤销课程
    RabbitMQ struct {
        Host             string
        Port             int
        User             string
        Pass             string
        PayQueue         string
        RefundQueue      string
        PayExchange      string
        RefundExchange   string
        PayRoutingKey    string
        RefundRoutingKey string
    }
}
```

> **注意**：learning RPC **不配置任何下游 RpcClient**。它是交易链路的终点，只被 learning-api 调用，以及被 MQ 事件驱动。

---

## 配置项说明

### 数据库连接串

```
root:0000@tcp(127.0.0.1:3306)/tj_learning?charset=utf8mb4&parseTime=true&loc=Local
```

- 用户名: `root`
- 密码: `0000`
- 地址: `127.0.0.1:3306`
- 数据库: `tj_learning`
- 字符集: `utf8mb4`
- 时区: `Local`

对应 `sql/ddl/tj_learning.sql` 中的**唯一一张表**：`learning_lesson`。

### 端口分配

| 服务 | 端口 | 协议 |
|------|------|------|
| `learning-api` | 8813 | HTTP |
| `learning.rpc` | 8085 | gRPC |

### RabbitMQ 事件拓扑

learning 是事件**消费者**，trade 是**生产者**，两端交换机与路由键必须一致：

| 项 | trade（生产端 `trade.yaml`） | learning（消费端 `learning.yaml`） |
|----|---------------------------|----------------------------------|
| 支付交换机 | `Exchange: order.exchange` | `PayExchange: order.exchange` |
| 退款交换机 | `Exchange: order.exchange` | `RefundExchange: order.exchange` |
| 支付路由键 | `PayRoutingKey: order.pay` | `PayRoutingKey: order.pay` |
| 退款路由键 | `RefundRoutingKey: order.refund` | `RefundRoutingKey: order.refund` |
| 支付队列 | —（生产端不声明队列） | `PayQueue: learning.lesson.pay.queue` |
| 退款队列 | —（生产端不声明队列） | `RefundQueue: learning.lesson.refund.queue` |

**事件到业务方法的映射（设计意图）**：

| 队列 | 触发的 RPC 方法 | 效果 |
|------|----------------|------|
| `learning.lesson.pay.queue` | `GrantCourses(user_id, course_ids)` | 幂等写入 `learning_lesson`，`status=0` |
| `learning.lesson.refund.queue` | `RevokeCourses(user_id, course_ids)` | `status=3`（失效）、`plan_status=0` |

### 缓存策略

`Cache` 配置传给 `model.NewLearningLessonModel(conn, c.Cache)`，但 `learninglessonmodel.go` 中 **11 个自定义扩展方法全部使用 `ExecNoCacheCtx` / `QueryRowNoCacheCtx` / `QueryRowsNoCacheCtx` 绕过缓存**。缓存实际只在 gen 层的 `FindOne` / `FindOneByUserIdCourseId` / `Insert` / `Update` / `Delete` 生效。

### 与外部服务的连接

| 服务 | 连接方式 | 说明 |
|------|---------|------|
| MySQL (tj_learning) | DataSource 配置 | 自建存储：`learning_lesson` 单表 |
| Redis | Cache 配置 | goctl model 缓存层（自定义方法绕过） |
| RabbitMQ | RabbitMQ 配置 | 消费 `order.pay` / `order.refund` 事件 |
| course.rpc | API 层 RpcClient 配置 | 补全课程名/封面/章节数 |
| etcd | Etcd 配置 | 服务注册与发现 |

### JWT 密钥

生产环境**必须修改**默认值 `change-me-in-production`，否则 JWT 可被伪造。

`learning-api` 只做 token 校验（不签发），`Auth.AccessSecret` 需与 `apps/auth/rpc/etc/auth.yaml` 的 `Jwt.AccessSecret` 一致。

---

## 已知配置缺口

| 缺口 | 说明 |
|------|------|
| **MQ 消费端未接线** | `learning.yaml` 完整配置了 6 项 RabbitMQ 参数，`config.go` 也声明了对应结构体，但 `apps/learning/rpc/internal/svc/servicecontext.go` 只初始化了 `LearningLessonModel` 与 `LearningService`，**无 Consumer 创建代码**；`apps/learning/rpc/learning.go` 主入口也未启动消费协程。配置项目前处于「已声明未使用」状态。 |
| **无 MQ 容错策略** | 对照 trade 的 `servicecontext.go:42-46`（MQ 初始化失败仅记日志、`MQProducer` 置 nil、不阻塞启动），learning 侧尚无对应的降级设计。 |
| **RPC 层无 CourseRpc** | 课程侧字段补全只在 API 层配置了 `CourseRpc`。若后续需要在 RPC 层直接返回完整 `LearningLessonVO`，需在 `learning.yaml` 与 `config.go` 中补 `CourseRpc zrpc.RpcClientConf`。 |
| **无积分服务配置** | `PlanPageReply.week_points` 需要积分数据，`learning.yaml` 与 `learning-api.yaml` 均无相关服务连接配置。 |
