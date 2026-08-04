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

type ApplyRefundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyRefundLogic {
	return &ApplyRefundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyRefundLogic) ApplyRefund(req *types.ApplyRefundReq) (resp *types.Result, err error) {
	rpcResp, err := l.svcCtx.PayRpc.ApplyRefund(l.ctx, &payclient.ApplyRefundRequest{
		BizOrderNo:       req.BizOrderNo,
		BizRefundOrderNo: req.BizRefundOrderNo,
		RefundAmount:     req.RefundAmount,
	})
	if err != nil {
		return nil, err
	}
	return &types.Result{Code: 200, Msg: "OK", Data: rpcResp}, nil
}