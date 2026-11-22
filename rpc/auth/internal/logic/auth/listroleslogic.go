package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"
	"tjxt/pkg/utils/page"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListRolesLogic {
	return &ListRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListRoles 分页查询角色列表。
func (l *ListRolesLogic) ListRoles(in *pb.PageReq) (*pb.RoleListReply, error) {
	offset, limit := page.Normalize(int64(in.PageNo), int64(in.PageSize))

	roles, total, err := l.svcCtx.RoleModel.FindPage(l.ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	list := make([]*pb.RoleVO, 0, len(roles))
	for _, r := range roles {
		list = append(list, toRoleVO(r))
	}
	return &pb.RoleListReply{Total: total, List: list}, nil
}
