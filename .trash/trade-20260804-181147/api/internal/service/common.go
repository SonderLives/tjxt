package service

import (
	"database/sql"
	"time"

	"tjxt/apps/trade/api/internal/model"
	"tjxt/apps/trade/api/internal/types"
)

// formatTime 格式化时间，零值返回空串。
func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

// formatNullTime 格式化可空时间。
func formatNullTime(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(time.RFC3339)
}

// formatNullInt 格式化可空整型，无效返回 0。
func formatNullInt(n sql.NullInt64) int64 {
	if !n.Valid {
		return 0
	}
	return n.Int64
}

// statusDesc 取订单状态描述。
func orderStatusDesc(status int64) string {
	if desc, ok := model.OrderStatusDesc[status]; ok {
		return desc
	}
	return "未知状态"
}

// detailStatusDesc 取订单明细状态描述。
func detailStatusDesc(status int64) string {
	if desc, ok := model.OrderDetailStatusDesc[status]; ok {
		return desc
	}
	return "未知状态"
}

// refundStatusDesc 取退款状态描述，0/NULL 视为未退款。
func refundStatusDesc(status int64) string {
	if status == 0 {
		return "未退款"
	}
	if desc, ok := model.RefundStatusDesc[status]; ok {
		return desc
	}
	return "未知状态"
}

// canRefund 判断订单明细是否可发起退款。
func canRefund(status, refundStatus int64) bool {
	if status != model.OrderDetailStatusPaid {
		return false
	}
	switch refundStatus {
	case 0, // 从未退款
		model.RefundStatusCancelled,
		model.RefundStatusRejected,
		model.RefundStatusFailed:
		return true
	default: // 待审批、同意退款、退款成功等不允许重复申请
		return false
	}
}

// detailToVO 订单明细 -> 用户端 VO。
func detailToVO(d *model.OrderDetail) types.OrderDetailVO {
	return types.OrderDetailVO{
		Id:            d.Id,
		OrderId:       d.OrderId,
		CourseId:      d.CourseId,
		Name:          d.Name,
		CoverUrl:      d.CoverUrl,
		Price:         d.Price,
		RealPayAmount: d.RealPayAmount,
		RefundStatus:  formatNullInt(d.RefundStatus),
		CouponDesc:    "",
		CanRefund:     canRefund(d.Status, formatNullInt(d.RefundStatus)),
	}
}

// orderProgressNodes 构造订单进度节点。
func orderProgressNodes(o *model.Order) []types.OrderProgressNodeVO {
	nodes := make([]types.OrderProgressNodeVO, 0, 4)
	add := func(name string, t time.Time) {
		if t.IsZero() {
			return
		}
		nodes = append(nodes, types.OrderProgressNodeVO{Name: name, Time: t.Format(time.RFC3339)})
	}
	add("创建订单", o.CreateTime)
	add("支付成功", nullTimeVal(o.PayTime))
	add("订单完成", nullTimeVal(o.FinishTime))
	if o.CloseTime.Valid {
		add("订单关闭", o.CloseTime.Time)
	}
	return nodes
}

func nullTimeVal(t sql.NullTime) time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}
