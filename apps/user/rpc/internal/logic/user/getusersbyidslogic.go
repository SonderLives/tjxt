package userlogic

import (
	"context"

	"tjxt/apps/user/rpc/internal/svc"
	"tjxt/apps/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUsersByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUsersByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersByIdsLogic {
	return &GetUsersByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUsersByIdsLogic) GetUsersByIds(in *pb.UserIdsRequest) (*pb.UserListResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.UserListResponse{}, nil
}
