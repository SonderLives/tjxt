package logic

import (
	"context"

	"tjxt/service/course/api/internal/svc"
	"tjxt/service/course/api/internal/types"
	"tjxt/service/course/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryListOneLevelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryListOneLevelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryListOneLevelLogic {
	return &CategoryListOneLevelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CategoryListOneLevel 查询全部一级分类。
func (l *CategoryListOneLevelLogic) CategoryListOneLevel() (resp []types.CategoryVO, err error) {
	list, gerr := l.svcCtx.CourseRpc.CategoryListOneLevel(l.ctx, &pb.Empty{})
	if gerr != nil {
		return nil, gerr
	}
	resp = make([]types.CategoryVO, 0, len(list.Items))
	for _, c := range list.Items {
		resp = append(resp, toCategoryVO(c))
	}
	return resp, nil
}
