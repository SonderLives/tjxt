// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notice

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePublicNoticeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePublicNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePublicNoticeLogic {
	return &DeletePublicNoticeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePublicNoticeLogic) DeletePublicNotice(req *types.IdPathReq) (resp *types.OkVO, err error) {
	// todo: add your logic here and delete this line

	return
}
