package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserCouponModel = (*customUserCouponModel)(nil)

type (
	// UserCouponModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserCouponModel.
	UserCouponModel interface {
		userCouponModel
		// FindByUserAndStatus 查询指定用户在某状态下的优惠券，status 为空表示不限。
		FindByUserAndStatus(ctx context.Context, userId int64, status string) ([]*UserCoupon, error)
		// FindPageByUser 按用户+状态分页查询，offset/limit 由 pkg/utils/page 归一化后传入。
		FindPageByUser(ctx context.Context, userId int64, status string, offset, limit int64) ([]*UserCoupon, int64, error)
		// FindByIdsAndUser 批量按 id 查询当前用户的优惠券，防止越权使用他人券。
		FindByIdsAndUser(ctx context.Context, ids []int64, userId int64) ([]*UserCoupon, error)
		// CountByUserAndCoupon 统计用户已领取某券的数量，用于 user_limit 校验。
		CountByUserAndCoupon(ctx context.Context, userId, couponId int64) (int64, error)
		// UpdateStatusByIds 批量更新状态（使用/退还），带条件状态防并发覆盖。
		UpdateStatusByIds(ctx context.Context, ids []int64, userId int64, fromStatus, toStatus string, orderId int64) (int64, error)
	}

	customUserCouponModel struct {
		*defaultUserCouponModel
	}
)

// NewUserCouponModel returns a model for the database table.
func NewUserCouponModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserCouponModel {
	return &customUserCouponModel{
		defaultUserCouponModel: newUserCouponModel(conn, c, opts...),
	}
}

func (m *customUserCouponModel) FindByUserAndStatus(ctx context.Context, userId int64, status string) ([]*UserCoupon, error) {
	where := []string{"`deleted` = 0", "`user_id` = ?"}
	args := []any{userId}
	if status != "" {
		where = append(where, "`status` = ?")
		args = append(args, status)
	}
	query := fmt.Sprintf("select %s from %s where %s order by `id` desc",
		userCouponRows, m.table, strings.Join(where, " and "))

	var resp []*UserCoupon
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}

func (m *customUserCouponModel) FindPageByUser(ctx context.Context, userId int64, status string, offset, limit int64) ([]*UserCoupon, int64, error) {
	where := []string{"`deleted` = 0", "`user_id` = ?"}
	args := []any{userId}
	if status != "" {
		where = append(where, "`status` = ?")
		args = append(args, status)
	}
	whereClause := strings.Join(where, " and ")

	query := fmt.Sprintf("select %s from %s where %s order by `id` desc limit ?,?",
		userCouponRows, m.table, whereClause)
	countQuery := fmt.Sprintf("select count(*) from %s where %s", m.table, whereClause)

	var (
		resp  []*UserCoupon
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

func (m *customUserCouponModel) FindByIdsAndUser(ctx context.Context, ids []int64, userId int64) ([]*UserCoupon, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	holders := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, 0, len(ids)+1)
	for _, id := range ids {
		args = append(args, id)
	}
	args = append(args, userId)

	query := fmt.Sprintf("select %s from %s where `deleted` = 0 and `id` in (%s) and `user_id` = ?",
		userCouponRows, m.table, holders)

	var resp []*UserCoupon
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}

func (m *customUserCouponModel) CountByUserAndCoupon(ctx context.Context, userId, couponId int64) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where `deleted` = 0 and `user_id` = ? and `coupon_id` = ?", m.table)
	var total int64
	err := m.QueryRowNoCacheCtx(ctx, &total, query, userId, couponId)
	return total, err
}

func (m *customUserCouponModel) UpdateStatusByIds(ctx context.Context, ids []int64, userId int64,
	fromStatus, toStatus string, orderId int64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	holders := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")

	var (
		setClause string
		args      []any
	)
	switch toStatus {
	case UserCouponStatusUsed:
		setClause = "`status` = ?, `use_time` = now(), `order_id` = ?, `update_time` = now()"
		args = append(args, toStatus, orderId)
	default:
		setClause = "`status` = ?, `use_time` = null, `order_id` = null, `update_time` = now()"
		args = append(args, toStatus)
	}
	for _, id := range ids {
		args = append(args, id)
	}
	args = append(args, userId, fromStatus)

	query := fmt.Sprintf("update %s set %s where `id` in (%s) and `user_id` = ? and `status` = ? and `deleted` = 0",
		m.table, setClause, holders)

	res, err := m.ExecNoCacheCtx(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	// 状态变更后缓存需失效，避免读到旧数据。
	for _, id := range ids {
		_ = m.DelCacheCtx(ctx, fmt.Sprintf("%s%v", cacheUserCouponIdPrefix, id))
	}
	return res.RowsAffected()
}
