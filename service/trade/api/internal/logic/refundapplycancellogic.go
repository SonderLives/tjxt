// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	backendlogic "tjxt/api/trade/backend/logic"
	"tjxt/api/trade/backend/pb"
	"tjxt/api/trade/internal/svc"
	"tjxt/api/trade/internal/types"

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
	if _, err = backendlogic.NewRefundApplyCancelLogic(l.ctx, l.svcCtx.Backend).RefundApplyCancel(&pb.RefundCancelRequest{
		Id:            req.Id,
		OrderDetailId: req.OrderDetailId,
	}); err != nil {
		return nil, err
	}
	return &types.NamePlaceVO{Existed: true, Message: "ok"}, nil
}
