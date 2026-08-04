package event

import "time"

// OrderBasic 订单基础事件
// 对应 Java: OrderBasicDTO
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

// OrderRefundEvent 订单退款事件
type OrderRefundEvent struct {
	OrderBasic
}
