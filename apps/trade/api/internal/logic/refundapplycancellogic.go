// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyCancelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundApplyCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyCancelLogic {
	return &RefundApplyCancelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefundApplyCancelLogic) RefundApplyCancel(req *types.RefundCancelReq) (resp *types.NamePlaceVO, err error) {
	if _, err = l.svcCtx.TradeRpc.RefundApplyCancel(l.ctx, &pb.RefundCancelRequest{
		Id:            req.Id,
		OrderDetailId: req.OrderDetailId,
	}); err != nil {
		return nil, err
	}
	return &types.NamePlaceVO{Existed: true, Message: "ok"}, nil
}
