// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package noticetask

import (
	"context"

	backendlogic "tjxt/api/message/backend/logic"
	"tjxt/api/message/backend/pb"
	"tjxt/api/message/internal/svc"
	"tjxt/api/message/internal/types"

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
	if _, err := backendlogic.NewDeleteNoticeTaskLogic(l.ctx, l.svcCtx.Backend).DeleteNoticeTask(&pb.IdReq{Id: req.Id}); err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
