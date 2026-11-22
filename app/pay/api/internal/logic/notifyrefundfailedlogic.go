// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"fmt"

	"tjxt/app/pay/api/internal/svc"
	"tjxt/app/pay/api/internal/types"
	payclient "tjxt/app/pay/rpc/client/pay"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyRefundFailedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNotifyRefundFailedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyRefundFailedLogic {
	return &NotifyRefundFailedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NotifyRefundFailedLogic) NotifyRefundFailed(req *types.NotifyRefundFailedReq) (resp *types.Result, err error) {
	if !verifySign(l.svcCtx.Config.PayNotifySecret, fmt.Sprintf("refundOrderNo=%d", req.RefundOrderNo), req.Sign) {
		return nil, errorx.Unauthorized("退款回调签名校验失败")
	}
	if _, err := l.svcCtx.PayRpc.NotifyRefundFailed(l.ctx, &payclient.NotifyRefundFailedRequest{
		RefundOrderNo: req.RefundOrderNo,
		ResultCode:    req.ResultCode,
		ResultMsg:     req.ResultMsg,
	}); err != nil {
		return nil, err
	}
	return &types.Result{Code: 200, Msg: "OK"}, nil
}
