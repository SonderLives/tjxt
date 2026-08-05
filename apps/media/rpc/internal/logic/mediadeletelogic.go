package logic

import (
	"context"

	"tjxt/apps/media/rpc/internal/svc"
	"tjxt/apps/media/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MediaDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMediaDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MediaDeleteLogic {
	return &MediaDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MediaDeleteLogic) MediaDelete(in *pb.MediaIdRequest) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
