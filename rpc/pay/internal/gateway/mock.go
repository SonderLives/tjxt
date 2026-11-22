package gateway

import (
	"context"
	"fmt"
)

// MockGateway demo 渠道：生成本地可识别的 mock 二维码，退款同步成功。
// 真实项目替换为 wechat/alipay 实现（调渠道下单/退款 API）。
type MockGateway struct{}

// CreatePayment 生成 mock 支付二维码（tjxt:// 协议头，本地可识别）。
func (MockGateway) CreatePayment(_ context.Context, req CreatePaymentRequest) (CreatePaymentResult, error) {
	return CreatePaymentResult{
		QrCodeUrl: fmt.Sprintf("tjxt://mock-pay?order_no=%d&amount=%d", req.PayOrderNo, req.Amount),
	}, nil
}

// Refund mock 渠道退款同步成功。
func (MockGateway) Refund(_ context.Context, _ RefundRequest) (RefundResult, error) {
	return RefundResult{
		Synchronous: true,
		Success:     true,
		ChannelNo:   "MOCK_OK",
		Msg:         "mock 退款成功",
	}, nil
}
