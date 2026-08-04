// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/pay/api/internal/svc"
	"tjxt/apps/pay/api/internal/types"
	payclient "tjxt/apps/pay/rpc/pay"
	"tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyPayOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyPayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyPayOrderLogic {
	return &ApplyPayOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyPayOrderLogic) ApplyPayOrder(req *types.ApplyPayOrderReq) (resp *types.Result, err error) {
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.PayRpc.ApplyPayOrder(l.ctx, &payclient.ApplyPayOrderRequest{
		BizUserId:      userId,
		BizOrderNo:     req.BizOrderNo,
		Amount:         req.Amount,
		PayChannelCode: req.PayChannelCode,
		PayType:        int32(req.PayType),
		NotifyUrl:      req.NotifyUrl,
		ExpandJson:     req.ExpandJson,
		PayOverSeconds: req.PayOverSeconds,
	})
	if err != nil {
		return nil, err
	}
	return &types.Result{Code: 200, Msg: "OK", Data: rpcResp.QrCodeUrl}, nil
}