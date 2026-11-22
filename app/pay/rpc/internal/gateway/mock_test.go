package gateway

import (
	"context"
	"strings"
	"testing"
)

func TestMockGatewayCreatePayment(t *testing.T) {
	got, err := MockGateway{}.CreatePayment(context.Background(), CreatePaymentRequest{
		PayOrderNo: 990010, BizOrderNo: 990001, Amount: 10000,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got.QrCodeUrl, "tjxt://mock-pay?order_no=990010&amount=10000") {
		t.Fatalf("unexpected qr url: %s", got.QrCodeUrl)
	}
}

func TestMockGatewayRefundSynchronousSuccess(t *testing.T) {
	got, err := MockGateway{}.Refund(context.Background(), RefundRequest{RefundAmount: 100})
	if err != nil {
		t.Fatal(err)
	}
	// mock 渠道契约：退款同步返回成功终态，logic 层据此直接落成功
	if !got.Synchronous || !got.Success || got.ChannelNo != "MOCK_OK" {
		t.Fatalf("unexpected refund result: %+v", got)
	}
}
