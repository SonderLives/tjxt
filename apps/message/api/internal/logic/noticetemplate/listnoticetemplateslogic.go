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

type ListNoticeTemplatesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListNoticeTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListNoticeTemplatesLogic {
	return &ListNoticeTemplatesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListNoticeTemplatesLogic) ListNoticeTemplates(req *types.PageReq) (resp *types.NoticeTemplateListVO, err error) {
	r, err := l.svcCtx.MessageRpc.ListNoticeTemplates(l.ctx, &messageclient.PageReq{
		PageNo:   int32(req.PageNo),
		PageSize: int32(req.PageSize),
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.NoticeTemplateVO, 0, len(r.List))
	for _, item := range r.List {
		list = append(list, types.NoticeTemplateVO{
			Id:            item.Id,
			Name:          item.Name,
			Code:          item.Code,
			Type:          item.Type,
			Status:        item.Status,
			Title:         item.Title,
			Content:       item.Content,
			IsSmsTemplate: item.IsSmsTemplate,
			CreateTime:    item.CreateTime,
		})
	}
	return &types.NoticeTemplateListVO{Total: r.Total, List: list}, nil
}
