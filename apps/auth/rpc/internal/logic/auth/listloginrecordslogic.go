package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLoginRecordsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLoginRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLoginRecordsLogic {
	return &ListLoginRecordsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListLoginRecordsLogic) ListLoginRecords(in *pb.LoginRecordPageReq) (*pb.LoginRecordListReply, error) {
	// todo: add your logic here and delete this line

	return &pb.LoginRecordListReply{}, nil
}
