// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package privilege

import (
	"context"

	backendlogic "tjxt/api/auth/backend/logic"
	authpb "tjxt/api/auth/backend/pb"
	"tjxt/api/auth/internal/svc"
	"tjxt/api/auth/internal/types"

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
	_, err = backendlogic.NewDeletePrivilegeLogic(l.ctx, l.svcCtx.Backend).DeletePrivilege(&authpb.IdReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
