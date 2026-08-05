package logic

import (
	"context"

	"tjxt/apps/media/rpc/internal/svc"
	"tjxt/apps/media/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MediaGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMediaGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MediaGetLogic {
	return &MediaGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 媒资管理
func (l *MediaGetLogic) MediaGet(in *pb.MediaIdRequest) (*pb.MediaVO, error) {
	// todo: add your logic here and delete this line

	return &pb.MediaVO{}, nil
}
