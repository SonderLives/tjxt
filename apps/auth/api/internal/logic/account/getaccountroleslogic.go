// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package account

import (
	"context"

	authclient "tjxt/apps/auth/rpc/client/auth"
	"tjxt/apps/auth/api/internal/svc"
	"tjxt/apps/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAccountRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAccountRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAccountRolesLogic {
	return &GetAccountRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetAccountRoles 查询账户已分配的角色 id 列表。
func (l *GetAccountRolesLogic) GetAccountRoles(req *types.IdPathReq) (resp *types.IdListVO, err error) {
	reply, err := l.svcCtx.AuthRpc.GetAccountRoles(l.ctx, &authclient.IdReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	ids := reply.Ids
	if ids == nil {
		ids = []int64{}
	}
	return &types.IdListVO{Ids: ids}, nil
}
