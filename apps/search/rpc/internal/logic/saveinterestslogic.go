package logic

import (
	"context"

	"tjxt/apps/search/rpc/internal/svc"
	"tjxt/apps/search/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveInterestsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveInterestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveInterestsLogic {
	return &SaveInterestsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户兴趣
func (l *SaveInterestsLogic) SaveInterests(in *pb.SaveInterestsReq) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
