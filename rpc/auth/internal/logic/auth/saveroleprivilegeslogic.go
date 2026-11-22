package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/model"
	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveRolePrivilegesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveRolePrivilegesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveRolePrivilegesLogic {
	return &SaveRolePrivilegesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SaveRolePrivileges 全量覆盖角色的权限分配。传空列表等价于收回该角色的全部权限。
func (l *SaveRolePrivilegesLogic) SaveRolePrivileges(in *pb.RolePrivilegeReq) (*pb.Empty, error) {
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

	// 校验权限均有效，避免授出指向已删除权限的空壳授权。
	if len(in.PrivilegeIds) > 0 {
		privileges, err := l.svcCtx.PrivilegeModel.FindByIds(l.ctx, in.PrivilegeIds)
		if err != nil {
			return nil, err
		}
		valid := make(map[int64]struct{}, len(privileges))
		for _, p := range privileges {
			valid[p.Id] = struct{}{}
		}
		for _, id := range in.PrivilegeIds {
			if _, ok := valid[id]; !ok {
				return nil, xerr.BadRequestf("权限 %d 不存在或已删除", id)
			}
		}
	}

	if err := l.svcCtx.RolePrivilegeModel.ReplaceByRoleId(l.ctx, in.RoleId, in.PrivilegeIds); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
