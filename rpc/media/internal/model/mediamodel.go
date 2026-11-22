package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MediaModel = (*customMediaModel)(nil)

type (
	// MediaModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMediaModel.
	MediaModel interface {
		mediaModel

		// SoftDelete 逻辑删除（deleted = 1），幂等
		SoftDelete(ctx context.Context, id int64) error
		// FindOneNotDeleted 查询未删除的媒资
		FindOneNotDeleted(ctx context.Context, id int64) (*Media, error)
		// FindPage 分页查询媒资（deleted = 0），支持名称模糊与排序
		FindPage(ctx context.Context, name, sortBy, isAsc string, offset, limit int64) ([]*Media, error)
		// Count 统计媒资总数（deleted = 0）
		Count(ctx context.Context, name string) (int64, error)
	}

	customMediaModel struct {
		*defaultMediaModel
	}
)

// NewMediaModel returns a model for the database table.
func NewMediaModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) MediaModel {
	return &customMediaModel{
		defaultMediaModel: newMediaModel(conn, c, opts...),
	}
}

// mediaSortField sortBy 白名单映射，非法值由调用方兜底为 create_time
var mediaSortField = map[string]string{
	"createTime": "create_time",
	"duration":   "duration",
	"size":       "size",
}

// SoftDelete 逻辑删除：将 deleted 置为 1，并清理缓存
func (m *customMediaModel) SoftDelete(ctx context.Context, id int64) error {
	mediaIdKey := fmt.Sprintf("%s%v", cacheMediaIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `deleted` = 1 where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, mediaIdKey)
	return err
}

// FindOneNotDeleted 按 id 查询未删除的媒资
func (m *customMediaModel) FindOneNotDeleted(ctx context.Context, id int64) (*Media, error) {
	mediaIdKey := fmt.Sprintf("%s%v", cacheMediaIdPrefix, id)
	var resp Media
	err := m.QueryRowCtx(ctx, &resp, mediaIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? and `deleted` = 0 limit 1", mediaRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindPage 分页查询媒资列表，deleted = 0，名称 LIKE 模糊匹配
func (m *customMediaModel) FindPage(ctx context.Context, name, sortBy, isAsc string, offset, limit int64) ([]*Media, error) {
	sortField := mediaSortField[sortBy]
	if sortField == "" {
		sortField = "create_time"
	}
	order := "desc"
	if isAsc == "asc" || isAsc == "ASC" || isAsc == "1" || isAsc == "true" {
		order = "asc"
	}

	var (
		query = fmt.Sprintf("select %s from %s where `deleted` = 0", mediaRows, m.table)
		args  []any
	)
	if name != "" {
		query += " and `filename` like ?"
		args = append(args, "%"+name+"%")
	}
	query += fmt.Sprintf(" order by `%s` %s limit ? offset ?", sortField, order)
	args = append(args, limit, offset)

	var list []*Media
	err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...)
	return list, err
}

// Count 统计未删除的媒资总数
func (m *customMediaModel) Count(ctx context.Context, name string) (int64, error) {
	var (
		query = "select count(1) from " + m.table + " where `deleted` = 0"
		args  []any
		total int64
	)
	if name != "" {
		query += " and `filename` like ?"
		args = append(args, "%"+name+"%")
	}
	err := m.QueryRowNoCacheCtx(ctx, &total, query, args...)
	return total, err
}
