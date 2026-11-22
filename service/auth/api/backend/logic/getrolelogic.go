package authlogic

import (
	"context"

	"tjxt/api/auth/backend/model"
	"tjxt/api/auth/backend/pb"
	"tjxt/api/auth/backend/svc"
	"tjxt/pkg/errorx"

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
		return nil, errorx.BadRequestf("角色 id 无效")
	}

	role, err := l.svcCtx.RoleModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errorx.NotFound("角色不存在")
		}
		return nil, err
	}
	if role.Deleted != 0 {
		return nil, errorx.NotFound("角色不存在")
	}
	return toRoleVO(role), nil
}
