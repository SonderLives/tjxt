package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PayOrderModel = (*customPayOrderModel)(nil)

type (
	// PayOrderModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPayOrderModel.
	PayOrderModel interface {
		payOrderModel

		// MarkToPaying 待提交 -> 待支付，写入二维码
		MarkToPaying(ctx context.Context, id int64, qrCodeUrl string) error
		// MarkToSuccess 待支付 -> 支付成功（条件更新），返回是否真的发生状态流转
		MarkToSuccess(ctx context.Context, id int64, resultCode, resultMsg string) (bool, error)
		// MarkToClosed 待提交/待支付 -> 关闭（含超时/取消/失败原因）
		MarkToClosed(ctx context.Context, id int64, resultCode, resultMsg string) error
		// IncrNotifyTimes 增加通知次数
		IncrNotifyTimes(ctx context.Context, id int64) error
		// SetNotifyStatus 设置回调状态
		SetNotifyStatus(ctx context.Context, id int64, status int64) error
	}

	customPayOrderModel struct {
		*defaultPayOrderModel
	}
)

// NewPayOrderModel returns a model for the database table.
func NewPayOrderModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PayOrderModel {
	return &customPayOrderModel{
		defaultPayOrderModel: newPayOrderModel(conn, c, opts...),
	}
}

func (m *customPayOrderModel) MarkToPaying(ctx context.Context, id int64, qrCodeUrl string) error {
	_, err := m.ExecNoCacheCtx(ctx,
		fmt.Sprintf("update %s set `status` = 1, `qr_code_url` = ?, `update_time` = ? where `id` = ?", m.table),
		qrCodeUrl, time.Now(), id)
	return err
}

// MarkToSuccess 待支付 -> 支付成功。带 `status in (0,1)` 条件更新，
// 并发回调下只有一条请求真正流转成功（幂等）；同时失效主键/业务单号/支付单号三路缓存。
// 返回 false 表示状态已被其他请求处理，调用方按幂等成功处理且不应重复发事件。
func (m *customPayOrderModel) MarkToSuccess(ctx context.Context, id int64, resultCode, resultMsg string) (bool, error) {
	payOrder, err := m.FindOne(ctx, id)
	if err != nil {
		return false, err
	}
	payOrderIdKey := fmt.Sprintf("%s%v", cachePayOrderIdPrefix, id)
	payOrderBizOrderNoKey := fmt.Sprintf("%s%v", cachePayOrderBizOrderNoPrefix, payOrder.BizOrderNo)
	payOrderPayOrderNoKey := fmt.Sprintf("%s%v", cachePayOrderPayOrderNoPrefix, payOrder.PayOrderNo)

	res, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (sql.Result, error) {
		return conn.ExecCtx(ctx,
			fmt.Sprintf("update %s set `status` = 3, `result_code` = ?, `result_msg` = ?, `pay_success_time` = ?, `update_time` = ? where `id` = ? and `status` in (0, 1)", m.table),
			resultCode, resultMsg, time.Now(), time.Now(), id)
	}, payOrderIdKey, payOrderBizOrderNoKey, payOrderPayOrderNoKey)
	if err != nil {
		return false, err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func (m *customPayOrderModel) MarkToClosed(ctx context.Context, id int64, resultCode, resultMsg string) error {
	_, err := m.ExecNoCacheCtx(ctx,
		fmt.Sprintf("update %s set `status` = 2, `result_code` = ?, `result_msg` = ?, `update_time` = ? where `id` = ?", m.table),
		resultCode, resultMsg, time.Now(), id)
	return err
}

func (m *customPayOrderModel) IncrNotifyTimes(ctx context.Context, id int64) error {
	_, err := m.ExecNoCacheCtx(ctx,
		fmt.Sprintf("update %s set `notify_times` = `notify_times` + 1, `update_time` = ? where `id` = ?", m.table),
		time.Now(), id)
	return err
}

func (m *customPayOrderModel) SetNotifyStatus(ctx context.Context, id int64, status int64) error {
	_, err := m.ExecNoCacheCtx(ctx,
		fmt.Sprintf("update %s set `notify_status` = ?, `update_time` = ? where `id` = ?", m.table),
		status, time.Now(), id)
	return err
}