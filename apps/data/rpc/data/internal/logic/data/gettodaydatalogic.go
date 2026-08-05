package datalogic

import (
	"context"

	"tjxt/apps/data/rpc/data/internal/svc"
	"tjxt/apps/data/rpc/data/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTodayDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetTodayDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTodayDataLogic {
	return &GetTodayDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetTodayDataLogic) GetTodayData(in *pb.Empty) (*pb.TodayDataVO, error) {
	// todo: add your logic here and delete this line

	return &pb.TodayDataVO{}, nil
}
