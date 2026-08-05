// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package assign

import (
	"context"

	"tjxt/apps/auth/api/internal/svc"
	"tjxt/apps/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRolePrivilegesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRolePrivilegesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRolePrivilegesLogic {
	return &GetRolePrivilegesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRolePrivilegesLogic) GetRolePrivileges(req *types.IdPathReq) (resp *types.IdListVO, err error) {
	// todo: add your logic here and delete this line

	return
}
