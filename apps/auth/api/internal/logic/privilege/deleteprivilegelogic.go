// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package privilege

import (
	"context"

	authclient "tjxt/apps/auth/rpc/client/auth"
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

// DeletePrivilege 删除权限（软删），并清理角色-权限分配。
func (l *DeletePrivilegeLogic) DeletePrivilege(req *types.IdPathReq) (resp *types.OkVO, err error) {
	_, err = l.svcCtx.AuthRpc.DeletePrivilege(l.ctx, &authclient.IdReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
