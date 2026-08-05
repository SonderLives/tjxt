package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkInboxReadLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMarkInboxReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkInboxReadLogic {
	return &MarkInboxReadLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MarkInboxReadLogic) MarkInboxRead(in *pb.IdReq) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
