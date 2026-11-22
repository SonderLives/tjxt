package logic

import (
	"context"

	backendlogic "tjxt/app/trade/api/backend/logic"
	"tjxt/app/trade/api/backend/pb"
	"tjxt/app/trade/api/internal/svc"
	"tjxt/app/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CartAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCartAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartAddLogic {
	return &CartAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CartAddLogic) CartAdd(req *types.CartAddReq) (resp *types.NamePlaceVO, err error) {
	if _, err = backendlogic.NewCartAddLogic(l.ctx, l.svcCtx.Backend).CartAdd(&pb.CartAddRequest{CourseId: req.CourseId}); err != nil {
		return nil, err
	}
	return &types.NamePlaceVO{Existed: true, Message: "ok"}, nil
}
