package logic

import (
	"context"

	"tjxt/api/course/internal/svc"
	"tjxt/api/course/internal/types"
	"tjxt/rpc/course/pb"
	"tjxt/pkg/errorx"

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
		Id:     req.Id,
		Status: int32(req.Status),
	})
	if err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "更新分类状态失败")
	}
	return &types.NameExistVO{Existed: false}, nil
}
