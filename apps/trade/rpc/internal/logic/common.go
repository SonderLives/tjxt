package logic

import (
	"context"
	"database/sql"
	"time"

	courseclient "tjxt/apps/course/rpc/course"
	"tjxt/apps/trade/rpc/internal/model"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"tjxt/pkg/utils/idgen"
)

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
	RefundStatusPending  int64 = 1
	RefundStatusCancel   int64 = 2
	RefundStatusApprove  int64 = 3
	RefundStatusReject   int64 = 4
	RefundStatusSuccess  int64 = 5
	RefundStatusFailed   int64 = 6
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

// ===================== 工具函数 =====================

func nextID() int64 {
	return idgen.NextID()
}

func now() time.Time {
	return time.Now()
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func formatNullTime(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}
	return formatTime(t.Time)
}

func nullInt64Value(v sql.NullInt64) int64 {
	if v.Valid {
		return v.Int64
	}
	return 0
}

func nullStringValue(v sql.NullString) string {
	if v.Valid {
		return v.String
	}
	return ""
}

func calcPages(total, pageSize int64) int64 {
	if pageSize <= 0 || total <= 0 {
		return 0
	}
	pages := total / pageSize
	if total%pageSize != 0 {
		pages++
	}
	return pages
}

// ===================== 课程信息回填 =====================

// fetchCourseMap 批量查询课程简况，返回 course_id -> 课程信息
func fetchCourseMap(ctx context.Context, svcCtx *svc.ServiceContext, ids []int64) map[int64]*courseclient.CourseSimpleInfoItem {
	res := make(map[int64]*courseclient.CourseSimpleInfoItem)
	if len(ids) == 0 {
		return res
	}
	reply, err := svcCtx.CourseRpc.CourseSimpleInfoList(ctx, &courseclient.CourseSimpleInfoQueryRequest{Ids: ids})
	if err != nil {
		return res
	}
	for _, it := range reply.Items {
		res[it.Id] = it
	}
	return res
}

func courseName(m map[int64]*courseclient.CourseSimpleInfoItem, id int64) string {
	if it, ok := m[id]; ok {
		return it.Name
	}
	return ""
}

func courseCover(m map[int64]*courseclient.CourseSimpleInfoItem, id int64) string {
	if it, ok := m[id]; ok {
		return it.CoverUrl
	}
	return ""
}

func coursePrice(m map[int64]*courseclient.CourseSimpleInfoItem, id int64) int64 {
	if it, ok := m[id]; ok {
		return it.Price
	}
	return 0
}

// ===================== VO 构建 =====================

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
		CreateTime:     formatTime(order.CreateTime),
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
		{Name: "提交订单", Status: 1, Time: formatTime(order.CreateTime)},
	}
	if order.Status >= OrderStatusPaid {
		nodes = append(nodes, &pb.OrderProgressNodeVO{Name: "支付成功", Status: 1, Time: formatNullTime(order.PayTime)})
	} else if order.Status == OrderStatusClosed {
		nodes = append(nodes, &pb.OrderProgressNodeVO{Name: "订单关闭", Status: 3, Time: formatNullTime(order.CloseTime), Desc: order.Message})
	} else {
		nodes = append(nodes, &pb.OrderProgressNodeVO{Name: "支付成功", Status: 0, Time: ""})
	}

	switch order.Status {
	case OrderStatusFinished, OrderStatusEnrolled:
		nodes = append(nodes, &pb.OrderProgressNodeVO{Name: "已完成", Status: 1, Time: formatNullTime(order.FinishTime)})
	case OrderStatusRefunding:
		nodes = append(nodes, &pb.OrderProgressNodeVO{Name: "退款中", Status: 0, Time: formatNullTime(order.RefundTime)})
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
		StudyValidTime: formatNullTime(d.CourseExpireTime),
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
		Id:              ra.Id,
		OrderId:         ra.OrderId,
		OrderDetailId:   ra.OrderDetailId,
		Price:           nullInt64Value(detailPrice(detail)),
		RefundAmount:    ra.RefundAmount,
		RefundStatus:    int32(ra.Status),
		RefundOrderNo:   nullInt64Value(ra.RefundOrderNo),
		PayOrderNo:      nullInt64Value(ra.PayOrderNo),
		PayChannel:      nullStringValue(ra.RefundChannel),
		RefundChannel:   nullStringValue(ra.RefundChannel),
		RefundReason:    ra.RefundReason,
		RefundMessage:   ra.Message,
		FailedReason:    nullStringValue(ra.FailedReason),
		ApproveOpinion:  nullStringValue(ra.ApproveOpinion),
		ApproveType:     0,
		Remark:          nullStringValue(ra.Remark),
		CreateTime:      formatTime(ra.CreateTime),
		OrderTime:       formatNullTime(orderCreateTime(order)),
		PaySuccessTime:  formatNullTime(paySuccessTime(order)),
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
		CreateTime:    formatTime(ra.CreateTime),
		ApproveTime:   formatNullTime(ra.ApproveTime),
	}
}
