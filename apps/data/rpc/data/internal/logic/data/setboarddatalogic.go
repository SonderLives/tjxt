package datalogic

import (
	"context"

	"tjxt/apps/data/rpc/data/internal/svc"
	"tjxt/apps/data/rpc/data/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetBoardDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetBoardDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetBoardDataLogic {
	return &SetBoardDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetBoardDataLogic) SetBoardData(in *pb.BoardDataSetReq) (*pb.OkReply, error) {
	// todo: add your logic here and delete this line

	return &pb.OkReply{}, nil
}
