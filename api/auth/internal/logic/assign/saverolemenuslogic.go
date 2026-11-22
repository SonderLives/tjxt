// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package assign

import (
	"context"

	backendlogic "tjxt/api/auth/backend/logic"
	authpb "tjxt/api/auth/backend/pb"
	"tjxt/api/auth/internal/svc"
	"tjxt/api/auth/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveRoleMenusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveRoleMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveRoleMenusLogic {
	return &SaveRoleMenusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SaveRoleMenus 全量替换角色下的菜单分配；路径 :id 为权威角色标识。
func (l *SaveRoleMenusLogic) SaveRoleMenus(req *types.RoleMenuReq) (resp *types.OkVO, err error) {
	_, err = backendlogic.NewSaveRoleMenusLogic(l.ctx, l.svcCtx.Backend).SaveRoleMenus(&authpb.RoleMenuReq{
		RoleId:  req.Id,
		MenuIds: req.MenuIds,
	})
	if err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
