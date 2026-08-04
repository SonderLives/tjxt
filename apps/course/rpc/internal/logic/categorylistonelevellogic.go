package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryListOneLevelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCategoryListOneLevelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryListOneLevelLogic {
	return &CategoryListOneLevelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CategoryListOneLevelLogic) CategoryListOneLevel(in *pb.Empty) (*pb.CategoryList, error) {
	// todo: add your logic here and delete this line

	return &pb.CategoryList{}, nil
}
