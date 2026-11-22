package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/model"
	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePrivilegeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePrivilegeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePrivilegeLogic {
	return &DeletePrivilegeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeletePrivilege 逻辑删除权限，并清理其角色分配。
func (l *DeletePrivilegeLogic) DeletePrivilege(in *pb.IdReq) (*pb.Empty, error) {
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("权限 id 无效")
	}

	privilege, err := l.svcCtx.PrivilegeModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NotFound("权限不存在")
		}
		return nil, err
	}
	if privilege.Deleted != 0 {
		return nil, xerr.NotFound("权限不存在")
	}

	// 先断开角色引用，再删权限，避免鉴权时命中已删除的权限行。
	if err := l.svcCtx.RolePrivilegeModel.DeleteByPrivilegeId(l.ctx, in.Id); err != nil {
		return nil, err
	}
	if err := l.svcCtx.PrivilegeModel.SoftDelete(l.ctx, in.Id, 0); err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}
