package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAccountRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAccountRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAccountRolesLogic {
	return &GetAccountRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAccountRolesLogic) GetAccountRoles(in *pb.IdReq) (*pb.IdListReply, error) {
	// todo: add your logic here and delete this line

	return &pb.IdListReply{}, nil
}
