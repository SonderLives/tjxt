package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// 订单状态常量
const (
	OrderStatusPendingPay int64 = 1 // 待支付
	OrderStatusPaid       int64 = 2 // 已支付
	OrderStatusClosed     int64 = 3 // 已关闭
	OrderStatusFinished   int64 = 4 // 已完成
	OrderStatusEnrolled   int64 = 5 // 已报名（免费课）
	OrderStatusRefunding  int64 = 6 // 已申请退款
)

// OrderStatusDesc 订单状态描述
var OrderStatusDesc = map[int64]string{
	OrderStatusPendingPay: "待支付",
	OrderStatusPaid:       "已支付",
	OrderStatusClosed:     "已关闭",
	OrderStatusFinished:   "已完成",
	OrderStatusEnrolled:   "已报名",
	OrderStatusRefunding:  "已申请退款",
}

// Order 订单（order 表）
type Order struct {
	Id             int64
	UserId         int64
	PayOrderNo     sql.NullInt64
	Status         int64
	Message        string
	TotalAmount    int64
	RealAmount     int64
	DiscountAmount int64
	PayChannel     string
	CouponIds      sql.NullString
	CreateTime     time.Time
	PayTime        sql.NullTime
	CloseTime      sql.NullTime
	FinishTime     sql.NullTime
	RefundTime     sql.NullTime
	UpdateTime     time.Time
	Creater        int64
	Updater        int64
	Deleted        int64
}

// OrderModel 订单数据访问
type OrderModel struct {
	conn  sqlx.SqlConn
	table string
}

// execer 兼容 sqlx.SqlConn 与 sqlx.Session 的写执行器。
type execer interface {
	ExecCtx(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// NewOrderModel 创建订单数据访问对象
func NewOrderModel(conn sqlx.SqlConn) *OrderModel {
	return &OrderModel{conn: conn, table: "`order`"}
}

// Insert 创建订单。
func (m *OrderModel) Insert(ctx context.Context, o *Order) error {
	return m.insert(ctx, m.conn, o)
}

// InsertTx 在事务中创建订单。
func (m *OrderModel) InsertTx(ctx context.Context, session sqlx.Session, o *Order) error {
	return m.insert(ctx, session, o)
}

func (m *OrderModel) insert(ctx context.Context, e execer, o *Order) error {
	_, err := e.ExecCtx(ctx, fmt.Sprintf(
		`INSERT INTO %s (id, user_id, pay_order_no, status, message, total_amount, real_amount,
		 discount_amount, pay_channel, coupon_ids, create_time, pay_time, close_time, finish_time,
		 refund_time, update_time, creater, updater, deleted)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), ?, ?, ?, ?, NOW(), ?, ?, 0)`, m.table),
		o.Id, o.UserId, nullInt64(o.PayOrderNo), o.Status, o.Message, o.TotalAmount, o.RealAmount,
		o.DiscountAmount, o.PayChannel, nullString(o.CouponIds), nullTime(o.PayTime),
		nullTime(o.CloseTime), nullTime(o.FinishTime), nullTime(o.RefundTime), o.Creater, o.Updater)
	return err
}

// FindById 按主键查询订单。
func (m *OrderModel) FindById(ctx context.Context, id int64) (*Order, error) {
	var o Order
	err := m.conn.QueryRowCtx(ctx, &o, fmt.Sprintf(
		`SELECT id, user_id, pay_order_no, status, message, total_amount, real_amount, discount_amount,
		 pay_channel, coupon_ids, create_time, pay_time, close_time, finish_time, refund_time,
		 update_time, creater, updater, deleted FROM %s WHERE id = ? AND deleted = 0`, m.table), id)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

// UpdateStatus 更新订单状态及状态备注。
func (m *OrderModel) UpdateStatus(ctx context.Context, id, status int64, message string) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET status = ?, message = ?, update_time = NOW() WHERE id = ?`, m.table),
		status, message, id)
	return err
}

// UpdateAmount 更新订单金额（下单确认时）。
func (m *OrderModel) UpdateAmount(ctx context.Context, id, realAmount, discountAmount int64) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET real_amount = ?, discount_amount = ?, update_time = NOW() WHERE id = ?`, m.table),
		realAmount, discountAmount, id)
	return err
}

// MarkPaid 标记订单支付成功。
func (m *OrderModel) MarkPaid(ctx context.Context, id, payOrderNo int64, channel string, finishTime time.Time) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET status = ?, message = ?, pay_order_no = ?, pay_channel = ?, pay_time = NOW(),
		 finish_time = ?, update_time = NOW() WHERE id = ?`, m.table),
		OrderStatusPaid, "用户支付成功", payOrderNo, channel, finishTime, id)
	return err
}

// MarkRefunding 标记订单进入退款申请状态。
func (m *OrderModel) MarkRefunding(ctx context.Context, id int64) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET status = ?, message = ?, refund_time = NOW(), update_time = NOW() WHERE id = ?`, m.table),
		OrderStatusRefunding, id)
	return err
}

// CancelPending 取消待支付订单。
func (m *OrderModel) CancelPending(ctx context.Context, id, userId int64) (int64, error) {
	result, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET status = ?, message = ?, close_time = NOW(), update_time = NOW()
		 WHERE id = ? AND user_id = ? AND status = ?`, m.table),
		OrderStatusClosed, "用户取消订单", id, userId, OrderStatusPendingPay)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// SoftDelete 逻辑删除订单。
func (m *OrderModel) SoftDelete(ctx context.Context, id, userId int64) (int64, error) {
	result, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET deleted = 1, update_time = NOW() WHERE id = ? AND user_id = ? AND status <> ?`, m.table),
		id, userId, OrderStatusPendingPay)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// ListByUser 分页查询用户订单。
func (m *OrderModel) ListByUser(ctx context.Context, userId int64, status int64, offset, limit int64, asc bool) ([]Order, int64, error) {
	cond := "user_id = ? AND deleted = 0"
	args := []any{userId}
	if status > 0 {
		cond += " AND status = ?"
		args = append(args, status)
	}

	var total int64
	if err := m.conn.QueryRowCtx(ctx, &total, fmt.Sprintf(
		`SELECT COUNT(1) FROM %s WHERE %s`, m.table, cond), args...); err != nil {
		return nil, 0, err
	}

	order := "DESC"
	if asc {
		order = "ASC"
	}
	args = append(args, limit, offset)
	var list []Order
	err := m.conn.QueryRowsCtx(ctx, &list, fmt.Sprintf(
		`SELECT id, user_id, pay_order_no, status, message, total_amount, real_amount, discount_amount,
		 pay_channel, coupon_ids, create_time, pay_time, close_time, finish_time, refund_time,
		 update_time, creater, updater, deleted FROM %s WHERE %s ORDER BY create_time %s LIMIT ? OFFSET ?`,
		m.table, cond, order), args...)
	return list, total, err
}

// ListByIds 批量查询订单。
func (m *OrderModel) ListByIds(ctx context.Context, ids []int64) ([]Order, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	placeholders := strings.TrimRight(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	var list []Order
	err := m.conn.QueryRowsCtx(ctx, &list, fmt.Sprintf(
		`SELECT id, user_id, pay_order_no, status, message, total_amount, real_amount, discount_amount,
		 pay_channel, coupon_ids, create_time, pay_time, close_time, finish_time, refund_time,
		 update_time, creater, updater, deleted FROM %s WHERE id IN (%s) AND deleted = 0`,
		m.table, placeholders), args...)
	return list, err
}
