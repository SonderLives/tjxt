package paylogic

import "time"

// 支付渠道状态
const (
	PayChannelStatusEnabled  = 1 // 使用中
	PayChannelStatusDisabled = 2 // 停用
)

// 支付单状态（与表注释一致）
const (
	PayOrderStatusPending = 0 // 待提交
	PayOrderStatusPaying  = 1 // 待支付
	PayOrderStatusClosed  = 2 // 支付超时或取消
	PayOrderStatusSuccess = 3 // 支付成功
)

// 退款单状态
const (
	RefundStatusInit       = 0 // 未提交
	RefundStatusProcessing = 1 // 退款中
	RefundStatusFailed     = 2 // 退款失败
	RefundStatusSuccess    = 3 // 退款成功
)

// 回调状态（PayOrder.NotifyStatus）
const (
	NotifyStatusPending = 0 // 待回调
	NotifyStatusOK      = 1 // 回调成功
	NotifyStatusFail    = 2 // 回调失败
)

// 退款通知状态（RefundOrder.NotifyStatus）
const (
	RefundNotifyStatusPending    = 0 // 待通知
	RefundNotifyStatusSuccess    = 1 // 通知成功
	RefundNotifyStatusProcessing = 2 // 通知中
	RefundNotifyStatusFailed     = 3 // 通知失败
)

// 渠道类型（与表注释一致，1=h5 2=小程序 3=公众号 4=扫码）
const (
	PayTypeH5     = 1
	PayTypeMini   = 2
	PayTypeMp     = 3
	PayTypeNative = 4
)

// PayOrderTimeout 支付单默认有效期（二维码超时未支付即关闭）。
// 与 trade 服务 OrderPayTimeout 同口径：均为"支付有效期"，此处支付单默认 30 分钟，
// 调用方可通过 ApplyPayOrderRequest.PayOverSeconds 覆盖。
const PayOrderTimeout = 30 * time.Minute
