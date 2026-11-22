package messagelogic

import (
	"context"

	"tjxt/service/message/api/backend/pb"
	"tjxt/service/message/api/backend/svc"
	"tjxt/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMessageTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMessageTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMessageTemplateLogic {
	return &DeleteMessageTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 短信模板 物理删除
func (l *DeleteMessageTemplateLogic) DeleteMessageTemplate(in *pb.IdReq) (*pb.Empty, error) {
	if err := l.svcCtx.MessageTemplateModel.Delete(l.ctx, in.Id); err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "删除短信模板失败")
	}
	return &pb.Empty{}, nil
}
