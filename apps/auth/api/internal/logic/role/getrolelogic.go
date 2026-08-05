// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"context"

	authclient "tjxt/apps/auth/rpc/client/auth"
	"tjxt/apps/auth/api/internal/svc"
	"tjxt/apps/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleLogic {
	return &GetRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetRole 查询单个角色详情。
func (l *GetRoleLogic) GetRole(req *types.IdPathReq) (resp *types.RoleVO, err error) {
	reply, err := l.svcCtx.AuthRpc.GetRole(l.ctx, &authclient.IdReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.RoleVO{
		Id:         reply.Id,
		Code:       reply.Code,
		Name:       reply.Name,
		Type:       reply.Type,
		CreateTime: reply.CreateTime,
	}, nil
}
