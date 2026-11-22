package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseDraftModel = (*customCourseDraftModel)(nil)

type (
	// CourseDraftModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseDraftModel.
	CourseDraftModel interface {
		courseDraftModel
		FindByNameExceptId(ctx context.Context, name string, id int64) (*CourseDraft, error)
		CountByThirdCateId(ctx context.Context, thirdCateId int64) (int64, error)
	}
	customCourseDraftModel struct {
		*defaultCourseDraftModel
	}
)

// NewCourseDraftModel returns a model for the database table.
func NewCourseDraftModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseDraftModel {
	return &customCourseDraftModel{
		defaultCourseDraftModel: newCourseDraftModel(conn, c, opts...),
	}
}

func (m *customCourseDraftModel) FindByNameExceptId(ctx context.Context, name string, id int64) (*CourseDraft, error) {
	var list []*CourseDraft
	query := fmt.Sprintf("select * from %s where `deleted` = 0 and `name` = ? and `id` <> ? limit 1", m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, name, id); err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, ErrNotFound
	}
	return list[0], nil
}

func (m *customCourseDraftModel) CountByThirdCateId(ctx context.Context, thirdCateId int64) (int64, error) {
	var total int64
	query := fmt.Sprintf("select count(*) from %s where `deleted` = 0 and `third_cate_id` = ?", m.table)
	if err := m.QueryRowNoCacheCtx(ctx, &total, query, thirdCateId); err != nil {
		return 0, err
	}
	return total, nil
}
