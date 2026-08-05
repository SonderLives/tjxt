// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package assign

import (
	"context"

	"tjxt/apps/auth/api/internal/svc"
	"tjxt/apps/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveRoleMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveRoleMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveRoleMenusLogic {
	return &SaveRoleMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveRoleMenusLogic) SaveRoleMenus(req *types.RoleMenuReq) (resp *types.OkVO, err error) {
	// todo: add your logic here and delete this line

	return
}
