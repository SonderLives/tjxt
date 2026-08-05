// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package noticetemplate

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
