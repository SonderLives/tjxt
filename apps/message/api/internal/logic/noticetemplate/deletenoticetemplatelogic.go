// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package noticetemplate

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteNoticeTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteNoticeTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNoticeTemplateLogic {
	return &DeleteNoticeTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteNoticeTemplateLogic) DeleteNoticeTemplate(req *types.IdPathReq) (resp *types.OkVO, err error) {
	// todo: add your logic here and delete this line

	return
}
