// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package messagetemplate

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveMessageTemplateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveMessageTemplateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveMessageTemplateLogic {
	return &SaveMessageTemplateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveMessageTemplateLogic) SaveMessageTemplate(req *types.MessageTemplateSaveReq) (resp *types.IdVO, err error) {
	// todo: add your logic here and delete this line

	return
}
