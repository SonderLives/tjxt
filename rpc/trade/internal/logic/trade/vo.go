package tradelogic

import (
	"database/sql"

	"tjxt/rpc/trade/internal/model"
	"tjxt/rpc/trade/pb"
	"tjxt/pkg/timex"
)

// ===================== VO 构建 =====================

// nullInt64Value 解包 sql.NullInt64：无效值返回 0
func nullInt64Value(v sql.NullInt64) int64 {
	if v.Valid {
		return v.Int64
	}
	return 0
}

// nullStringValue 解包 sql.NullString：无效值返回空串
func nullStringValue(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}

func toCartVO(c *model.Cart) *pb.CartVO {
	price := c.Price
	name := c.CourseName
	cover := c.CoverUrl
	expired := false
	return &pb.CartVO{
		Id:         c.Id,
		CourseId:   c.CourseId,
		CourseName: name,
		CoverUrl:   cover,
		Price:      price,
		NowPrice:   price,
		Expired:    expired,
	}
}

func toOrderDetailItemVO(d *model.OrderDetail) *pb.OrderDetailItemVO {
	canRefund := (d.Status == DetailStatusPaid || d.Status == DetailStatusFinished || d.Status == DetailStatusEnrolled) &&
		(!d.RefundStatus.Valid || d.RefundStatus.Int64 == 0 || d.RefundStatus.Int64 == RefundStatusPending)
	return &pb.OrderDetailItemVO{
		Id:            d.Id,
		CourseId:      d.CourseId,
		CourseName:    d.Name,
		CoverUrl:      d.CoverUrl,
		Price:         d.Price,
		RealPayAmount: d.RealPayAmount,
		Status:        int32(d.Status),
		RefundStatus:  int32(nullInt64Value(d.RefundStatus)),
		CanRefund:     canRefund,
	}
}

func toOrderVO(order *model.Order, details []*model.OrderDetail) *pb.OrderVO {
	vo := &pb.OrderVO{
		Id:             order.Id,
		CreateTime:     timex.Format(order.CreateTime),
		TotalAmount:    order.TotalAmount,
		RealAmount:     order.RealAmount,
		DiscountAmount: order.DiscountAmount,
		Status:         int32(order.Status),
		StatusDesc:     orderStatusDesc(order.Status),
		Message:        order.Message,
	}
	for _, d := range details {
		vo.Details = append(vo.Details, toOrderDetailItemVO(d))
	}
	vo.ProgressNodes = buildOrderProgressNodes(order)
	return vo
}

// buildOrderProgressNodes 构建订单进度节点（下单→支付→完成/关闭）
func buildOrderProgressNodes(order *model.Order) []*pb.OrderProgressNodeVO {
	nodes := []*pb.OrderProgressNodeVO{
		{Name: "提交订单", Status: 1, Time: timex.Format(order.CreateTime)},
	}
	if order.Status >= OrderStatusPaid {
		nodes = append(nodes, &pb.OrderProgressNodeVO{Name: "支付成功", Status: 1, Time: timex.FormatNullTime(order.PayTime)})
	} else if order.Status == OrderStatusClosed {
		nodes = append(nodes, &pb.OrderProgressNodeVO{Name: "订单关闭", Status: 3, Time: timex.FormatNullTime(order.CloseTime), Desc: order.Message})
	} else {
		nodes = append(nodes, &pb.OrderProgressNodeVO{Name: "支付成功", Status: 0, Time: ""})
	}

	switch order.Status {
	case OrderStatusFinished, OrderStatusEnrolled:
		nodes = append(nodes, &pb.OrderProgressNodeVO{Name: "已完成", Status: 1, Time: timex.FormatNullTime(order.FinishTime)})
	case OrderStatusRefunding:
		nodes = append(nodes, &pb.OrderProgressNodeVO{Name: "退款中", Status: 0, Time: timex.FormatNullTime(order.RefundTime)})
	case OrderStatusClosed:
		nodes = append(nodes, &pb.OrderProgressNodeVO{Name: "已完成", Status: 3, Time: ""})
	default:
		nodes = append(nodes, &pb.OrderProgressNodeVO{Name: "已完成", Status: 0, Time: ""})
	}
	return nodes
}

func toOrderDetailAdminVO(d *model.OrderDetail, order *model.Order, ra *model.RefundApply) *pb.OrderDetailAdminVO {
	vo := &pb.OrderDetailAdminVO{
		Id:             d.Id,
		OrderId:        d.OrderId,
		Name:           d.Name,
		Price:          d.Price,
		RealPayAmount:  d.RealPayAmount,
		DiscountAmount: d.DiscountAmount,
		Status:         int32(d.Status),
		PayChannel:     d.PayChannel,
		PayOrderNo:     nullInt64Value(order.PayOrderNo),
		StudyValidTime: timex.FormatNullTime(d.CourseExpireTime),
		CanRefund:      (d.Status == DetailStatusPaid || d.Status == DetailStatusFinished || d.Status == DetailStatusEnrolled) && (!d.RefundStatus.Valid || d.RefundStatus.Int64 == 0 || d.RefundStatus.Int64 == RefundStatusPending),
	}
	if ra != nil {
		vo.RefundApplyId = ra.Id
		vo.RefundOrderNo = nullInt64Value(ra.RefundOrderNo)
		vo.RefundStatus = int32(ra.Status)
		vo.RefundReason = ra.RefundReason
		vo.RefundMessage = ra.Message
		vo.RefundChannel = nullStringValue(ra.RefundChannel)
		vo.RefundFailedReason = nullStringValue(ra.FailedReason)
	}
	return vo
}

func toRefundApplyVO(ra *model.RefundApply, order *model.Order, detail *model.OrderDetail) *pb.RefundApplyVO {
	vo := &pb.RefundApplyVO{
		Id:             ra.Id,
		OrderId:        ra.OrderId,
		OrderDetailId:  ra.OrderDetailId,
		Price:          nullInt64Value(detailPrice(detail)),
		RefundAmount:   ra.RefundAmount,
		RefundStatus:   int32(ra.Status),
		RefundOrderNo:  nullInt64Value(ra.RefundOrderNo),
		PayOrderNo:     nullInt64Value(ra.PayOrderNo),
		PayChannel:     nullStringValue(ra.RefundChannel),
		RefundChannel:  nullStringValue(ra.RefundChannel),
		RefundReason:   ra.RefundReason,
		RefundMessage:  ra.Message,
		FailedReason:   nullStringValue(ra.FailedReason),
		ApproveOpinion: nullStringValue(ra.ApproveOpinion),
		ApproveType:    0,
		Remark:         nullStringValue(ra.Remark),
		CreateTime:     timex.Format(ra.CreateTime),
		OrderTime:      timex.FormatNullTime(orderCreateTime(order)),
		PaySuccessTime: timex.FormatNullTime(paySuccessTime(order)),
	}
	return vo
}

func detailPrice(d *model.OrderDetail) sql.NullInt64 {
	if d == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: d.RealPayAmount, Valid: true}
}

func orderCreateTime(order *model.Order) sql.NullTime {
	if order == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: order.CreateTime, Valid: true}
}

func paySuccessTime(order *model.Order) sql.NullTime {
	if order == nil {
		return sql.NullTime{}
	}
	return order.PayTime
}

func toRefundApplyPageVO(ra *model.RefundApply) *pb.RefundApplyPageVO {
	return &pb.RefundApplyPageVO{
		Id:            ra.Id,
		OrderId:       ra.OrderId,
		OrderDetailId: ra.OrderDetailId,
		Price:         ra.RefundAmount,
		RefundAmount:  ra.RefundAmount,
		Status:        int32(ra.Status),
		StatusDesc:    refundStatusDesc(ra.Status),
		RefundReason:  ra.RefundReason,
		CreateTime:    timex.Format(ra.CreateTime),
		ApproveTime:   timex.FormatNullTime(ra.ApproveTime),
	}
}
