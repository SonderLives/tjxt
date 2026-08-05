package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/model"
	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleLogic {
	return &GetRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetRole 按主键查询角色详情。
func (l *GetRoleLogic) GetRole(in *pb.IdReq) (*pb.RoleVO, error) {
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
	return toRoleVO(role), nil
}
