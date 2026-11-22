package messagelogic

import (
	"context"
	"errors"

	"tjxt/pkg/errorx"
	"tjxt/rpc/message/internal/model"
	"tjxt/rpc/message/internal/svc"
	"tjxt/rpc/message/pb"

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
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("通知模板不存在")
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询通知模板失败")
	}
	return toNoticeTemplateVO(data), nil
}
