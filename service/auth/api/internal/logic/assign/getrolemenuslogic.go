// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package assign

import (
	"context"

	backendlogic "tjxt/service/auth/api/backend/logic"
	authpb "tjxt/service/auth/api/backend/pb"
	"tjxt/service/auth/api/internal/svc"
	"tjxt/service/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoleMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleMenusLogic {
	return &GetRoleMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetRoleMenus 查询角色已分配的菜单 id 列表。
func (l *GetRoleMenusLogic) GetRoleMenus(req *types.IdPathReq) (resp *types.IdListVO, err error) {
	reply, err := backendlogic.NewGetRoleMenusLogic(l.ctx, l.svcCtx.Backend).GetRoleMenus(&authpb.IdReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	ids := reply.Ids
	if ids == nil {
		ids = []int64{}
	}
	return &types.IdListVO{Ids: ids}, nil
}
