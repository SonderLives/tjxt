package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

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

// CategoryListOneLevel 查询全部一级分类。
func (l *CategoryListOneLevelLogic) CategoryListOneLevel(in *pb.Empty) (*pb.CategoryList, error) {
	list, err := l.svcCtx.CategoryModel.FindByLevel(l.ctx, 1)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询一级分类失败")
	}
	items := make([]*pb.CategoryInfo, 0, len(list))
	for _, c := range list {
		items = append(items, toCategoryInfo(c, 0, 0))
	}
	return &pb.CategoryList{Items: items}, nil
}
