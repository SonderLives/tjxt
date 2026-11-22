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
