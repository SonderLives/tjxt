package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayChannelGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayChannelGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayChannelGetLogic {
	return &PayChannelGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PayChannelGetLogic) PayChannelGet(in *pb.IdRequest) (*pb.PayChannelDTO, error) {
	// todo: add your logic here and delete this line

	return &pb.PayChannelDTO{}, nil
}
