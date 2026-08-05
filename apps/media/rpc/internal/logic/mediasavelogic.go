package logic

import (
	"context"

	"tjxt/apps/media/rpc/internal/svc"
	"tjxt/apps/media/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MediaSaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMediaSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MediaSaveLogic {
	return &MediaSaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MediaSaveLogic) MediaSave(in *pb.MediaSaveRequest) (*pb.MediaIdReply, error) {
	// todo: add your logic here and delete this line

	return &pb.MediaIdReply{}, nil
}
