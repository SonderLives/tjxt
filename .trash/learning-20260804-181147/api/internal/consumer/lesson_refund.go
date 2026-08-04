package consumer

import (
	"context"

	"tjxt/apps/learning/api/internal/mq"
	"tjxt/apps/learning/api/internal/mq/event"
	"tjxt/apps/learning/api/internal/service"

	"github.com/zeromicro/go-zero/core/logx"
)

// LessonRefundConsumer 课程退款消费者
type LessonRefundConsumer struct {
	lessonService service.LessonService
}

// NewLessonRefundConsumer 创建课程退款消费者
func NewLessonRefundConsumer(lessonService service.LessonService) *LessonRefundConsumer {
	return &LessonRefundConsumer{
		lessonService: lessonService,
	}
}

// Register 注册消费者到 MQ
func (c *LessonRefundConsumer) Register(mqClient *mq.Client, refundQueue, refundExchange, refundRoutingKey string) {
	mq.Register(
		mqClient,
		mq.Binding{
			Queue:      refundQueue,
			Exchange:   refundExchange,
			RoutingKey: refundRoutingKey,
			Kind:       "direct",
		},
		c.handleRefund,
	)
}

// handleRefund 处理退款事件
func (c *LessonRefundConsumer) handleRefund(ctx context.Context, order *event.OrderRefundEvent) error {
	if order == nil || order.UserID == 0 || len(order.CourseIDs) == 0 {
		logx.Info("skip order refund event, invalid params")
		return nil
	}

	logx.Infof("receive order refund event: order_id=%d user_id=%d courses=%v", order.OrderID, order.UserID, order.CourseIDs)

	// 退款通常只处理第一个课程，或全部课程
	return c.lessonService.DeleteCourseFromLesson(ctx, order.UserID, order.CourseIDs[0])
}
