package logic

import (
	"context"
	"database/sql"
	"time"

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
//   - 实际项目中应调用第三方支付渠道（微信/支付宝），此处为 demo 生成 mock qr url
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
			String: mockQrCodeUrl(poNo, in.Amount),
			Valid:  true,
		},
	}
	if _, err := l.svcCtx.PayOrderModel.Insert(l.ctx, po); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "创建支付单失败")
	}
	return &pb.ApplyPayOrderResponse{QrCodeUrl: po.QrCodeUrl.String}, nil
}

// mockQrCodeUrl 在真实项目里应该是调用微信/支付宝下单接口拿到的 code_url/prepay_id，
// demo 中生成一个本地可识别的占位 url。
func mockQrCodeUrl(payOrderNo, amount int64) string {
	return "tjxt://mock-pay?order_no=" + itoa(payOrderNo) + "&amount=" + itoa(amount)
}

func itoa(v int64) string {
	// 避免引入 strconv 的小工具，特殊情况下仍走 fmt
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	buf := [20]byte{}
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}