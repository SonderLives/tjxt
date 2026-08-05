package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseCatalogueModel = (*customCourseCatalogueModel)(nil)

// MediaUseCount 媒资在课程目录中的引用计数（用于 CourseMediaUseInfo）。
type MediaUseCount struct {
	MediaId  int64 `db:"media_id"`
	QuoteNum int64 `db:"quote_num"`
}

type (
	// CourseCatalogueModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseCatalogueModel.
	CourseCatalogueModel interface {
		courseCatalogueModel
		ListByCourseId(ctx context.Context, courseId int64) ([]*CourseCatalogue, error)
		ListByIdIn(ctx context.Context, ids []int64) ([]*CourseCatalogue, error)
		CountByMediaIds(ctx context.Context, mediaIds []int64) ([]*MediaUseCount, error)
	}
	customCourseCatalogueModel struct {
		*defaultCourseCatalogueModel
	}
)

// NewCourseCatalogueModel returns a model for the database table.
func NewCourseCatalogueModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseCatalogueModel {
	return &customCourseCatalogueModel{
		defaultCourseCatalogueModel: newCourseCatalogueModel(conn, c, opts...),
	}
}

func (m *customCourseCatalogueModel) ListByCourseId(ctx context.Context, courseId int64) ([]*CourseCatalogue, error) {
	var list []*CourseCatalogue
	query := fmt.Sprintf("select * from %s where `deleted` = 0 and `course_id` = ? order by `c_index` asc", m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, courseId); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCourseCatalogueModel) ListByIdIn(ctx context.Context, ids []int64) ([]*CourseCatalogue, error) {
	if len(ids) == 0 {
		return []*CourseCatalogue{}, nil
	}
	var list []*CourseCatalogue
	ph := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	query := fmt.Sprintf("select * from %s where `deleted` = 0 and `id` in (%s) order by `c_index` asc", m.table, ph)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCourseCatalogueModel) CountByMediaIds(ctx context.Context, mediaIds []int64) ([]*MediaUseCount, error) {
	if len(mediaIds) == 0 {
		return []*MediaUseCount{}, nil
	}
	var list []*MediaUseCount
	ph := strings.TrimSuffix(strings.Repeat("?,", len(mediaIds)), ",")
	args := make([]any, 0, len(mediaIds))
	for _, id := range mediaIds {
		args = append(args, id)
	}
	query := fmt.Sprintf("select `media_id`, count(*) as `quote_num` from %s where `deleted` = 0 and `media_id` in (%s) group by `media_id`", m.table, ph)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}
