// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"context"

	"tjxt/api/auth/internal/svc"
	"tjxt/api/auth/internal/types"
	authclient "tjxt/rpc/auth/client/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DeleteRole 删除角色（软删），并级联清理角色-菜单、角色-权限分配。
func (l *DeleteRoleLogic) DeleteRole(req *types.IdPathReq) (resp *types.OkVO, err error) {
	_, err = l.svcCtx.AuthRpc.DeleteRole(l.ctx, &authclient.IdReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
