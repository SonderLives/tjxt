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

type PayChannelGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayChannelGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayChannelGetLogic {
	return &PayChannelGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayChannelGetLogic) PayChannelGet(req *types.PayChannelIdReq) (resp *types.PayChannelDTO, err error) {
	reply, err := l.svcCtx.TradeRpc.PayChannelGet(l.ctx, &pb.IdRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}

	return &types.PayChannelDTO{
		Id:              reply.Id,
		Name:            reply.Name,
		ChannelCode:     reply.ChannelCode,
		ChannelIcon:     reply.ChannelIcon,
		ChannelPriority: int64(reply.ChannelPriority),
		Status:          int64(reply.Status),
	}, nil
}
