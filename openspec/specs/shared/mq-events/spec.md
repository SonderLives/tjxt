# mq-events（RabbitMQ 事件契约）Specification

## Purpose

定义跨服务 RabbitMQ 事件总线的拓扑、事件清单、结构体规范与消费约定（`pkg/mq`、`pkg/mq/event/`）。事件用于订单支付、消息发送、积分变更等异步解耦场景。

## Requirements

### Requirement: 事件总线命名约定

事件拓扑 SHALL 遵循统一命名：Exchange 为 `<domain>.events`，RoutingKey 为 `<domain>.<action>`，Queue 为 `<svc>.<domain>.<action>`。

#### Scenario: 交易域事件发布

- **WHEN** trade-rpc 发布 `order.created` 事件
- **THEN** 事件发往 `trade.events` Exchange，RoutingKey 为 `order.created`
- **AND** learning-rpc 通过 `learning.order.created` 队列消费

### Requirement: 交易域事件（trade.events）

trade-rpc SHALL 发布以下事件：

| 事件 | RoutingKey | 消费者 | 说明 |
|------|------------|--------|------|
| 订单创建 | `order.created` | learning-rpc, message-rpc | 下单成功，解锁课程+发通知 |
| 订单支付 | `order.paid` | promotion-rpc, learning-rpc | 支付成功，核销优惠券、更新学习权限 |
| 订单取消 | `order.cancelled` | promotion-rpc | 释放占用的优惠券/库存 |
| 退款发起 | `refund.initiated` | message-rpc | 发起退款通知 |
| 退款完成 | `refund.completed` | promotion-rpc, learning-rpc | 退款成功，回收权益 |

#### Scenario: 订单支付成功触发核销与解锁

- **WHEN** trade-rpc 发布 `order.paid`
- **THEN** promotion-rpc 消费该事件核销优惠券，learning-rpc 消费该事件更新学习权限

### Requirement: 支付域事件（pay.events）

pay-rpc SHALL 发布以下事件：

| 事件 | RoutingKey | 消费者 | 说明 |
|------|------------|--------|------|
| 支付成功 | `payment.paid` | trade-rpc | 回调确认，驱动订单状态流转 |
| 支付失败 | `payment.failed` | trade-rpc, message-rpc | 失败通知 |
| 对账差异 | `reconciliation.mismatch` | 无（告警人工处理） | 对账差异告警 |

#### Scenario: 支付回调驱动订单流转

- **WHEN** pay-rpc 收到渠道支付成功回调并发布 `payment.paid`
- **THEN** trade-rpc 消费该事件，将订单状态推进为已支付

### Requirement: 学习域事件（learning.events）

learning-rpc SHALL 发布以下事件：

| 事件 | RoutingKey | 消费者 | 说明 |
|------|------------|--------|------|
| 课程完成 | `course.completed` | message-rpc, promotion-rpc | 结业证书、积分奖励 |
| 进度更新 | `progress.updated` | 无（内部统计用） | 学习进度统计 |
| 打卡完成 | `checkin.completed` | message-rpc | 连续打卡奖励 |

#### Scenario: 结业触发通知与奖励

- **WHEN** learning-rpc 发布 `course.completed`
- **THEN** message-rpc 消费该事件发送结业通知，promotion-rpc 消费该事件发放积分奖励

### Requirement: 优惠券域事件（promotion.events）

promotion-rpc SHALL 发布以下事件：

| 事件 | RoutingKey | 消费者 | 说明 |
|------|------------|--------|------|
| 券发放 | `coupon.issued` | message-rpc | 发放通知 |
| 券核销 | `coupon.used` | 无（统计用） | 核销统计 |
| 券过期 | `coupon.expired` | message-rpc | 过期提醒 |

#### Scenario: 领券触发通知

- **WHEN** promotion-rpc 向用户发放优惠券并发布 `coupon.issued`
- **THEN** message-rpc 消费该事件，向用户推送领券成功通知

### Requirement: 消息域事件（message.events）

message-rpc SHALL 发布以下事件：`inbox.send`（站内信内部投递）、`sms.send`（第三方短信通道）。

#### Scenario: 站内信内部投递

- **WHEN** message-rpc 需要向用户投递站内信
- **THEN** 通过 `message.events` Exchange 的 `inbox.send` RoutingKey 发布事件，由内部消费者完成投递

### Requirement: 事件结构体规范

所有事件结构体 SHALL 定义在 `pkg/mq/event/` 下，JSON 序列化，字段命名 camelCase，必须携带 `IdempotencyKey`（幂等键）与 `Timestamp`。

#### Scenario: 订单创建事件结构

- **WHEN** trade-rpc 构造 `OrderCreatedEvent`
- **THEN** 结构体包含 `OrderId`、`UserId`、`CourseIds`、`Amount`（单位：分）、`CouponId`（omitempty）、`IdempotencyKey`、`Timestamp` 字段
- **AND** Schema 变更只增字段不删字段，消费者忽略未知字段

### Requirement: 消费者实现约定

消费者 SHALL 遵循以下约定：

1. **幂等**：基于 `IdempotencyKey` 或业务主键去重，已处理的事件直接 ack
2. **手动 ACK**：处理成功才 ack；失败 requeue 或进死信队列
3. **重试策略**：最大 3 次，指数退避，最终进 DLQ 人工处理
4. **监控**：消费延迟、积压量、失败率、DLQ 堆积纳入监控指标

#### Scenario: 重复投递幂等处理

- **WHEN** 同一 `IdempotencyKey` 的 `OrderCreatedEvent` 被投递两次
- **THEN** 第一次正常处理并 ack；第二次被幂等校验拦截，直接 ack 不重复执行业务

#### Scenario: 处理失败重试与死信

- **WHEN** 消费处理器返回 error
- **THEN** 消息 requeue 重试，最多 3 次指数退避
- **AND** 超过重试上限后进入死信队列等待人工处理

