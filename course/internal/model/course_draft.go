package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// CourseDraft 课程草稿表（course_draft 表）
type CourseDraft struct {
	Id                int64
	Name              string
	CourseType        int64
	CoverUrl          string
	FirstCateId       int64
	SecondCateId      int64
	ThirdCateId       int64
	Free              int64
	Price             int64
	TemplateType      int64
	TemplateUrl       string
	Status            int64
	PurchaseStartTime sql.NullTime
	PurchaseEndTime   time.Time
	Step              int64
	Score             sql.NullInt64
	MediaDuration     int64
	ValidDuration     int64
	SectionNum        int64
	CanUpdate         int64
	DepId             int64
	PublishTime       sql.NullTime
	CreateTime        time.Time
	UpdateTime        time.Time
	Creater           int64
	Updater           int64
	Deleted           int64
	CVersion          sql.NullInt64
}

// CourseDraftModel 课程草稿数据访问
type CourseDraftModel struct {
	conn  sqlx.SqlConn
	table string
}

// NewCourseDraftModel 创建课程草稿数据访问对象
func NewCourseDraftModel(conn sqlx.SqlConn) *CourseDraftModel {
	return &CourseDraftModel{conn: conn, table: "course_draft"}
}

const courseDraftColumns = `id, name, course_type, cover_url, first_cate_id, second_cate_id, third_cate_id,
	free, price, template_type, template_url, status, purchase_start_time, purchase_end_time, step,
	score, media_duration, valid_duration, section_num, can_update, dep_id, publish_time,
	create_time, update_time, COALESCE(creater, 0), COALESCE(updater, 0), deleted, c_version`

// FindById 按主键查询草稿。
func (m *CourseDraftModel) FindById(ctx context.Context, id int64) (*CourseDraft, error) {
	var c CourseDraft
	err := m.conn.QueryRowCtx(ctx, &c, fmt.Sprintf(
		`SELECT %s FROM %s WHERE id = ? AND deleted = 0`, courseDraftColumns, m.table), id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &c, err
}

// Insert 新增草稿。
func (m *CourseDraftModel) Insert(ctx context.Context, c *CourseDraft) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`INSERT INTO %s (id, name, course_type, cover_url, first_cate_id, second_cate_id, third_cate_id,
			free, price, template_type, template_url, status, purchase_start_time, purchase_end_time,
			step, score, media_duration, valid_duration, section_num, can_update, dep_id, publish_time,
			create_time, update_time, creater, updater, c_version)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), ?, ?, ?)`,
		m.table),
		c.Id, c.Name, c.CourseType, c.CoverUrl, c.FirstCateId, c.SecondCateId, c.ThirdCateId,
		c.Free, c.Price, c.TemplateType, c.TemplateUrl, c.Status, c.PurchaseStartTime, c.PurchaseEndTime,
		c.Step, c.Score, c.MediaDuration, c.ValidDuration, c.SectionNum, c.CanUpdate, c.DepId, c.PublishTime,
		c.Creater, c.Updater, c.CVersion)
	return err
}

// UpdateById 更新草稿指定字段。
func (m *CourseDraftModel) UpdateById(ctx context.Context, c *CourseDraft, fields ...string) error {
	if len(fields) == 0 {
		return nil
	}
	sets := make([]string, 0, len(fields))
	args := make([]any, 0, len(fields))
	for _, f := range fields {
		switch f {
		case "name":
			sets = append(sets, "name = ?")
			args = append(args, c.Name)
		case "cover_url":
			sets = append(sets, "cover_url = ?")
			args = append(args, c.CoverUrl)
		case "price":
			sets = append(sets, "price = ?")
			args = append(args, c.Price)
		case "free":
			sets = append(sets, "free = ?")
			args = append(args, c.Free)
		case "first_cate_id":
			sets = append(sets, "first_cate_id = ?")
			args = append(args, c.FirstCateId)
		case "second_cate_id":
			sets = append(sets, "second_cate_id = ?")
			args = append(args, c.SecondCateId)
		case "third_cate_id":
			sets = append(sets, "third_cate_id = ?")
			args = append(args, c.ThirdCateId)
		case "course_type":
			sets = append(sets, "course_type = ?")
			args = append(args, c.CourseType)
		case "status":
			sets = append(sets, "status = ?")
			args = append(args, c.Status)
		case "step":
			sets = append(sets, "step = ?")
			args = append(args, c.Step)
		case "section_num":
			sets = append(sets, "section_num = ?")
			args = append(args, c.SectionNum)
		case "c_version":
			sets = append(sets, "c_version = ?")
			args = append(args, c.CVersion)
		case "purchase_end_time":
			sets = append(sets, "purchase_end_time = ?")
			args = append(args, c.PurchaseEndTime)
		case "purchase_start_time":
			sets = append(sets, "purchase_start_time = ?")
			args = append(args, c.PurchaseStartTime)
		case "can_update":
			sets = append(sets, "can_update = ?")
			args = append(args, c.CanUpdate)
		case "media_duration":
			sets = append(sets, "media_duration = ?")
			args = append(args, c.MediaDuration)
		case "valid_duration":
			sets = append(sets, "valid_duration = ?")
			args = append(args, c.ValidDuration)
		}
	}
	if len(sets) == 0 {
		return nil
	}
	sets = append(sets, "update_time = NOW()")
	args = append(args, c.Id)
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET %s WHERE id = ? AND deleted = 0`, m.table, strings.Join(sets, ", ")), args...)
	return err
}

// DeleteById 逻辑删除草稿。
func (m *CourseDraftModel) DeleteById(ctx context.Context, id int64) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET deleted = 1, update_time = NOW() WHERE id = ? AND deleted = 0`, m.table), id)
	return err
}

// CountByName 统计同名草稿数量，id>0 时排除自身。
func (m *CourseDraftModel) CountByName(ctx context.Context, name string, excludeId int64) (int64, error) {
	where := "deleted = 0 AND name = ?"
	args := []any{name}
	if excludeId > 0 {
		where += " AND id != ?"
		args = append(args, excludeId)
	}
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, fmt.Sprintf(
		`SELECT COUNT(*) FROM %s WHERE %s`, m.table, where), args...)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// PageCond 草稿课程分页查询条件
type CourseDraftPageCond struct {
	Keyword     string
	FirstCateId int64
	SecondCateId int64
	ThirdCateId int64
	CourseType  int64
	Status      int64
	BeginTime   string
	EndTime     string
	OrderBy     string
	Offset      int64
	Limit       int64
}

// Page 分页查询草稿课程。
func (m *CourseDraftModel) Page(ctx context.Context, cond *CourseDraftPageCond) ([]*CourseDraft, int64, error) {
	where := []string{"deleted = 0"}
	args := []any{}
	if cond.Keyword != "" {
		where = append(where, "name LIKE ?")
		args = append(args, "%"+cond.Keyword+"%")
	}
	if cond.FirstCateId > 0 {
		where = append(where, "first_cate_id = ?")
		args = append(args, cond.FirstCateId)
	}
	if cond.SecondCateId > 0 {
		where = append(where, "second_cate_id = ?")
		args = append(args, cond.SecondCateId)
	}
	if cond.ThirdCateId > 0 {
		where = append(where, "third_cate_id = ?")
		args = append(args, cond.ThirdCateId)
	}
	if cond.CourseType > 0 {
		where = append(where, "course_type = ?")
		args = append(args, cond.CourseType)
	}
	if cond.Status > 0 {
		where = append(where, "status = ?")
		args = append(args, cond.Status)
	}
	if cond.BeginTime != "" && cond.EndTime != "" {
		where = append(where, "update_time BETWEEN ? AND ?")
		args = append(args, cond.BeginTime, cond.EndTime)
	}
	whereSQL := strings.Join(where, " AND ")

	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, fmt.Sprintf(
		`SELECT COUNT(*) FROM %s WHERE %s`, m.table, whereSQL), args...)
	if err != nil {
		return nil, 0, err
	}

	order := cond.OrderBy
	if order == "" {
		order = "update_time DESC"
	}
	args = append(args, cond.Limit, cond.Offset)
	var rows []*CourseDraft
	err = m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE %s ORDER BY %s LIMIT ? OFFSET ?`,
		courseDraftColumns, m.table, whereSQL, order), args...)
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// ListAll 查询全部草稿课程。
func (m *CourseDraftModel) ListAll(ctx context.Context) ([]*CourseDraft, error) {
	var rows []*CourseDraft
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE deleted = 0`, courseDraftColumns, m.table))
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// CountTeacherCourse 统计老师名下待上架课程数量（course_teacher_draft 联表）。
func (m *CourseDraftModel) CountTeacherCourse(ctx context.Context, teacherIds []int64) (map[int64]int64, error) {
	result := make(map[int64]int64)
	if len(teacherIds) == 0 {
		return result, nil
	}
	ph := make([]string, 0, len(teacherIds))
	for range teacherIds {
		ph = append(ph, "?")
	}
	args := make([]any, 0, len(teacherIds))
	for _, id := range teacherIds {
		args = append(args, id)
	}
	var list []struct {
		TeacherId int64 `db:"teacher_id"`
		Num       int64 `db:"num"`
	}
	err := m.conn.QueryRowsCtx(ctx, &list, fmt.Sprintf(
		`SELECT ct.teacher_id AS teacher_id, COUNT(*) AS num
		 FROM course_teacher_draft ct
		 JOIN %s cd ON cd.id = ct.course_id AND cd.deleted = 0
		 WHERE ct.teacher_id IN (%s) AND ct.deleted = 0
		 GROUP BY ct.teacher_id`, m.table, strings.Join(ph, ",")),
		args...)
	if err != nil {
		return result, err
	}
	for _, r := range list {
		result[r.TeacherId] = r.Num
	}
	return result, nil
}
