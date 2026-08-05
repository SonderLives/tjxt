package datalogic

import (
	"context"

	"tjxt/apps/data/rpc/data/internal/svc"
	"tjxt/apps/data/rpc/data/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetTodayDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetTodayDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetTodayDataLogic {
	return &SetTodayDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetTodayDataLogic) SetTodayData(in *pb.TodayDataSetReq) (*pb.OkReply, error) {
	// todo: add your logic here and delete this line

	return &pb.OkReply{}, nil
}
