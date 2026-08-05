package logic

import (
	"context"

	"tjxt/apps/search/rpc/internal/svc"
	"tjxt/apps/search/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetInterestsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetInterestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInterestsLogic {
	return &GetInterestsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetInterestsLogic) GetInterests(in *pb.IdReq) (*pb.InterestsVO, error) {
	// todo: add your logic here and delete this line

	return &pb.InterestsVO{}, nil
}
