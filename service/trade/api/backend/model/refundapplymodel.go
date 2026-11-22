package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// placeholders 生成 n 个 "?" 占位符，逗号分隔
func placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.TrimSuffix(strings.Repeat("?,", n), ",")
}

var _ RefundApplyModel = (*customRefundApplyModel)(nil)

type (
	RefundApplyModel interface {
		refundApplyModel
		ListByUserId(ctx context.Context, userId int64) ([]*RefundApply, error)
		FindByOrderDetailId(ctx context.Context, orderDetailId int64) (*RefundApply, error)
		ListByOrderId(ctx context.Context, orderId int64) ([]*RefundApply, error)
		FindNextPending(ctx context.Context) (*RefundApply, error)
		PageQuery(ctx context.Context, f RefundApplyPageFilter) ([]*RefundApply, int64, error)
		UpdateStatus(ctx context.Context, id, status int64) error
		UpdateApprove(ctx context.Context, id, status, approver int64, opinion, remark string, approveTime sql.NullTime, refundOrderNo, payOrderNo int64) error
	}
	customRefundApplyModel struct {
		*defaultRefundApplyModel
	}
)

func NewRefundApplyModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RefundApplyModel {
	return &customRefundApplyModel{
		defaultRefundApplyModel: newRefundApplyModel(conn, c, opts...),
	}
}

// ListByUserId 查询某用户的退款申请
func (m *customRefundApplyModel) ListByUserId(ctx context.Context, userId int64) ([]*RefundApply, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? order by `create_time` desc", refundApplyRows, m.table)
	var resp []*RefundApply
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindByOrderDetailId 查询某订单明细关联的退款申请
func (m *customRefundApplyModel) FindByOrderDetailId(ctx context.Context, orderDetailId int64) (*RefundApply, error) {
	query := fmt.Sprintf("select %s from %s where `order_detail_id` = ? order by `create_time` desc limit 1", refundApplyRows, m.table)
	var resp RefundApply
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, orderDetailId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// ListByOrderId 查询某订单的全部退款申请
func (m *customRefundApplyModel) ListByOrderId(ctx context.Context, orderId int64) ([]*RefundApply, error) {
	query := fmt.Sprintf("select %s from %s where `order_id` = ? order by `create_time` desc", refundApplyRows, m.table)
	var resp []*RefundApply
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, orderId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindNextPending 取出下一条待审批（status=1）的退款申请，用于轮询/审批台
func (m *customRefundApplyModel) FindNextPending(ctx context.Context) (*RefundApply, error) {
	query := fmt.Sprintf("select %s from %s where `status` = 1 order by `create_time` asc limit 1", refundApplyRows, m.table)
	var resp RefundApply
	err := m.QueryRowNoCacheCtx(ctx, &resp, query)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// RefundApplyPageFilter 退款申请管理分页过滤条件
type RefundApplyPageFilter struct {
	Id            int64
	OrderDetailId int64
	OrderId       int64
	RefundStatus  int64
	StartTime     string
	EndTime       string
	PageNo        int64
	PageSize      int64
	IsAsc         bool
	SortBy        string
}

// PageQuery 退款申请管理分页（mobile 过滤因无 user 服务接入，暂不支持）
func (m *customRefundApplyModel) PageQuery(ctx context.Context, f RefundApplyPageFilter) ([]*RefundApply, int64, error) {
	conditions := []string{}
	args := []any{}
	if f.Id > 0 {
		conditions = append(conditions, "`id` = ?")
		args = append(args, f.Id)
	}
	if f.OrderDetailId > 0 {
		conditions = append(conditions, "`order_detail_id` = ?")
		args = append(args, f.OrderDetailId)
	}
	if f.OrderId > 0 {
		conditions = append(conditions, "`order_id` = ?")
		args = append(args, f.OrderId)
	}
	if f.RefundStatus > 0 {
		conditions = append(conditions, "`status` = ?")
		args = append(args, f.RefundStatus)
	}
	if f.StartTime != "" {
		conditions = append(conditions, "`create_time` >= ?")
		args = append(args, f.StartTime)
	}
	if f.EndTime != "" {
		conditions = append(conditions, "`create_time` <= ?")
		args = append(args, f.EndTime)
	}
	where := "1=1"
	if len(conditions) > 0 {
		where = strings.Join(conditions, " and ")
	}

	var total int64
	if err := m.QueryRowNoCacheCtx(ctx, &total, fmt.Sprintf("select count(*) from %s where %s", m.table, where), args...); err != nil {
		return nil, 0, err
	}

	orderBy := "create_time desc"
	if f.SortBy != "" {
		dir := "desc"
		if f.IsAsc {
			dir = "asc"
		}
		orderBy = f.SortBy + " " + dir
	}
	pageNo, pageSize := f.PageNo, f.PageSize
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (pageNo - 1) * pageSize

	query := fmt.Sprintf("select %s from %s where %s order by %s limit ?, ?", refundApplyRows, m.table, where, orderBy)
	var list []*RefundApply
	err := m.QueryRowsNoCacheCtx(ctx, &list, query, append(args, offset, pageSize)...)
	if err == sqlc.ErrNotFound {
		return nil, total, nil
	}
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// UpdateStatus 更新退款申请状态（终态写入 finish_time）
func (m *customRefundApplyModel) UpdateStatus(ctx context.Context, id, status int64) error {
	ra, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}
	ra.Status = status
	if status == 5 || status == 6 {
		ra.FinishTime = sql.NullTime{Time: ra.UpdateTime, Valid: true}
	}
	return m.Update(ctx, ra)
}

// UpdateApprove 写入审批结果（同意/拒绝）
func (m *customRefundApplyModel) UpdateApprove(ctx context.Context, id, status, approver int64, opinion, remark string, approveTime sql.NullTime, refundOrderNo, payOrderNo int64) error {
	ra, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}
	ra.Status = status
	ra.Approver = sql.NullInt64{Int64: approver, Valid: approver > 0}
	ra.ApproveOpinion = sql.NullString{String: opinion, Valid: opinion != ""}
	ra.Remark = sql.NullString{String: remark, Valid: remark != ""}
	ra.ApproveTime = approveTime
	if refundOrderNo > 0 {
		ra.RefundOrderNo = sql.NullInt64{Int64: refundOrderNo, Valid: true}
	}
	if payOrderNo > 0 {
		ra.PayOrderNo = sql.NullInt64{Int64: payOrderNo, Valid: true}
	}
	return m.Update(ctx, ra)
}
