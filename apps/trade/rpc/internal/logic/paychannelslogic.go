package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayChannelsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayChannelsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayChannelsLogic {
	return &PayChannelsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PayChannelsLogic) PayChannels(in *pb.Empty) (*pb.PayChannelVOList, error) {
	// todo: add your logic here and delete this line

	return &pb.PayChannelVOList{}, nil
}
