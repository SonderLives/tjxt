package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/model"
	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuLogic {
	return &DeleteMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteMenu 逻辑删除菜单，并级联清理其下权限与角色分配。
// 存在子菜单时拒绝删除，要求调用方自下而上处理，避免整棵子树被隐式摘除。
func (l *DeleteMenuLogic) DeleteMenu(in *pb.IdReq) (*pb.Empty, error) {
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("菜单 id 无效")
	}

	menu, err := l.svcCtx.MenuModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NotFound("菜单不存在")
		}
		return nil, err
	}
	if menu.Deleted != 0 {
		return nil, xerr.NotFound("菜单不存在")
	}

	children, err := l.svcCtx.MenuModel.CountChildren(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if children > 0 {
		return nil, xerr.Conflict("存在子菜单，请先删除子菜单")
	}

	// 先摘除该菜单下的权限及其分配，再删菜单本身，保证不残留可用权限。
	privileges, err := l.svcCtx.PrivilegeModel.FindByMenuId(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	for _, p := range privileges {
		if err := l.svcCtx.RolePrivilegeModel.DeleteByPrivilegeId(l.ctx, p.Id); err != nil {
			return nil, err
		}
		if err := l.svcCtx.PrivilegeModel.SoftDelete(l.ctx, p.Id, 0); err != nil {
			return nil, err
		}
	}

	if err := l.svcCtx.RoleMenuModel.DeleteByMenuId(l.ctx, in.Id); err != nil {
		return nil, err
	}
	if err := l.svcCtx.MenuModel.SoftDelete(l.ctx, in.Id, 0); err != nil {
		return nil, err
	}
	// 删除后父节点可能已无子菜单，重算标记。
	if err := l.svcCtx.MenuModel.SyncHasChildren(l.ctx, menu.ParentId); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
