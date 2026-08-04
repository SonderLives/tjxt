package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// 订单明细状态常量
const (
	OrderDetailStatusPendingPay int64 = 1 // 待支付
	OrderDetailStatusPaid       int64 = 2 // 已支付
	OrderDetailStatusClosed     int64 = 3 // 已关闭
	OrderDetailStatusFinished   int64 = 4 // 已完成
	OrderDetailStatusEnrolled   int64 = 5 // 已报名（免费课）
	OrderDetailStatusRefunding  int64 = 6 // 已申请退款
)

// OrderDetailStatusDesc 订单明细状态描述
var OrderDetailStatusDesc = map[int64]string{
	OrderDetailStatusPendingPay: "待支付",
	OrderDetailStatusPaid:       "已支付",
	OrderDetailStatusClosed:     "已关闭",
	OrderDetailStatusFinished:   "已完成",
	OrderDetailStatusEnrolled:   "已报名",
	OrderDetailStatusRefunding:  "已申请退款",
}

// 退款状态常量
const (
	RefundStatusPending   int64 = 1 // 待审批
	RefundStatusCancelled int64 = 2 // 取消退款
	RefundStatusApproved  int64 = 3 // 同意退款
	RefundStatusRejected  int64 = 4 // 拒绝退款
	RefundStatusSuccess   int64 = 5 // 退款成功
	RefundStatusFailed    int64 = 6 // 退款失败
)

// RefundStatusDesc 退款状态描述
var RefundStatusDesc = map[int64]string{
	RefundStatusPending:   "待审批",
	RefundStatusCancelled: "取消退款",
	RefundStatusApproved:  "同意退款",
	RefundStatusRejected:  "拒绝退款",
	RefundStatusSuccess:   "退款成功",
	RefundStatusFailed:    "退款失败",
}

// OrderDetail 订单明细（order_detail 表）
type OrderDetail struct {
	Id               int64
	OrderId          int64
	UserId           int64
	CourseId         int64
	Price            int64
	Name             string
	CoverUrl         string
	ValidDuration    sql.NullInt64
	CourseExpireTime sql.NullTime
	DiscountAmount   int64
	RealPayAmount    int64
	Status           int64
	RefundStatus     sql.NullInt64
	PayChannel       string
	CreateTime       time.Time
	UpdateTime       time.Time
	Creater          int64
	Updater          int64
}

// OrderDetailModel 订单明细数据访问
type OrderDetailModel struct {
	conn  sqlx.SqlConn
	table string
}

// NewOrderDetailModel 创建订单明细数据访问对象
func NewOrderDetailModel(conn sqlx.SqlConn) *OrderDetailModel {
	return &OrderDetailModel{conn: conn, table: "order_detail"}
}

const orderDetailColumns = `id, order_id, user_id, course_id, price, name, cover_url,
	 valid_duration, course_expire_time, discount_amount, real_pay_amount, status,
	 refund_status, pay_channel, create_time, update_time, creater, updater`

// Insert 写入订单明细。
func (m *OrderDetailModel) Insert(ctx context.Context, d *OrderDetail) error {
	return m.insert(ctx, m.conn, d)
}

// InsertTx 在事务中写入订单明细。
func (m *OrderDetailModel) InsertTx(ctx context.Context, session sqlx.Session, d *OrderDetail) error {
	return m.insert(ctx, session, d)
}

func (m *OrderDetailModel) insert(ctx context.Context, e execer, d *OrderDetail) error {
	_, err := e.ExecCtx(ctx, fmt.Sprintf(
		`INSERT INTO %s (id, order_id, user_id, course_id, price, name, cover_url,
		 valid_duration, course_expire_time, discount_amount, real_pay_amount, status,
		 refund_status, pay_channel, create_time, update_time, creater, updater)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), ?, ?)`, m.table),
		d.Id, d.OrderId, d.UserId, d.CourseId, d.Price, d.Name, d.CoverUrl,
		nullInt64(d.ValidDuration), nullTime(d.CourseExpireTime), d.DiscountAmount, d.RealPayAmount,
		d.Status, nullInt64(d.RefundStatus), d.PayChannel, d.Creater, d.Updater)
	return err
}

// InsertBatch 批量写入订单明细（单条 SQL，用于下单事务）。
func (m *OrderDetailModel) InsertBatch(ctx context.Context, details []*OrderDetail) error {
	return m.insertBatch(ctx, m.conn, details)
}

// InsertBatchTx 在事务中批量写入订单明细。
func (m *OrderDetailModel) InsertBatchTx(ctx context.Context, session sqlx.Session, details []*OrderDetail) error {
	return m.insertBatch(ctx, session, details)
}

func (m *OrderDetailModel) insertBatch(ctx context.Context, e execer, details []*OrderDetail) error {
	if len(details) == 0 {
		return nil
	}
	var sb strings.Builder
	args := make([]any, 0, len(details)*16)
	sb.WriteString(fmt.Sprintf(
		`INSERT INTO %s (id, order_id, user_id, course_id, price, name, cover_url,
		 valid_duration, course_expire_time, discount_amount, real_pay_amount, status,
		 refund_status, pay_channel, create_time, update_time, creater, updater) VALUES `, m.table))
	for i, d := range details {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString("(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), ?, ?)")
		args = append(args, d.Id, d.OrderId, d.UserId, d.CourseId, d.Price, d.Name, d.CoverUrl,
			nullInt64(d.ValidDuration), nullTime(d.CourseExpireTime), d.DiscountAmount, d.RealPayAmount,
			d.Status, nullInt64(d.RefundStatus), d.PayChannel, d.Creater, d.Updater)
	}
	_, err := e.ExecCtx(ctx, sb.String(), args...)
	return err
}

// FindById 按主键查询订单明细。
func (m *OrderDetailModel) FindById(ctx context.Context, id int64) (*OrderDetail, error) {
	var d OrderDetail
	err := m.conn.QueryRowCtx(ctx, &d, fmt.Sprintf(
		`SELECT %s FROM %s WHERE id = ?`, orderDetailColumns, m.table), id)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// ListByOrderId 查询订单下所有明细。
func (m *OrderDetailModel) ListByOrderId(ctx context.Context, orderId int64) ([]OrderDetail, error) {
	var list []OrderDetail
	err := m.conn.QueryRowsCtx(ctx, &list, fmt.Sprintf(
		`SELECT %s FROM %s WHERE order_id = ?`, orderDetailColumns, m.table), orderId)
	return list, err
}

// ListByOrderIds 批量查询多个订单的明细。
func (m *OrderDetailModel) ListByOrderIds(ctx context.Context, orderIds []int64) ([]OrderDetail, error) {
	if len(orderIds) == 0 {
		return nil, nil
	}
	placeholders := strings.TrimRight(strings.Repeat("?,", len(orderIds)), ",")
	args := make([]any, 0, len(orderIds))
	for _, id := range orderIds {
		args = append(args, id)
	}
	var list []OrderDetail
	err := m.conn.QueryRowsCtx(ctx, &list, fmt.Sprintf(
		`SELECT %s FROM %s WHERE order_id IN (%s)`, orderDetailColumns, m.table, placeholders), args...)
	return list, err
}

// FindPaidByUserCourse 查询用户某课程的已支付明细。
func (m *OrderDetailModel) FindPaidByUserCourse(ctx context.Context, userId, courseId int64) (*OrderDetail, error) {
	var d OrderDetail
	err := m.conn.QueryRowCtx(ctx, &d, fmt.Sprintf(
		`SELECT %s FROM %s WHERE user_id = ? AND course_id = ? AND status IN (?, ?, ?) LIMIT 1`,
		orderDetailColumns, m.table),
		userId, courseId, OrderDetailStatusPaid, OrderDetailStatusFinished, OrderDetailStatusEnrolled)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

// MarkPaid 标记明细已支付。
func (m *OrderDetailModel) MarkPaid(ctx context.Context, id, status int64, channel string, expireTime *time.Time) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET status = ?, pay_channel = ?, course_expire_time = ?, update_time = NOW() WHERE id = ?`,
		m.table), status, channel, expireTime, id)
	return err
}

// MarkRefundStatus 更新明细退款状态。
func (m *OrderDetailModel) MarkRefundStatus(ctx context.Context, id, refundStatus int64) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET refund_status = ?, update_time = NOW() WHERE id = ?`, m.table), refundStatus, id)
	return err
}

// CountEnrolledByUser 统计学员已报名课程数（含免费课）。
func (m *OrderDetailModel) CountEnrolledByUser(ctx context.Context, userId int64) (int64, error) {
	var count int64
	err := m.conn.QueryRowCtx(ctx, &count, fmt.Sprintf(
		`SELECT COUNT(1) FROM %s WHERE user_id = ? AND status IN (?, ?, ?)`, m.table),
		userId, OrderDetailStatusPaid, OrderDetailStatusFinished, OrderDetailStatusEnrolled)
	return count, err
}

// CountEnrolledByUsers 批量统计学员已报名课程数，返回 userId -> 数量。
func (m *OrderDetailModel) CountEnrolledByUsers(ctx context.Context, userIds []int64) (map[int64]int64, error) {
	result := make(map[int64]int64, len(userIds))
	if len(userIds) == 0 {
		return result, nil
	}
	placeholders := strings.TrimRight(strings.Repeat("?,", len(userIds)), ",")
	args := make([]any, 0, len(userIds))
	for _, id := range userIds {
		args = append(args, id)
	}
	type row struct {
		UserId int64 `db:"user_id"`
		Count  int64 `db:"cnt"`
	}
	var rows []row
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT user_id, COUNT(1) AS cnt FROM %s WHERE user_id IN (%s) AND status IN (?, ?, ?)
		 GROUP BY user_id`, m.table, placeholders),
		append(args, OrderDetailStatusPaid, OrderDetailStatusFinished, OrderDetailStatusEnrolled)...)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		result[r.UserId] = r.Count
	}
	return result, nil
}

// CountEnrolledByCourses 批量统计课程报名人数，返回 courseId -> 数量。
func (m *OrderDetailModel) CountEnrolledByCourses(ctx context.Context, courseIds []int64) (map[int64]int64, error) {
	result := make(map[int64]int64, len(courseIds))
	if len(courseIds) == 0 {
		return result, nil
	}
	placeholders := strings.TrimRight(strings.Repeat("?,", len(courseIds)), ",")
	args := make([]any, 0, len(courseIds))
	for _, id := range courseIds {
		args = append(args, id)
	}
	type row struct {
		CourseId int64 `db:"course_id"`
		Count    int64 `db:"cnt"`
	}
	var rows []row
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT course_id, COUNT(1) AS cnt FROM %s WHERE course_id IN (%s) AND status IN (?, ?, ?)
		 GROUP BY course_id`, m.table, placeholders),
		append(args, OrderDetailStatusPaid, OrderDetailStatusFinished, OrderDetailStatusEnrolled)...)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		result[r.CourseId] = r.Count
	}
	return result, nil
}

// PurchaseInfo 课程购买统计
type PurchaseInfo struct {
	EnrollNum     int64 `db:"enrollNum"`
	RealPayAmount int64 `db:"realPayAmount"`
	RefundNum     int64 `db:"refundNum"`
}

// SumPurchaseInfo 统计课程购买信息（报名数、实付总额、退款数）。
func (m *OrderDetailModel) SumPurchaseInfo(ctx context.Context, courseId int64) (*PurchaseInfo, error) {
	var info PurchaseInfo
	err := m.conn.QueryRowCtx(ctx, &info, fmt.Sprintf(
		`SELECT COUNT(1) AS enrollNum,
		        COALESCE(SUM(CASE WHEN status IN (?, ?, ?) THEN real_pay_amount ELSE 0 END), 0) AS realPayAmount,
		        COUNT(CASE WHEN status IN (?, ?, ?) AND refund_status = ? THEN 1 END) AS refundNum
		 FROM %s WHERE course_id = ?`, m.table),
		OrderDetailStatusPaid, OrderDetailStatusFinished, OrderDetailStatusEnrolled,
		OrderDetailStatusPaid, OrderDetailStatusFinished, OrderDetailStatusEnrolled,
		RefundStatusSuccess, courseId)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// ListPage 管理端订单明细分页查询。
func (m *OrderDetailModel) ListPage(ctx context.Context, cond *DetailPageCond) ([]OrderDetail, int64, error) {
	where := "1 = 1"
	args := make([]any, 0)
	if cond.Id > 0 {
		where += " AND id = ?"
		args = append(args, cond.Id)
	}
	if cond.UserId > 0 {
		where += " AND user_id = ?"
		args = append(args, cond.UserId)
	}
	if cond.Status > 0 {
		where += " AND status = ?"
		args = append(args, cond.Status)
	}
	if cond.RefundStatus > 0 {
		where += " AND refund_status = ?"
		args = append(args, cond.RefundStatus)
	}
	if cond.PayChannel != "" {
		where += " AND pay_channel = ?"
		args = append(args, cond.PayChannel)
	}
	if !cond.StartTime.IsZero() {
		where += " AND create_time >= ?"
		args = append(args, cond.StartTime)
	}
	if !cond.EndTime.IsZero() {
		where += " AND create_time <= ?"
		args = append(args, cond.EndTime)
	}

	var total int64
	if err := m.conn.QueryRowCtx(ctx, &total, fmt.Sprintf(
		`SELECT COUNT(1) FROM %s WHERE %s`, m.table, where), args...); err != nil {
		return nil, 0, err
	}

	order := "DESC"
	if cond.IsAsc {
		order = "ASC"
	}
	args = append(args, cond.Limit, cond.Offset)
	var list []OrderDetail
	err := m.conn.QueryRowsCtx(ctx, &list, fmt.Sprintf(
		`SELECT %s FROM %s WHERE %s ORDER BY create_time %s LIMIT ? OFFSET ?`,
		orderDetailColumns, m.table, where, order), args...)
	return list, total, err
}

// DetailPageCond 订单明细分页条件
type DetailPageCond struct {
	Id           int64
	UserId       int64
	Status       int64
	RefundStatus int64
	PayChannel   string
	StartTime    time.Time
	EndTime      time.Time
	Offset       int64
	Limit        int64
	IsAsc        bool
}
