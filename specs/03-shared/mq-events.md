# RabbitMQ 事件定义与约定

> 版本：v1.0 | 更新：2026-08-05 | 来源：`pkg/mq/event/`

## 事件总线拓扑

```
┌─────────────┐     trade.events      ┌──────────────┐
│ trade-rpc   │ ─────────────────────▶ │ learning-rpc │ (解锁课程)
│ (发布)      │   order.created       │ (消费)       │
└─────────────┘                       └──────────────┘
       │                                      │
       │ payment.paid                         │ course.completed
       ▼                                      ▼
┌─────────────┐     pay.events        ┌──────────────┐
│ pay-rpc     │ ─────────────────────▶ │ message-rpc  │ (发送通知)
│ (发布)      │   payment.paid        │ (消费)       │
└─────────────┘                       └──────────────┘
       │
       │ refund.processed
       ▼
┌─────────────┐
│ promotion-  │ (回收优惠券)
│ rpc         │
└─────────────┘
```

## 核心事件定义

### 交易域 (trade.events)

| 事件 | RoutingKey | 发布者 | 消费者 | 说明 |
|------|------------|--------|--------|------|
| 订单创建 | `order.created` | trade-rpc | learning-rpc, message-rpc | 下单成功，解锁课程+发通知 |
| 订单支付 | `order.paid` | trade-rpc | promotion-rpc, learning-rpc | 支付成功，核销优惠券、更新学习权限 |
| 订单取消 | `order.cancelled` | trade-rpc | promotion-rpc | 释放占用的优惠券/库存 |
| 退款发起 | `refund.initiated` | trade-rpc | message-rpc | 发起退款通知 |
| 退款完成 | `refund.completed` | trade-rpc | promotion-rpc, learning-rpc | 退款成功，回收权益 |

### 支付域 (pay.events)

| 事件 | RoutingKey | 发布者 | 消费者 | 说明 |
|------|------------|--------|--------|------|
| 支付成功 | `payment.paid` | pay-rpc | trade-rpc | 回调确认，驱动订单状态流转 |
| 支付失败 | `payment.failed` | pay-rpc | trade-rpc, message-rpc | 失败通知 |
| 对账差异 | `reconciliation.mismatch` | pay-rpc | - | 告警，人工处理 |

### 学习域 (learning.events)

| 事件 | RoutingKey | 发布者 | 消费者 | 说明 |
|------|------------|--------|--------|------|
| 课程完成 | `course.completed` | learning-rpc | message-rpc, promotion-rpc | 结业证书、积分奖励 |
| 进度更新 | `progress.updated` | learning-rpc | - | 内部统计用 |
| 打卡完成 | `checkin.completed` | learning-rpc | message-rpc | 连续打卡奖励 |

### 优惠券域 (promotion.events)

| 事件 | RoutingKey | 发布者 | 消费者 | 说明 |
|------|------------|--------|--------|------|
| 券发放 | `coupon.issued` | promotion-rpc | message-rpc | 发放通知 |
| 券核销 | `coupon.used` | promotion-rpc | - | 统计用 |
| 券过期 | `coupon.expired` | promotion-rpc | message-rpc | 过期提醒 |

### 消息域 (message.events)

| 事件 | RoutingKey | 发布者 | 消费者 | 说明 |
|------|------------|--------|--------|------|
| 站内信发送 | `inbox.send` | message-rpc | - | 内部投递 |
| 短信发送 | `sms.send` | message-rpc | - | 第三方短信通道 |

## 事件结构体规范 (Go)

所有事件定义在 `pkg/mq/event/`：

```go
// pkg/mq/event/order.go
type OrderCreatedEvent struct {
    OrderId     int64   `json:"orderId"`
    UserId      int64   `json:"userId"`
    CourseIds   []int64 `json:"courseIds"`   // 涉及的课程
    Amount      int64   `json:"amount"`      // 分
    CouponId    int64   `json:"couponId,omitempty"`
    IdempotencyKey string `json:"idempotencyKey"` // 幂等键
    Timestamp   int64   `json:"timestamp"`
}

type OrderPaidEvent struct {
    OrderId         int64   `json:"orderId"`
    UserId          int64   `json:"userId"`
    CourseIds       []int64 `json:"courseIds"`
    PaymentId       int64   `json:"paymentId"`
    PaidAt          int64   `json:"paidAt"`
    IdempotencyKey  string  `json:"idempotencyKey"`
}
```

## 消费者实现约定

```go
// Logic 中注册消费者
func (l *OrderLogic) RegisterConsumers() {
    l.svcCtx.MqClient.Consume("learning.order.created", l.handleOrderCreated)
}

// 处理器签名
func (l *OrderLogic) handleOrderCreated(ctx context.Context, body []byte) error {
    var evt event.OrderCreatedEvent
    if err := json.Unmarshal(body, &evt); err != nil {
        return err
    }
    // 幂等校验
    if l.checkProcessed(ctx, evt.IdempotencyKey) {
        return nil // 已处理，直接 ack
    }
    // 业务处理
    return l.processOrderCreated(ctx, &evt)
}
```

## 运维规范

| 规则 | 说明 |
|------|------|
| 消费者必须幂等 | 基于 `IdempotencyKey` 或业务主键去重 |
| 手动 ACK | 处理成功才 ack，失败 requeue 或进死信队列 |
| 重试策略 | 最大 3 次，指数退避，最终进 DLQ 人工处理 |
| 监控指标 | 消费延迟、积压量、失败率、DLQ 堆积 |
| Schema 变更 | 只增字段不删字段，消费者忽略未知字段 |