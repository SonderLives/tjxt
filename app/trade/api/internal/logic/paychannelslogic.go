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

type PayChannelsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayChannelsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayChannelsLogic {
	return &PayChannelsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayChannelsLogic) PayChannels() (resp []types.PayChannelVO, err error) {
	reply, err := backendlogic.NewPayChannelsLogic(l.ctx, l.svcCtx.Backend).PayChannels(&pb.Empty{})
	if err != nil {
		return nil, err
	}

	resp = make([]types.PayChannelVO, 0)
	for _, item := range reply.Items {
		resp = append(resp, types.PayChannelVO{
			Id:              item.Id,
			Name:            item.Name,
			ChannelCode:     item.ChannelCode,
			ChannelIcon:     item.ChannelIcon,
			ChannelPriority: int64(item.ChannelPriority),
		})
	}
	return resp, nil
}
