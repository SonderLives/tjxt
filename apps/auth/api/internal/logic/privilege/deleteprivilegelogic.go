// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package privilege

import (
	"context"

	"tjxt/apps/auth/api/internal/svc"
	"tjxt/apps/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePrivilegeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePrivilegeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePrivilegeLogic {
	return &DeletePrivilegeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePrivilegeLogic) DeletePrivilege(req *types.IdPathReq) (resp *types.OkVO, err error) {
	// todo: add your logic here and delete this line

	return
}
