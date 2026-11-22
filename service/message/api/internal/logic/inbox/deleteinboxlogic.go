// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package inbox

import (
	"context"

	backendlogic "tjxt/service/message/api/backend/logic"
	"tjxt/service/message/api/backend/pb"
	"tjxt/service/message/api/internal/svc"
	"tjxt/service/message/api/internal/types"

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
	if _, err := backendlogic.NewDeleteInboxLogic(l.ctx, l.svcCtx.Backend).DeleteInbox(&pb.IdReq{Id: req.Id}); err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
