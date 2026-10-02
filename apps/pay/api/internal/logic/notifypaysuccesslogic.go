// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"tjxt/apps/pay/api/internal/svc"
	"tjxt/apps/pay/api/internal/types"
	payclient "tjxt/apps/pay/rpc/pay"
	"tjxt/pkg/xerr"

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

// verifySign 校验支付网关回调签名。
// 规则：sign = hex(hmac_sha256(PayNotifySecret, "payOrderNo=<PayOrderNo>")).
func verifySign(secret string, payOrderNo int64, sign string) bool {
	if secret == "" || sign == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = fmt.Fprintf(mac, "payOrderNo=%d", payOrderNo)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(sign))
}

func (l *NotifyPaySuccessLogic) NotifyPaySuccess(req *types.NotifyPaySuccessReq) (resp *types.Result, err error) {
	if !verifySign(l.svcCtx.Config.PayNotifySecret, req.PayOrderNo, req.Sign) {
		return nil, xerr.Unauthorized("支付回调签名校验失败")
	}
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
