// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"context"

	"tjxt/api/auth/internal/svc"
	"tjxt/api/auth/internal/types"
	authclient "tjxt/rpc/auth/client/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DeleteMenu 删除菜单（软删），并清理其下权限与角色-菜单分配。
func (l *DeleteMenuLogic) DeleteMenu(req *types.IdPathReq) (resp *types.OkVO, err error) {
	_, err = l.svcCtx.AuthRpc.DeleteMenu(l.ctx, &authclient.IdReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
