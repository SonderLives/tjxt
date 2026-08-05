package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"
	"tjxt/apps/course/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryDisableOrEnableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryDisableOrEnableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryDisableOrEnableLogic {
	return &CategoryDisableOrEnableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CategoryDisableOrEnable 启用/禁用分类。
func (l *CategoryDisableOrEnableLogic) CategoryDisableOrEnable(req *types.CategoryDisableReq) (resp *types.NameExistVO, err error) {
	_, err = l.svcCtx.CourseRpc.CategoryDisableOrEnable(l.ctx, &pb.CategoryStatusRequest{
		Id:    req.Id,
		Status: int32(req.Status),
	})
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "更新分类状态失败")
	}
	return &types.NameExistVO{Existed: false}, nil
}
