package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCategoryGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryGetLogic {
	return &CategoryGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CategoryGetLogic) CategoryGet(in *pb.IdRequest) (*pb.CategoryInfo, error) {
	// todo: add your logic here and delete this line

	return &pb.CategoryInfo{}, nil
}
