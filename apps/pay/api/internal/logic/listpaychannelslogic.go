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

type ListPayChannelsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListPayChannelsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPayChannelsLogic {
	return &ListPayChannelsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPayChannelsLogic) ListPayChannels() (resp *types.Result, err error) {
	rpcResp, err := l.svcCtx.PayRpc.ListPayChannels(l.ctx, &payclient.ListPayChannelsRequest{})
	if err != nil {
		return nil, err
	}
	return &types.Result{Code: 200, Msg: "OK", Data: rpcResp.List}, nil
}