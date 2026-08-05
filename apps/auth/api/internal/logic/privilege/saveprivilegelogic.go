// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package privilege

import (
	"context"

	"tjxt/apps/auth/api/internal/svc"
	"tjxt/apps/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SavePrivilegeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSavePrivilegeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SavePrivilegeLogic {
	return &SavePrivilegeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SavePrivilegeLogic) SavePrivilege(req *types.PrivilegeSaveReq) (resp *types.IdVO, err error) {
	// todo: add your logic here and delete this line

	return
}
