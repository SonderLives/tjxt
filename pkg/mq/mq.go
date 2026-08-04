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
