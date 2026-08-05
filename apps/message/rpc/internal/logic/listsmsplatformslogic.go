package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSmsPlatformsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListSmsPlatformsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSmsPlatformsLogic {
	return &ListSmsPlatformsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 短信平台
func (l *ListSmsPlatformsLogic) ListSmsPlatforms(in *pb.Empty) (*pb.SmsPlatformListReply, error) {
	// todo: add your logic here and delete this line

	return &pb.SmsPlatformListReply{}, nil
}
