package logic

import (
	"context"

	"tjxt/apps/pay/rpc/internal/svc"
	"tjxt/apps/pay/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

// NotifyPayFailedLogic 支付渠道通知支付失败：把单置为已关闭。
type NotifyPayFailedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyPayFailedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyPayFailedLogic {
	return &NotifyPayFailedLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *NotifyPayFailedLogic) NotifyPayFailed(in *pb.NotifyPayFailedRequest) (*pb.EmptyResponse, error) {
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
		return nil, xerr.Conflict("订单已支付成功，不能再标记失败")
	case PayOrderStatusClosed:
		return &pb.EmptyResponse{}, nil
	}
	if err := l.svcCtx.PayOrderModel.MarkToClosed(l.ctx, m.Id, in.ResultCode, in.ResultMsg); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "更新支付单状态失败")
	}
	return &pb.EmptyResponse{}, nil
}