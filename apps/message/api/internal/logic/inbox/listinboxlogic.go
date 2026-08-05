// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package inbox

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListInboxLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListInboxLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListInboxLogic {
	return &ListInboxLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListInboxLogic) ListInbox(req *types.InboxListReq) (resp *types.InboxListVO, err error) {
	// todo: add your logic here and delete this line

	return
}
