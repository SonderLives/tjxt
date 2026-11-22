// Package gateway 支付渠道抽象层。
//
// logic 层只依赖 PaymentGateway 接口；接入真实渠道（微信/支付宝）时实现该接口
// 并在 etc/pay.yaml 的 PayProvider 切换实现，业务代码零改动。
// demo 默认使用 MockGateway。
package gateway

import "context"

// CreatePaymentRequest 渠道下单请求。
type CreatePaymentRequest struct {
	PayOrderNo  int64  // 支付单号（渠道侧透传标识）
	BizOrderNo  int64  // 业务订单号
	BizUserID   int64  // 支付用户
	Amount      int64  // 金额，单位分
	Subject     string // 商品标题
	ChannelCode string // 渠道编码
}

// CreatePaymentResult 渠道下单结果。
type CreatePaymentResult struct {
	QrCodeUrl string // 二维码/跳转地址（native 扫码支付）
}

// RefundRequest 渠道退款请求。
type RefundRequest struct {
	PayOrderNo       int64
	RefundOrderNo    int64
	BizRefundOrderNo int64
	RefundAmount     int64
	ChannelCode      string
}

// RefundResult 渠道退款结果。
type RefundResult struct {
	// Synchronous 渠道是否同步返回终态（mock 渠道 true；微信/支付宝为异步回调驱动）
	Synchronous bool
	Success     bool
	ChannelNo   string // 渠道流水号
	Msg         string
}

// PaymentGateway 支付渠道接口。
type PaymentGateway interface {
	// CreatePayment 渠道下单，返回支付二维码。
	CreatePayment(ctx context.Context, req CreatePaymentRequest) (CreatePaymentResult, error)
	// Refund 发起退款。
	Refund(ctx context.Context, req RefundRequest) (RefundResult, error)
}
