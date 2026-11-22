package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseModel = (*customCourseModel)(nil)

// CoursePageFilter 课程分页查询条件（status/free/courseType 为 0 表示不限制）。
type CoursePageFilter struct {
	Keyword      string
	Status       int64
	Free         int64 // 0 不限, 1 免费, 2 付费
	CourseType   int64
	FirstCateId  int64
	SecondCateId int64
	ThirdCateId  int64
	BeginTime    string
	EndTime      string
}

type (
	// CourseModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseModel.
	CourseModel interface {
		courseModel
		PageQuery(ctx context.Context, f CoursePageFilter, offset, limit int64) ([]*Course, int64, error)
		ListByIds(ctx context.Context, ids []int64) ([]*Course, error)
		ListByThirdCateIds(ctx context.Context, ids []int64) ([]*Course, error)
		FindByNameLike(ctx context.Context, name string) ([]*Course, error)
		CountByThirdCateId(ctx context.Context, thirdCateId int64) (int64, error)
		InsertTx(ctx context.Context, session sqlx.Session, data *Course) (sql.Result, error)
		UpdateTx(ctx context.Context, session sqlx.Session, data *Course) error
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

func (m *customCourseModel) PageQuery(ctx context.Context, f CoursePageFilter, offset, limit int64) ([]*Course, int64, error) {
	conds := []string{"`deleted` = 0"}
	args := []any{}
	if f.Keyword != "" {
		conds = append(conds, "`name` like ?")
		args = append(args, "%"+f.Keyword+"%")
	}
	if f.Status != 0 {
		conds = append(conds, "`status` = ?")
		args = append(args, f.Status)
	}
	switch f.Free {
	case 1:
		conds = append(conds, "`free` = 1")
	case 2:
		conds = append(conds, "`free` = 0")
	}
	if f.CourseType != 0 {
		conds = append(conds, "`course_type` = ?")
		args = append(args, f.CourseType)
	}
	if f.FirstCateId != 0 {
		conds = append(conds, "`first_cate_id` = ?")
		args = append(args, f.FirstCateId)
	}
	if f.SecondCateId != 0 {
		conds = append(conds, "`second_cate_id` = ?")
		args = append(args, f.SecondCateId)
	}
	if f.ThirdCateId != 0 {
		conds = append(conds, "`third_cate_id` = ?")
		args = append(args, f.ThirdCateId)
	}
	if f.BeginTime != "" {
		conds = append(conds, "`create_time` >= ?")
		args = append(args, f.BeginTime)
	}
	if f.EndTime != "" {
		conds = append(conds, "`create_time` <= ?")
		args = append(args, f.EndTime)
	}
	where := strings.Join(conds, " and ")

	var total int64
	countQuery := fmt.Sprintf("select count(*) from %s where %s", m.table, where)
	if err := m.QueryRowNoCacheCtx(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	var list []*Course
	query := fmt.Sprintf("select * from %s where %s order by `create_time` desc limit ?, ?", m.table, where)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, append(args, offset, limit)...); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (m *customCourseModel) ListByIds(ctx context.Context, ids []int64) ([]*Course, error) {
	if len(ids) == 0 {
		return []*Course{}, nil
	}
	var list []*Course
	ph := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	query := fmt.Sprintf("select * from %s where `deleted` = 0 and `id` in (%s)", m.table, ph)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCourseModel) ListByThirdCateIds(ctx context.Context, ids []int64) ([]*Course, error) {
	if len(ids) == 0 {
		return []*Course{}, nil
	}
	var list []*Course
	ph := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	query := fmt.Sprintf("select * from %s where `deleted` = 0 and `third_cate_id` in (%s)", m.table, ph)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCourseModel) FindByNameLike(ctx context.Context, name string) ([]*Course, error) {
	var list []*Course
	query := fmt.Sprintf("select * from %s where `deleted` = 0 and `name` like ?", m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, "%"+name+"%"); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCourseModel) CountByThirdCateId(ctx context.Context, thirdCateId int64) (int64, error) {
	var total int64
	query := fmt.Sprintf("select count(*) from %s where `deleted` = 0 and `third_cate_id` = ?", m.table)
	if err := m.QueryRowNoCacheCtx(ctx, &total, query, thirdCateId); err != nil {
		return 0, err
	}
	return total, nil
}

// InsertTx 事务内插入课程行：SQL 与缓存键失效语义同 Insert（生成模型），写入走传入的事务会话，
// 供跨模型事务（如课程上架发布）复用。
func (m *customCourseModel) InsertTx(ctx context.Context, session sqlx.Session, data *Course) (sql.Result, error) {
	cc := m.WithSession(session)
	courseIdKey := fmt.Sprintf("%s%v", cacheCourseIdPrefix, data.Id)
	return cc.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values %s", m.table, courseRowsExpectAutoSet, rowPlaceholders(insertColumns(courseRowsExpectAutoSet)))
		return conn.ExecCtx(ctx, query,
			data.Id, data.Name, data.CourseType, data.CoverUrl, data.FirstCateId, data.SecondCateId,
			data.ThirdCateId, data.Free, data.Price, data.TemplateType, data.TemplateUrl, data.Status,
			data.PurchaseStartTime, data.PurchaseEndTime, data.Step, data.Score, data.MediaDuration,
			data.ValidDuration, data.SectionNum, data.DepId, data.PublishTimes, data.PublishTime,
			data.Creater, data.Updater, data.Deleted)
	}, courseIdKey)
}

// UpdateTx 事务内更新课程行：SQL 与缓存键失效语义同 Update（生成模型），写入走传入的事务会话，
// 供跨模型事务（如课程上架发布）复用。
func (m *customCourseModel) UpdateTx(ctx context.Context, session sqlx.Session, data *Course) error {
	cc := m.WithSession(session)
	courseIdKey := fmt.Sprintf("%s%v", cacheCourseIdPrefix, data.Id)
	_, err := cc.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, courseRowsWithPlaceHolder)
		return conn.ExecCtx(ctx, query,
			data.Name, data.CourseType, data.CoverUrl, data.FirstCateId, data.SecondCateId, data.ThirdCateId,
			data.Free, data.Price, data.TemplateType, data.TemplateUrl, data.Status, data.PurchaseStartTime,
			data.PurchaseEndTime, data.Step, data.Score, data.MediaDuration, data.ValidDuration, data.SectionNum,
			data.DepId, data.PublishTimes, data.PublishTime, data.Creater, data.Updater, data.Deleted, data.Id)
	}, courseIdKey)
	return err
}
