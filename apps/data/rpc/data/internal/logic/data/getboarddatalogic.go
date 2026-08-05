package datalogic

import (
	"context"

	"tjxt/apps/data/rpc/data/internal/svc"
	"tjxt/apps/data/rpc/data/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBoardDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBoardDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBoardDataLogic {
	return &GetBoardDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBoardDataLogic) GetBoardData(in *pb.BoardDataReq) (*pb.EchartsVO, error) {
	// todo: add your logic here and delete this line

	return &pb.EchartsVO{}, nil
}
