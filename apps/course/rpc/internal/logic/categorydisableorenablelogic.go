package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryDisableOrEnableLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCategoryDisableOrEnableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryDisableOrEnableLogic {
	return &CategoryDisableOrEnableLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CategoryDisableOrEnableLogic) CategoryDisableOrEnable(in *pb.CategoryStatusRequest) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
