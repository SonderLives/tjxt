package logic

import (
	"context"
	"database/sql"
	"time"

	"tjxt/apps/pay/rpc/internal/gateway"
	"tjxt/apps/pay/rpc/internal/model"
	"tjxt/apps/pay/rpc/internal/svc"
	"tjxt/apps/pay/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

// ApplyPayOrderLogic 申请支付单（创建 or 返回已有）
//
// 业务规则：
//   - 同一 biz_order_no 幂等：若已有 pay_order，直接返回原二维码
//   - 若原单已支付/已关闭，返回错误，让上游重新发起业务订单
//   - 默认 30 分钟支付超时
//   - 渠道交互经 gateway.PaymentGateway 抽象（demo 为 MockGateway，真实渠道实现该接口即可）
type ApplyPayOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyPayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyPayOrderLogic {
	return &ApplyPayOrderLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ApplyPayOrderLogic) ApplyPayOrder(in *pb.ApplyPayOrderRequest) (*pb.ApplyPayOrderResponse, error) {
	if in.BizOrderNo <= 0 || in.BizUserId <= 0 || in.Amount <= 0 {
		return nil, xerr.BadRequestf("biz_order_no/biz_user_id/amount 非法")
	}
	if in.PayChannelCode == "" {
		return nil, xerr.BadRequestf("支付渠道编码不能为空")
	}

	// 渠道校验
	channel, err := l.svcCtx.PayChannelModel.FindByCode(l.ctx, in.PayChannelCode)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("支付渠道不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询支付渠道失败")
	}
	if channel.Status != PayChannelStatusEnabled {
		return nil, xerr.Conflict("支付渠道已停用")
	}

	// 幂等：biz_order_no 已有单则直接复用
	existing, err := l.svcCtx.PayOrderModel.FindOneByBizOrderNo(l.ctx, in.BizOrderNo)
	if err == nil {
		switch existing.Status {
		case PayOrderStatusPaying:
			return &pb.ApplyPayOrderResponse{QrCodeUrl: formatNullString(existing.QrCodeUrl)}, nil
		case PayOrderStatusSuccess:
			return nil, xerr.Conflict("订单已支付，请勿重复支付")
		default:
			return nil, xerr.Conflict("订单已关闭，请重新下单")
		}
	}
	if !isNotFound(err) {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询已有支付单失败")
	}

	// 默认 30 分钟超时
	overSeconds := in.PayOverSeconds
	if overSeconds <= 0 {
		overSeconds = 30 * 60
	}
	payOverTime := time.Now().Add(time.Duration(overSeconds) * time.Second)
	payType := in.PayType
	if payType <= 0 {
		payType = PayTypeNative
	}

	poNo := nextID()
	// 渠道下单（mock 渠道返回本地可识别二维码；真实渠道在此调用微信/支付宝下单接口）
	gp, err := l.svcCtx.Gateway.CreatePayment(l.ctx, gateway.CreatePaymentRequest{
		PayOrderNo:  poNo,
		BizOrderNo:  in.BizOrderNo,
		BizUserID:   in.BizUserId,
		Amount:      in.Amount,
		ChannelCode: in.PayChannelCode,
	})
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "渠道下单失败")
	}
	po := &model.PayOrder{
		BizOrderNo:     in.BizOrderNo,
		PayOrderNo:     poNo,
		BizUserId:      in.BizUserId,
		PayChannelCode: in.PayChannelCode,
		Amount:         in.Amount,
		PayType:        int64(payType),
		Status:         PayOrderStatusPaying,
		ExpandJson:     in.ExpandJson,
		NotifyUrl:      in.NotifyUrl,
		NotifyTimes:    0,
		NotifyStatus:   NotifyStatusPending,
		PayOverTime:    payOverTime,
		QrCodeUrl: sql.NullString{
			String: gp.QrCodeUrl,
			Valid:  true,
		},
	}
	if _, err := l.svcCtx.PayOrderModel.Insert(l.ctx, po); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "创建支付单失败")
	}
	return &pb.ApplyPayOrderResponse{QrCodeUrl: po.QrCodeUrl.String}, nil
}
