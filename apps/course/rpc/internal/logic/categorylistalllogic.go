package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/model"
	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

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

// CategoryListAll 查询全部分类并以树形结构返回（前端新增课程时选择分类用）。
func (l *CategoryListAllLogic) CategoryListAll(in *pb.CategoryListAllRequest) (*pb.CategoryNodeList, error) {
	list, err := l.svcCtx.CategoryModel.ListAll(l.ctx)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询分类列表失败")
	}
	if in.Name != "" {
		filtered := make([]*model.Category, 0, len(list))
		for _, c := range list {
			if containsStr(c.Name, in.Name) {
				filtered = append(filtered, c)
			}
		}
		list = filtered
	}
	return &pb.CategoryNodeList{Items: buildCategoryTree(list)}, nil
}

// containsStr 简单包含判断（分类名称模糊匹配）。
func containsStr(s, sub string) bool {
	if sub == "" {
		return true
	}
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
