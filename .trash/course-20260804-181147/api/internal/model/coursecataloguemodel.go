package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseCatalogueModel = (*customCourseCatalogueModel)(nil)

type (
	// CourseCatalogueModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseCatalogueModel.
	CourseCatalogueModel interface {
		courseCatalogueModel
		ListByIds(ctx context.Context, ids []int64) ([]*CourseCatalogue, error)
		FindById(ctx context.Context, id int64) (*CourseCatalogue, error)
		ListByMediaIds(ctx context.Context, mediaIds []int64) ([]*CourseCatalogue, error)
		ListByCourseId(ctx context.Context, courseId int64, withPractice bool) ([]*CourseCatalogue, error)
		SaveAll(ctx context.Context, list []*CourseCatalogue) error
		DeleteByCourseId(ctx context.Context, courseId int64, types []int64) error
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

// ListByIds 批量查询
func (m *customCourseCatalogueModel) ListByIds(ctx context.Context, ids []int64) ([]*CourseCatalogue, error) {
	if len(ids) == 0 {
		return []*CourseCatalogue{}, nil
	}
	ph := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		ph = append(ph, "?")
		args = append(args, id)
	}
	var rows []*CourseCatalogue
	query := fmt.Sprintf("select %s from %s where `deleted` = 0 and `id` in (%s)", courseCatalogueRows, m.table, strings.Join(ph, ","))
	err := m.QueryRowsNoCacheCtx(ctx, &rows, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// FindById 根据ID查询
func (m *customCourseCatalogueModel) FindById(ctx context.Context, id int64) (*CourseCatalogue, error) {
	return m.FindOne(ctx, id)
}

// ListByMediaIds 根据媒资ID列表查询
func (m *customCourseCatalogueModel) ListByMediaIds(ctx context.Context, mediaIds []int64) ([]*CourseCatalogue, error) {
	if len(mediaIds) == 0 {
		return []*CourseCatalogue{}, nil
	}
	ph := make([]string, 0, len(mediaIds))
	args := make([]any, 0, len(mediaIds))
	for _, id := range mediaIds {
		ph = append(ph, "?")
		args = append(args, id)
	}
	var rows []*CourseCatalogue
	query := fmt.Sprintf("select %s from %s where `deleted` = 0 and `media_id` in (%s)", courseCatalogueRows, m.table, strings.Join(ph, ","))
	err := m.QueryRowsNoCacheCtx(ctx, &rows, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// ListByCourseId 根据课程ID查询目录
func (m *customCourseCatalogueModel) ListByCourseId(ctx context.Context, courseId int64, withPractice bool) ([]*CourseCatalogue, error) {
	where := "`course_id` = ? and `deleted` = 0"
	args := []any{courseId}
	if !withPractice {
		where += " and `type` != 3"
	}
	query := fmt.Sprintf("select %s from %s where %s order by `parent_id` asc, `index` asc", courseCatalogueRows, m.table, where)
	var rows []*CourseCatalogue
	err := m.QueryRowsNoCacheCtx(ctx, &rows, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SaveAll 批量保存（用于上架/下架同步）
func (m *customCourseCatalogueModel) SaveAll(ctx context.Context, list []*CourseCatalogue) error {
	if len(list) == 0 {
		return nil
	}
	for _, c := range list {
		if c.Id > 0 {
			if err := m.Update(ctx, c); err != nil {
				return err
			}
		} else {
			if _, err := m.Insert(ctx, c); err != nil {
				return err
			}
		}
	}
	return nil
}

// DeleteByCourseId 删除课程的所有目录
func (m *customCourseCatalogueModel) DeleteByCourseId(ctx context.Context, courseId int64, types []int64) error {
	if len(types) == 0 {
		return nil
	}
	ph := make([]string, 0, len(types))
	args := make([]any, 0, len(types)+1)
	args = append(args, courseId)
	for _, t := range types {
		ph = append(ph, "?")
		args = append(args, t)
	}
	query := fmt.Sprintf("delete from %s where `course_id` = ? and `type` in (%s)", m.table, strings.Join(ph, ","))
	courseCatalogueIdKey := fmt.Sprintf("%s%v", cacheCourseCatalogueIdPrefix, courseId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, courseCatalogueIdKey)
	return err
}
