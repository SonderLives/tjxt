package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListInboxLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListInboxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListInboxLogic {
	return &ListInboxLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListInboxLogic) ListInbox(in *pb.InboxPageReq) (*pb.InboxListReply, error) {
	// todo: add your logic here and delete this line

	return &pb.InboxListReply{}, nil
}
