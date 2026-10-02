// Package mq 提供基于 RabbitMQ 的通用发布/订阅能力。
//
// 支付链路事件契约：
//  1. pay 服务在支付单置为成功后，向 pay.exchange（routing key=pay.success）
//     发布 PaySuccessEvent；trade 服务消费它回写订单状态。
//  2. trade 服务在订单置为已支付后，向 order.exchange（routing key=order.pay）
//     发布 OrderPayEvent；learning 服务消费它为用户开通课程。
//  3. trade 服务在退款成功后，向 order.exchange（routing key=order.refund）
//     发布 OrderRefundEvent；learning 服务消费它撤销课程。
//
// 课程上下架事件契约：course 服务发布到 course.events，search 服务消费
// 用于同步 ES 课程索引（见 pkg/mq/event.CourseEvent）。
package mq

// 事件契约常量
const (
	// 支付相关交换机（发布方：pay 服务，消费方：trade 服务）
	ExchangePay = "pay.exchange"
	// 支付成功路由键
	RoutingKeyPaySuccess = "pay.success"
)

// 事件契约常量
const (
	// 订单相关交换机（发布方：trade 服务，消费方：learning 服务）
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
