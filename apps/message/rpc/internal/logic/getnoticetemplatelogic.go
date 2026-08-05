package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNoticeTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNoticeTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoticeTemplateLogic {
	return &GetNoticeTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNoticeTemplateLogic) GetNoticeTemplate(in *pb.IdReq) (*pb.NoticeTemplateVO, error) {
	// todo: add your logic here and delete this line

	return &pb.NoticeTemplateVO{}, nil
}
