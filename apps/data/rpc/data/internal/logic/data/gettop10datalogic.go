package datalogic

import (
	"context"

	"tjxt/apps/data/rpc/data/internal/svc"
	"tjxt/apps/data/rpc/data/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTop10DataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTop10DataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTop10DataLogic {
	return &GetTop10DataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTop10DataLogic) GetTop10Data(in *pb.Empty) (*pb.Top10DataVO, error) {
	// todo: add your logic here and delete this line

	return &pb.Top10DataVO{}, nil
}
