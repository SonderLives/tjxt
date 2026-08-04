package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseDraftModel = (*customCourseDraftModel)(nil)

type (
	// CourseDraftModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseDraftModel.
	CourseDraftModel interface {
		courseDraftModel
		ListAll(ctx context.Context) ([]*CourseDraft, error)
		FindById(ctx context.Context, id int64) (*CourseDraft, error)
		CountByName(ctx context.Context, name string, excludeId int64) (int64, error)
		UpdateById(ctx context.Context, data *CourseDraft, fields ...string) error
		DeleteById(ctx context.Context, id int64) error
		Page(ctx context.Context, cond *CourseDraftPageCond) ([]*CourseDraft, int64, error)
	}

	customCourseDraftModel struct {
		*defaultCourseDraftModel
	}
)

// NewCourseDraftModel returns a model for the database table.
func NewCourseDraftModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseDraftModel {
	return &customCourseDraftModel{
		defaultCourseDraftModel: newCourseDraftModel(conn, c, opts...),
	}
}

// ListAll 查询所有草稿
func (m *customCourseDraftModel) ListAll(ctx context.Context) ([]*CourseDraft, error) {
	var resp []*CourseDraft
	query := fmt.Sprintf("select %s from %s where `deleted` = 0 order by `id` desc", courseDraftRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// FindById 根据ID查询
func (m *customCourseDraftModel) FindById(ctx context.Context, id int64) (*CourseDraft, error) {
	return m.FindOne(ctx, id)
}

// CountByName 统计同名草稿（排除指定ID）
func (m *customCourseDraftModel) CountByName(ctx context.Context, name string, excludeId int64) (int64, error) {
	var count int64
	var err error
	if excludeId > 0 {
		query := fmt.Sprintf("select count(1) from %s where `name` = ? and `id` != ? and `deleted` = 0", m.table)
		err = m.QueryRowNoCacheCtx(ctx, &count, query, name, excludeId)
	} else {
		query := fmt.Sprintf("select count(1) from %s where `name` = ? and `deleted` = 0", m.table)
		err = m.QueryRowNoCacheCtx(ctx, &count, query, name)
	}
	return count, err
}

// UpdateById 更新指定字段
func (m *customCourseDraftModel) UpdateById(ctx context.Context, data *CourseDraft, fields ...string) error {
	if len(fields) == 0 {
		return m.Update(ctx, data)
	}
	courseDraftIdKey := fmt.Sprintf("%s%v", cacheCourseDraftIdPrefix, data.Id)
	setParts := make([]string, 0, len(fields))
	args := make([]any, 0, len(fields)+1)
	for _, f := range fields {
		setParts = append(setParts, fmt.Sprintf("`%s` = ?", f))
		switch f {
		case "name":
			args = append(args, data.Name)
		case "cover_url":
			args = append(args, data.CoverUrl)
		case "price":
			args = append(args, data.Price)
		case "valid_duration":
			args = append(args, data.ValidDuration)
		case "free":
			args = append(args, data.Free)
		case "status":
			args = append(args, data.Status)
		case "first_cate_id":
			args = append(args, data.FirstCateId)
		case "second_cate_id":
			args = append(args, data.SecondCateId)
		case "third_cate_id":
			args = append(args, data.ThirdCateId)
		case "purchase_end_time":
			args = append(args, data.PurchaseEndTime)
		case "step":
			args = append(args, data.Step)
		case "section_num":
			args = append(args, data.SectionNum)
		}
	}
	args = append(args, data.Id)
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(setParts, ", "))
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, courseDraftIdKey)
	return err
}

// DeleteById 根据ID删除
func (m *customCourseDraftModel) DeleteById(ctx context.Context, id int64) error {
	return m.Delete(ctx, id)
}

// CourseDraftPageCond 草稿分页查询条件
type CourseDraftPageCond struct {
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
func (m *customCourseDraftModel) Page(ctx context.Context, cond *CourseDraftPageCond) ([]*CourseDraft, int64, error) {
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
		return []*CourseDraft{}, 0, nil
	}
	orderBy := cond.OrderBy
	if orderBy == "" {
		orderBy = "`id` desc"
	}
	dataQuery := fmt.Sprintf("select %s from %s where %s order by %s limit ? offset ?", courseDraftRows, m.table, where, orderBy)
	args = append(args, cond.Limit, cond.Offset)
	var rows []*CourseDraft
	err = m.QueryRowsNoCacheCtx(ctx, &rows, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}
