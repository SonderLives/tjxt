package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// RefundApply 退款申请（refund_apply 表）
type RefundApply struct {
	Id             int64
	OrderDetailId  int64
	OrderId        int64
	PayOrderNo     sql.NullInt64
	RefundOrderNo  sql.NullInt64
	UserId         int64
	RefundAmount   int64
	Status         int64
	RefundReason   string
	Message        string
	Approver       sql.NullInt64
	ApproveOpinion sql.NullString
	Remark         sql.NullString
	FailedReason   sql.NullString
	QuestionDesc   sql.NullString
	RefundChannel  sql.NullString
	CreateTime     time.Time
	ApproveTime    sql.NullTime
	FinishTime     sql.NullTime
	UpdateTime     time.Time
	Creater        int64
	Updater        int64
}

// RefundApplyModel 退款申请数据访问
type RefundApplyModel struct {
	conn  sqlx.SqlConn
	table string
}

// NewRefundApplyModel 创建退款申请数据访问对象
func NewRefundApplyModel(conn sqlx.SqlConn) *RefundApplyModel {
	return &RefundApplyModel{conn: conn, table: "refund_apply"}
}

const refundApplyColumns = `id, order_detail_id, order_id, pay_order_no, refund_order_no, user_id,
	 refund_amount, status, refund_reason, message, approver, approve_opinion, remark, failed_reason,
	 question_desc, refund_channel, create_time, approve_time, finish_time, update_time, creater, updater`

// Insert 创建退款申请。
func (m *RefundApplyModel) Insert(ctx context.Context, r *RefundApply) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`INSERT INTO %s (id, order_detail_id, order_id, pay_order_no, refund_order_no, user_id,
		 refund_amount, status, refund_reason, message, approver, approve_opinion, remark,
		 failed_reason, question_desc, refund_channel, create_time, approve_time, finish_time,
		 update_time, creater, updater)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), ?, ?, NOW(), ?, ?)`, m.table),
		r.Id, r.OrderDetailId, r.OrderId, nullInt64(r.PayOrderNo), nullInt64(r.RefundOrderNo),
		r.UserId, r.RefundAmount, r.Status, r.RefundReason, r.Message, nullInt64(r.Approver),
		nullString(r.ApproveOpinion), nullString(r.Remark), nullString(r.FailedReason),
		nullString(r.QuestionDesc), nullString(r.RefundChannel), nullTime(r.ApproveTime),
		nullTime(r.FinishTime), r.Creater, r.Updater)
	return err
}

// FindById 按主键查询退款申请。
func (m *RefundApplyModel) FindById(ctx context.Context, id int64) (*RefundApply, error) {
	var r RefundApply
	err := m.conn.QueryRowCtx(ctx, &r, fmt.Sprintf(
		`SELECT %s FROM %s WHERE id = ?`, refundApplyColumns, m.table), id)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// FindByOrderDetail 查询某明细的退款申请（最新一条）。
func (m *RefundApplyModel) FindByOrderDetail(ctx context.Context, orderDetailId int64) (*RefundApply, error) {
	var r RefundApply
	err := m.conn.QueryRowCtx(ctx, &r, fmt.Sprintf(
		`SELECT %s FROM %s WHERE order_detail_id = ? ORDER BY create_time DESC LIMIT 1`,
		refundApplyColumns, m.table), orderDetailId)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// UpdateStatus 更新退款状态、描述及时间。
func (m *RefundApplyModel) UpdateStatus(ctx context.Context, id, status int64, message string, approveTime, finishTime *time.Time) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET status = ?, message = ?, approve_time = ?, finish_time = ?, update_time = NOW()
		 WHERE id = ?`, m.table),
		status, message, nullTimePtr(approveTime), nullTimePtr(finishTime), id)
	return err
}

// Approve 审批退款申请（同意/拒绝）。
func (m *RefundApplyModel) Approve(ctx context.Context, id, status, approver int64, opinion, remark string, approveTime, finishTime *time.Time) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET status = ?, approver = ?, approve_opinion = ?, remark = ?,
		 approve_time = ?, finish_time = ?, update_time = NOW() WHERE id = ?`, m.table),
		status, approver, opinion, remark, nullTimePtr(approveTime), nullTimePtr(finishTime), id)
	return err
}

// MarkRefundDone 标记退款成功/失败（回填退款单号与渠道）。
func (m *RefundApplyModel) MarkRefundDone(ctx context.Context, id, status int64, refundOrderNo int64, channel, failedReason string) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET status = ?, refund_order_no = ?, refund_channel = ?, failed_reason = ?,
		 finish_time = NOW(), update_time = NOW() WHERE id = ?`, m.table),
		status, refundOrderNo, channel, failedReason, id)
	return err
}

// ListPendingOrder 查询最早一条待审批的退款申请。
func (m *RefundApplyModel) ListPendingOrder(ctx context.Context) (*RefundApply, error) {
	var r RefundApply
	err := m.conn.QueryRowCtx(ctx, &r, fmt.Sprintf(
		`SELECT %s FROM %s WHERE status = ? ORDER BY create_time ASC LIMIT 1`,
		refundApplyColumns, m.table), RefundStatusPending)
	if err != nil {
		return nil, err
	}
	return &r, nil
}

// ListPage 管理端退款申请分页查询。
func (m *RefundApplyModel) ListPage(ctx context.Context, cond *RefundApplyPageCond) ([]RefundApply, int64, error) {
	where := "1 = 1"
	args := make([]any, 0)
	if cond.Id > 0 {
		where += " AND id = ?"
		args = append(args, cond.Id)
	}
	if cond.OrderId > 0 {
		where += " AND order_id = ?"
		args = append(args, cond.OrderId)
	}
	if cond.OrderDetailId > 0 {
		where += " AND order_detail_id = ?"
		args = append(args, cond.OrderDetailId)
	}
	if cond.UserId > 0 {
		where += " AND user_id = ?"
		args = append(args, cond.UserId)
	}
	if cond.Status > 0 {
		where += " AND status = ?"
		args = append(args, cond.Status)
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
	var list []RefundApply
	err := m.conn.QueryRowsCtx(ctx, &list, fmt.Sprintf(
		`SELECT %s FROM %s WHERE %s ORDER BY create_time %s LIMIT ? OFFSET ?`,
		refundApplyColumns, m.table, where, order), args...)
	return list, total, err
}

// RefundApplyPageCond 退款申请分页条件
type RefundApplyPageCond struct {
	Id            int64
	OrderId       int64
	OrderDetailId int64
	UserId        int64
	Status        int64
	StartTime     time.Time
	EndTime       time.Time
	Offset        int64
	Limit         int64
	IsAsc         bool
}

func nullTimePtr(t *time.Time) any {
	if t == nil {
		return nil
	}
	return *t
}
