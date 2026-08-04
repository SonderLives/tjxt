package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// Course 课程正式表（course 表）
type Course struct {
	Id               int64
	Name             string
	CourseType       int64
	CoverUrl         string
	FirstCateId      int64
	SecondCateId     int64
	ThirdCateId      int64
	Free             int64
	Price            int64
	TemplateType     int64
	TemplateUrl      string
	Status           int64
	PurchaseStartTime sql.NullTime
	PurchaseEndTime  time.Time
	Step             int64
	Score            sql.NullInt64
	MediaDuration    sql.NullInt64
	ValidDuration    int64
	SectionNum       sql.NullInt64
	DepId            int64
	PublishTimes     sql.NullInt64
	PublishTime      sql.NullTime
	CreateTime       time.Time
	UpdateTime       time.Time
	Creater          int64
	Updater          int64
	Deleted          int64
}

// CourseModel 课程正式数据访问
type CourseModel struct {
	conn  sqlx.SqlConn
	table string
}

// NewCourseModel 创建课程正式数据访问对象
func NewCourseModel(conn sqlx.SqlConn) *CourseModel {
	return &CourseModel{conn: conn, table: "course"}
}

const courseColumns = `id, name, course_type, cover_url, first_cate_id, second_cate_id, third_cate_id,
	free, price, template_type, template_url, status, purchase_start_time, purchase_end_time, step,
	score, media_duration, valid_duration, section_num, dep_id, publish_times, publish_time,
	create_time, update_time, creater, updater, deleted`

const courseScanColumns = `id, name, course_type, cover_url, first_cate_id, second_cate_id, third_cate_id,
	free, price, template_type, template_url, status, purchase_start_time, purchase_end_time, step,
	score, media_duration, valid_duration, section_num, dep_id, publish_times, publish_time,
	create_time, update_time, COALESCE(creater, 0), COALESCE(updater, 0), deleted`

// FindById 按主键查询课程。
func (m *CourseModel) FindById(ctx context.Context, id int64) (*Course, error) {
	var c Course
	err := m.conn.QueryRowCtx(ctx, &c, fmt.Sprintf(
		`SELECT %s FROM %s WHERE id = ? AND deleted = 0`, courseScanColumns, m.table), id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &c, err
}

// FindByIds 按主键列表批量查询。
func (m *CourseModel) FindByIds(ctx context.Context, ids []int64) ([]*Course, error) {
	if len(ids) == 0 {
		return []*Course{}, nil
	}
	ph := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		ph = append(ph, "?")
		args = append(args, id)
	}
	var rows []*Course
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE deleted = 0 AND id IN (%s)`,
		courseScanColumns, m.table, strings.Join(ph, ",")), args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// FindByName 按名称模糊查询。
func (m *CourseModel) FindByName(ctx context.Context, name string) ([]*Course, error) {
	var rows []*Course
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE deleted = 0 AND name LIKE ?`,
		courseScanColumns, m.table), "%"+name+"%")
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// CountSameName 统计正式表中同名课程数量。
func (m *CourseModel) CountSameName(ctx context.Context, name string) (int64, error) {
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, fmt.Sprintf(
		`SELECT COUNT(*) FROM %s WHERE deleted = 0 AND name = ?`, m.table), name)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// Insert 新增课程。
func (m *CourseModel) Insert(ctx context.Context, c *Course) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`INSERT INTO %s (id, name, course_type, cover_url, first_cate_id, second_cate_id, third_cate_id,
			free, price, template_type, template_url, status, purchase_start_time, purchase_end_time,
			step, score, media_duration, valid_duration, section_num, dep_id, publish_times, publish_time,
			create_time, update_time, creater, updater)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW(), ?, ?)`,
		m.table),
		c.Id, c.Name, c.CourseType, c.CoverUrl, c.FirstCateId, c.SecondCateId, c.ThirdCateId,
		c.Free, c.Price, c.TemplateType, c.TemplateUrl, c.Status, c.PurchaseStartTime, c.PurchaseEndTime,
		c.Step, c.Score, c.MediaDuration, c.ValidDuration, c.SectionNum, c.DepId, c.PublishTimes, c.PublishTime,
		c.Creater, c.Updater)
	return err
}

// UpdateById 更新课程指定字段，支持可空列透传。
func (m *CourseModel) UpdateById(ctx context.Context, c *Course, fields ...string) error {
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
		case "status":
			sets = append(sets, "status = ?")
			args = append(args, c.Status)
		case "media_duration":
			sets = append(sets, "media_duration = ?")
			args = append(args, c.MediaDuration)
		case "valid_duration":
			sets = append(sets, "valid_duration = ?")
			args = append(args, c.ValidDuration)
		case "purchase_end_time":
			sets = append(sets, "purchase_end_time = ?")
			args = append(args, c.PurchaseEndTime)
		case "purchase_start_time":
			sets = append(sets, "purchase_start_time = ?")
			args = append(args, c.PurchaseStartTime)
		case "publish_time":
			sets = append(sets, "publish_time = ?")
			args = append(args, c.PublishTime)
		case "publish_times":
			sets = append(sets, "publish_times = ?")
			args = append(args, c.PublishTimes)
		case "score":
			sets = append(sets, "score = ?")
			args = append(args, c.Score)
		case "section_num":
			sets = append(sets, "section_num = ?")
			args = append(args, c.SectionNum)
		case "step":
			sets = append(sets, "step = ?")
			args = append(args, c.Step)
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

// ListByCates 按课程分类筛选（用于分类删除联动与统计）。
func (m *CourseModel) ListByCates(ctx context.Context, categoryId int64, level int64) ([]*Course, error) {
	col := "first_cate_id"
	switch level {
	case 2:
		col = "second_cate_id"
	case 3:
		col = "third_cate_id"
	}
	var rows []*Course
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE deleted = 0 AND %s = ?`,
		courseScanColumns, m.table, col), categoryId)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// CountByCategoryId 统计某分类下课程数量（任一分类级别匹配即计入）。
func (m *CourseModel) CountByCategoryId(ctx context.Context, categoryId int64) (int64, error) {
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, fmt.Sprintf(
		`SELECT COUNT(*) FROM %s WHERE deleted = 0 AND
		 (first_cate_id = ? OR second_cate_id = ? OR third_cate_id = ?)`,
		m.table), categoryId, categoryId, categoryId)
	if err != nil {
		return 0, err
	}
	return total, nil
}

// PageCond 课程分页查询条件（正式表）
type CoursePageCond struct {
	Keyword      string
	FirstCateId  int64
	SecondCateId int64
	ThirdCateId  int64
	CourseType   int64
	Status       int64
	BeginTime    string
	EndTime      string
	OrderBy      string
	Offset       int64
	Limit        int64
}

// Page 分页查询课程（正式表）。
func (m *CourseModel) Page(ctx context.Context, cond *CoursePageCond) ([]*Course, int64, error) {
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
	var rows []*Course
	err = m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE %s ORDER BY %s LIMIT ? OFFSET ?`,
		courseScanColumns, m.table, whereSQL, order), args...)
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

// ListAll 查询全部课程。
func (m *CourseModel) ListAll(ctx context.Context) ([]*Course, error) {
	var rows []*Course
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE deleted = 0`, courseScanColumns, m.table))
	if err != nil {
		return nil, err
	}
	return rows, nil
}
