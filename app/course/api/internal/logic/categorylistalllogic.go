package logic

import (
	"context"

	"tjxt/service/course/api/internal/svc"
	"tjxt/service/course/api/internal/types"
	"tjxt/service/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryListAllLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryListAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryListAllLogic {
	return &CategoryListAllLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CategoryListAll 查询全部分类（树形）。
func (l *CategoryListAllLogic) CategoryListAll(req *types.CategoryListAllReq) (resp []types.SimpleCategoryVO, err error) {
	list, gerr := l.svcCtx.CourseRpc.CategoryListAll(l.ctx, &pb.CategoryListAllRequest{
		Admin: req.Admin,
		Name:  req.Name,
	})
	if gerr != nil {
		return nil, gerr
	}
	resp = make([]types.SimpleCategoryVO, 0, len(list.Items))
	for _, n := range list.Items {
		resp = append(resp, *toSimpleCategoryVO(n))
	}
	return resp, nil
}
