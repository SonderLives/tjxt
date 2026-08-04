// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/pay/api/internal/svc"
	"tjxt/apps/pay/api/internal/types"
	payclient "tjxt/apps/pay/rpc/pay"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePayChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePayChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePayChannelLogic {
	return &UpdatePayChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePayChannelLogic) UpdatePayChannel(req *types.PayChannelUpdateReq) (resp *types.Result, err error) {
	if _, err := l.svcCtx.PayRpc.UpdatePayChannel(l.ctx, &payclient.PayChannelRequest{
		Id:              req.Id,
		Name:            req.Name,
		ChannelPriority: int32(req.ChannelPriority),
		ChannelIcon:     req.ChannelIcon,
	}); err != nil {
		return nil, err
	}
	return &types.Result{Code: 200, Msg: "OK"}, nil
}