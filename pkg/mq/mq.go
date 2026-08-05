// Package mq 提供基于 RabbitMQ 的通用发布/订阅能力。
//
// 事件契约与 learning 服务对齐：支付成功事件发布到
// order.exchange（routing key=order.pay），退款事件发布到
// order.exchange（routing key=order.refund）。交换机构建后即可被
// learning 等下游服务消费。
package mq

// 事件契约常量
const (
	// 订单相关交换机
	ExchangeOrder = "order.exchange"
	// 支付成功路由键
	RoutingKeyOrderPay = "order.pay"
	// 退款成功路由键
	RoutingKeyOrderRefund = "order.refund"
)

// 课程相关交换机与路由键。
// 事件契约：course 服务在课程上架/下架时向 course.events 交换机
// 发布 payload 为 {"courseId": 123} 的 JSON 消息，search 服务消费
// 这两个事件用于同步 ES 课程索引（见 pkg/mq/event.CourseEvent）。
const (
	// 课程相关交换机
	ExchangeCourse = "course.events"
	// 课程上架路由键（消费方：search 服务，用于写入/更新 ES 索引）
	RoutingKeyCourseUp = "course.up"
	// 课程下架路由键（消费方：search 服务，用于删除 ES 索引）
	RoutingKeyCourseDown = "course.down"
)
