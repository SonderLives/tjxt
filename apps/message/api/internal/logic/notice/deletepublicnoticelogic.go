// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notice

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"
	messageclient "tjxt/apps/message/rpc/message"

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
	if _, err := l.svcCtx.MessageRpc.DeletePublicNotice(l.ctx, &messageclient.IdReq{Id: req.Id}); err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
