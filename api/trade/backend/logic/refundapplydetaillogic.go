package tradelogic

import (
	"context"
	"errors"

	"tjxt/api/trade/backend/model"
	"tjxt/api/trade/backend/svc"
	"tjxt/rpc/trade/pb"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyDetailLogic {
	return &RefundApplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefundApplyDetailLogic) RefundApplyDetail(in *pb.IdRequest) (*pb.RefundApplyVO, error) {
	if in.Id <= 0 {
		return nil, errorx.BadRequestf("退款申请ID不能为空")
	}

	ra, err := l.svcCtx.RefundApplyModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("退款申请不存在")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询退款申请失败")
	}

	order, _ := l.svcCtx.OrderModel.FindOne(l.ctx, ra.OrderId)
	detail, _ := l.svcCtx.OrderDetailModel.FindOne(l.ctx, ra.OrderDetailId)

	return toRefundApplyVO(ra, order, detail), nil
}
