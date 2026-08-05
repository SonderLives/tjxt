package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRoleMenusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRoleMenusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRoleMenusLogic {
	return &GetRoleMenusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRoleMenusLogic) GetRoleMenus(in *pb.IdReq) (*pb.IdListReply, error) {
	// todo: add your logic here and delete this line

	return &pb.IdListReply{}, nil
}
