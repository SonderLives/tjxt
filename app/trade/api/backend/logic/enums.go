package tradelogic

import "time"

// 订单支付超时时间：下单后 15 分钟内未支付则关闭（与前端倒计时、
// 订单超时关闭任务的过期口径一致）。
const OrderPayTimeout = 15 * time.Minute

// ===================== 状态枚举 =====================

// 订单状态：1待支付 2已支付 3已关闭 4已完成 5已报名 6申请退款
const (
	OrderStatusPending   int64 = 1
	OrderStatusPaid      int64 = 2
	OrderStatusClosed    int64 = 3
	OrderStatusFinished  int64 = 4
	OrderStatusEnrolled  int64 = 5
	OrderStatusRefunding int64 = 6
)

// 订单明细状态：1待支付 2已支付 3已关闭 4已完成 5已报名
const (
	DetailStatusPending  int64 = 1
	DetailStatusPaid     int64 = 2
	DetailStatusClosed   int64 = 3
	DetailStatusFinished int64 = 4
	DetailStatusEnrolled int64 = 5
)

// 退款状态：1待审批 2取消退款 3同意退款 4拒绝退款 5退款成功 6退款失败
const (
	RefundStatusPending int64 = 1
	RefundStatusCancel  int64 = 2
	RefundStatusApprove int64 = 3
	RefundStatusReject  int64 = 4
	RefundStatusSuccess int64 = 5
	RefundStatusFailed  int64 = 6
)

// ===================== 订单状态描述 =====================

func orderStatusDesc(status int64) string {
	switch status {
	case OrderStatusPending:
		return "待支付"
	case OrderStatusPaid:
		return "已支付"
	case OrderStatusClosed:
		return "已关闭"
	case OrderStatusFinished:
		return "已完成"
	case OrderStatusEnrolled:
		return "已报名"
	case OrderStatusRefunding:
		return "申请退款"
	default:
		return "未知"
	}
}

func detailStatusDesc(status int64) string {
	switch status {
	case DetailStatusPending:
		return "待支付"
	case DetailStatusPaid:
		return "已支付"
	case DetailStatusClosed:
		return "已关闭"
	case DetailStatusFinished:
		return "已完成"
	case DetailStatusEnrolled:
		return "已报名"
	default:
		return "未知"
	}
}

func refundStatusDesc(status int64) string {
	switch status {
	case RefundStatusPending:
		return "待审批"
	case RefundStatusCancel:
		return "取消退款"
	case RefundStatusApprove:
		return "同意退款"
	case RefundStatusReject:
		return "拒绝退款"
	case RefundStatusSuccess:
		return "退款成功"
	case RefundStatusFailed:
		return "退款失败"
	default:
		return "未知"
	}
}
