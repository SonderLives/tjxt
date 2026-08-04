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

type QueryPayOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryPayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryPayOrderLogic {
	return &QueryPayOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryPayOrderLogic) QueryPayOrder(req *types.QueryPayResultReq) (resp *types.Result, err error) {
	rpcResp, err := l.svcCtx.PayRpc.QueryPayOrderByBizOrderNo(l.ctx, &payclient.QueryPayOrderRequest{BizOrderNo: req.BizOrderNo})
	if err != nil {
		return nil, err
	}
	return &types.Result{Code: 200, Msg: "OK", Data: rpcResp}, nil
}