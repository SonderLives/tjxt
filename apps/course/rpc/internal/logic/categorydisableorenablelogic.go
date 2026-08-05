package logic

import (
	"context"

	"tjxt/apps/course/rpc/internal/svc"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryDisableOrEnableLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCategoryDisableOrEnableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryDisableOrEnableLogic {
	return &CategoryDisableOrEnableLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CategoryDisableOrEnable 启用/禁用分类（1 启用 2 禁用）。
func (l *CategoryDisableOrEnableLogic) CategoryDisableOrEnable(in *pb.CategoryStatusRequest) (*pb.Empty, error) {
	c, err := l.svcCtx.CategoryModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("分类不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询分类失败")
	}
	c.Status = int64(in.Status)
	if err := l.svcCtx.CategoryModel.Update(l.ctx, c); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "更新分类状态失败")
	}
	return &pb.Empty{}, nil
}
