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

type GetNoticeTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetNoticeTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoticeTemplateLogic {
	return &GetNoticeTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetNoticeTemplateLogic) GetNoticeTemplate(req *types.IdPathReq) (resp *types.NoticeTemplateVO, err error) {
	r, err := l.svcCtx.MessageRpc.GetNoticeTemplate(l.ctx, &messageclient.IdReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.NoticeTemplateVO{
		Id:            r.Id,
		Name:          r.Name,
		Code:          r.Code,
		Type:          r.Type,
		Status:        r.Status,
		Title:         r.Title,
		Content:       r.Content,
		IsSmsTemplate: r.IsSmsTemplate,
		CreateTime:    r.CreateTime,
	}, nil
}
