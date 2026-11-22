package tradelogic

import (
	"context"
	"errors"

	"tjxt/app/trade/api/backend/model"
	"tjxt/app/trade/api/backend/svc"
	"tjxt/app/trade/api/backend/pb"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyNextLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundApplyNextLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyNextLogic {
	return &RefundApplyNextLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefundApplyNextLogic) RefundApplyNext(in *pb.Empty) (*pb.RefundApplyVO, error) {
	ra, err := l.svcCtx.RefundApplyModel.FindNextPending(l.ctx)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("无待审批退款")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询待审批退款失败")
	}

	order, _ := l.svcCtx.OrderModel.FindOne(l.ctx, ra.OrderId)
	detail, _ := l.svcCtx.OrderDetailModel.FindOne(l.ctx, ra.OrderDetailId)

	return toRefundApplyVO(ra, order, detail), nil
}
