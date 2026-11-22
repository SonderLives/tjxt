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

type PayChannelAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayChannelAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayChannelAddLogic {
	return &PayChannelAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayChannelAddLogic) PayChannelAdd(req *types.PayChannelDTO) (resp *types.NamePlaceVO, err error) {
	if _, err = backendlogic.NewPayChannelAddLogic(l.ctx, l.svcCtx.Backend).PayChannelAdd(&pb.PayChannelDTO{
		Name:            req.Name,
		ChannelCode:     req.ChannelCode,
		ChannelIcon:     req.ChannelIcon,
		ChannelPriority: int32(req.ChannelPriority),
		Status:          int32(req.Status),
	}); err != nil {
		return nil, err
	}
	return &types.NamePlaceVO{Existed: true, Message: "ok"}, nil
}
