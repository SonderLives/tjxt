package admin

import (
	"context"

	"tjxt/service/user/api/internal/logic/convert"
	"tjxt/service/user/api/internal/svc"
	"tjxt/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageQueryStaffsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPageQueryStaffsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageQueryStaffsLogic {
	return &PageQueryStaffsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// PageQueryStaffs 员工分页查询。
func (l *PageQueryStaffsLogic) PageQueryStaffs(req *types.UserPageReq) (resp *types.StaffVOList, err error) {
	out, err := l.svcCtx.UserRpc.PageQueryStaffs(l.ctx, convert.FromUserPageReq(req))
	if err != nil {
		return nil, err
	}
	return convert.ToStaffVOList(out), nil
}
