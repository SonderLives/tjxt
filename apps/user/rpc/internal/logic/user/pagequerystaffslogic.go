package userlogic

import (
	"context"

	"tjxt/apps/user/rpc/internal/svc"
	"tjxt/apps/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageQueryStaffsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageQueryStaffsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageQueryStaffsLogic {
	return &PageQueryStaffsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PageQueryStaffsLogic) PageQueryStaffs(in *pb.PageQueryUsersRequest) (*pb.PageQueryUsersResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.PageQueryUsersResponse{}, nil
}
