// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"context"

	authclient "tjxt/apps/auth/rpc/client/auth"
	"tjxt/apps/auth/api/internal/svc"
	"tjxt/apps/auth/api/internal/types"
	"tjxt/apps/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMenuTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuTreeLogic {
	return &GetMenuTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetMenuTree 获取完整菜单树（含子节点）。
func (l *GetMenuTreeLogic) GetMenuTree() (resp *types.MenuTreeVO, err error) {
	reply, err := l.svcCtx.AuthRpc.GetMenuTree(l.ctx, &authclient.Empty{})
	if err != nil {
		return nil, err
	}
	return &types.MenuTreeVO{List: toMenuVOList(reply.List)}, nil
}

func toMenuVOList(in []*pb.MenuVO) []types.MenuVO {
	out := make([]types.MenuVO, 0, len(in))
	for _, v := range in {
		out = append(out, toMenuVO(v))
	}
	return out
}

func toMenuVO(v *pb.MenuVO) types.MenuVO {
	children := make([]types.MenuVO, 0, len(v.Children))
	for _, c := range v.Children {
		children = append(children, toMenuVO(c))
	}
	return types.MenuVO{
		Id:          v.Id,
		ParentId:    v.ParentId,
		HasChildren: v.HasChildren,
		Label:       v.Label,
		Path:        v.Path,
		Icon:        v.Icon,
		Priority:    v.Priority,
		Children:    children,
	}
}
