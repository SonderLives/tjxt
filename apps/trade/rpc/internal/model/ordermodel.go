package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ OrderModel = (*customOrderModel)(nil)

type (
	OrderModel interface {
		orderModel
		PageQueryByUser(ctx context.Context, userId, pageNo, pageSize, status, noNo int64, sortBy string, isAsc bool) ([]*Order, int64, error)
		UpdateStatus(ctx context.Context, id, status int64, message string) error
		MarkPaid(ctx context.Context, id, payOrderNo int64, payChannel string, payTime time.Time, realAmount int64) error
		MarkClosed(ctx context.Context, id int64, message string) error
	}
	customOrderModel struct {
		*defaultOrderModel
	}
)

func NewOrderModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) OrderModel {
	return &customOrderModel{
		defaultOrderModel: newOrderModel(conn, c, opts...),
	}
}

// PageQueryByUser 分页查询某用户的订单。
// status>0 时按状态过滤；noNo>0 时按订单 id 精确过滤（proto 字段 no_no 语义为订单号）。
func (m *customOrderModel) PageQueryByUser(ctx context.Context, userId, pageNo, pageSize, status, noNo int64, sortBy string, isAsc bool) ([]*Order, int64, error) {
	conditions := []string{"`user_id` = ?"}
	args := []any{userId}
	if status > 0 {
		conditions = append(conditions, "`status` = ?")
		args = append(args, status)
	}
	if noNo > 0 {
		conditions = append(conditions, "`id` = ?")
		args = append(args, noNo)
	}
	where := strings.Join(conditions, " and ")

	var total int64
	countQuery := fmt.Sprintf("select count(*) from %s where %s", m.table, where)
	if err := m.QueryRowNoCacheCtx(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	orderBy := "create_time desc"
	if sortBy != "" {
		dir := "desc"
		if isAsc {
			dir = "asc"
		}
		orderBy = fmt.Sprintf("%s %s", sortBy, dir)
	}
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (pageNo - 1) * pageSize

	listQuery := fmt.Sprintf("select %s from %s where %s order by %s limit ?, ?", orderRows, m.table, where, orderBy)
	var list []*Order
	err := m.QueryRowsNoCacheCtx(ctx, &list, listQuery, append(args, offset, pageSize)...)
	if err == sqlc.ErrNotFound {
		return nil, total, nil
	}
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// UpdateStatus 更新订单状态与状态备注
func (m *customOrderModel) UpdateStatus(ctx context.Context, id, status int64, message string) error {
	order, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}
	order.Status = status
	order.Message = message
	return m.Update(ctx, order)
}

// MarkPaid 标记订单为已支付，写入支付流水号、渠道、支付时间与实付金额
func (m *customOrderModel) MarkPaid(ctx context.Context, id, payOrderNo int64, payChannel string, payTime time.Time, realAmount int64) error {
	order, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}
	order.Status = 2
	order.PayOrderNo = sql.NullInt64{Int64: payOrderNo, Valid: true}
	order.PayChannel = payChannel
	order.PayTime = sql.NullTime{Time: payTime, Valid: true}
	if realAmount > 0 {
		order.RealAmount = realAmount
	}
	return m.Update(ctx, order)
}

// MarkClosed 标记订单关闭（取消/超时）
func (m *customOrderModel) MarkClosed(ctx context.Context, id int64, message string) error {
	return m.UpdateStatus(ctx, id, 3, message)
}
