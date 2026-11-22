package tradelogic

import (
	"context"

	"tjxt/app/trade/api/backend/svc"
	"tjxt/app/trade/api/backend/pb"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundApplyDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyDeleteLogic {
	return &RefundApplyDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefundApplyDeleteLogic) RefundApplyDelete(in *pb.IdRequest) (*pb.Empty, error) {
	if in.Id <= 0 {
		return nil, errorx.BadRequestf("退款申请ID不能为空")
	}

	// 未找到时忽略，直接执行删除保证幂等
	if ra, err := l.svcCtx.RefundApplyModel.FindOne(l.ctx, in.Id); err == nil && ra != nil {
		if e := l.svcCtx.OrderDetailModel.UpdateRefundStatus(l.ctx, ra.OrderDetailId, 0); e != nil {
			return nil, errorx.Wrap(e, errorx.CodeInternal, "更新订单明细退款状态失败")
		}
	}

	if err := l.svcCtx.RefundApplyModel.Delete(l.ctx, in.Id); err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "删除退款申请失败")
	}
	return &pb.Empty{}, nil
}
