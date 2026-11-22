// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	backendlogic "tjxt/service/trade/api/backend/logic"
	"tjxt/service/trade/api/backend/pb"
	"tjxt/service/trade/api/internal/svc"
	"tjxt/service/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDeleteLogic {
	return &OrderDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderDeleteLogic) OrderDelete(req *types.OrderIdReq) (resp *types.NamePlaceVO, err error) {
	if _, err = backendlogic.NewOrderDeleteLogic(l.ctx, l.svcCtx.Backend).OrderDelete(&pb.IdRequest{Id: req.Id}); err != nil {
		return nil, err
	}
	return &types.NamePlaceVO{Existed: true, Message: "ok"}, nil
}
