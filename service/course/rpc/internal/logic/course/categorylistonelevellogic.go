package courselogic

import (
	"context"

	"tjxt/service/course/rpc/internal/svc"
	"tjxt/service/course/rpc/pb"
	"tjxt/common/errorx"

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
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询一级分类失败")
	}
	items := make([]*pb.CategoryInfo, 0, len(list))
	for _, c := range list {
		items = append(items, toCategoryInfo(c, 0, 0))
	}
	return &pb.CategoryList{Items: items}, nil
}
