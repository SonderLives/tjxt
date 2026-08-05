package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"

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

// 权限
func (l *GetPrivilegesByMenuLogic) GetPrivilegesByMenu(in *pb.IdReq) (*pb.PrivilegeListReply, error) {
	// todo: add your logic here and delete this line

	return &pb.PrivilegeListReply{}, nil
}
