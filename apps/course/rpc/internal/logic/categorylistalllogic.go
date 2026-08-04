package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryListAllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCategoryListAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryListAllLogic {
	return &CategoryListAllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CategoryListAllLogic) CategoryListAll(in *pb.CategoryListAllRequest) (*pb.CategoryNodeList, error) {
	// todo: add your logic here and delete this line

	return &pb.CategoryNodeList{}, nil
}
