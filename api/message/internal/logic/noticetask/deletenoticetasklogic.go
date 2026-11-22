// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package noticetask

import (
	"context"

	"tjxt/api/message/internal/svc"
	"tjxt/api/message/internal/types"
	messageclient "tjxt/rpc/message/client/message"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteNoticeTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteNoticeTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNoticeTaskLogic {
	return &DeleteNoticeTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteNoticeTaskLogic) DeleteNoticeTask(req *types.IdPathReq) (resp *types.OkVO, err error) {
	if _, err := l.svcCtx.MessageRpc.DeleteNoticeTask(l.ctx, &messageclient.IdReq{Id: req.Id}); err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
