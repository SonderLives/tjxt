package event

import "time"

// OrderBasic 订单基础事件
// 与 learning/internal/mq/event 中的结构保持一致（json 字段名必须匹配）。
type OrderBasic struct {
	OrderID    int64     `json:"orderId"`
	UserID     int64     `json:"userId"`
	CourseIDs  []int64   `json:"courseIds"`
	FinishTime time.Time `json:"finishTime"`
}

// OrderPayEvent 订单支付成功事件
type OrderPayEvent struct {
	OrderBasic
}

// OrderRefundEvent 订单退款成功事件
type OrderRefundEvent struct {
	OrderBasic
}
