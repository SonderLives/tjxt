package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/model"
	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryAddLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCategoryAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryAddLogic {
	return &CategoryAddLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CategoryAdd 新增分类：根据父分类自动计算级别，支持多级分类。
func (l *CategoryAddLogic) CategoryAdd(in *pb.CategoryAddRequest) (*pb.Empty, error) {
	level := int64(1)
	if in.ParentId != 0 {
		parent, err := l.svcCtx.CategoryModel.FindOne(l.ctx, in.ParentId)
		if err != nil && !isNotFound(err) {
			return nil, xerr.Wrap(err, xerr.CodeInternal, "查询父分类失败")
		}
		if parent != nil {
			level = int64(parent.Level) + 1
		}
	}
	data := &model.Category{
		Name:     in.Name,
		ParentId: in.ParentId,
		Level:    level,
		Priority: int64(in.Index),
		Status:   int64(CategoryStatusNormal),
		Creater:  0,
		Updater:  0,
	}
	if _, err := l.svcCtx.CategoryModel.Insert(l.ctx, data); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "新增分类失败")
	}
	return &pb.Empty{}, nil
}
