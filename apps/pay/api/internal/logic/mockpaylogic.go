// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/pay/api/internal/svc"
	"tjxt/apps/pay/api/internal/types"
	payclient "tjxt/apps/pay/rpc/pay"
	"tjxt/pkg/auth"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type MockPayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMockPayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MockPayLogic {
	return &MockPayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// MockPay 模拟支付网关回调（demo）：登录用户本人确认支付自己的订单。
// 真实环境由支付渠道回调 /pay-notify/pay-success（带 HMAC 签名），不走本接口。
func (l *MockPayLogic) MockPay(req *types.MockPayReq) (resp *types.Result, err error) {
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, xerr.New(xerr.CodeUnauthorized, "未登录")
	}
	if req.BizOrderNo <= 0 {
		return nil, xerr.BadRequestf("bizOrderNo 非法")
	}

	// 属主校验：支付单必须属于当前登录用户
	payOrder, err := l.svcCtx.PayRpc.QueryPayOrderByBizOrderNo(l.ctx, &payclient.QueryPayOrderRequest{
		BizOrderNo: req.BizOrderNo,
	})
	if err != nil {
		return nil, err
	}
	if payOrder.BizUserId != userId {
		return nil, xerr.Forbidden("无权操作该订单")
	}
	switch payOrder.Status {
	case 3: // 已支付：幂等成功
		return &types.Result{Code: 200, Msg: "OK"}, nil
	case 2: // 已关闭
		return nil, xerr.Conflict("支付单已关闭，无法支付")
	}

	if _, err := l.svcCtx.PayRpc.NotifyPaySuccess(l.ctx, &payclient.NotifyPaySuccessRequest{
		PayOrderNo: payOrder.PayOrderNo,
		ResultCode: "MOCK",
		ResultMsg:  "mock 支付成功",
	}); err != nil {
		return nil, err
	}
	return &types.Result{Code: 200, Msg: "OK"}, nil
}
