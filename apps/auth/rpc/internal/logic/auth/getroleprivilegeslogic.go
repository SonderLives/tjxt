package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"

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

func (l *GetRolePrivilegesLogic) GetRolePrivileges(in *pb.IdReq) (*pb.IdListReply, error) {
	// todo: add your logic here and delete this line

	return &pb.IdListReply{}, nil
}
