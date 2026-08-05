package userlogic

import (
	"context"

	"tjxt/apps/user/rpc/internal/svc"
	"tjxt/apps/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageQueryTeachersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageQueryTeachersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageQueryTeachersLogic {
	return &PageQueryTeachersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PageQueryTeachersLogic) PageQueryTeachers(in *pb.PageQueryUsersRequest) (*pb.PageQueryUsersResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.PageQueryUsersResponse{}, nil
}
