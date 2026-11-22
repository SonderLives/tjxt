package logic

import (
	"context"

	"tjxt/api/course/internal/svc"
	"tjxt/api/course/internal/types"
	"tjxt/rpc/course/pb"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryDeleteLogic {
	return &CategoryDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CategoryDelete 删除分类。
func (l *CategoryDeleteLogic) CategoryDelete(req *types.IdPathReq) (resp *types.NameExistVO, err error) {
	_, err = l.svcCtx.CourseRpc.CategoryDelete(l.ctx, &pb.IdRequest{Id: req.Id})
	if err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "删除分类失败")
	}
	return &types.NameExistVO{Existed: false}, nil
}
