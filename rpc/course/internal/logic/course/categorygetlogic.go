package courselogic

import (
	"context"
	"errors"

	"tjxt/pkg/errorx"
	"tjxt/rpc/course/internal/model"
	"tjxt/rpc/course/internal/svc"
	"tjxt/rpc/course/pb"

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

// CategoryGet 查询单个分类详情，附带课程数、三级分类数及一/二级分类名称。
func (l *CategoryGetLogic) CategoryGet(in *pb.IdRequest) (*pb.CategoryInfo, error) {
	c, err := l.svcCtx.CategoryModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("分类不存在")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询分类失败")
	}
	courseNum, _ := l.svcCtx.CourseModel.CountByThirdCateId(l.ctx, c.Id)
	thirdNum, _ := l.svcCtx.CategoryModel.CountByParentId(l.ctx, c.Id)

	info := toCategoryInfo(c, courseNum, thirdNum)
	// 回溯一/二级分类名称
	all, aerr := l.svcCtx.CategoryModel.ListAll(l.ctx)
	if aerr == nil {
		m := categoryNameMap(all)
		if c.Level >= 2 && m[c.ParentId] != nil {
			info.SecondCategoryName = m[c.ParentId].Name
			if m[c.ParentId].ParentId != 0 && m[m[c.ParentId].ParentId] != nil {
				info.FirstCategoryName = m[m[c.ParentId].ParentId].Name
			}
		}
	}
	return info, nil
}
