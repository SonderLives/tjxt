// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package account

import (
	"context"

	backendlogic "tjxt/api/auth/backend/logic"
	authpb "tjxt/api/auth/backend/pb"
	"tjxt/api/auth/internal/svc"
	"tjxt/api/auth/internal/types"

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

// SaveAccountRoles 全量替换账户下的角色分配；路径 :id 为权威账户标识。
func (l *SaveAccountRolesLogic) SaveAccountRoles(req *types.AccountRoleReq) (resp *types.OkVO, err error) {
	_, err = backendlogic.NewSaveAccountRolesLogic(l.ctx, l.svcCtx.Backend).SaveAccountRoles(&authpb.AccountRoleReq{
		AccountId: req.Id,
		RoleIds:   req.RoleIds,
	})
	if err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
