// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	backendlogic "tjxt/app/trade/api/backend/logic"
	"tjxt/app/trade/api/backend/pb"
	"tjxt/app/trade/api/internal/svc"
	"tjxt/app/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderCancelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderCancelLogic {
	return &OrderCancelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderCancelLogic) OrderCancel(req *types.OrderIdReq) (resp *types.NamePlaceVO, err error) {
	if _, err = backendlogic.NewOrderCancelLogic(l.ctx, l.svcCtx.Backend).OrderCancel(&pb.IdRequest{Id: req.Id}); err != nil {
		return nil, err
	}
	return &types.NamePlaceVO{Existed: true, Message: "ok"}, nil
}
