// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package assign

import (
	"context"

	authclient "tjxt/apps/auth/rpc/client/auth"
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

// SaveRolePrivileges 全量替换角色下的权限分配；路径 :id 为权威角色标识。
func (l *SaveRolePrivilegesLogic) SaveRolePrivileges(req *types.RolePrivilegeReq) (resp *types.OkVO, err error) {
	_, err = l.svcCtx.AuthRpc.SaveRolePrivileges(l.ctx, &authclient.RolePrivilegeReq{
		RoleId:       req.Id,
		PrivilegeIds: req.PrivilegeIds,
	})
	if err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
