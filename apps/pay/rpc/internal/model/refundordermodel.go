package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RefundOrderModel = (*customRefundOrderModel)(nil)

type (
	// RefundOrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRefundOrderModel.
	RefundOrderModel interface {
		refundOrderModel

		// FindOneByBizRefundOrderNo 按业务退款单号查询
		FindOneByBizRefundOrderNo(ctx context.Context, bizRefundOrderNo int64) (*RefundOrder, error)
		// FindOneByRefundOrderNo 按退款单号查询
		FindOneByRefundOrderNo(ctx context.Context, refundOrderNo int64) (*RefundOrder, error)
		// FindListByBizOrderNo 列出某业务订单下所有退款单
		FindListByBizOrderNo(ctx context.Context, bizOrderNo int64) ([]*RefundOrder, error)
		// MarkToProcessing 未提交 -> 退款中
		MarkToProcessing(ctx context.Context, id int64) error
		// MarkToSuccess 退款中 -> 退款成功
		MarkToSuccess(ctx context.Context, id int64, resultCode, resultMsg, refundChannel string) error
		// MarkToFailed 退款中 -> 退款失败
		MarkToFailed(ctx context.Context, id int64, resultCode, resultMsg string) error
		// SetNotifyStatus 设置退款回调状态
		SetNotifyStatus(ctx context.Context, id int64, status int64) error
		// IncrNotifyFailedTimes 累计通知失败次数
		IncrNotifyFailedTimes(ctx context.Context, id int64) error
	}

	customRefundOrderModel struct {
		*defaultRefundOrderModel
	}
)

// NewRefundOrderModel returns a model for the database table.
func NewRefundOrderModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RefundOrderModel {
	return &customRefundOrderModel{
		defaultRefundOrderModel: newRefundOrderModel(conn, c, opts...),
	}
}

func (m *customRefundOrderModel) FindOneByBizRefundOrderNo(ctx context.Context, bizRefundOrderNo int64) (*RefundOrder, error) {
	var resp RefundOrder
	query := fmt.Sprintf("select %s from %s where `biz_refund_order_no` = ? order by `id` desc limit 1", refundOrderRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, bizRefundOrderNo)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *customRefundOrderModel) FindOneByRefundOrderNo(ctx context.Context, refundOrderNo int64) (*RefundOrder, error) {
	var resp RefundOrder
	query := fmt.Sprintf("select %s from %s where `refund_order_no` = ? limit 1", refundOrderRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, refundOrderNo)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *customRefundOrderModel) FindListByBizOrderNo(ctx context.Context, bizOrderNo int64) ([]*RefundOrder, error) {
	var list []*RefundOrder
	query := fmt.Sprintf("select %s from %s where `biz_order_no` = ? and `deleted` = 0 order by `id` desc", refundOrderRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &list, query, bizOrderNo)
	return list, err
}

func (m *customRefundOrderModel) MarkToProcessing(ctx context.Context, id int64) error {
	_, err := m.ExecNoCacheCtx(ctx,
		fmt.Sprintf("update %s set `status` = 1, `update_time` = ? where `id` = ?", m.table),
		time.Now(), id)
	return err
}

func (m *customRefundOrderModel) MarkToSuccess(ctx context.Context, id int64, resultCode, resultMsg, refundChannel string) error {
	_, err := m.ExecNoCacheCtx(ctx,
		fmt.Sprintf("update %s set `status` = 3, `result_code` = ?, `result_msg` = ?, `refund_channel` = ?, `update_time` = ? where `id` = ?", m.table),
		resultCode, resultMsg, sql.NullString{String: refundChannel, Valid: refundChannel != ""}, time.Now(), id)
	return err
}

func (m *customRefundOrderModel) MarkToFailed(ctx context.Context, id int64, resultCode, resultMsg string) error {
	_, err := m.ExecNoCacheCtx(ctx,
		fmt.Sprintf("update %s set `status` = 2, `result_code` = ?, `result_msg` = ?, `update_time` = ? where `id` = ?", m.table),
		resultCode, resultMsg, time.Now(), id)
	return err
}

func (m *customRefundOrderModel) SetNotifyStatus(ctx context.Context, id int64, status int64) error {
	_, err := m.ExecNoCacheCtx(ctx,
		fmt.Sprintf("update %s set `notify_status` = ?, `update_time` = ? where `id` = ?", m.table),
		status, time.Now(), id)
	return err
}

func (m *customRefundOrderModel) IncrNotifyFailedTimes(ctx context.Context, id int64) error {
	_, err := m.ExecNoCacheCtx(ctx,
		fmt.Sprintf("update %s set `notify_failed_times` = `notify_failed_times` + 1, `update_time` = ? where `id` = ?", m.table),
		time.Now(), id)
	return err
}