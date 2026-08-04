package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// CourseCataSubjectDraft 课程-题目关系草稿表（course_cata_subject_draft 表）
type CourseCataSubjectDraft struct {
	Id        int64
	CourseId  int64
	CataId    int64
	SubjectId int64
}

// CourseCataSubjectDraftModel 课程-题目关系草稿数据访问
type CourseCataSubjectDraftModel struct {
	conn  sqlx.SqlConn
	table string
}

// NewCourseCataSubjectDraftModel 创建课程-题目关系草稿数据访问对象
func NewCourseCataSubjectDraftModel(conn sqlx.SqlConn) *CourseCataSubjectDraftModel {
	return &CourseCataSubjectDraftModel{conn: conn, table: "course_cata_subject_draft"}
}

const courseCataSubjectDraftColumns = `id, course_id, cata_id, subject_id`

// ListByCourseId 按课程查询题目关系。
func (m *CourseCataSubjectDraftModel) ListByCourseId(ctx context.Context, courseId int64) ([]*CourseCataSubjectDraft, error) {
	var rows []*CourseCataSubjectDraft
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE course_id = ?`, courseCataSubjectDraftColumns, m.table), courseId)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Insert 新增关系。
func (m *CourseCataSubjectDraftModel) Insert(ctx context.Context, c *CourseCataSubjectDraft) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`INSERT INTO %s (id, course_id, cata_id, subject_id) VALUES (?, ?, ?, ?)`, m.table),
		c.Id, c.CourseId, c.CataId, c.SubjectId)
	return err
}

// DeleteByCourseId 删除课程全部题目关系。
func (m *CourseCataSubjectDraftModel) DeleteByCourseId(ctx context.Context, courseId int64) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`DELETE FROM %s WHERE course_id = ?`, m.table), courseId)
	return err
}

// DeleteNotInCataIds 删除不在指定目录 id 列表中的关系。
func (m *CourseCataSubjectDraftModel) DeleteNotInCataIds(ctx context.Context, courseId int64, cataIds []int64) error {
	if len(cataIds) == 0 {
		return m.DeleteByCourseId(ctx, courseId)
	}
	ph := make([]string, 0, len(cataIds))
	args := make([]any, 0, len(cataIds))
	for _, id := range cataIds {
		ph = append(ph, "?")
		args = append(args, id)
	}
	args = append(args, courseId)
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`DELETE FROM %s WHERE course_id = ? AND cata_id NOT IN (%s)`,
		m.table, strings.Join(ph, ",")), args...)
	return err
}
