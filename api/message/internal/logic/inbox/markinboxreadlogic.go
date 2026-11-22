// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package inbox

import (
	"context"

	backendlogic "tjxt/api/message/backend/logic"
	"tjxt/api/message/backend/pb"
	"tjxt/api/message/internal/svc"
	"tjxt/api/message/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MarkInboxReadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMarkInboxReadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MarkInboxReadLogic {
	return &MarkInboxReadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MarkInboxReadLogic) MarkInboxRead(req *types.IdPathReq) (resp *types.OkVO, err error) {
	if _, err := backendlogic.NewMarkInboxReadLogic(l.ctx, l.svcCtx.Backend).MarkInboxRead(&pb.IdReq{Id: req.Id}); err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
