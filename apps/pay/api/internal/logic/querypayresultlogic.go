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

type QueryPayResultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryPayResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryPayResultLogic {
	return &QueryPayResultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryPayResultLogic) QueryPayResult(req *types.QueryPayResultReq) (resp *types.Result, err error) {
	rpcResp, err := l.svcCtx.PayRpc.QueryPayResult(l.ctx, &payclient.QueryPayResultRequest{BizOrderNo: req.BizOrderNo})
	if err != nil {
		return nil, err
	}
	return &types.Result{Code: 200, Msg: "OK", Data: rpcResp}, nil
}