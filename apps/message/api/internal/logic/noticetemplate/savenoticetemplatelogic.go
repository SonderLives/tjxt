// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package noticetemplate

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
