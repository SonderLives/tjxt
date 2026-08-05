package logic

import (
	"context"
	"errors"

	payclient "tjxt/apps/pay/rpc/pay"
	"tjxt/apps/trade/rpc/internal/model"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

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
		return nil, xerr.BadRequestf("订单ID不能为空")
	}
	if in.PayChannelCode == "" {
		return nil, xerr.BadRequestf("支付渠道不能为空")
	}

	// 支付人与金额取自本地订单
	order, err := l.svcCtx.OrderModel.FindOne(l.ctx, in.OrderId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NotFound("订单不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询订单失败")
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
		return nil, xerr.Wrap(err, xerr.CodeInternal, "发起支付失败")
	}
	return &pb.PayApplyReply{QrUrl: resp.QrCodeUrl}, nil
}
