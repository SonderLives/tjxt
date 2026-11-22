package svc

import (
	"context"

	"tjxt/pkg/mq"
	"tjxt/pkg/mq/event"

	"github.com/zeromicro/go-zero/core/logx"
)

// initMQ 注册订单支付/退款事件消费者（order.exchange，队列/路由键来自配置）。
// MQ 不可用时跳过注册，不阻塞服务启动。
func initMQ(svcCtx *ServiceContext, dsn string) {
	c := svcCtx.Config.RabbitMQ
	if c.PayQueue == "" || c.RefundQueue == "" {
		logx.Error("skip learning mq consumer: rabbitmq queues not configured")
		return
	}

	client := mq.NewClient(dsn)
	consumer := &lessonConsumer{svcCtx: svcCtx}
	mq.Register(client, mq.Binding{
		Queue:      c.PayQueue,
		Exchange:   c.PayExchange,
		RoutingKey: c.PayRoutingKey,
	}, consumer.handleOrderPay)
	mq.Register(client, mq.Binding{
		Queue:      c.RefundQueue,
		Exchange:   c.RefundExchange,
		RoutingKey: c.RefundRoutingKey,
	}, consumer.handleOrderRefund)
	svcCtx.MQClient = client
}

// lessonConsumer 消费订单支付/退款事件，为用户开通/撤销课程。
// 幂等性：GrantCourses 为 ON DUPLICATE KEY 幂等插入，RevokeCourses 置失效可重复执行。
type lessonConsumer struct {
	svcCtx *ServiceContext
}

// handleOrderPay 订单支付成功：为用户开通订单内全部课程。
func (c *lessonConsumer) handleOrderPay(ctx context.Context, evt *event.OrderPayEvent) error {
	if evt.UserID <= 0 || len(evt.CourseIDs) == 0 {
		logx.Errorf("ignore order.pay event with empty payload: %+v", evt)
		return nil
	}
	if err := c.svcCtx.LearningLessonModel.GrantCourses(ctx, evt.UserID, evt.CourseIDs); err != nil {
		return err // 返回 error 让消息重回队列重试
	}
	logx.Infof("granted courses for paid order, orderId=%d userId=%d courseIds=%v",
		evt.OrderID, evt.UserID, evt.CourseIDs)
	return nil
}

// handleOrderRefund 订单退款成功：撤销用户订单内全部课程（status 置失效）。
func (c *lessonConsumer) handleOrderRefund(ctx context.Context, evt *event.OrderRefundEvent) error {
	if evt.UserID <= 0 || len(evt.CourseIDs) == 0 {
		logx.Errorf("ignore order.refund event with empty payload: %+v", evt)
		return nil
	}
	if err := c.svcCtx.LearningLessonModel.RevokeCourses(ctx, evt.UserID, evt.CourseIDs); err != nil {
		return err
	}
	logx.Infof("revoked courses for refunded order, orderId=%d userId=%d courseIds=%v",
		evt.OrderID, evt.UserID, evt.CourseIDs)
	return nil
}
