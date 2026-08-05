package userlogic

import (
	"context"

	"tjxt/apps/user/rpc/internal/svc"
	"tjxt/apps/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageQueryStudentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageQueryStudentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageQueryStudentsLogic {
	return &PageQueryStudentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 管理后台分页查询 =====
func (l *PageQueryStudentsLogic) PageQueryStudents(in *pb.PageQueryUsersRequest) (*pb.PageQueryUsersResponse, error) {
	// todo: add your logic here and delete this line

	return &pb.PageQueryUsersResponse{}, nil
}
