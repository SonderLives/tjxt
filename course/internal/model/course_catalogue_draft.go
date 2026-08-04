package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// CourseCatalogueDraft 课程目录草稿表（course_catalogue_draft 表）
type CourseCatalogueDraft struct {
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
	CIndex            int64
	MediaDuration     int64
	CanUpdate         int64
	DepId             int64
	CreateTime        time.Time
	UpdateTime        time.Time
	Creater           int64
	Updater           int64
	Deleted           int64
}

// CourseCatalogueDraftModel 课程目录草稿数据访问
type CourseCatalogueDraftModel struct {
	conn  sqlx.SqlConn
	table string
}

// NewCourseCatalogueDraftModel 创建课程目录草稿数据访问对象
func NewCourseCatalogueDraftModel(conn sqlx.SqlConn) *CourseCatalogueDraftModel {
	return &CourseCatalogueDraftModel{conn: conn, table: "course_catalogue_draft"}
}

const courseCatalogueDraftColumns = `id, name, trailer, course_id, type, parent_catalogue_id, media_id,
	video_id, video_name, living_start_time, living_end_time, play_back, c_index, media_duration,
	can_update, dep_id, create_time, update_time, COALESCE(creater, 0), COALESCE(updater, 0), deleted`

// FindById 按主键查询草稿目录。
func (m *CourseCatalogueDraftModel) FindById(ctx context.Context, id int64) (*CourseCatalogueDraft, error) {
	var c CourseCatalogueDraft
	err := m.conn.QueryRowCtx(ctx, &c, fmt.Sprintf(
		`SELECT %s FROM %s WHERE id = ? AND deleted = 0`, courseCatalogueDraftColumns, m.table), id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &c, err
}

// ListByCourseId 按课程查询草稿目录（按 type, c_index 排序）。
func (m *CourseCatalogueDraftModel) ListByCourseId(ctx context.Context, courseId int64, withPractice bool) ([]*CourseCatalogueDraft, error) {
	var rows []*CourseCatalogueDraft
	where := "course_id = ? AND deleted = 0"
	if !withPractice {
		where += " AND type IN (1, 2)"
	}
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE %s ORDER BY type, c_index`,
		courseCatalogueDraftColumns, m.table, where), courseId)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// ListByTypes 按课程与目录类型查询。
func (m *CourseCatalogueDraftModel) ListByTypes(ctx context.Context, courseId int64, types []int64) ([]*CourseCatalogueDraft, error) {
	if len(types) == 0 {
		return []*CourseCatalogueDraft{}, nil
	}
	ph := make([]string, 0, len(types))
	args := make([]any, 0, len(types))
	for _, t := range types {
		ph = append(ph, "?")
		args = append(args, t)
	}
	args = append(args, courseId)
	var rows []*CourseCatalogueDraft
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE deleted = 0 AND type IN (%s) AND course_id = ? ORDER BY type, c_index`,
		courseCatalogueDraftColumns, m.table, strings.Join(ph, ",")), args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// ListAllByCourseId 查询课程全部草稿目录（含练习）。
func (m *CourseCatalogueDraftModel) ListAllByCourseId(ctx context.Context, courseId int64) ([]*CourseCatalogueDraft, error) {
	var rows []*CourseCatalogueDraft
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE course_id = ? AND deleted = 0 ORDER BY type, c_index`,
		courseCatalogueDraftColumns, m.table), courseId)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// DeleteByCourseId 逻辑删除课程草稿目录（可按类型过滤）。
func (m *CourseCatalogueDraftModel) DeleteByCourseId(ctx context.Context, courseId int64, types []int64) error {
	where := "course_id = ? AND deleted = 0"
	args := []any{courseId}
	if len(types) > 0 {
		ph := make([]string, 0, len(types))
		for _, t := range types {
			ph = append(ph, "?")
			args = append(args, t)
		}
		where += fmt.Sprintf(" AND type IN (%s)", strings.Join(ph, ","))
	}
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET deleted = 1, update_time = NOW() WHERE %s`, m.table, where), args...)
	return err
}

// ReplaceAll 全删重插（目录保存/上架时使用）。
func (m *CourseCatalogueDraftModel) ReplaceAll(ctx context.Context, courseId int64, list []*CourseCatalogueDraft) error {
	return m.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := session.ExecCtx(ctx, fmt.Sprintf(
			`UPDATE %s SET deleted = 1, update_time = NOW() WHERE course_id = ? AND deleted = 0`,
			m.table), courseId); err != nil {
			return err
		}
		for _, c := range list {
			if _, err := session.ExecCtx(ctx, fmt.Sprintf(
				`INSERT INTO %s (id, name, trailer, course_id, type, parent_catalogue_id, media_id,
					video_id, video_name, living_start_time, living_end_time, play_back, c_index,
					media_duration, can_update, dep_id, create_time, update_time, creater, updater)
				 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), ?, ?)`,
				m.table),
				c.Id, c.Name, c.Trailer, c.CourseId, c.Type, c.ParentCatalogueId, c.MediaId,
				c.VideoId, c.VideoName, c.LivingStartTime, c.LivingEndTime, c.PlayBack, c.CIndex,
				c.MediaDuration, c.CanUpdate, c.DepId, c.Creater, c.Updater); err != nil {
				return err
			}
		}
		return nil
	})
}

// SaveMediaInfo 批量更新小节媒资信息。
func (m *CourseCatalogueDraftModel) SaveMediaInfo(ctx context.Context, list []*CourseCatalogueDraft) error {
	for _, c := range list {
		if _, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
			`UPDATE %s SET media_id = ?, video_id = ?, video_name = ?, trailer = ?, media_duration = ?,
				update_time = NOW() WHERE id = ? AND deleted = 0`,
			m.table), c.MediaId, c.VideoId, c.VideoName, c.Trailer, c.MediaDuration, c.Id); err != nil {
			return err
		}
	}
	return nil
}

// UpdateMediaDuration 批量更新章级媒资总时长。
func (m *CourseCatalogueDraftModel) UpdateMediaDuration(ctx context.Context, items []*CourseCatalogueDraft) error {
	for _, c := range items {
		if _, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
			`UPDATE %s SET media_duration = ?, update_time = NOW() WHERE id = ? AND deleted = 0`,
			m.table), c.MediaDuration, c.Id); err != nil {
			return err
		}
	}
	return nil
}

// SumMediaDurationByCourse 统计课程小节媒资总时长（按章分组）。
func (m *CourseCatalogueDraftModel) SumMediaDurationByCourse(ctx context.Context, courseId int64) (map[int64]int64, error) {
	result := make(map[int64]int64)
	var rows []struct {
		ParentId int64 `db:"parent_catalogue_id"`
		Sum      int64 `db:"sum"`
	}
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT parent_catalogue_id, SUM(media_duration) AS sum FROM %s
		 WHERE deleted = 0 AND course_id = ? AND type = 2 GROUP BY parent_catalogue_id`,
		m.table), courseId)
	if err != nil {
		return nil, err
	}
	for _, r := range rows {
		result[r.ParentId] = r.Sum
	}
	return result, nil
}

// CountSections 统计课程小节与练习总数。
func (m *CourseCatalogueDraftModel) CountSections(ctx context.Context, courseId int64) (int64, error) {
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, fmt.Sprintf(
		`SELECT COUNT(*) FROM %s WHERE deleted = 0 AND course_id = ? AND type IN (2, 3)`, m.table), courseId)
	if err != nil {
		return 0, err
	}
	return total, nil
}
