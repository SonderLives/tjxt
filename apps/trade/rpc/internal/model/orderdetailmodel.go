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

var _ OrderDetailModel = (*customOrderDetailModel)(nil)

type (
	OrderDetailModel interface {
		orderDetailModel
		ListByOrderId(ctx context.Context, orderId int64) ([]*OrderDetail, error)
		ListByCourseId(ctx context.Context, courseId int64) ([]*OrderDetail, error)
		ListByCourseIds(ctx context.Context, courseIds []int64) ([]*OrderDetail, error)
		CountPaidByCourseIds(ctx context.Context, courseIds []int64) (map[int64]int64, error)
		StatByCourseId(ctx context.Context, courseId int64) (enrollNum, realPayAmount, refundNum int64, err error)
		CountPaidByUserIds(ctx context.Context, userIds []int64) (map[int64]int64, error)
		UpdateStatus(ctx context.Context, id, status int64) error
		UpdateRefundStatus(ctx context.Context, id, refundStatus int64) error
		PageQuery(ctx context.Context, f OrderDetailPageFilter) ([]*OrderDetail, int64, error)
	}
	customOrderDetailModel struct {
		*defaultOrderDetailModel
	}
)

func NewOrderDetailModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderDetailModel {
	return &customOrderDetailModel{
		defaultOrderDetailModel: newOrderDetailModel(conn, c, opts...),
	}
}

// 已支付/已报名/已完成的订单明细状态（可作为"报名"口径）
const paidDetailStatusSQL = "2,4,5"
// 已同意退款/退款成功的退款状态（可作为"退款"口径）
const refundedDetailStatusSQL = "3,5"

// ListByOrderId 查询某订单的全部明细
func (m *customOrderDetailModel) ListByOrderId(ctx context.Context, orderId int64) ([]*OrderDetail, error) {
	query := fmt.Sprintf("select %s from %s where `order_id` = ? order by `id` asc", orderDetailRows, m.table)
	var resp []*OrderDetail
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

// ListByCourseId 查询某课程的全部明细
func (m *customOrderDetailModel) ListByCourseId(ctx context.Context, courseId int64) ([]*OrderDetail, error) {
	query := fmt.Sprintf("select %s from %s where `course_id` = ?", orderDetailRows, m.table)
	var resp []*OrderDetail
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, courseId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// ListByCourseIds 批量查询课程明细（用于报名人数统计）
func (m *customOrderDetailModel) ListByCourseIds(ctx context.Context, courseIds []int64) ([]*OrderDetail, error) {
	if len(courseIds) == 0 {
		return nil, nil
	}
	query := fmt.Sprintf("select %s from %s where `course_id` in (%s)", orderDetailRows, m.table, placeholders(len(courseIds)))
	var resp []*OrderDetail
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, toAnySlice(courseIds)...)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// CountPaidByCourseIds 统计各课程的已支付报名人数，返回 course_id -> 人数
func (m *customOrderDetailModel) CountPaidByCourseIds(ctx context.Context, courseIds []int64) (map[int64]int64, error) {
	res := make(map[int64]int64)
	if len(courseIds) == 0 {
		return res, nil
	}
	query := fmt.Sprintf("select `course_id` as `id`, count(*) as `cnt` from %s where `course_id` in (%s) and `status` in (%s) group by `course_id`",
		m.table, placeholders(len(courseIds)), paidDetailStatusSQL)
	rows, err := m.queryGroupCount(ctx, query, toAnySlice(courseIds))
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		res[r.Id] = r.Cnt
	}
	return res, nil
}

// StatByCourseId 统计某课程的报名人数、实付总额、退款人数
func (m *customOrderDetailModel) StatByCourseId(ctx context.Context, courseId int64) (enrollNum, realPayAmount, refundNum int64, err error) {
	statQuery := fmt.Sprintf("select count(*) as `cnt`, coalesce(sum(`real_pay_amount`),0) as `sum` from %s where `course_id` = ? and `status` in (%s)", m.table, paidDetailStatusSQL)
	var row struct {
		Cnt int64         `db:"cnt"`
		Sum  sql.NullInt64 `db:"sum"`
	}
	if e := m.QueryRowNoCacheCtx(ctx, &row, statQuery, courseId); e != nil {
		return 0, 0, 0, e
	}
	enrollNum = row.Cnt
	if row.Sum.Valid {
		realPayAmount = row.Sum.Int64
	}

	refundQuery := fmt.Sprintf("select count(*) from %s where `course_id` = ? and `refund_status` in (%s)", m.table, refundedDetailStatusSQL)
	if e := m.QueryRowNoCacheCtx(ctx, &refundNum, refundQuery, courseId); e != nil {
		return 0, 0, 0, e
	}
	return enrollNum, realPayAmount, refundNum, nil
}

// CountPaidByUserIds 统计各学员的报名课程数，返回 user_id -> 课程数
func (m *customOrderDetailModel) CountPaidByUserIds(ctx context.Context, userIds []int64) (map[int64]int64, error) {
	res := make(map[int64]int64)
	if len(userIds) == 0 {
		return res, nil
	}
	query := fmt.Sprintf("select `user_id` as `id`, count(distinct `course_id`) as `cnt` from %s where `user_id` in (%s) and `status` in (%s) group by `user_id`",
		m.table, placeholders(len(userIds)), paidDetailStatusSQL)
	rows, err := m.queryGroupCount(ctx, query, toAnySlice(userIds))
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		res[r.Id] = r.Cnt
	}
	return res, nil
}

// UpdateStatus 更新订单明细状态
func (m *customOrderDetailModel) UpdateStatus(ctx context.Context, id, status int64) error {
	d, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}
	d.Status = status
	return m.Update(ctx, d)
}

// UpdateRefundStatus 更新订单明细退款状态
func (m *customOrderDetailModel) UpdateRefundStatus(ctx context.Context, id, refundStatus int64) error {
	d, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}
	d.RefundStatus = sql.NullInt64{Int64: refundStatus, Valid: true}
	return m.Update(ctx, d)
}

// OrderDetailPageFilter 订单明细管理分页过滤条件
type OrderDetailPageFilter struct {
	Id           int64
	OrderId      int64
	Status       int64
	RefundStatus int64
	PayChannel   string
	StartTime    string
	EndTime      string
	PageNo       int64
	PageSize     int64
	IsAsc        bool
	SortBy       string
}

// PageQuery 订单明细管理分页（mobile 过滤因无 user 服务接入，暂不支持）
func (m *customOrderDetailModel) PageQuery(ctx context.Context, f OrderDetailPageFilter) ([]*OrderDetail, int64, error) {
	conditions := []string{}
	args := []any{}
	if f.Id > 0 {
		conditions = append(conditions, "`id` = ?")
		args = append(args, f.Id)
	}
	if f.OrderId > 0 {
		conditions = append(conditions, "`order_id` = ?")
		args = append(args, f.OrderId)
	}
	if f.Status > 0 {
		conditions = append(conditions, "`status` = ?")
		args = append(args, f.Status)
	}
	if f.RefundStatus > 0 {
		conditions = append(conditions, "`refund_status` = ?")
		args = append(args, f.RefundStatus)
	}
	if f.PayChannel != "" {
		conditions = append(conditions, "`pay_channel` = ?")
		args = append(args, f.PayChannel)
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

	query := fmt.Sprintf("select %s from %s where %s order by %s limit ?, ?", orderDetailRows, m.table, where, orderBy)
	var list []*OrderDetail
	err := m.QueryRowsNoCacheCtx(ctx, &list, query, append(args, offset, pageSize)...)
	if err == sqlc.ErrNotFound {
		return nil, total, nil
	}
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// queryGroupCount 执行 group by 查询并扫描为 id->cnt 列表
func (m *customOrderDetailModel) queryGroupCount(ctx context.Context, q string, args []any) ([]idCount, error) {
	var rows []idCount
	err := m.QueryRowsNoCacheCtx(ctx, &rows, q, args...)
	switch err {
	case nil:
		return rows, nil
	case sqlc.ErrNotFound:
		return nil, nil
	default:
		return nil, err
	}
}

type idCount struct {
	Id  int64 `db:"id"`
	Cnt int64 `db:"cnt"`
}

func toAnySlice(ids []int64) []any {
	out := make([]any, 0, len(ids))
	for _, id := range ids {
		out = append(out, id)
	}
	return out
}
