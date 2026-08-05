package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SavePrivilegeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSavePrivilegeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SavePrivilegeLogic {
	return &SavePrivilegeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SavePrivilegeLogic) SavePrivilege(in *pb.PrivilegeSaveReq) (*pb.IdReply, error) {
	// todo: add your logic here and delete this line

	return &pb.IdReply{}, nil
}
