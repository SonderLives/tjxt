package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LoginRecordModel = (*customLoginRecordModel)(nil)

type (
	// LoginRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLoginRecordModel.
	LoginRecordModel interface {
		loginRecordModel
		// FindPage 按用户分页查询登录记录，userId <= 0 表示查询全部用户。
		FindPage(ctx context.Context, userId, offset, limit int64) ([]*LoginRecord, int64, error)
		// MarkLogout 标记最近一次未登出的记录为已登出，并回填登录时长（秒）。
		MarkLogout(ctx context.Context, userId int64) error
	}

	customLoginRecordModel struct {
		*defaultLoginRecordModel
	}
)

// NewLoginRecordModel returns a model for the database table.
func NewLoginRecordModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LoginRecordModel {
	return &customLoginRecordModel{
		defaultLoginRecordModel: newLoginRecordModel(conn, c, opts...),
	}
}

// FindPage 分页查询登录记录，按登录时间倒序。
func (m *customLoginRecordModel) FindPage(ctx context.Context, userId, offset, limit int64) ([]*LoginRecord, int64, error) {
	where := "where 1 = 1"
	args := make([]any, 0, 3)
	if userId > 0 {
		where = "where `user_id` = ?"
		args = append(args, userId)
	}

	query := fmt.Sprintf("select %s from %s %s order by `login_time` desc limit ?, ?", loginRecordRows, m.table, where)
	var list []*LoginRecord
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, append(append([]any{}, args...), offset, limit)...); err != nil {
		return nil, 0, err
	}

	var total int64
	countQuery := fmt.Sprintf("select count(*) from %s %s", m.table, where)
	if err := m.QueryRowNoCacheCtx(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// MarkLogout 结算用户最近一次登录会话。
func (m *customLoginRecordModel) MarkLogout(ctx context.Context, userId int64) error {
	var id int64
	query := fmt.Sprintf("select `id` from %s where `user_id` = ? and `logout_time` is null order by `login_time` desc limit 1", m.table)
	if err := m.QueryRowNoCacheCtx(ctx, &id, query, userId); err != nil {
		if err == ErrNotFound {
			return nil
		}
		return err
	}

	update := fmt.Sprintf("update %s set `logout_time` = now(), `duration` = timestampdiff(second, `login_time`, now()) where `id` = ?", m.table)
	if _, err := m.ExecNoCacheCtx(ctx, update, id); err != nil {
		return err
	}
	return m.DelCacheCtx(ctx, fmt.Sprintf("%s%v", cacheLoginRecordIdPrefix, id))
}
