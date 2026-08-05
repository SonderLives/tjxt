package logic

import (
	"context"

	"tjxt/apps/media/rpc/internal/svc"
	"tjxt/apps/media/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFileSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileSaveLogic {
	return &FileSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FileSaveLogic) FileSave(in *pb.FileSaveRequest) (*pb.FileIdReply, error) {
	// todo: add your logic here and delete this line

	return &pb.FileIdReply{}, nil
}
