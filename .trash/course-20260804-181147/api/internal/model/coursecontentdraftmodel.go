package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseContentDraftModel = (*customCourseContentDraftModel)(nil)

type (
	// CourseContentDraftModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseContentDraftModel.
	CourseContentDraftModel interface {
		courseContentDraftModel
		FindById(ctx context.Context, id int64) (*CourseContentDraft, error)
		UpdateById(ctx context.Context, data *CourseContentDraft) error
	}

	customCourseContentDraftModel struct {
		*defaultCourseContentDraftModel
	}
)

// NewCourseContentDraftModel returns a model for the database table.
func NewCourseContentDraftModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseContentDraftModel {
	return &customCourseContentDraftModel{
		defaultCourseContentDraftModel: newCourseContentDraftModel(conn, c, opts...),
	}
}

// FindById 根据ID查询
func (m *customCourseContentDraftModel) FindById(ctx context.Context, id int64) (*CourseContentDraft, error) {
	return m.FindOne(ctx, id)
}

// UpdateById 更新课程内容草稿
func (m *customCourseContentDraftModel) UpdateById(ctx context.Context, data *CourseContentDraft) error {
	courseContentDraftIdKey := fmt.Sprintf("%s%v", cacheCourseContentDraftIdPrefix, data.Id)
	query := fmt.Sprintf("update %s set `course_introduce` = ?, `course_detail` = ?, `use_people` = ?, `updater` = ? where `id` = ?", m.table)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, data.CourseIntroduce, data.CourseDetail, data.UsePeople, data.Updater, data.Id)
	}, courseContentDraftIdKey)
	return err
}
