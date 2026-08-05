package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"

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

// 账户角色
func (l *SaveAccountRolesLogic) SaveAccountRoles(in *pb.AccountRoleReq) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
