package event

import "time"

// PaySuccessEvent 支付单支付成功事件（pay → trade）。
// trade 消费后按 BizOrderNo 回写订单状态为已支付，并向 learning
// 转发 OrderPayEvent（见 event.OrderPayEvent）。
type PaySuccessEvent struct {
	BizOrderNo int64     `json:"bizOrderNo"`  // 业务订单号（trade 订单 id）
	PayOrderNo int64     `json:"payOrderNo"`  // 支付单号
	Amount     int64     `json:"amount"`      // 实付金额，单位分
	PayChannel string    `json:"payChannel"`  // 支付渠道编码
	PayTime    time.Time `json:"payTime"`     // 支付成功时间
}
