package logic

import (
	"context"

	"tjxt/apps/pay/rpc/internal/svc"
	"tjxt/apps/pay/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

// NotifyRefundFailedLogic 第三方退款失败回调
type NotifyRefundFailedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyRefundFailedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyRefundFailedLogic {
	return &NotifyRefundFailedLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *NotifyRefundFailedLogic) NotifyRefundFailed(in *pb.NotifyRefundFailedRequest) (*pb.EmptyResponse, error) {
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
		return nil, xerr.Conflict("退款单已成功，不能改为失败")
	case RefundStatusFailed:
		return &pb.EmptyResponse{}, nil
	}
	if err := l.svcCtx.RefundOrderModel.MarkToFailed(l.ctx, m.Id, in.ResultCode, in.ResultMsg); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "更新退款状态失败")
	}
	return &pb.EmptyResponse{}, nil
}