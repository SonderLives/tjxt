package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseCatalogueDraftModel = (*customCourseCatalogueDraftModel)(nil)

type (
	// CourseCatalogueDraftModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseCatalogueDraftModel.
	CourseCatalogueDraftModel interface {
		courseCatalogueDraftModel
		ListAllByCourseId(ctx context.Context, courseId int64) ([]*CourseCatalogueDraft, error)
		ListByCourseId(ctx context.Context, courseId int64, withPractice bool) ([]*CourseCatalogueDraft, error)
		CountSections(ctx context.Context, courseId int64) (int64, error)
		ListByTypes(ctx context.Context, courseId int64, types []int64) ([]*CourseCatalogueDraft, error)
		SumMediaDurationByCourse(ctx context.Context, courseId int64) (int64, error)
		SumMediaDurationByChapter(ctx context.Context, courseId int64) (map[int64]int64, error)
		SaveMediaInfo(ctx context.Context, items []*CourseCatalogueDraft) error
		UpdateMediaDuration(ctx context.Context, updates map[int64]int64) error
		ReplaceAll(ctx context.Context, courseId int64, list []*CourseCatalogueDraft) error
		DeleteByCourseId(ctx context.Context, courseId int64, types []int64) error
	}

	customCourseCatalogueDraftModel struct {
		*defaultCourseCatalogueDraftModel
	}
)

// NewCourseCatalogueDraftModel returns a model for the database table.
func NewCourseCatalogueDraftModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseCatalogueDraftModel {
	return &customCourseCatalogueDraftModel{
		defaultCourseCatalogueDraftModel: newCourseCatalogueDraftModel(conn, c, opts...),
	}
}

// ListAllByCourseId 查询课程所有目录草稿
func (m *customCourseCatalogueDraftModel) ListAllByCourseId(ctx context.Context, courseId int64) ([]*CourseCatalogueDraft, error) {
	var rows []*CourseCatalogueDraft
	query := fmt.Sprintf("select %s from %s where `course_id` = ? and `deleted` = 0 order by `parent_id` asc, `index` asc", courseCatalogueDraftRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &rows, query, courseId)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// ListByCourseId 查询课程目录（可选是否包含练习）
func (m *customCourseCatalogueDraftModel) ListByCourseId(ctx context.Context, courseId int64, withPractice bool) ([]*CourseCatalogueDraft, error) {
	where := "`course_id` = ? and `deleted` = 0"
	args := []any{courseId}
	if !withPractice {
		where += " and `type` != 3"
	}
	query := fmt.Sprintf("select %s from %s where %s order by `parent_id` asc, `index` asc", courseCatalogueDraftRows, m.table, where)
	var rows []*CourseCatalogueDraft
	err := m.QueryRowsNoCacheCtx(ctx, &rows, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// CountSections 统计小节数量
func (m *customCourseCatalogueDraftModel) CountSections(ctx context.Context, courseId int64) (int64, error) {
	query := fmt.Sprintf("select count(1) from %s where `course_id` = ? and `type` = 2 and `deleted` = 0", m.table)
	var count int64
	err := m.QueryRowNoCacheCtx(ctx, &count, query, courseId)
	return count, err
}

// ListByTypes 根据类型列表查询
func (m *customCourseCatalogueDraftModel) ListByTypes(ctx context.Context, courseId int64, types []int64) ([]*CourseCatalogueDraft, error) {
	if len(types) == 0 {
		return []*CourseCatalogueDraft{}, nil
	}
	ph := make([]string, 0, len(types))
	args := make([]any, 0, len(types)+1)
	args = append(args, courseId)
	for _, t := range types {
		ph = append(ph, "?")
		args = append(args, t)
	}
	query := fmt.Sprintf("select %s from %s where `course_id` = ? and `type` in (%s) and `deleted` = 0 order by `parent_id` asc, `index` asc", courseCatalogueDraftRows, m.table, strings.Join(ph, ","))
	var rows []*CourseCatalogueDraft
	err := m.QueryRowsNoCacheCtx(ctx, &rows, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SumMediaDurationByCourse 汇总课程媒资时长
func (m *customCourseCatalogueDraftModel) SumMediaDurationByCourse(ctx context.Context, courseId int64) (int64, error) {
	query := fmt.Sprintf("select coalesce(sum(`media_duration`), 0) from %s where `course_id` = ? and `deleted` = 0", m.table)
	var total int64
	err := m.QueryRowNoCacheCtx(ctx, &total, query, courseId)
	return total, err
}

// SumMediaDurationByChapter 汇总每章的媒资时长
func (m *customCourseCatalogueDraftModel) SumMediaDurationByChapter(ctx context.Context, courseId int64) (map[int64]int64, error) {
	query := fmt.Sprintf("select `parent_id`, coalesce(sum(`media_duration`), 0) from %s where `course_id` = ? and `deleted` = 0 and `type` = 2 group by `parent_id`", m.table)
	var rows []struct {
		ParentId  int64 `db:"parent_id"`
		Duration  int64 `db:"duration"`
	}
	err := m.QueryRowsNoCacheCtx(ctx, &rows, query, courseId)
	if err != nil {
		return nil, err
	}
	result := make(map[int64]int64)
	for _, row := range rows {
		result[row.ParentId] = row.Duration
	}
	return result, nil
}

// SaveMediaInfo 批量保存媒资信息
func (m *customCourseCatalogueDraftModel) SaveMediaInfo(ctx context.Context, items []*CourseCatalogueDraft) error {
	if len(items) == 0 {
		return nil
	}
	for _, item := range items {
		if err := m.Update(ctx, item); err != nil {
			return err
		}
	}
	return nil
}

// UpdateMediaDuration 批量更新媒资时长
func (m *customCourseCatalogueDraftModel) UpdateMediaDuration(ctx context.Context, updates map[int64]int64) error {
	if len(updates) == 0 {
		return nil
	}
	for id, duration := range updates {
		courseCatalogueDraftIdKey := fmt.Sprintf("%s%v", cacheCourseCatalogueDraftIdPrefix, id)
		query := fmt.Sprintf("update %s set `media_duration` = ? where `id` = ?", m.table)
		_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
			return conn.ExecCtx(ctx, query, duration, id)
		}, courseCatalogueDraftIdKey)
		if err != nil {
			return err
		}
	}
	return nil
}

// ReplaceAll 全量替换课程目录
func (m *customCourseCatalogueDraftModel) ReplaceAll(ctx context.Context, courseId int64, list []*CourseCatalogueDraft) error {
	// 先删除旧数据
	if err := m.DeleteByCourseId(ctx, courseId, []int64{1, 2, 3}); err != nil {
		return err
	}
	// 批量插入新数据
	for _, c := range list {
		if _, err := m.Insert(ctx, c); err != nil {
			return err
		}
	}
	return nil
}

// DeleteByCourseId 删除课程的所有目录
func (m *customCourseCatalogueDraftModel) DeleteByCourseId(ctx context.Context, courseId int64, types []int64) error {
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
	courseCatalogueDraftIdKey := fmt.Sprintf("%s%v", cacheCourseCatalogueDraftIdPrefix, courseId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, courseCatalogueDraftIdKey)
	return err
}
