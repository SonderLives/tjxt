package logic

import (
	"context"

	"tjxt/apps/pay/rpc/internal/svc"
	"tjxt/apps/pay/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

// ClosePayOrderLogic 业务端主动取消订单。
type ClosePayOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClosePayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClosePayOrderLogic {
	return &ClosePayOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ClosePayOrderLogic) ClosePayOrder(in *pb.ClosePayOrderRequest) (*pb.EmptyResponse, error) {
	if in.PayOrderNo <= 0 {
		return nil, xerr.BadRequestf("pay_order_no 非法")
	}
	m, err := l.svcCtx.PayOrderModel.FindOneByPayOrderNo(l.ctx, in.PayOrderNo)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("支付单不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询支付单失败")
	}
	switch m.Status {
	case PayOrderStatusSuccess:
		return nil, xerr.Conflict("订单已支付成功，无法关单")
	case PayOrderStatusClosed:
		return &pb.EmptyResponse{}, nil
	}
	if err := l.svcCtx.PayOrderModel.MarkToClosed(l.ctx, m.Id, "MANUAL_CLOSE", "业务端主动关单"); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "关单失败")
	}
	return &pb.EmptyResponse{}, nil
}