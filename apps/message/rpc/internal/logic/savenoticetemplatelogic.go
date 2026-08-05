package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveNoticeTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveNoticeTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveNoticeTemplateLogic {
	return &SaveNoticeTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通知模板
func (l *SaveNoticeTemplateLogic) SaveNoticeTemplate(in *pb.NoticeTemplateSaveReq) (*pb.IdReply, error) {
	// todo: add your logic here and delete this line

	return &pb.IdReply{}, nil
}
