// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"fmt"

	"tjxt/apps/pay/api/internal/svc"
	"tjxt/apps/pay/api/internal/types"
	payclient "tjxt/apps/pay/rpc/pay"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type NotifyRefundSuccessLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNotifyRefundSuccessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyRefundSuccessLogic {
	return &NotifyRefundSuccessLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NotifyRefundSuccessLogic) NotifyRefundSuccess(req *types.NotifyRefundSuccessReq) (resp *types.Result, err error) {
	if !verifySign(l.svcCtx.Config.PayNotifySecret, fmt.Sprintf("refundOrderNo=%d", req.RefundOrderNo), req.Sign) {
		return nil, xerr.Unauthorized("退款回调签名校验失败")
	}
	if _, err := l.svcCtx.PayRpc.NotifyRefundSuccess(l.ctx, &payclient.NotifyRefundSuccessRequest{
		RefundOrderNo: req.RefundOrderNo,
		ResultCode:    req.ResultCode,
		ResultMsg:     req.ResultMsg,
		RefundChannel: req.RefundChannel,
	}); err != nil {
		return nil, err
	}
	return &types.Result{Code: 200, Msg: "OK"}, nil
}
