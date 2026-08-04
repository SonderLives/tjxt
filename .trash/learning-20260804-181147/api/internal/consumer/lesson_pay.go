package consumer

import (
	"context"

	"tjxt/apps/learning/api/internal/mq"
	"tjxt/apps/learning/api/internal/mq/event"
	"tjxt/apps/learning/api/internal/service"

	"github.com/zeromicro/go-zero/core/logx"
)

// LessonPayConsumer 课程支付消费者
type LessonPayConsumer struct {
	lessonService service.LessonService
}

// NewLessonPayConsumer 创建课程支付消费者
func NewLessonPayConsumer(lessonService service.LessonService) *LessonPayConsumer {
	return &LessonPayConsumer{
		lessonService: lessonService,
	}
}

// Register 注册消费者到 MQ
func (c *LessonPayConsumer) Register(mqClient *mq.Client, payQueue, payExchange, payRoutingKey string) {
	mq.Register(
		mqClient,
		mq.Binding{
			Queue:      payQueue,
			Exchange:   payExchange,
			RoutingKey: payRoutingKey,
			Kind:       "direct",
		},
		c.handlePay,
	)
}

// handlePay 处理支付事件
func (c *LessonPayConsumer) handlePay(ctx context.Context, order *event.OrderPayEvent) error {
	if order == nil || order.UserID == 0 || len(order.CourseIDs) == 0 {
		logx.Info("接受订单支付事件，但参数无效")
		return nil
	}
	logx.Infof("receive order pay event: order_id=%d user_id=%d courses=%v", order.OrderID, order.UserID, order.CourseIDs)

	return c.lessonService.AddUserLessons(ctx, order.UserID, order.CourseIDs)
}
