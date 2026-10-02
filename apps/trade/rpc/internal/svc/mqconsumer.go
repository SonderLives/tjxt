package svc

import (
	"context"
	"errors"

	"tjxt/apps/trade/rpc/internal/model"
	"tjxt/pkg/mq"
	"tjxt/pkg/mq/event"

	"github.com/zeromicro/go-zero/core/logx"
)

// initMQ 注册支付成功事件消费者（pay.exchange / pay.success）。
// MQ 不可用时跳过注册，不阻塞服务启动。
func initMQ(svcCtx *ServiceContext, dsn string) {
	client := mq.NewClient(dsn)
	consumer := &paySuccessConsumer{svcCtx: svcCtx}
	mq.Register(client, mq.Binding{
		Queue:      "trade.order.paid",
		Exchange:   mq.ExchangePay,
		RoutingKey: mq.RoutingKeyPaySuccess,
	}, consumer.handlePaySuccess)
	svcCtx.MQClient = client
}

// paySuccessConsumer 消费支付成功事件：回写订单状态为已支付，
// 并向 learning 转发 OrderPayEvent（order.exchange / order.pay）开通课程。
// 消费幂等：订单非待支付态直接 ack；数据库条件更新兜底并发。
type paySuccessConsumer struct {
	svcCtx *ServiceContext
}

func (c *paySuccessConsumer) handlePaySuccess(ctx context.Context, evt *event.PaySuccessEvent) error {
	if evt.BizOrderNo <= 0 {
		logx.Errorf("ignore pay.success event with invalid bizOrderNo: %+v", evt)
		return nil
	}

	order, err := c.svcCtx.OrderModel.FindOne(ctx, evt.BizOrderNo)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			logx.Errorf("order not found for pay.success event, bizOrderNo=%d, discard", evt.BizOrderNo)
			return nil
		}
		return err // 返回 error 让消息重回队列重试
	}
	// 仅待支付订单需要回写；其余状态（已支付/已关闭/已完成等）视为终态，幂等丢弃
	if order.Status != 1 {
		logx.Infof("skip pay.success event for non-pending order, orderId=%d status=%d", order.Id, order.Status)
		return nil
	}

	transit, err := c.svcCtx.OrderModel.MarkPaid(ctx, order.Id, evt.PayOrderNo, evt.PayChannel, evt.PayTime, evt.Amount)
	if err != nil {
		return err
	}
	if !transit {
		// 并发消费下已被其他消费者处理
		return nil
	}

	c.publishOrderPay(ctx, order, evt)
	return nil
}

// publishOrderPay 向 learning 发布订单支付成功事件（best-effort：失败仅告警）。
func (c *paySuccessConsumer) publishOrderPay(ctx context.Context, order *model.Order, evt *event.PaySuccessEvent) {
	if c.svcCtx.MQProducer == nil {
		logx.Errorf("mq producer unavailable, skip order.pay event, orderId=%d", order.Id)
		return
	}
	details, err := c.svcCtx.OrderDetailModel.ListByOrderId(ctx, order.Id)
	if err != nil {
		logx.Errorf("list order details for order.pay event failed, orderId=%d: %v", order.Id, err)
		return
	}
	courseIds := make([]int64, 0, len(details))
	for _, d := range details {
		courseIds = append(courseIds, d.CourseId)
	}
	if err := c.svcCtx.MQProducer.Publish(ctx, mq.ExchangeOrder, mq.RoutingKeyOrderPay, event.OrderPayEvent{
		OrderBasic: event.OrderBasic{
			OrderID:    order.Id,
			UserID:     order.UserId,
			CourseIDs:  courseIds,
			FinishTime: evt.PayTime,
		},
	}); err != nil {
		logx.Errorf("publish order.pay event failed, orderId=%d: %v", order.Id, err)
	}
}
