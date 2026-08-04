package logic

import (
	"context"

	"tjxt/apps/pay/rpc/internal/svc"
	"tjxt/apps/pay/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

// NotifyRefundSuccessLogic 第三方退款成功回调
type NotifyRefundSuccessLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyRefundSuccessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyRefundSuccessLogic {
	return &NotifyRefundSuccessLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *NotifyRefundSuccessLogic) NotifyRefundSuccess(in *pb.NotifyRefundSuccessRequest) (*pb.EmptyResponse, error) {
	if in.RefundOrderNo <= 0 {
		return nil, xerr.BadRequestf("refund_order_no 非法")
	}
	m, err := l.svcCtx.RefundOrderModel.FindOneByRefundOrderNo(l.ctx, in.RefundOrderNo)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("退款单不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询退款单失败")
	}
	switch m.Status {
	case RefundStatusSuccess:
		return &pb.EmptyResponse{}, nil
	case RefundStatusFailed:
		return nil, xerr.Conflict("退款单已标记失败，不允许改为成功")
	case RefundStatusInit, RefundStatusProcessing:
	}
	if err := l.svcCtx.RefundOrderModel.MarkToSuccess(l.ctx, m.Id, in.ResultCode, in.ResultMsg, in.RefundChannel); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "更新退款状态失败")
	}
	return &pb.EmptyResponse{}, nil
}