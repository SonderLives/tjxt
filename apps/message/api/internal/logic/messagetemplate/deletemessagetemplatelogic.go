// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package messagetemplate

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMessageTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMessageTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMessageTemplateLogic {
	return &DeleteMessageTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMessageTemplateLogic) DeleteMessageTemplate(req *types.IdPathReq) (resp *types.OkVO, err error) {
	// todo: add your logic here and delete this line

	return
}
