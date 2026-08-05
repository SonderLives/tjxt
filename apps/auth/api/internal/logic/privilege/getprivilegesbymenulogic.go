// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package privilege

import (
	"context"

	"tjxt/apps/auth/api/internal/svc"
	"tjxt/apps/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPrivilegesByMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPrivilegesByMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPrivilegesByMenuLogic {
	return &GetPrivilegesByMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPrivilegesByMenuLogic) GetPrivilegesByMenu(req *types.IdPathReq) (resp *types.PrivilegeListVO, err error) {
	// todo: add your logic here and delete this line

	return
}
