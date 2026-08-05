package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyCancelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundApplyCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyCancelLogic {
	return &RefundApplyCancelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefundApplyCancelLogic) RefundApplyCancel(in *pb.RefundCancelRequest) (*pb.Empty, error) {
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("退款申请ID不能为空")
	}

	if err := l.svcCtx.RefundApplyModel.UpdateStatus(l.ctx, in.Id, RefundStatusCancel); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "取消退款申请失败")
	}

	// 取消后清空订单明细的退款状态，允许重新申请
	if in.OrderDetailId > 0 {
		if err := l.svcCtx.OrderDetailModel.UpdateRefundStatus(l.ctx, in.OrderDetailId, 0); err != nil {
			return nil, xerr.Wrap(err, xerr.CodeInternal, "更新订单明细退款状态失败")
		}
	}
	return &pb.Empty{}, nil
}
