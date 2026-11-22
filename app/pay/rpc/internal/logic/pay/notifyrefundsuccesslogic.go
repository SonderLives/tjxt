package paylogic

import (
	"context"
	"errors"

	"tjxt/pkg/errorx"
	"tjxt/app/pay/rpc/internal/model"
	"tjxt/app/pay/rpc/internal/svc"
	"tjxt/app/pay/rpc/pb"

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
		return nil, errorx.BadRequestf("refund_order_no 非法")
	}
	m, err := l.svcCtx.RefundOrderModel.FindOneByRefundOrderNo(l.ctx, in.RefundOrderNo)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("退款单不存在")
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询退款单失败")
	}
	switch m.Status {
	case RefundStatusSuccess:
		return &pb.EmptyResponse{}, nil
	case RefundStatusFailed:
		return nil, errorx.Conflict("退款单已标记失败，不允许改为成功")
	case RefundStatusInit, RefundStatusProcessing:
	}
	if err := l.svcCtx.RefundOrderModel.MarkToSuccess(l.ctx, m.Id, in.ResultCode, in.ResultMsg, in.RefundChannel); err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "更新退款状态失败")
	}
	return &pb.EmptyResponse{}, nil
}
