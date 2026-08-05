package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteInboxLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteInboxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteInboxLogic {
	return &DeleteInboxLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteInboxLogic) DeleteInbox(in *pb.IdReq) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
