package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	courseScanColumns         = "`id`, `name`, `course_type`, `cover_url`, `first_cate_id`, `second_cate_id`, `third_cate_id`, `free`, `price`, `template_type`, `template_url`, `status`, `purchase_start_time`, `purchase_end_time`, `step`, `score`, `media_duration`, `valid_duration`, `section_num`, `dep_id`, `publish_times`, `publish_time`, `create_time`, `update_time`, COALESCE(`creater`, 0), COALESCE(`updater`, 0), `deleted`"
)

var _ CourseModel = (*customCourseModel)(nil)

type (
	// CourseModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseModel.
	CourseModel interface {
		courseModel
		ListAll(ctx context.Context) ([]*Course, error)
		FindById(ctx context.Context, id int64) (*Course, error)
		FindByIds(ctx context.Context, ids []int64) ([]*Course, error)
		CountByCategoryId(ctx context.Context, categoryId int64) (int64, error)
		CountSameName(ctx context.Context, name string) (int64, error)
		FindByName(ctx context.Context, name string) ([]*Course, error)
		Page(ctx context.Context, cond *CoursePageCond) ([]*Course, int64, error)
	}
	customCourseModel struct {
		*defaultCourseModel
	}
)

// NewCourseModel returns a model for the database table.
func NewCourseModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseModel {
	return &customCourseModel{
		defaultCourseModel: newCourseModel(conn, c, opts...),
	}
}

// ListAll 查询所有课程
func (m *customCourseModel) ListAll(ctx context.Context) ([]*Course, error) {
	var resp []*Course
	query := fmt.Sprintf("select %s from %s where `deleted` = 0 order by `id` desc", courseScanColumns, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// FindById 根据ID查询
func (m *customCourseModel) FindById(ctx context.Context, id int64) (*Course, error) {
	return m.FindOne(ctx, id)
}

// FindByIds 批量查询
func (m *customCourseModel) FindByIds(ctx context.Context, ids []int64) ([]*Course, error) {
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
	query := fmt.Sprintf("select %s from %s where `deleted` = 0 and `id` in (%s)", courseScanColumns, m.table, strings.Join(ph, ","))
	err := m.QueryRowsNoCacheCtx(ctx, &rows, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// CountByCategoryId 统计分类下的课程数量
func (m *customCourseModel) CountByCategoryId(ctx context.Context, categoryId int64) (int64, error) {
	query := fmt.Sprintf("select count(1) from %s where (`first_cate_id` = ? or `second_cate_id` = ? or `third_cate_id` = ?) and `deleted` = 0", m.table)
	var count int64
	err := m.QueryRowNoCacheCtx(ctx, &count, query, categoryId, categoryId, categoryId)
	return count, err
}

// CountSameName 统计同名课程
func (m *customCourseModel) CountSameName(ctx context.Context, name string) (int64, error) {
	query := fmt.Sprintf("select count(1) from %s where `name` = ? and `deleted` = 0", m.table)
	var count int64
	err := m.QueryRowNoCacheCtx(ctx, &count, query, name)
	return count, err
}

// FindByName 模糊查询
func (m *customCourseModel) FindByName(ctx context.Context, name string) ([]*Course, error) {
	var rows []*Course
	query := fmt.Sprintf("select %s from %s where `name` like ? and `deleted` = 0", courseScanColumns, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &rows, query, "%"+name+"%")
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// CoursePageCond 分页查询条件
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

// Page 分页查询
func (m *customCourseModel) Page(ctx context.Context, cond *CoursePageCond) ([]*Course, int64, error) {
	where := "`deleted` = 0"
	args := make([]any, 0)
	if cond.Keyword != "" {
		where += " and `name` like ?"
		args = append(args, "%"+cond.Keyword+"%")
	}
	if cond.FirstCateId > 0 {
		where += " and `first_cate_id` = ?"
		args = append(args, cond.FirstCateId)
	}
	if cond.SecondCateId > 0 {
		where += " and `second_cate_id` = ?"
		args = append(args, cond.SecondCateId)
	}
	if cond.ThirdCateId > 0 {
		where += " and `third_cate_id` = ?"
		args = append(args, cond.ThirdCateId)
	}
	if cond.CourseType > 0 {
		where += " and `course_type` = ?"
		args = append(args, cond.CourseType)
	}
	if cond.Status > 0 {
		where += " and `status` = ?"
		args = append(args, cond.Status)
	}
	if cond.BeginTime != "" {
		where += " and `create_time` >= ?"
		args = append(args, cond.BeginTime)
	}
	if cond.EndTime != "" {
		where += " and `create_time` <= ?"
		args = append(args, cond.EndTime)
	}
	// 总数
	countQuery := fmt.Sprintf("select count(1) from %s where %s", m.table, where)
	var total int64
	err := m.QueryRowNoCacheCtx(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []*Course{}, 0, nil
	}
	// 分页数据
	orderBy := cond.OrderBy
	if orderBy == "" {
		orderBy = "`id` desc"
	}
	dataQuery := fmt.Sprintf("select %s from %s where %s order by %s limit ? offset ?", courseScanColumns, m.table, where, orderBy)
	args = append(args, cond.Limit, cond.Offset)
	var rows []*Course
	err = m.QueryRowsNoCacheCtx(ctx, &rows, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}
