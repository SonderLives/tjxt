package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/model"
	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteRoleLogic {
	return &DeleteRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteRole 逻辑删除角色，并级联清理其菜单/权限分配。
// 固定角色与仍被账户引用的角色不允许删除。
func (l *DeleteRoleLogic) DeleteRole(in *pb.IdReq) (*pb.Empty, error) {
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("角色 id 无效")
	}

	role, err := l.svcCtx.RoleModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NotFound("角色不存在")
		}
		return nil, err
	}
	if role.Deleted != 0 {
		return nil, xerr.NotFound("角色不存在")
	}
	if role.Type == fixedRoleType {
		return nil, xerr.Conflict("固定角色不允许删除")
	}

	// 仍被账户引用时拒绝删除，避免出现权限悬空的账号。
	used, err := l.svcCtx.AccountRoleModel.CountByRoleId(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if used > 0 {
		return nil, xerr.Conflict("该角色已分配给账户，请先解除分配")
	}

	if err := l.svcCtx.RoleModel.SoftDelete(l.ctx, in.Id, 0); err != nil {
		return nil, err
	}
	// 级联清理分配关系，防止角色 id 复用时继承到旧权限。
	if err := l.svcCtx.RoleMenuModel.DeleteByRoleId(l.ctx, in.Id); err != nil {
		return nil, err
	}
	if err := l.svcCtx.RolePrivilegeModel.DeleteByRoleId(l.ctx, in.Id); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
