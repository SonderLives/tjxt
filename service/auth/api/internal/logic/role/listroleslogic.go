// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"context"

	backendlogic "tjxt/service/auth/api/backend/logic"
	authpb "tjxt/service/auth/api/backend/pb"
	"tjxt/service/auth/api/internal/svc"
	"tjxt/service/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRolesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRolesLogic {
	return &ListRolesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListRoles 分页查询角色列表。
func (l *ListRolesLogic) ListRoles(req *types.RoleListReq) (resp *types.RoleListVO, err error) {
	reply, err := backendlogic.NewListRolesLogic(l.ctx, l.svcCtx.Backend).ListRoles(&authpb.PageReq{
		PageNo:   int32(req.PageNo),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.RoleVO, 0, len(reply.List))
	for _, v := range reply.List {
		list = append(list, types.RoleVO{
			Id:         v.Id,
			Code:       v.Code,
			Name:       v.Name,
			Type:       v.Type,
			CreateTime: v.CreateTime,
		})
	}
	return &types.RoleListVO{Total: reply.Total, List: list}, nil
}
