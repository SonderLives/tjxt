package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

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

// CategoryListQuery 按名称/状态过滤查询分类列表。
func (l *CategoryListQueryLogic) CategoryListQuery(in *pb.CategoryListQueryRequest) (*pb.CategoryList, error) {
	list, err := l.svcCtx.CategoryModel.FindByCondition(l.ctx, in.Name, int64(in.Status))
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询分类列表失败")
	}
	items := make([]*pb.CategoryInfo, 0, len(list))
	for _, c := range list {
		items = append(items, toCategoryInfo(c, 0, 0))
	}
	return &pb.CategoryList{Items: items}, nil
}
