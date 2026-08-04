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

type NotifyPaySuccessLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNotifyPaySuccessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyPaySuccessLogic {
	return &NotifyPaySuccessLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NotifyPaySuccessLogic) NotifyPaySuccess(req *types.NotifyPaySuccessReq) (resp *types.Result, err error) {
	if _, err := l.svcCtx.PayRpc.NotifyPaySuccess(l.ctx, &payclient.NotifyPaySuccessRequest{
		PayOrderNo: req.PayOrderNo,
		ResultCode: req.ResultCode,
		ResultMsg:  req.ResultMsg,
		QrCodeUrl:  req.QrCodeUrl,
	}); err != nil {
		return nil, err
	}
	return &types.Result{Code: 200, Msg: "OK"}, nil
}