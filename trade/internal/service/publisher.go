package service

import (
	"context"
	"time"

	"common/mq"
	"common/mq/event"
	"common/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

// mqEventPublisher 基于 RabbitMQ 的订单事件发布实现。
type mqEventPublisher struct {
	producer         *mq.Producer
	exchange         string
	payRoutingKey    string
	refundRoutingKey string
}

// NewMQEventPublisher 创建 MQ 事件发布器。
func NewMQEventPublisher(producer *mq.Producer, exchange, payRoutingKey, refundRoutingKey string) EventPublisher {
	return &mqEventPublisher{
		producer:         producer,
		exchange:         exchange,
		payRoutingKey:    payRoutingKey,
		refundRoutingKey: refundRoutingKey,
	}
}

// PublishPay 发布订单支付成功事件。
func (p *mqEventPublisher) PublishPay(ctx context.Context, orderID, userID int64, courseIDs []int64, finishTime time.Time) error {
	return p.publish(ctx, p.payRoutingKey, &event.OrderPayEvent{
		OrderBasic: event.OrderBasic{
			OrderID:    orderID,
			UserID:     userID,
			CourseIDs:  courseIDs,
			FinishTime: finishTime,
		},
	})
}

// PublishRefund 发布订单退款成功事件。
func (p *mqEventPublisher) PublishRefund(ctx context.Context, orderID, userID int64, courseIDs []int64, finishTime time.Time) error {
	return p.publish(ctx, p.refundRoutingKey, &event.OrderRefundEvent{
		OrderBasic: event.OrderBasic{
			OrderID:    orderID,
			UserID:     userID,
			CourseIDs:  courseIDs,
			FinishTime: finishTime,
		},
	})
}

func (p *mqEventPublisher) publish(ctx context.Context, routingKey string, body any) error {
	if p.producer == nil {
		logx.Errorf("mq producer is nil, skip publish, routingKey=%s", routingKey)
		return xerr.ServiceUnavailable("消息队列暂不可用")
	}
	if err := p.producer.Publish(ctx, p.exchange, routingKey, body); err != nil {
		logx.Errorf("publish order event failed, routingKey=%s err=%v", routingKey, err)
		return xerr.Wrap(err, xerr.CodeServiceUnavailable, "发布订单事件失败")
	}
	return nil
}
