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
	reply, err := l.svcCtx.TradeRpc.PayChannels(l.ctx, &pb.Empty{})
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
