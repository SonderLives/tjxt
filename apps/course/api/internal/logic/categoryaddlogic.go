package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryAddLogic {
	return &CategoryAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CategoryAdd 新增分类（透传 RPC）。
func (l *CategoryAddLogic) CategoryAdd(req *types.CategoryAddReq) (resp *types.NameExistVO, err error) {
	_, err = l.svcCtx.CourseRpc.CategoryAdd(l.ctx, &pb.CategoryAddRequest{
		Name:     req.Name,
		ParentId: req.ParentId,
		Index:    int32(req.Index),
	})
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "新增分类失败")
	}
	return &types.NameExistVO{Existed: false}, nil
}
