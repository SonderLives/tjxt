package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseTeacherModel = (*customCourseTeacherModel)(nil)

type (
	// CourseTeacherModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseTeacherModel.
	CourseTeacherModel interface {
		courseTeacherModel
		ListByCourseId(ctx context.Context, courseId int64) ([]*CourseTeacher, error)
		ListByTeacherIds(ctx context.Context, teacherIds []int64) ([]*CourseTeacher, error)
		DeleteByCourseId(ctx context.Context, courseId int64) error
		DeleteByIdsTx(ctx context.Context, session sqlx.Session, ids []int64) error
		BatchInsertTx(ctx context.Context, session sqlx.Session, list []*CourseTeacher) error
	}
	customCourseTeacherModel struct {
		*defaultCourseTeacherModel
	}
)

// NewCourseTeacherModel returns a model for the database table.
func NewCourseTeacherModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseTeacherModel {
	return &customCourseTeacherModel{
		defaultCourseTeacherModel: newCourseTeacherModel(conn, c, opts...),
	}
}

func (m *customCourseTeacherModel) ListByCourseId(ctx context.Context, courseId int64) ([]*CourseTeacher, error) {
	var list []*CourseTeacher
	query := fmt.Sprintf("select * from %s where `deleted` = 0 and `course_id` = ? order by `c_index` asc", m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, courseId); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCourseTeacherModel) ListByTeacherIds(ctx context.Context, teacherIds []int64) ([]*CourseTeacher, error) {
	if len(teacherIds) == 0 {
		return []*CourseTeacher{}, nil
	}
	var list []*CourseTeacher
	ph := strings.TrimSuffix(strings.Repeat("?,", len(teacherIds)), ",")
	args := make([]any, 0, len(teacherIds))
	for _, id := range teacherIds {
		args = append(args, id)
	}
	query := fmt.Sprintf("select * from %s where `deleted` = 0 and `teacher_id` in (%s)", m.table, ph)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCourseTeacherModel) DeleteByCourseId(ctx context.Context, courseId int64) error {
	query := fmt.Sprintf("delete from %s where `course_id` = ?", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, courseId)
	return err
}

// DeleteByIdsTx 事务内按 id 批量物理删除老师行（一条 DELETE ... WHERE id IN），并失效对应行缓存，
// 供跨模型事务（如课程上架发布时清理旧老师）复用；ids 为空时为空操作。
func (m *customCourseTeacherModel) DeleteByIdsTx(ctx context.Context, session sqlx.Session, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	query := fmt.Sprintf("delete from %s where `id` in (%s)", m.table, strings.TrimSuffix(strings.Repeat("?,", len(ids)), ","))
	args := make([]any, 0, len(ids))
	keys := make([]string, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
		keys = append(keys, fmt.Sprintf("%s%v", cacheCourseTeacherIdPrefix, id))
	}
	cc := m.WithSession(session)
	_, err := cc.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, keys...)
	return err
}

// BatchInsertTx 事务内批量插入老师行（单条 INSERT 多 VALUES，列与生成模型 Insert 一致），
// 并逐行失效行缓存，供跨模型事务（如课程上架发布）复用；list 为空时为空操作。
func (m *customCourseTeacherModel) BatchInsertTx(ctx context.Context, session sqlx.Session, list []*CourseTeacher) error {
	if len(list) == 0 {
		return nil
	}
	cols := insertColumns(courseTeacherRowsExpectAutoSet)
	values := make([]string, 0, len(list))
	args := make([]any, 0, len(list)*cols)
	keys := make([]string, 0, len(list))
	for _, d := range list {
		values = append(values, rowPlaceholders(cols))
		args = append(args, d.Id, d.CourseId, d.TeacherId, d.IsShow, d.CIndex, d.DepId,
			d.Creater, d.Updater, d.Deleted)
		keys = append(keys, fmt.Sprintf("%s%v", cacheCourseTeacherIdPrefix, d.Id))
	}
	query := fmt.Sprintf("insert into %s (%s) values %s", m.table, courseTeacherRowsExpectAutoSet, strings.Join(values, ","))
	cc := m.WithSession(session)
	_, err := cc.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, keys...)
	return err
}
