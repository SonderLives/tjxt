// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package account

import (
	"context"

	"tjxt/apps/auth/api/internal/svc"
	"tjxt/apps/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAccountRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAccountRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAccountRolesLogic {
	return &GetAccountRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAccountRolesLogic) GetAccountRoles(req *types.IdPathReq) (resp *types.IdListVO, err error) {
	// todo: add your logic here and delete this line

	return
}
