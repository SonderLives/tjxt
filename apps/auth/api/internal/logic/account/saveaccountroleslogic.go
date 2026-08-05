// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package account

import (
	"context"

	"tjxt/apps/auth/api/internal/svc"
	"tjxt/apps/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveAccountRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveAccountRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveAccountRolesLogic {
	return &SaveAccountRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveAccountRolesLogic) SaveAccountRoles(req *types.AccountRoleReq) (resp *types.OkVO, err error) {
	// todo: add your logic here and delete this line

	return
}
