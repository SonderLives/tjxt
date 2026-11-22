package logic

import (
	"context"

	"tjxt/api/course/internal/svc"
	"tjxt/api/course/internal/types"
	"tjxt/rpc/course/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryUpdateLogic {
	return &CategoryUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CategoryUpdate 更新分类。
func (l *CategoryUpdateLogic) CategoryUpdate(req *types.CategoryUpdateReq) (resp *types.NameExistVO, err error) {
	_, err = l.svcCtx.CourseRpc.CategoryUpdate(l.ctx, &pb.CategoryUpdateRequest{
		Id:    req.Id,
		Name:  req.Name,
		Index: int32(req.Index),
	})
	if err != nil {
		return nil, err
	}
	return &types.NameExistVO{Existed: false}, nil
}
