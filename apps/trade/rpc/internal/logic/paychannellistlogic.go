package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayChannelListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayChannelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayChannelListLogic {
	return &PayChannelListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PayChannelListLogic) PayChannelList(in *pb.Empty) (*pb.PayChannelListReply, error) {
	// todo: add your logic here and delete this line

	return &pb.PayChannelListReply{}, nil
}
