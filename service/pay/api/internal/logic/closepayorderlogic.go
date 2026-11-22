// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/service/pay/api/internal/svc"
	"tjxt/service/pay/api/internal/types"
	payclient "tjxt/service/pay/rpc/client/pay"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClosePayOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewClosePayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClosePayOrderLogic {
	return &ClosePayOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ClosePayOrderLogic) ClosePayOrder(req *types.PayOrderNoReq) (resp *types.Result, err error) {
	if _, err := l.svcCtx.PayRpc.ClosePayOrder(l.ctx, &payclient.ClosePayOrderRequest{PayOrderNo: req.PayOrderNo}); err != nil {
		return nil, err
	}
	return &types.Result{Code: 200, Msg: "OK"}, nil
}
