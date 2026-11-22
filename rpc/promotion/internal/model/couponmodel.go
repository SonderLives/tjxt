package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CouponModel = (*customCouponModel)(nil)

type (
	// CouponModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCouponModel.
	CouponModel interface {
		couponModel
		// FindList 查询全部未删除优惠券。
		FindList(ctx context.Context) ([]*Coupon, error)
		// FindPage 按条件分页查询优惠券，offset/limit 由 pkg/utils/page 归一化后传入。
		FindPage(ctx context.Context, name, status, discountType string, offset, limit int64) ([]*Coupon, int64, error)
		// FindByIds 批量查询优惠券。
		FindByIds(ctx context.Context, ids []int64) ([]*Coupon, error)
		// IncrIssueNum 原子扣减库存（已领数量+1），返回受影响行数，0 表示已抢光或非发放中。
		IncrIssueNum(ctx context.Context, id int64) (int64, error)
		// DecrIssueNum 回滚库存（已领数量-1），用于领取流程后续步骤失败时补偿。
		DecrIssueNum(ctx context.Context, id int64) error
		// AddUsedNum 增减已使用数量，delta 可为负（退还场景）。
		AddUsedNum(ctx context.Context, id, delta int64) error
		// UpdateStatus 更新优惠券状态。
		UpdateStatus(ctx context.Context, id int64, status string, updater int64) error
		// SoftDelete 逻辑删除，仅允许删除未发放的券。
		SoftDelete(ctx context.Context, id, updater int64) (int64, error)
	}

	customCouponModel struct {
		*defaultCouponModel
	}
)

// NewCouponModel returns a model for the database table.
func NewCouponModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CouponModel {
	return &customCouponModel{
		defaultCouponModel: newCouponModel(conn, c, opts...),
	}
}

// FindList 查询全部未删除优惠券（管理后台列表）。
func (m *customCouponModel) FindList(ctx context.Context) ([]*Coupon, error) {
	query := fmt.Sprintf("select %s from %s where `deleted` = 0 order by `create_time` desc", couponRows, m.table)
	var resp []*Coupon
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	return resp, err
}

// FindPage 按条件分页查询优惠券。
func (m *customCouponModel) FindPage(ctx context.Context, name, status, discountType string, offset, limit int64) ([]*Coupon, int64, error) {
	where := []string{"`deleted` = 0"}
	args := []any{}
	if name != "" {
		where = append(where, "`name` like ?")
		args = append(args, "%"+name+"%")
	}
	if status != "" {
		where = append(where, "`status` = ?")
		args = append(args, status)
	}
	if discountType != "" {
		where = append(where, "`discount_type` = ?")
		args = append(args, discountType)
	}
	whereClause := strings.Join(where, " and ")
	query := fmt.Sprintf("select %s from %s where %s order by `create_time` desc limit ?,?", couponRows, m.table, whereClause)
	countQuery := fmt.Sprintf("select count(*) from %s where %s", m.table, whereClause)

	var (
		resp  []*Coupon
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

// FindByIds 批量查询优惠券，用于折扣计算时一次性加载券规则。
func (m *customCouponModel) FindByIds(ctx context.Context, ids []int64) ([]*Coupon, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	holders := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	query := fmt.Sprintf("select %s from %s where `deleted` = 0 and `id` in (%s)", couponRows, m.table, holders)

	var resp []*Coupon
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...)
	return resp, err
}

// IncrIssueNum 原子扣减库存：仅在发放中且未超发行总量时 +1，靠 SQL 条件保证并发安全。
func (m *customCouponModel) IncrIssueNum(ctx context.Context, id int64) (int64, error) {
	query := fmt.Sprintf("update %s set `issue_num` = `issue_num` + 1, `update_time` = now() "+
		"where `id` = ? and `deleted` = 0 and `status` = ? and (`total_num` = 0 or `issue_num` < `total_num`)", m.table)
	res, err := m.ExecNoCacheCtx(ctx, query, id, CouponStatusIssued)
	if err != nil {
		return 0, err
	}
	_ = m.DelCacheCtx(ctx, fmt.Sprintf("%s%v", cacheCouponIdPrefix, id))
	return res.RowsAffected()
}

// DecrIssueNum 库存回滚，保证不会把已领数量减成负数。
func (m *customCouponModel) DecrIssueNum(ctx context.Context, id int64) error {
	query := fmt.Sprintf("update %s set `issue_num` = `issue_num` - 1, `update_time` = now() "+
		"where `id` = ? and `issue_num` > 0", m.table)
	if _, err := m.ExecNoCacheCtx(ctx, query, id); err != nil {
		return err
	}
	return m.DelCacheCtx(ctx, fmt.Sprintf("%s%v", cacheCouponIdPrefix, id))
}

// AddUsedNum 已使用数量增减，delta 为负时用于退还，并保证不会减成负数。
func (m *customCouponModel) AddUsedNum(ctx context.Context, id, delta int64) error {
	query := fmt.Sprintf("update %s set `used_num` = `used_num` + ?, `update_time` = now() "+
		"where `id` = ? and `deleted` = 0 and `used_num` + ? >= 0", m.table)
	if _, err := m.ExecNoCacheCtx(ctx, query, delta, id, delta); err != nil {
		return err
	}
	return m.DelCacheCtx(ctx, fmt.Sprintf("%s%v", cacheCouponIdPrefix, id))
}

// UpdateStatus 更新优惠券状态（发放/暂停/结束）。
func (m *customCouponModel) UpdateStatus(ctx context.Context, id int64, status string, updater int64) error {
	query := fmt.Sprintf("update %s set `status` = ?, `updater` = ?, `update_time` = now() "+
		"where `id` = ? and `deleted` = 0", m.table)
	if _, err := m.ExecNoCacheCtx(ctx, query, status, updater, id); err != nil {
		return err
	}
	return m.DelCacheCtx(ctx, fmt.Sprintf("%s%v", cacheCouponIdPrefix, id))
}

// SoftDelete 逻辑删除，仅允许删除草稿或暂停状态的券，避免影响已发放数据。
func (m *customCouponModel) SoftDelete(ctx context.Context, id, updater int64) (int64, error) {
	query := fmt.Sprintf("update %s set `deleted` = 1, `updater` = ?, `update_time` = now() "+
		"where `id` = ? and `deleted` = 0 and `status` in (?, ?)", m.table)
	res, err := m.ExecNoCacheCtx(ctx, query, updater, id, CouponStatusDraft, CouponStatusPaused)
	if err != nil {
		return 0, err
	}
	_ = m.DelCacheCtx(ctx, fmt.Sprintf("%s%v", cacheCouponIdPrefix, id))
	return res.RowsAffected()
}
