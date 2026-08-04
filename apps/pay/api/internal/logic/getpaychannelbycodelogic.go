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

type GetPayChannelByCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPayChannelByCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPayChannelByCodeLogic {
	return &GetPayChannelByCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPayChannelByCodeLogic) GetPayChannelByCode(req *types.PayChannelCodeReq) (resp *types.Result, err error) {
	r, err := l.svcCtx.PayRpc.QueryPayChannelByCode(l.ctx, &payclient.QueryPayChannelByCodeRequest{
		ChannelCode: req.ChannelCode,
	})
	if err != nil {
		return nil, err
	}
	return &types.Result{Code: 200, Msg: "OK", Data: r}, nil
}