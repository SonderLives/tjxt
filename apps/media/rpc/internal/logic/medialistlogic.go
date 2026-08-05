package logic

import (
	"context"

	"tjxt/apps/media/rpc/internal/svc"
	"tjxt/apps/media/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MediaListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMediaListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MediaListLogic {
	return &MediaListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MediaListLogic) MediaList(in *pb.MediaListRequest) (*pb.MediaListReply, error) {
	// todo: add your logic here and delete this line

	return &pb.MediaListReply{}, nil
}
