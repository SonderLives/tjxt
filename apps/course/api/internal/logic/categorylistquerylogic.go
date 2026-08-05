package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

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
		return nil, xerr.Wrap(gerr, xerr.CodeInternal, "查询分类列表失败")
	}
	resp = make([]types.CategoryVO, 0, len(list.Items))
	for _, c := range list.Items {
		resp = append(resp, toCategoryVO(c))
	}
	return resp, nil
}
