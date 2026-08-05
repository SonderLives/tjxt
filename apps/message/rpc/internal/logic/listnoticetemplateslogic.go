package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListNoticeTemplatesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListNoticeTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNoticeTemplatesLogic {
	return &ListNoticeTemplatesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListNoticeTemplatesLogic) ListNoticeTemplates(in *pb.PageReq) (*pb.NoticeTemplateListReply, error) {
	// todo: add your logic here and delete this line

	return &pb.NoticeTemplateListReply{}, nil
}
