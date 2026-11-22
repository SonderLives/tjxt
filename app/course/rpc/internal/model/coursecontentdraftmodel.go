package model

import (
	"context"
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
		Upsert(ctx context.Context, data *CourseContentDraft) error
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

// Upsert 课程草稿内容按课程 id 主键替换插入或更新。
func (m *customCourseContentDraftModel) Upsert(ctx context.Context, data *CourseContentDraft) error {
	query := fmt.Sprintf("replace into %s (`id`,`course_introduce`,`use_people`,`course_detail`,`dep_id`,`creater`,`updater`,`deleted`,`create_time`,`update_time`) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query,
		data.Id, data.CourseIntroduce, data.UsePeople, data.CourseDetail,
		data.DepId, data.Creater, data.Updater, data.Deleted, data.CreateTime, data.UpdateTime)
	return err
}
