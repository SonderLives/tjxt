// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package noticetemplate

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"
	messageclient "tjxt/apps/message/rpc/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveNoticeTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveNoticeTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveNoticeTemplateLogic {
	return &SaveNoticeTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveNoticeTemplateLogic) SaveNoticeTemplate(req *types.NoticeTemplateSaveReq) (resp *types.IdVO, err error) {
	r, err := l.svcCtx.MessageRpc.SaveNoticeTemplate(l.ctx, &messageclient.NoticeTemplateSaveReq{
		Id:            req.Id,
		Name:          req.Name,
		Code:          req.Code,
		Type:          req.Type,
		Title:         req.Title,
		Content:       req.Content,
		IsSmsTemplate: req.IsSmsTemplate,
	})
	if err != nil {
		return nil, err
	}
	return &types.IdVO{Id: r.Id}, nil
}
