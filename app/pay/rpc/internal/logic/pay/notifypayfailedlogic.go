package paylogic

import (
	"context"
	"errors"

	"tjxt/common/errorx"
	"tjxt/service/pay/rpc/internal/model"
	"tjxt/service/pay/rpc/internal/svc"
	"tjxt/service/pay/rpc/pb"

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
		return nil, errorx.BadRequestf("pay_order_no 非法")
	}
	m, err := l.svcCtx.PayOrderModel.FindOneByPayOrderNo(l.ctx, in.PayOrderNo)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("支付单不存在")
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询支付单失败")
	}
	switch m.Status {
	case PayOrderStatusSuccess:
		return nil, errorx.Conflict("订单已支付成功，不能再标记失败")
	case PayOrderStatusClosed:
		return &pb.EmptyResponse{}, nil
	}
	transit, err := l.svcCtx.PayOrderModel.MarkToClosed(l.ctx, m.Id, in.ResultCode, in.ResultMsg)
	if err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "更新支付单状态失败")
	}
	if !transit {
		logx.WithContext(l.ctx).Infof("notify pay failed no-op, already terminal, payOrderId=%d", m.Id)
	}
	return &pb.EmptyResponse{}, nil
}
