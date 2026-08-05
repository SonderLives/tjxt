package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"
	"tjxt/pkg/xerr"

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

// 通知模板 详情
func (l *GetNoticeTemplateLogic) GetNoticeTemplate(in *pb.IdReq) (*pb.NoticeTemplateVO, error) {
	data, err := l.svcCtx.NoticeTemplateModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("通知模板不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询通知模板失败")
	}
	return toNoticeTemplateVO(data), nil
}
