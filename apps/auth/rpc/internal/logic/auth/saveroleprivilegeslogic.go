package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"

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

func (l *SaveRolePrivilegesLogic) SaveRolePrivileges(in *pb.RolePrivilegeReq) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
