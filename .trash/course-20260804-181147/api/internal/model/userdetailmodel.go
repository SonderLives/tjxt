package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserDetailModel = (*customUserDetailModel)(nil)

type (
	// UserDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserDetailModel.
	UserDetailModel interface {
		userDetailModel
		FindByIds(ctx context.Context, ids []int64) (map[int64]*UserDetail, error)
	}

	customUserDetailModel struct {
		*defaultUserDetailModel
	}
)

// NewUserDetailModel returns a model for the database table.
func NewUserDetailModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserDetailModel {
	return &customUserDetailModel{
		defaultUserDetailModel: newUserDetailModel(conn, c, opts...),
	}
}

// FindByIds 批量查询用户详情，返回 map[id]*UserDetail
func (m *customUserDetailModel) FindByIds(ctx context.Context, ids []int64) (map[int64]*UserDetail, error) {
	if len(ids) == 0 {
		return map[int64]*UserDetail{}, nil
	}
	ph := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		ph = append(ph, "?")
		args = append(args, id)
	}
	var rows []*UserDetail
	query := fmt.Sprintf("select %s from %s where `id` in (%s)", userDetailRows, m.table, strings.Join(ph, ","))
	err := m.QueryRowsNoCacheCtx(ctx, &rows, query, args...)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]*UserDetail, len(rows))
	for _, row := range rows {
		result[row.Id] = row
	}
	return result, nil
}
