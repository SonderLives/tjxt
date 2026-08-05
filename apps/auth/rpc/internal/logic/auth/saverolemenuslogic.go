package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/model"
	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveRoleMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveRoleMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveRoleMenusLogic {
	return &SaveRoleMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SaveRoleMenus 全量覆盖角色的菜单分配。传空列表等价于清空该角色的全部菜单。
func (l *SaveRoleMenusLogic) SaveRoleMenus(in *pb.RoleMenuReq) (*pb.Empty, error) {
	if in.RoleId <= 0 {
		return nil, xerr.BadRequestf("角色 id 无效")
	}

	role, err := l.svcCtx.RoleModel.FindOne(l.ctx, in.RoleId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NotFound("角色不存在")
		}
		return nil, err
	}
	if role.Deleted != 0 {
		return nil, xerr.NotFound("角色不存在")
	}

	// 校验菜单均有效，避免写入指向已删除菜单的脏分配。
	if len(in.MenuIds) > 0 {
		menus, err := l.svcCtx.MenuModel.FindByIds(l.ctx, in.MenuIds)
		if err != nil {
			return nil, err
		}
		valid := make(map[int64]struct{}, len(menus))
		for _, m := range menus {
			valid[m.Id] = struct{}{}
		}
		for _, id := range in.MenuIds {
			if _, ok := valid[id]; !ok {
				return nil, xerr.BadRequestf("菜单 %d 不存在或已删除", id)
			}
		}
	}

	if err := l.svcCtx.RoleMenuModel.ReplaceByRoleId(l.ctx, in.RoleId, in.MenuIds); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
