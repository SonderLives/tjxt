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

type PayChannelListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayChannelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayChannelListLogic {
	return &PayChannelListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayChannelListLogic) PayChannelList() (resp []types.PayChannelDTO, err error) {
	reply, err := l.svcCtx.TradeRpc.PayChannelList(l.ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}

	resp = make([]types.PayChannelDTO, 0)
	for _, item := range reply.Items {
		resp = append(resp, types.PayChannelDTO{
			Id:              item.Id,
			Name:            item.Name,
			ChannelCode:     item.ChannelCode,
			ChannelIcon:     item.ChannelIcon,
			ChannelPriority: int64(item.ChannelPriority),
			Status:          int64(item.Status),
		})
	}
	return resp, nil
}
