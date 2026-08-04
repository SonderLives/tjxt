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

type UpdatePayChannelStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdatePayChannelStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePayChannelStatusLogic {
	return &UpdatePayChannelStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdatePayChannelStatusLogic) UpdatePayChannelStatus(req *types.PayChannelStatusReq) (resp *types.Result, err error) {
	if _, err := l.svcCtx.PayRpc.UpdatePayChannelStatus(l.ctx, &payclient.UpdatePayChannelStatusRequest{
		Id:     req.Id,
		Status: int32(req.Status),
	}); err != nil {
		return nil, err
	}
	return &types.Result{Code: 200, Msg: "OK"}, nil
}