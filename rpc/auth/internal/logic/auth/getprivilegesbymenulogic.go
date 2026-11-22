package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPrivilegesByMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPrivilegesByMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPrivilegesByMenuLogic {
	return &GetPrivilegesByMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetPrivilegesByMenu 查询指定菜单下的权限列表。
func (l *GetPrivilegesByMenuLogic) GetPrivilegesByMenu(in *pb.IdReq) (*pb.PrivilegeListReply, error) {
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("菜单 id 无效")
	}

	privileges, err := l.svcCtx.PrivilegeModel.FindByMenuId(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}

	list := make([]*pb.PrivilegeVO, 0, len(privileges))
	for _, p := range privileges {
		list = append(list, toPrivilegeVO(p))
	}
	return &pb.PrivilegeListReply{List: list}, nil
}
