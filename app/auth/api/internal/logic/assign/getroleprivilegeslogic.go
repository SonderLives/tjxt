// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package assign

import (
	"context"

	backendlogic "tjxt/app/auth/api/backend/logic"
	authpb "tjxt/app/auth/api/backend/pb"
	"tjxt/app/auth/api/internal/svc"
	"tjxt/app/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRolePrivilegesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRolePrivilegesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRolePrivilegesLogic {
	return &GetRolePrivilegesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetRolePrivileges 查询角色已分配的权限 id 列表。
func (l *GetRolePrivilegesLogic) GetRolePrivileges(req *types.IdPathReq) (resp *types.IdListVO, err error) {
	reply, err := backendlogic.NewGetRolePrivilegesLogic(l.ctx, l.svcCtx.Backend).GetRolePrivileges(&authpb.IdReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	ids := reply.Ids
	if ids == nil {
		ids = []int64{}
	}
	return &types.IdListVO{Ids: ids}, nil
}
