package logic

import (
	"context"

	"tjxt/api/course/internal/svc"
	"tjxt/api/course/internal/types"
	"tjxt/rpc/course/pb"
	"tjxt/pkg/errorx"

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
		return nil, errorx.Wrap(err, errorx.CodeInternal, "新增分类失败")
	}
	return &types.NameExistVO{Existed: false}, nil
}
