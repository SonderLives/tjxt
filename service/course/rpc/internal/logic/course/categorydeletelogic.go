package courselogic

import (
	"context"

	"tjxt/rpc/course/internal/svc"
	"tjxt/rpc/course/pb"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCategoryDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryDeleteLogic {
	return &CategoryDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CategoryDelete 递归删除分类及其全部子分类。
func (l *CategoryDeleteLogic) CategoryDelete(in *pb.IdRequest) (*pb.Empty, error) {
	if err := l.deleteRecursive(in.Id); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (l *CategoryDeleteLogic) deleteRecursive(id int64) error {
	children, err := l.svcCtx.CategoryModel.FindByParentId(l.ctx, id)
	if err != nil {
		return errorx.Wrap(err, errorx.CodeInternal, "查询子分类失败")
	}
	for _, ch := range children {
		if err := l.deleteRecursive(ch.Id); err != nil {
			return err
		}
	}
	if err := l.svcCtx.CategoryModel.Delete(l.ctx, id); err != nil {
		return errorx.Wrap(err, errorx.CodeInternal, "删除分类失败")
	}
	return nil
}
