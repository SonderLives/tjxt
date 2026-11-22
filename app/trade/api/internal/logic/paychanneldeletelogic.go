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

type PayChannelDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayChannelDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayChannelDeleteLogic {
	return &PayChannelDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayChannelDeleteLogic) PayChannelDelete(req *types.PayChannelIdReq) (resp *types.NamePlaceVO, err error) {
	if _, err = backendlogic.NewPayChannelDeleteLogic(l.ctx, l.svcCtx.Backend).PayChannelDelete(&pb.IdRequest{Id: req.Id}); err != nil {
		return nil, err
	}
	return &types.NamePlaceVO{Existed: true, Message: "ok"}, nil
}
