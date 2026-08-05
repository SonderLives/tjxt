// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package noticetemplate

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
