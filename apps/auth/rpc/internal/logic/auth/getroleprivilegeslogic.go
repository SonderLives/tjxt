package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRolePrivilegesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRolePrivilegesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRolePrivilegesLogic {
	return &GetRolePrivilegesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetRolePrivileges 查询角色已分配的权限 id 列表。
func (l *GetRolePrivilegesLogic) GetRolePrivileges(in *pb.IdReq) (*pb.IdListReply, error) {
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("角色 id 无效")
	}

	ids, err := l.svcCtx.RolePrivilegeModel.FindPrivilegeIdsByRoleId(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if ids == nil {
		ids = []int64{}
	}
	return &pb.IdListReply{Ids: ids}, nil
}
