package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/model"
	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveAccountRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveAccountRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveAccountRolesLogic {
	return &SaveAccountRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SaveAccountRoles 全量覆盖账户的角色分配。传空列表等价于收回该账户的全部角色。
func (l *SaveAccountRolesLogic) SaveAccountRoles(in *pb.AccountRoleReq) (*pb.Empty, error) {
	if in.AccountId <= 0 {
		return nil, xerr.BadRequestf("账户 id 无效")
	}

	// 逐个校验角色有效性：角色数量少，直接走带缓存的 FindOne。
	for _, roleId := range in.RoleIds {
		if roleId <= 0 {
			return nil, xerr.BadRequestf("角色 id 无效")
		}
		role, err := l.svcCtx.RoleModel.FindOne(l.ctx, roleId)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, xerr.BadRequestf("角色 %d 不存在", roleId)
			}
			return nil, err
		}
		if role.Deleted != 0 {
			return nil, xerr.BadRequestf("角色 %d 已删除", roleId)
		}
	}

	if err := l.svcCtx.AccountRoleModel.ReplaceByAccountId(l.ctx, in.AccountId, in.RoleIds); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
