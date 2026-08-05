package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMessageTemplatesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMessageTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMessageTemplatesLogic {
	return &ListMessageTemplatesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListMessageTemplatesLogic) ListMessageTemplates(in *pb.PageReq) (*pb.MessageTemplateListReply, error) {
	// todo: add your logic here and delete this line

	return &pb.MessageTemplateListReply{}, nil
}
