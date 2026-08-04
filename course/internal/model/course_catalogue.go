package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// CourseCatalogue 课程目录正式表（course_catalogue 表）
type CourseCatalogue struct {
	Id                int64
	Name              string
	Trailer           int64
	CourseId          int64
	Type              int64
	ParentCatalogueId int64
	MediaId           int64
	VideoId           sql.NullInt64
	VideoName         string
	LivingStartTime   sql.NullTime
	LivingEndTime     sql.NullTime
	PlayBack          int64
	MediaDuration     int64
	CIndex            int64
	DepId             int64
	CreateTime        time.Time
	UpdateTime        time.Time
	Creater           int64
	Updater           int64
	Deleted           int64
}

// CourseCatalogueModel 课程目录正式数据访问
type CourseCatalogueModel struct {
	conn  sqlx.SqlConn
	table string
}

// NewCourseCatalogueModel 创建课程目录正式数据访问对象
func NewCourseCatalogueModel(conn sqlx.SqlConn) *CourseCatalogueModel {
	return &CourseCatalogueModel{conn: conn, table: "course_catalogue"}
}

const courseCatalogueColumns = `id, name, trailer, course_id, type, parent_catalogue_id, media_id,
	video_id, video_name, living_start_time, living_end_time, play_back, media_duration, c_index,
	dep_id, create_time, update_time, creater, updater, deleted`

const courseCatalogueScanColumns = `id, name, trailer, course_id, type, parent_catalogue_id, media_id,
	video_id, video_name, living_start_time, living_end_time, play_back, media_duration, c_index,
	dep_id, create_time, update_time, COALESCE(creater, 0), COALESCE(updater, 0), deleted`

// FindById 按主键查询目录。
func (m *CourseCatalogueModel) FindById(ctx context.Context, id int64) (*CourseCatalogue, error) {
	var c CourseCatalogue
	err := m.conn.QueryRowCtx(ctx, &c, fmt.Sprintf(
		`SELECT %s FROM %s WHERE id = ? AND deleted = 0`, courseCatalogueScanColumns, m.table), id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &c, err
}

// ListByCourseId 按课程查询目录（按 type, c_index 排序）。
func (m *CourseCatalogueModel) ListByCourseId(ctx context.Context, courseId int64, withPractice bool) ([]*CourseCatalogue, error) {
	var rows []*CourseCatalogue
	where := "course_id = ? AND deleted = 0"
	if !withPractice {
		where += " AND type IN (1, 2)"
	}
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE %s ORDER BY type, c_index`,
		courseCatalogueScanColumns, m.table, where), courseId)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// ListByMediaIds 按媒资 id 列表查询目录（统计引用次数）。
func (m *CourseCatalogueModel) ListByMediaIds(ctx context.Context, mediaIds []int64) ([]*CourseCatalogue, error) {
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
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE deleted = 0 AND media_id IN (%s)`,
		courseCatalogueScanColumns, m.table, strings.Join(ph, ",")), args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// ListByIds 按 id 列表批量查询。
func (m *CourseCatalogueModel) ListByIds(ctx context.Context, ids []int64) ([]*CourseCatalogue, error) {
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
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE deleted = 0 AND id IN (%s)`,
		courseCatalogueScanColumns, m.table, strings.Join(ph, ",")), args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// DeleteByCourseId 逻辑删除课程目录。
func (m *CourseCatalogueModel) DeleteByCourseId(ctx context.Context, courseId int64) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET deleted = 1, update_time = NOW() WHERE course_id = ? AND deleted = 0`, m.table), courseId)
	return err
}

// SaveAll 批量插入/覆盖目录（上架时草稿→正式，全删重插）。
func (m *CourseCatalogueModel) SaveAll(ctx context.Context, list []*CourseCatalogue) error {
	if len(list) == 0 {
		return nil
	}
	// 事务中先逻辑删除再批量插入
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		courseId := list[0].CourseId
		if _, err := session.ExecCtx(ctx, fmt.Sprintf(
			`UPDATE %s SET deleted = 1, update_time = NOW() WHERE course_id = ? AND deleted = 0`,
			m.table), courseId); err != nil {
			return err
		}
		for _, c := range list {
			if _, err := session.ExecCtx(ctx, fmt.Sprintf(
				`INSERT INTO %s (id, name, trailer, course_id, type, parent_catalogue_id, media_id,
					video_id, video_name, living_start_time, living_end_time, play_back, media_duration,
					c_index, dep_id, create_time, update_time, creater, updater)
				 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), ?, ?)`,
				m.table),
				c.Id, c.Name, c.Trailer, c.CourseId, c.Type, c.ParentCatalogueId, c.MediaId,
				c.VideoId, c.VideoName, c.LivingStartTime, c.LivingEndTime, c.PlayBack, c.MediaDuration,
				c.CIndex, c.DepId, c.Creater, c.Updater); err != nil {
				return err
			}
		}
		return nil
	})
}
