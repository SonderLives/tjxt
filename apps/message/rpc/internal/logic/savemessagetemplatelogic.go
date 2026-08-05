package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveMessageTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveMessageTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveMessageTemplateLogic {
	return &SaveMessageTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 短信模板
func (l *SaveMessageTemplateLogic) SaveMessageTemplate(in *pb.MessageTemplateSaveReq) (*pb.IdReply, error) {
	// todo: add your logic here and delete this line

	return &pb.IdReply{}, nil
}
