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

type SaveRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveRoleLogic {
	return &SaveRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SaveRole 新增或修改角色，返回新建或更新目标的 id。
func (l *SaveRoleLogic) SaveRole(req *types.RoleSaveReq) (resp *types.IdVO, err error) {
	reply, err := l.svcCtx.AuthRpc.SaveRole(l.ctx, &authclient.RoleSaveReq{
		Id:   req.Id,
		Code: req.Code,
		Name: req.Name,
		Type: req.Type,
	})
	if err != nil {
		return nil, err
	}
	return &types.IdVO{Id: reply.Id}, nil
}
