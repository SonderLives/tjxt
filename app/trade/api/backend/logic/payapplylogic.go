package tradelogic

import (
	"context"
	"errors"

	payclient "tjxt/app/pay/rpc/client/pay"
	"tjxt/app/trade/api/backend/model"
	"tjxt/app/trade/api/backend/svc"
	"tjxt/app/trade/api/backend/pb"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayApplyLogic {
	return &PayApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 支付 =====
func (l *PayApplyLogic) PayApply(in *pb.PayApplyRequest) (*pb.PayApplyReply, error) {
	if in.OrderId <= 0 {
		return nil, errorx.BadRequestf("订单ID不能为空")
	}
	if in.PayChannelCode == "" {
		return nil, errorx.BadRequestf("支付渠道不能为空")
	}

	// 支付人与金额取自本地订单
	order, err := l.svcCtx.OrderModel.FindOne(l.ctx, in.OrderId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("订单不存在")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询订单失败")
	}

	resp, err := l.svcCtx.PayRpc.ApplyPayOrder(l.ctx, &payclient.ApplyPayOrderRequest{
		BizUserId:      order.UserId,
		BizOrderNo:     in.OrderId,
		Amount:         order.TotalAmount,
		PayChannelCode: in.PayChannelCode,
		PayType:        4, // native 扫码支付
		NotifyUrl:      "",
	})
	if err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "发起支付失败")
	}
	return &pb.PayApplyReply{QrUrl: resp.QrCodeUrl}, nil
}
