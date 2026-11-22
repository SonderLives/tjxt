// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	backendlogic "tjxt/service/trade/api/backend/logic"
	"tjxt/service/trade/api/backend/pb"
	"tjxt/service/trade/api/internal/svc"
	"tjxt/service/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayChannelGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayChannelGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayChannelGetLogic {
	return &PayChannelGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayChannelGetLogic) PayChannelGet(req *types.PayChannelIdReq) (resp *types.PayChannelDTO, err error) {
	reply, err := backendlogic.NewPayChannelGetLogic(l.ctx, l.svcCtx.Backend).PayChannelGet(&pb.IdRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}

	return &types.PayChannelDTO{
		Id:              reply.Id,
		Name:            reply.Name,
		ChannelCode:     reply.ChannelCode,
		ChannelIcon:     reply.ChannelIcon,
		ChannelPriority: int64(reply.ChannelPriority),
		Status:          int64(reply.Status),
	}, nil
}
