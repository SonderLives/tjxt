package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// CourseContent 课程内容正式表（course_content 表）
type CourseContent struct {
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

// CourseContentModel 课程内容正式数据访问
type CourseContentModel struct {
	conn  sqlx.SqlConn
	table string
}

// NewCourseContentModel 创建课程内容正式数据访问对象
func NewCourseContentModel(conn sqlx.SqlConn) *CourseContentModel {
	return &CourseContentModel{conn: conn, table: "course_content"}
}

const courseContentColumns = `id, course_introduce, use_people, course_detail, dep_id, create_time,
	update_time, creater, updater, deleted`

// FindById 按主键查询课程内容。
func (m *CourseContentModel) FindById(ctx context.Context, id int64) (*CourseContent, error) {
	var c CourseContent
	err := m.conn.QueryRowCtx(ctx, &c, fmt.Sprintf(
		`SELECT %s FROM %s WHERE id = ? AND deleted = 0`,
		courseContentColumns, m.table), id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &c, err
}

// Insert 新增课程内容。
func (m *CourseContentModel) Insert(ctx context.Context, c *CourseContent) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`INSERT INTO %s (id, course_introduce, use_people, course_detail, dep_id, create_time, update_time,
			creater, updater) VALUES (?, ?, ?, ?, ?, NOW(), NOW(), ?, ?)`, m.table),
		c.Id, c.CourseIntroduce, c.UsePeople, c.CourseDetail, c.DepId, c.Creater, c.Updater)
	return err
}

// UpdateById 更新课程内容。
func (m *CourseContentModel) UpdateById(ctx context.Context, c *CourseContent) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET course_introduce = ?, use_people = ?, course_detail = ?, update_time = NOW()
		 WHERE id = ? AND deleted = 0`, m.table),
		c.CourseIntroduce, c.UsePeople, c.CourseDetail, c.Id)
	return err
}

// DeleteById 逻辑删除课程内容。
func (m *CourseContentModel) DeleteById(ctx context.Context, id int64) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET deleted = 1, update_time = NOW() WHERE id = ? AND deleted = 0`, m.table), id)
	return err
}
