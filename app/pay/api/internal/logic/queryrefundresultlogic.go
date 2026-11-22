// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/app/pay/api/internal/svc"
	"tjxt/app/pay/api/internal/types"
	payclient "tjxt/app/pay/rpc/client/pay"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryRefundResultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryRefundResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRefundResultLogic {
	return &QueryRefundResultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryRefundResultLogic) QueryRefundResult(req *types.QueryRefundResultReq) (resp *types.Result, err error) {
	rpcResp, err := l.svcCtx.PayRpc.QueryRefundResult(l.ctx, &payclient.QueryRefundResultRequest{BizRefundOrderNo: req.BizRefundOrderNo})
	if err != nil {
		return nil, err
	}
	return &types.Result{Code: 200, Msg: "OK", Data: rpcResp}, nil
}
