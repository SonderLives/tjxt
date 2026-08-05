// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"context"

	authclient "tjxt/apps/auth/rpc/client/auth"
	"tjxt/apps/auth/api/internal/svc"
	"tjxt/apps/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveMenuLogic {
	return &SaveMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SaveMenu 新增或修改菜单，返回菜单 id。
func (l *SaveMenuLogic) SaveMenu(req *types.MenuSaveReq) (resp *types.IdVO, err error) {
	reply, err := l.svcCtx.AuthRpc.SaveMenu(l.ctx, &authclient.MenuSaveReq{
		Id:       req.Id,
		ParentId: req.ParentId,
		Label:    req.Label,
		Path:     req.Path,
		Icon:     req.Icon,
		Priority: req.Priority,
	})
	if err != nil {
		return nil, err
	}
	return &types.IdVO{Id: reply.Id}, nil
}
