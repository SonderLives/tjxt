package messagelogic

import (
	"context"

	"tjxt/service/message/api/backend/pb"
	"tjxt/service/message/api/backend/svc"
	"tjxt/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteNoticeTemplateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteNoticeTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNoticeTemplateLogic {
	return &DeleteNoticeTemplateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通知模板 物理删除
func (l *DeleteNoticeTemplateLogic) DeleteNoticeTemplate(in *pb.IdReq) (*pb.Empty, error) {
	if err := l.svcCtx.NoticeTemplateModel.Delete(l.ctx, in.Id); err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "删除通知模板失败")
	}
	return &pb.Empty{}, nil
}
