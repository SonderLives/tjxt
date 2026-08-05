package logic

import (
	"context"

	"tjxt/apps/media/rpc/internal/svc"
	"tjxt/apps/media/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFileGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileGetLogic {
	return &FileGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 文件管理
func (l *FileGetLogic) FileGet(in *pb.FileIdRequest) (*pb.FileVO, error) {
	// todo: add your logic here and delete this line

	return &pb.FileVO{}, nil
}
