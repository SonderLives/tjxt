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

type NotifyPayFailedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNotifyPayFailedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyPayFailedLogic {
	return &NotifyPayFailedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NotifyPayFailedLogic) NotifyPayFailed(req *types.NotifyPayFailedReq) (resp *types.Result, err error) {
	if _, err := l.svcCtx.PayRpc.NotifyPayFailed(l.ctx, &payclient.NotifyPayFailedRequest{
		PayOrderNo: req.PayOrderNo,
		ResultCode: req.ResultCode,
		ResultMsg:  req.ResultMsg,
	}); err != nil {
		return nil, err
	}
	return &types.Result{Code: 200, Msg: "OK"}, nil
}