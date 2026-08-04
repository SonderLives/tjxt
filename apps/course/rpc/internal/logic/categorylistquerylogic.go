package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryListQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCategoryListQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryListQueryLogic {
	return &CategoryListQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CategoryListQueryLogic) CategoryListQuery(in *pb.CategoryListQueryRequest) (*pb.CategoryList, error) {
	// todo: add your logic here and delete this line

	return &pb.CategoryList{}, nil
}
