package paylogic

import (
	"context"

	"tjxt/rpc/pay/internal/svc"
	"tjxt/rpc/pay/pb"
	"tjxt/pkg/errorx"

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
		return nil, errorx.BadRequestf("refund_order_no 非法")
	}
	m, err := l.svcCtx.RefundOrderModel.FindOneByRefundOrderNo(l.ctx, in.RefundOrderNo)
	if err != nil {
		if isNotFound(err) {
			return nil, errorx.NotFound("退款单不存在")
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询退款单失败")
	}
	switch m.Status {
	case RefundStatusSuccess:
		return nil, errorx.Conflict("退款单已成功，不能改为失败")
	case RefundStatusFailed:
		return &pb.EmptyResponse{}, nil
	}
	if err := l.svcCtx.RefundOrderModel.MarkToFailed(l.ctx, m.Id, in.ResultCode, in.ResultMsg); err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "更新退款状态失败")
	}
	return &pb.EmptyResponse{}, nil
}
