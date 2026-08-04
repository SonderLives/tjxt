package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// CourseContentDraft 课程内容草稿表（course_content_draft 表）
type CourseContentDraft struct {
	Id               int64
	CourseIntroduce  string
	UsePeople        string
	CourseDetail     string
	DepId            int64
	CreateTime       time.Time
	UpdateTime       time.Time
	Creater          int64
	Updater          int64
	Deleted          int64
}

// CourseContentDraftModel 课程内容草稿数据访问
type CourseContentDraftModel struct {
	conn  sqlx.SqlConn
	table string
}

// NewCourseContentDraftModel 创建课程内容草稿数据访问对象
func NewCourseContentDraftModel(conn sqlx.SqlConn) *CourseContentDraftModel {
	return &CourseContentDraftModel{conn: conn, table: "course_content_draft"}
}

const courseContentDraftColumns = `id, course_introduce, use_people, course_detail, dep_id, create_time,
	update_time, creater, updater, deleted`

// FindById 按主键查询课程内容草稿。
func (m *CourseContentDraftModel) FindById(ctx context.Context, id int64) (*CourseContentDraft, error) {
	var c CourseContentDraft
	err := m.conn.QueryRowCtx(ctx, &c, fmt.Sprintf(
		`SELECT %s FROM %s WHERE id = ? AND deleted = 0`,
		courseContentDraftColumns, m.table), id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &c, err
}

// Insert 新增课程内容草稿。
func (m *CourseContentDraftModel) Insert(ctx context.Context, c *CourseContentDraft) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`INSERT INTO %s (id, course_introduce, use_people, course_detail, dep_id, create_time, update_time,
			creater, updater) VALUES (?, ?, ?, ?, ?, NOW(), NOW(), ?, ?)`, m.table),
		c.Id, c.CourseIntroduce, c.UsePeople, c.CourseDetail, c.DepId, c.Creater, c.Updater)
	return err
}

// UpdateById 更新课程内容草稿。
func (m *CourseContentDraftModel) UpdateById(ctx context.Context, c *CourseContentDraft) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET course_introduce = ?, use_people = ?, course_detail = ?, update_time = NOW()
		 WHERE id = ? AND deleted = 0`, m.table),
		c.CourseIntroduce, c.UsePeople, c.CourseDetail, c.Id)
	return err
}

// DeleteById 逻辑删除课程内容草稿。
func (m *CourseContentDraftModel) DeleteById(ctx context.Context, id int64) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET deleted = 1, update_time = NOW() WHERE id = ? AND deleted = 0`, m.table), id)
	return err
}
