package logic

import (
	"context"

	"tjxt/app/course/api/internal/svc"
	"tjxt/app/course/api/internal/types"
	"tjxt/app/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryListQueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryListQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryListQueryLogic {
	return &CategoryListQueryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CategoryListQuery 按条件查询分类列表。
func (l *CategoryListQueryLogic) CategoryListQuery(req *types.CategoryListQueryReq) (resp []types.CategoryVO, err error) {
	list, gerr := l.svcCtx.CourseRpc.CategoryListQuery(l.ctx, &pb.CategoryListQueryRequest{
		Name:   req.Name,
		Status: int32(req.Status),
	})
	if gerr != nil {
		return nil, gerr
	}
	resp = make([]types.CategoryVO, 0, len(list.Items))
	for _, c := range list.Items {
		resp = append(resp, toCategoryVO(c))
	}
	return resp, nil
}
