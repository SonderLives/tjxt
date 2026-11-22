package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CouponCodeModel = (*customCouponCodeModel)(nil)

type (
	// CouponCodeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCouponCodeModel.
	CouponCodeModel interface {
		couponCodeModel
		// FindPageByCoupon 按优惠券+状态分页查询兑换码，offset/limit 由 pkg/utils/page 归一化后传入。
		FindPageByCoupon(ctx context.Context, couponId int64, status string, offset, limit int64) ([]*CouponCode, int64, error)
		// MarkUsed 将兑换码标记为已使用，返回受影响行数（0 表示已被抢兑）。
		MarkUsed(ctx context.Context, id, userId int64) (int64, error)
		// BatchInsert 批量写入兑换码，发放优惠券时使用。
		BatchInsert(ctx context.Context, couponId int64, codes []string, creater int64) error
	}

	customCouponCodeModel struct {
		*defaultCouponCodeModel
	}
)

// NewCouponCodeModel returns a model for the database table.
func NewCouponCodeModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CouponCodeModel {
	return &customCouponCodeModel{
		defaultCouponCodeModel: newCouponCodeModel(conn, c, opts...),
	}
}

func (m *customCouponCodeModel) FindPageByCoupon(ctx context.Context, couponId int64, status string,
	offset, limit int64) ([]*CouponCode, int64, error) {
	where := []string{"`deleted` = 0"}
	args := []any{}
	if couponId > 0 {
		where = append(where, "`coupon_id` = ?")
		args = append(args, couponId)
	}
	if status != "" {
		where = append(where, "`status` = ?")
		args = append(args, status)
	}
	whereClause := strings.Join(where, " and ")

	query := fmt.Sprintf("select %s from %s where %s order by `id` desc limit ?,?",
		couponCodeRows, m.table, whereClause)
	countQuery := fmt.Sprintf("select count(*) from %s where %s", m.table, whereClause)

	var (
		resp  []*CouponCode
		total int64
	)
	if err := m.QueryRowsNoCacheCtx(ctx, &resp, query, append(append([]any{}, args...), offset, limit)...); err != nil {
		return nil, 0, err
	}
	if err := m.QueryRowNoCacheCtx(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}
	return resp, total, nil
}

func (m *customCouponCodeModel) MarkUsed(ctx context.Context, id, userId int64) (int64, error) {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return 0, err
	}

	query := fmt.Sprintf("update %s set `status` = ?, `user_id` = ?, `update_time` = now() "+
		"where `id` = ? and `status` = ? and `deleted` = 0", m.table)
	res, err := m.ExecNoCacheCtx(ctx, query, CouponCodeStatusUsed, userId, id, CouponCodeStatusUnused)
	if err != nil {
		return 0, err
	}

	_ = m.DelCacheCtx(ctx,
		fmt.Sprintf("%s%v", cacheCouponCodeIdPrefix, id),
		fmt.Sprintf("%s%v", cacheCouponCodeCodePrefix, data.Code))
	return res.RowsAffected()
}

func (m *customCouponCodeModel) BatchInsert(ctx context.Context, couponId int64, codes []string, creater int64) error {
	if len(codes) == 0 {
		return nil
	}
	values := make([]string, 0, len(codes))
	args := make([]any, 0, len(codes)*5)
	for _, code := range codes {
		values = append(values, "(?, ?, ?, ?, ?)")
		args = append(args, couponId, code, CouponCodeStatusUnused, creater, creater)
	}
	query := fmt.Sprintf("insert into %s (`coupon_id`, `code`, `status`, `creater`, `updater`) values %s",
		m.table, strings.Join(values, ","))
	_, err := m.ExecNoCacheCtx(ctx, query, args...)
	return err
}
