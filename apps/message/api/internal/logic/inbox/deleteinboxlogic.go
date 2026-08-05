// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package inbox

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteInboxLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteInboxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteInboxLogic {
	return &DeleteInboxLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteInboxLogic) DeleteInbox(req *types.IdPathReq) (resp *types.OkVO, err error) {
	// todo: add your logic here and delete this line

	return
}
