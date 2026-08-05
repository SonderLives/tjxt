package datalogic

import (
	"context"

	"tjxt/apps/data/rpc/data/internal/svc"
	"tjxt/apps/data/rpc/data/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetTop10DataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetTop10DataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetTop10DataLogic {
	return &SetTop10DataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetTop10DataLogic) SetTop10Data(in *pb.Top10DataSetReq) (*pb.OkReply, error) {
	// todo: add your logic here and delete this line

	return &pb.OkReply{}, nil
}
