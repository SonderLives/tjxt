// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package messagetemplate

import (
	"context"

	backendlogic "tjxt/api/message/backend/logic"
	"tjxt/api/message/backend/pb"
	"tjxt/api/message/internal/svc"
	"tjxt/api/message/internal/types"

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
	if _, err := backendlogic.NewDeleteMessageTemplateLogic(l.ctx, l.svcCtx.Backend).DeleteMessageTemplate(&pb.IdReq{Id: req.Id}); err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
