package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseCataSubjectDraftModel = (*customCourseCataSubjectDraftModel)(nil)

type (
	// CourseCataSubjectDraftModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseCataSubjectDraftModel.
	CourseCataSubjectDraftModel interface {
		courseCataSubjectDraftModel
		ListByCourseId(ctx context.Context, courseId int64) ([]*CourseCataSubjectDraft, error)
		DeleteNotInCataIds(ctx context.Context, courseId int64, cataIds []int64) error
		DeleteByCourseId(ctx context.Context, courseId int64) error
	}

	customCourseCataSubjectDraftModel struct {
		*defaultCourseCataSubjectDraftModel
	}
)

// NewCourseCataSubjectDraftModel returns a model for the database table.
func NewCourseCataSubjectDraftModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseCataSubjectDraftModel {
	return &customCourseCataSubjectDraftModel{
		defaultCourseCataSubjectDraftModel: newCourseCataSubjectDraftModel(conn, c, opts...),
	}
}

// ListByCourseId 根据课程ID查询题目关系
func (m *customCourseCataSubjectDraftModel) ListByCourseId(ctx context.Context, courseId int64) ([]*CourseCataSubjectDraft, error) {
	var rows []*CourseCataSubjectDraft
	query := fmt.Sprintf("select %s from %s where `course_id` = ?", courseCataSubjectDraftRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &rows, query, courseId)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// DeleteNotInCataIds 删除不在指定目录ID列表中的题目关系
func (m *customCourseCataSubjectDraftModel) DeleteNotInCataIds(ctx context.Context, courseId int64, cataIds []int64) error {
	if len(cataIds) == 0 {
		return m.DeleteByCourseId(ctx, courseId)
	}
	ph := make([]string, 0, len(cataIds))
	args := make([]any, 0, len(cataIds)+1)
	args = append(args, courseId)
	for _, id := range cataIds {
		ph = append(ph, "?")
		args = append(args, id)
	}
	query := fmt.Sprintf("delete from %s where `course_id` = ? and `cata_id` not in (%s)", m.table, strings.Join(ph, ","))
	courseCataSubjectDraftIdKey := fmt.Sprintf("%s%v", cacheCourseCataSubjectDraftIdPrefix, courseId)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, courseCataSubjectDraftIdKey)
	return err
}

// DeleteByCourseId 删除课程的所有题目关系
func (m *customCourseCataSubjectDraftModel) DeleteByCourseId(ctx context.Context, courseId int64) error {
	courseCataSubjectDraftIdKey := fmt.Sprintf("%s%v", cacheCourseCataSubjectDraftIdPrefix, courseId)
	query := fmt.Sprintf("delete from %s where `course_id` = ?", m.table)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, courseId)
	}, courseCataSubjectDraftIdKey)
	return err
}
