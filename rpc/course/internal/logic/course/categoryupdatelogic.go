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

type CategoryUpdateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCategoryUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryUpdateLogic {
	return &CategoryUpdateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CategoryUpdate 更新分类名称、排序，父级变化时重算级别。
func (l *CategoryUpdateLogic) CategoryUpdate(in *pb.CategoryUpdateRequest) (*pb.Empty, error) {
	c, err := l.svcCtx.CategoryModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("分类不存在")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询分类失败")
	}
	c.Name = in.Name
	c.Priority = int64(in.Index)
	if c.ParentId != 0 {
		parent, perr := l.svcCtx.CategoryModel.FindOne(l.ctx, c.ParentId)
		if perr == nil && parent != nil {
			c.Level = parent.Level + 1
		}
	} else {
		c.Level = 1
	}
	if err := l.svcCtx.CategoryModel.Update(l.ctx, c); err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "更新分类失败")
	}
	return &pb.Empty{}, nil
}
