// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package assign

import (
	"context"

	"tjxt/apps/auth/api/internal/svc"
	"tjxt/apps/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveRolePrivilegesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveRolePrivilegesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveRolePrivilegesLogic {
	return &SaveRolePrivilegesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveRolePrivilegesLogic) SaveRolePrivileges(req *types.RolePrivilegeReq) (resp *types.OkVO, err error) {
	// todo: add your logic here and delete this line

	return
}
