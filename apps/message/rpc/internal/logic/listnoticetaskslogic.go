package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListNoticeTasksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListNoticeTasksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNoticeTasksLogic {
	return &ListNoticeTasksLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListNoticeTasksLogic) ListNoticeTasks(in *pb.PageReq) (*pb.NoticeTaskListReply, error) {
	// todo: add your logic here and delete this line

	return &pb.NoticeTaskListReply{}, nil
}
