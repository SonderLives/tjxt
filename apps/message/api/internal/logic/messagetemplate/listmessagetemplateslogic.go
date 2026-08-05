// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package messagetemplate

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMessageTemplatesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMessageTemplatesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMessageTemplatesLogic {
	return &ListMessageTemplatesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMessageTemplatesLogic) ListMessageTemplates(req *types.PageReq) (resp *types.MessageTemplateListVO, err error) {
	// todo: add your logic here and delete this line

	return
}
