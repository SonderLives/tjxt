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

type AddPayChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddPayChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPayChannelLogic {
	return &AddPayChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddPayChannelLogic) AddPayChannel(req *types.PayChannelAddReq) (resp *types.Result, err error) {
	r, err := l.svcCtx.PayRpc.AddPayChannel(l.ctx, &payclient.PayChannelRequest{
		Name:            req.Name,
		ChannelCode:     req.ChannelCode,
		ChannelPriority: int32(req.ChannelPriority),
		ChannelIcon:     req.ChannelIcon,
	})
	if err != nil {
		return nil, err
	}
	return &types.Result{Code: 200, Msg: "OK", Data: r.Id}, nil
}