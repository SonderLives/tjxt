package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"tjxt/pkg/utils/idgen"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// 业务状态常量
const (
	LessonStatusNotStart int64 = 0 // 未开始
	LessonStatusInLearn  int64 = 1 // 学习中
	LessonStatusDone     int64 = 2 // 完成
	LessonStatusExpired  int64 = 3 // 失效

	PlanStatusNone   int64 = 0 // 无计划
	PlanStatusInPlan int64 = 1 // 计划中
)

var _ LearningLessonModel = (*customLearningLessonModel)(nil)

type (
	// LearningLessonModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLearningLessonModel.
	LearningLessonModel interface {
		learningLessonModel

		// GrantCourses 为用户开通课程（幂等，存在则不动）
		GrantCourses(ctx context.Context, userID int64, courseIDs []int64) error
		// RevokeCourses 撤销课程（status 置 3 = 失效）
		RevokeCourses(ctx context.Context, userID int64, courseIDs []int64) error
		// FindByUserCourse 查询用户在指定课程的学习记录
		FindByUserCourse(ctx context.Context, userID, courseID int64) (*LearningLesson, error)
		// ListByUser 分页查询用户的学习记录（不含已失效）
		ListByUser(ctx context.Context, userID, pageNo, pageSize int64, asc bool) ([]*LearningLesson, int64, error)
		// ListPlansByUser 只看"已设置计划"的分页
		ListPlansByUser(ctx context.Context, userID, pageNo, pageSize int64, asc bool) ([]*LearningLesson, int64, error)
		// FindLatestLearnedByUser 最近学习的一条（用于"我正在学"）
		FindLatestLearnedByUser(ctx context.Context, userID int64) (*LearningLesson, error)
		// CountByCourse 该课程的注册人数
		CountByCourse(ctx context.Context, courseID int64) (int64, error)
		// UpdatePlan 设置/更新学习计划（每周章节数）
		UpdatePlan(ctx context.Context, userID, courseID, weekFreq int64) error
		// RemoveLesson 删除该用户对这门课的学习记录（status 置 3）
		RemoveLesson(ctx context.Context, userID, courseID int64) error
		// UpdateLatestLearn 提交学习记录时，更新最新学习进度
		UpdateLatestLearn(ctx context.Context, lessonID, sectionID, moment, duration int64) error
		// IncrLearnedSections +1 learned_sections，用于学习记录被确认完成时
		IncrLearnedSections(ctx context.Context, lessonID int64) error
	}

	customLearningLessonModel struct {
		*defaultLearningLessonModel
	}
)

// NewLearningLessonModel returns a model for the database table.
func NewLearningLessonModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LearningLessonModel {
	return &customLearningLessonModel{
		defaultLearningLessonModel: newLearningLessonModel(conn, c, opts...),
	}
}

const insertCourseSQL = `INSERT INTO %s
	(id, user_id, course_id, status, week_freq, plan_status, learned_sections,
	 latest_section_id, latest_learn_time, create_time, expire_time, update_time)
VALUES (?, ?, ?, 0, NULL, 0, 0, NULL, NULL, NOW(), NULL, NOW())
ON DUPLICATE KEY UPDATE update_time = NOW()`

func (m *customLearningLessonModel) GrantCourses(ctx context.Context, userID int64, courseIDs []int64) error {
	for _, cid := range courseIDs {
		if _, err := m.ExecNoCacheCtx(ctx, fmt.Sprintf(insertCourseSQL, m.table), idgen.NextID(), userID, cid); err != nil {
			return err
		}
	}
	return nil
}

func (m *customLearningLessonModel) RevokeCourses(ctx context.Context, userID int64, courseIDs []int64) error {
	for _, cid := range courseIDs {
		if _, err := m.ExecNoCacheCtx(ctx,
			fmt.Sprintf("UPDATE %s SET status = ?, plan_status = ?, update_time = NOW() WHERE user_id = ? AND course_id = ?", m.table),
			LessonStatusExpired, PlanStatusNone, userID, cid); err != nil {
			return err
		}
	}
	return nil
}

func (m *customLearningLessonModel) FindByUserCourse(ctx context.Context, userID, courseID int64) (*LearningLesson, error) {
	return m.FindOneByUserIdCourseId(ctx, userID, courseID)
}

func (m *customLearningLessonModel) ListByUser(ctx context.Context, userID, pageNo, pageSize int64, asc bool) ([]*LearningLesson, int64, error) {
	return m.listBy(ctx, "user_id = ? AND status <> ?", []any{userID, LessonStatusExpired}, pageNo, pageSize, asc, "create_time")
}

func (m *customLearningLessonModel) ListPlansByUser(ctx context.Context, userID, pageNo, pageSize int64, asc bool) ([]*LearningLesson, int64, error) {
	return m.listBy(ctx,
		"user_id = ? AND status <> ? AND plan_status = ?",
		[]any{userID, LessonStatusExpired, PlanStatusInPlan},
		pageNo, pageSize, asc, "create_time")
}

func (m *customLearningLessonModel) listBy(ctx context.Context, cond string, args []any, pageNo, pageSize int64, asc bool, orderBy string) ([]*LearningLesson, int64, error) {
	var total int64
	if err := m.QueryRowNoCacheCtx(ctx, &total,
		fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE %s", m.table, cond), args...); err != nil {
		return nil, 0, err
	}
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	order := "DESC"
	if asc {
		order = "ASC"
	}
	offset := (pageNo - 1) * pageSize
	listArgs := append(append([]any{}, args...), pageSize, offset)
	var list []*LearningLesson
	err := m.QueryRowsNoCacheCtx(ctx, &list,
		fmt.Sprintf("SELECT %s FROM %s WHERE %s ORDER BY %s %s LIMIT ? OFFSET ?", learningLessonRows, m.table, cond, orderBy, order),
		listArgs...)
	return list, total, err
}

func (m *customLearningLessonModel) FindLatestLearnedByUser(ctx context.Context, userID int64) (*LearningLesson, error) {
	var resp LearningLesson
	err := m.QueryRowNoCacheCtx(ctx, &resp,
		fmt.Sprintf("SELECT %s FROM %s WHERE user_id = ? AND status <> ? AND latest_learn_time IS NOT NULL ORDER BY latest_learn_time DESC LIMIT 1", learningLessonRows, m.table),
		userID, LessonStatusExpired)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (m *customLearningLessonModel) CountByCourse(ctx context.Context, courseID int64) (int64, error) {
	var count int64
	err := m.QueryRowNoCacheCtx(ctx, &count,
		fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE course_id = ? AND status <> ?", m.table),
		courseID, LessonStatusExpired)
	return count, err
}

func (m *customLearningLessonModel) UpdatePlan(ctx context.Context, userID, courseID, weekFreq int64) error {
	res, err := m.ExecNoCacheCtx(ctx,
		fmt.Sprintf("UPDATE %s SET week_freq = ?, plan_status = ?, update_time = NOW() WHERE user_id = ? AND course_id = ? AND status <> ?", m.table),
		sql.NullInt64{Int64: weekFreq, Valid: weekFreq > 0}, PlanStatusInPlan, userID, courseID, LessonStatusExpired)
	if err != nil {
		return err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (m *customLearningLessonModel) RemoveLesson(ctx context.Context, userID, courseID int64) error {
	_, err := m.ExecNoCacheCtx(ctx,
		fmt.Sprintf("UPDATE %s SET status = ?, update_time = NOW() WHERE user_id = ? AND course_id = ?", m.table),
		LessonStatusExpired, userID, courseID)
	return err
}

func (m *customLearningLessonModel) UpdateLatestLearn(ctx context.Context, lessonID, sectionID, moment, duration int64) error {
	_, err := m.ExecNoCacheCtx(ctx,
		fmt.Sprintf(`UPDATE %s SET
			latest_section_id = ?,
			latest_learn_time = ?,
			status = IF(status = ?, ?, status),
			update_time = NOW()
		WHERE id = ?`, m.table),
		sql.NullInt64{Int64: sectionID, Valid: sectionID > 0},
		time.Now(),
		LessonStatusNotStart, LessonStatusInLearn,
		lessonID)
	return err
}

func (m *customLearningLessonModel) IncrLearnedSections(ctx context.Context, lessonID int64) error {
	_, err := m.ExecNoCacheCtx(ctx,
		fmt.Sprintf("UPDATE %s SET learned_sections = learned_sections + 1, update_time = NOW() WHERE id = ?", m.table),
		lessonID)
	return err
}